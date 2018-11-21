package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "/",
		Domain:  MonadicOp(Function(nil)),
		doc:     "reduce, n-wise reduction",
		derived: reduction,
	})
	/* TODO APL2 p 220
	register(operator{
		symbol:  "/",
		Domain:  Left(Array(nil)), // scalar or vector, integer
		doc:     "replicate",
		derived: replicate,
	})
	*/
	/* TODO APL2 p 85
	register(operator{
		symbol:  "/",
		Domain:  Left(Array(nil)), // scalar or vector, bool
		doc:     "compress",
		derived: compress,
	})
	*/
}

// Reduction returns the derived function f over r.
func reduction(a *apl.Apl, f, _ apl.Value) apl.Function {

	// Special cases: left, right tack.
	if p, ok := f.(apl.Primitive); ok {
		if p == "⊣" {
			return reduceTack(true)
		} else if p == "⊢" {
			return reduceTack(false)
		}
	}

	derived := func(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
		if l != nil {
			return nwise(a, l, r)
		}

		d := f.(apl.Function)

		// If R is a scalar, the operation is not applied and Z←R
		ar, ok := r.(apl.Array)
		if ok == false {
			return r, nil
		}

		shape := ar.Shape()
		if len(shape) == 0 {
			return ar, nil // Not sure if this is ok.
		}

		n := shape[len(shape)-1]
		if n == 1 {
			// TODO: If the last axis is 1, the operation is not applied and Z ← (b1↓ρR)ρR
			return nil, fmt.Errorf("reduce on R, with last axis 1: TODO Z ← (b1↓ρR)ρR")
		}
		if n == 0 {
			// TODO: If the last axis is 0, apply an identity function
			return nil, fmt.Errorf("reduce on R, with last axis 0: TODO apply identity function")
		}

		// Reduce directly, if R is a vector.
		if len(shape) == 1 {
			vec := make([]apl.Value, shape[0])
			var err error
			for i := range vec {
				vec[i], err = ar.At(i)
				if err != nil {
					return nil, err
				}
			}
			v, err := reduce(a, vec, d)
			return v, err
		}

		// Create a new array with the last axis removed.
		dims := make([]int, len(shape)-1)
		copy(dims, shape[:len(shape)-1])
		values := make([]apl.Value, apl.ArraySize(apl.GeneralArray{Dims: dims}))
		v := apl.GeneralArray{
			Dims:   dims,
			Values: values,
		}

		var err error
		lastAxis := make([]apl.Value, n)
		m := 0
		for k := range v.Values {
			for i := range lastAxis {
				lastAxis[i], err = ar.At(m)
				if err != nil {
					return nil, err
				}
				m++
			}
			if res, err := reduce(a, lastAxis, d); err != nil {
				return nil, fmt.Errorf("cannot reduce: %s", err)
			} else {
				v.Values[k] = res
			}
		}
		return v, nil
	}

	return function(derived)
}

func reduce(a *apl.Apl, vec []apl.Value, d apl.Function) (apl.Value, error) {
	var err error
	v := vec[len(vec)-1]
	for i := len(vec) - 2; i >= 0; i-- {
		v, err = d.Call(a, vec[i], v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

// Nwise is the function handle for n-wise recution.
// l must be a scalar (integer) or a 1 element vector.
func nwise(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	return nil, fmt.Errorf("TODO: n-wise reduction")
}

// reduceTack is the derived function from ⊣/ or ⊢/ .
type reduceTack bool

func (first reduceTack) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if L != nil {
		return nil, fmt.Errorf("tack-reduce can only be applied monadically")
	}
	ar, ok := R.(apl.Array)
	if ok == false {
		return R, nil
	}
	shape := ar.Shape()
	if len(shape) == 0 {
		return R, nil
	}

	// Reduce a vector to a scalar.
	if len(shape) == 1 {
		if shape[0] <= 0 {
			return apl.EmptyArray{}, nil
		}
		var v apl.Value
		var err error
		if first {
			v, err = ar.At(0)
		} else {
			v, err = ar.At(shape[0] - 1)
		}
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	// Create a new array
	inner := shape[len(shape)-1]
	newshape := apl.CopyShape(ar)
	newshape = newshape[:len(newshape)-1]
	ret := apl.GeneralArray{Dims: newshape}
	ret.Values = make([]apl.Value, apl.ArraySize(ret))
	i := 0
	n := 0 // index over inner axis.
	for k := 0; k < apl.ArraySize(ar); k++ {
		if first && n == 0 {
			if v, err := ar.At(k); err != nil {
				return nil, err
			} else {
				ret.Values[i] = v
				i++
			}
		} else if first == false && n == inner-1 {
			if v, err := ar.At(k); err != nil {
				return nil, err
			} else {
				ret.Values[i] = v
				i++
			}
		}
		n++
		if n == inner {
			n = 0
		}
	}
	return ret, nil
}
