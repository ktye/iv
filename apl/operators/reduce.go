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
		derived: reduceLast,
	})
	register(operator{
		symbol:  "⌿",
		Domain:  MonadicOp(Function(nil)),
		doc:     "reduce first, n-wise reduction",
		derived: reduceFirst,
	})
	register(operator{
		symbol:  `\`,
		Domain:  MonadicOp(Function(nil)),
		doc:     "scan",
		derived: scanLast,
	})
	register(operator{
		symbol:  `⍀`,
		Domain:  MonadicOp(Function(nil)),
		doc:     "scan first axis",
		derived: scanFirst,
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

func reduceLast(a *apl.Apl, f, _ apl.Value) apl.Function {
	return reduction(a, f, -1)
}

func reduceFirst(a *apl.Apl, f, _ apl.Value) apl.Function {
	return reduction(a, f, 0)
}

func scanLast(a *apl.Apl, f, _ apl.Value) apl.Function {
	return scanArray(a, f, -1)
}

func scanFirst(a *apl.Apl, f, _ apl.Value) apl.Function {
	return scanArray(a, f, 0)
}

// Reduction returns the derived function f/ .
func reduction(a *apl.Apl, f apl.Value, axis int) apl.Function {

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

		if axis < 0 {
			axis = len(shape) + axis
		}
		if axis < 0 || axis >= len(shape) {
			return nil, fmt.Errorf("reduce: axis rank is %d but axis %d", len(shape), axis)
		}

		// n is the number of values being reduced (the length if the reduction axis).
		n := shape[axis]

		// Dims is the shape of the result.
		dims := make([]int, len(shape)-1)
		k := 0
		for i := range shape {
			if i != axis {
				dims[k] = shape[i]
				k++
			}
		}

		// If the length of the axis is 1, the result is a reshape.
		// TODO: if the length of any other axis is 0, this should be triggered as well.
		if n == 1 {
			if rs, ok := ar.(apl.Reshaper); ok {
				return rs.Reshape(dims), nil
			} else {
				return nil, fmt.Errorf("reduce with axis length 1: cannot reshape %T", ar)
			}
		}
		if n == 0 {
			// TODO: If the last axis is 0, apply an identity function, DyaRef p 169
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

		// Create a new array with the given axis removed.
		values := make([]apl.Value, apl.ArraySize(apl.GeneralArray{Dims: dims}))
		v := apl.GeneralArray{
			Dims:   dims,
			Values: values,
		}

		vec := make([]apl.Value, n)
		sidx := make([]int, len(dims)+1) // src index vector
		tidx := make([]int, len(dims))   // index vector in target array
		for k := range v.Values {
			// Copy target index over the source index,
			// leaving the reduced axis unset.
			copy(sidx, tidx[:axis])
			copy(sidx[axis+1:], tidx[axis:])
			// Iterate over the reduced axis
			for i := range vec {
				sidx[axis] = i
				// TODO: maybe this could be done more efficiently
				// e.g. by iteration with a fixed increase.
				if m, err := apl.ArrayIndex(shape, sidx); err != nil {
					return nil, err
				} else if val, err := ar.At(m); err != nil {
					return nil, err
				} else {
					vec[i] = val
				}
			}
			apl.IncArrayIndex(tidx, dims)

			if res, err := reduce(a, vec, d); err != nil {
				return nil, fmt.Errorf("cannot reduce: %s", err)
			} else {
				v.Values[k] = res
			}
		}
		return v, nil
	}

	return function(derived)
}

// ScanArray is the derived function f\ .
func scanArray(a *apl.Apl, f apl.Value, axis int) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		if L != nil {
			return nil, fmt.Errorf("scan: derived function is not defined for dyadic context")
		}

		ar, ok := R.(apl.Array)
		if ok == false {
			// If R is scalar, return unchanged.
			return R, nil
		}

		d := f.(apl.Function)

		// The result has the same shape as R.
		dims := apl.CopyShape(ar)
		res := apl.GeneralArray{
			Values: make([]apl.Value, apl.ArraySize(ar)),
			Dims:   dims,
		}

		if len(dims) == 0 {
			return apl.EmptyArray{}, nil
		}

		if axis < 0 {
			axis = len(dims) + axis
		}
		if axis < 0 || axis >= len(dims) {
			return nil, fmt.Errorf("scan: axis rank is %d but axis %d", len(dims), axis)
		}

		// Shortcut, if R is a vector
		if len(dims) == 1 {
			vec := make([]apl.Value, dims[0])
			for i := range vec {
				if v, err := ar.At(i); err != nil {
					return nil, err
				} else {
					vec[i] = v
				}
			}
			vec, err := scan(a, vec, d)
			if err != nil {
				return nil, err
			}
			return apl.GeneralArray{
				Values: vec,
				Dims:   []int{len(vec)},
			}, nil
		}

		// Loop over the indexes, with the scan axis length set to 1.
		lidx := apl.CopyShape(ar)
		lidx[axis] = 1
		idx := make([]int, len(lidx))
		vec := make([]apl.Value, dims[axis])
		for i := 0; i < apl.ArraySize(apl.GeneralArray{Dims: lidx}); i++ {
			// Build the scan vector, by iterating over the axis.
			for k := range vec {
				idx[axis] = k
				n, err := apl.ArrayIndex(dims, idx)
				if err != nil {
					fmt.Println("idx", idx, "dims", dims) // TODO rm
					panic("err1")
					return nil, err
				}
				val, err := ar.At(n)
				if err != nil {
					fmt.Println("n", n) // TODO rm
					panic("err2")
					return nil, err
				}
				vec[k] = val
			}
			vals, err := scan(a, vec, d)
			if err != nil {
				return nil, err
			}

			// Assign the values to the destination indexes.
			for k := range vals {
				idx[axis] = k
				n, err := apl.ArrayIndex(dims, idx)
				if err != nil {
					fmt.Println("n", n) // TODO rm
					panic("err3")
					return nil, err
				}
				res.Values[n] = vals[k]
			}

			// Reset the index vector and increment.
			idx[axis] = 0
			apl.IncArrayIndex(idx, lidx)
		}
		return res, nil
	}
	return function(derived)
}

func reduce(a *apl.Apl, vec []apl.Value, d apl.Function) (apl.Value, error) {
	var err error
	v := vec[len(vec)-1] // TODO: copy?
	for i := len(vec) - 2; i >= 0; i-- {
		v, err = d.Call(a, vec[i], v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func scan(a *apl.Apl, vec []apl.Value, d apl.Function) ([]apl.Value, error) {
	// The ith element of the result is: d/I↑V
	res := make([]apl.Value, len(vec))
	res[0] = vec[0] // TODO: copy?
	for i := 1; i < len(res); i++ {
		if v, err := reduce(a, vec[:i+1], d); err != nil {
			return nil, err
		} else {
			res[i] = v
		}
	}
	return res, nil
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
