package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "⌺",
		Domain:  DyadicOp(Split(Function(nil), ToIndexArray(nil))),
		doc:     "stencil",
		derived: stencil,
	})
}

// stencil: f is a function, RO an index array and R an array.
func stencil(a *apl.Apl, f, RO apl.Value) apl.Function {
	derived := func(a *apl.Apl, dummyL, R apl.Value) (apl.Value, error) {
		// Stencil derived function must be called monadically.
		if dummyL != nil {
			return nil, fmt.Errorf("stencil: derived function cannot be called dyadically")
		}

		// f is a Function
		f := f.(apl.Function)

		// RO is a 2 x rank R index array with rows that indicate stencil shape and movement.
		var ai apl.IndexArray
		if _, ok := RO.(apl.EmptyArray); ok {
			ai = apl.IndexArray{}
		} else {
			ai = RO.(apl.IndexArray)
		}
		is := ai.Shape()
		if len(is) > 2 {
			return nil, fmt.Errorf("stencil: rank of RO is > 2: %d", len(is))
		}

		// R is an array.
		ar, ok := R.(apl.Array)
		if ok == false {
			return nil, fmt.Errorf("stencil: right argument must be an array: %T", R)
		}
		rs := ar.Shape()
		if len(is) > 0 && is[len(is)-1] > len(rs) {
			return nil, fmt.Errorf("stencil: shape of RO is too large: %v, max: [2 %d]", is, len(rs))
		}

		// Default RO matrix.
		def := apl.IndexArray{Dims: []int{2, len(rs)}}
		def.Ints = make([]int, 2*len(rs))
		for i := range def.Ints {
			def.Ints[i] = 1
		}
		// Overwrite default matrix with given values and swap.
		if len(is) == 1 {
			is = append([]int{1}, is[0])
		}
		ic, idx := apl.NewIdxConverter(is)
		for i := 0; i < apl.ArraySize(ai); i++ {
			v, err := ai.At(i)
			if err != nil {
				return nil, err
			}
			n := ic.Index(idx)
			def.Ints[n] = int(v.(apl.Index))
			apl.IncArrayIndex(idx, is)
		}
		ai = def
		is = ai.Shape()

		// The result has the same shape as R.
		res := apl.GeneralArray{Dims: apl.CopyShape(ar), Values: make([]apl.Value, apl.ArraySize(ar))}

		// The temporary array has the requested stencil shape, the first row of RO.
		tmp := apl.GeneralArray{Dims: ai.Ints[:len(ai.Ints)/2]}
		tmp.Values = make([]apl.Value, apl.ArraySize(tmp))
		if apl.ArraySize(tmp) == 0 {
			return nil, fmt.Errorf("stencil: stencil size is 0")
		}

		// lvec is the left vector for the stencil function application,
		// which indicates the number of fill elements per axis.
		lvec := apl.IndexArray{Dims: []int{len(rs)}}
		vec := make([]int, len(rs))
		lvec.Ints = vec

		// Apply the stencil to all elements of R.
		ic, idx = apl.NewIdxConverter(rs)
		sdx := make([]int, len(tmp.Dims))
		dst := make([]int, len(idx))
		for i := 0; i < len(res.Values); i++ {

			// Center the stencil on idx.
			ic.Indexes(i, idx)
			copy(dst, idx)
			for k := range vec {
				vec[k] = 0
			}
			for k := range tmp.Values {
				out := false
				for d := range tmp.Dims {
					dst[d] = idx[d] + sdx[d] - tmp.Dims[d]/2
					if dst[d] < 0 || dst[d] >= res.Dims[d] {
						out = true
					}
					if v := dst[d]; v < 0 && -v > vec[d] {
						// Positive lvec value indicates the padding preceeds the array.
						vec[d] = -v
					} else if v = dst[d] + 1 - tmp.Dims[d]; v > 0 && -v < vec[d] {
						// Negative lvec value indicates padding follows the array values.
						vec[d] = -v
					}
				}
				if out {
					tmp.Values[k] = apl.Index(0)
				} else {
					v, err := ar.At(ic.Index(dst))
					if err != nil {
						return nil, err
					}
					tmp.Values[k] = v // TODO copy?
				}

				apl.IncArrayIndex(sdx, tmp.Dims)
			}

			// Apply the stencil and set the result.
			v, err := f.Call(a, lvec, tmp)
			if err != nil {
				return nil, err
			}
			// We only accept scalar results. Is that ok?
			if _, ok := v.(apl.Array); ok {
				return nil, fmt.Errorf("stencil function must return a scalar, not an array")
			}
			res.Values[i] = v

			apl.IncArrayIndex(idx, rs)
		}
		return res, nil
	}
	return function(derived)
}
