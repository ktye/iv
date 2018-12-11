package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "←",
		Domain:  MonadicOp(nil),
		doc:     "assign, variable assignment, specification, copula",
		derived: assign,
	})
}

func assign(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		v, ok := f.(apl.Identifier)
		if ok == false {
			return nil, fmt.Errorf("cannot assign to %T", f)
		}
		if L != nil {
			return nil, fmt.Errorf("assign cannot be called dyadically")
		}

		return R, a.Assign(string(v), R)
	}
	return function(derived)
}

// Index returns the indexes within the array A for the given index specification.
// The result may have a smaller size and shape as the input array.
// The indexes in the spec are origin dependend, but in IndexArray are always origin 0.
func Index(a *apl.Apl, spec apl.IdxSpec, shape []int) (apl.IndexArray, error) {
	intspec, err := spec2ints(a, spec, shape)
	if err != nil {
		return apl.IndexArray{}, nil
	}

	// Initially the rank is the same as spec.
	// Single element axis will be reduced later.
	res := apl.IndexArray{Dims: make([]int, len(intspec))}
	for i := range intspec {
		res.Dims[i] = len(intspec[i])
	}

	res.Ints = make([]int, apl.ArraySize(res))
	ic, src := apl.NewIdxConverter(shape)
	dst := make([]int, len(res.Dims))
	for i := range res.Ints {
		for k, n := range dst {
			src[k] = intspec[k][n]
		}
		res.Ints[i] = ic.Index(src)
		apl.IncArrayIndex(dst, res.Dims)
	}

	// Reduce rank by collapsing single element axis.
	rs := make([]int, 0, len(res.Dims))
	for _, v := range res.Dims {
		if v != 1 {
			rs = append(rs, v)
		}
	}
	res.Dims = rs
	return res, nil
}

// Spec2ints converts an index specification to [][]int for the given shape.
// spec is origin dependent, the result has always origin 0.
func spec2ints(a *apl.Apl, spec apl.IdxSpec, shape []int) ([][]int, error) {
	if len(spec) != len(shape) {
		return nil, fmt.Errorf("indexing: array and index specification have different rank")
	}

	to := ToIndexArray(nil)
	idx := make([][]int, len(shape))
	for i := range spec {
		v, ok := to.To(a, spec[i])
		if ok == false {
			return nil, fmt.Errorf("index specification for axis %d is illegal: %T", i+1, spec[i])
		}

		// Empty axis are expanded to all elements of the axis.
		if _, ok := v.(apl.EmptyArray); ok {
			idx[i] = make([]int, shape[i])
			for k := range idx[i] {
				idx[i][k] = k
			}
			continue
		}
		ia := v.(apl.IndexArray)
		idx[i] = make([]int, len(ia.Ints))
		for k := range ia.Ints {
			if n := ia.Ints[k] - a.Origin; n < 0 || n >= shape[i] {
				return nil, fmt.Errorf("index specification for axis %d is out of range", i+1)
			} else {
				idx[i][k] = n
			}
		}
	}
	return idx, nil
}
