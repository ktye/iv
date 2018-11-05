package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

func init() {
	register("/", reduction{})
	addDoc("/", `/ monadic operator: reduce, n-wise reduction, replacate
Z←L LO / R	
`)
}

// Reducer is an interface that custom types may implement.
// If the reduction operator finds an unknown type on it's left that
// implements this interface, it forwards the call.
type Reducer interface {
	Reduce() apl.FunctionHandle
}

type reduction struct {
	monadic
}

// OperateMonadic returns the derived function f over r (summation).
func (r reduction) Apply(f, dummy apl.Value) apl.FunctionHandle {

	// Forward custom implementations.
	if fn, ok := f.(Reducer); ok {
		return fn.Reduce()
	}

	return func(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
		if l != nil {
			return nwise(a, l, r)
		}

		// TODO compression f is an array.
		if _, ok := f.(apl.Array); ok {
			return true, nil, fmt.Errorf("TODO: compression (array/ )")
		}

		// Reduction needs a dyadic function to it's left.
		var d apl.Function
		if fn, ok := f.(apl.Function); ok == false {
			return true, nil, fmt.Errorf("left argument to / must be a function: %T", d)
		} else {
			d = fn
		}

		// If R is a scalar, the operation is not applied and Z←R
		ar, ok := r.(apl.Array)
		if ok == false {
			return true, r, nil
		}

		shape := ar.Shape()
		if len(shape) == 0 {
			return true, ar, nil // Not sure if this is ok.
		}

		n := shape[len(shape)-1]
		if n == 1 {
			// TODO: If the last axis is 1, the operation is not applied and Z ← (b1↓ρR)ρR
			return true, nil, fmt.Errorf("reduce on R, with last axis 1: TODO Z ← (b1↓ρR)ρR")
		}
		if n == 0 {
			// TODO: If the last axis is 0, apply an identity function
			return true, nil, fmt.Errorf("reduce on R, with last axis 0: TODO apply identity function")
		}

		// Reduce directly, if R is a vector.
		if len(shape) == 1 {
			vec := make([]apl.Value, shape[0])
			var err error
			for i := range vec {
				vec[i], err = ar.At(i)
				if err != nil {
					return true, nil, err
				}
			}
			v, err := reduce(a, vec, d)
			return true, v, err
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
					return true, nil, err
				}
				m++
			}
			if res, err := reduce(a, lastAxis, d); err != nil {
				return true, nil, fmt.Errorf("cannot reduce: %s", err)
			} else {
				v.Values[k] = res
			}
		}
		return true, v, nil
	}
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
func nwise(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	return true, nil, fmt.Errorf("TODO: n-wise reduction")
}
