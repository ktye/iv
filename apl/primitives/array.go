package primitives

import (
	"fmt"
	"sort"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

// arrays is the domain for binary arithmetic functions which
// may be scalars or arrays on both sides.
// If this function suceeds, only these cases are possible:
//	- one or both are empty (apl.EmptyArray)
//	- one is scalar and the other an array
//	- both are arrays of the same shape
// A single element array is converted to a scalar, if the other is a larger array.
type arrays struct{}

func (ars arrays) To(a *apl.Apl, L, R apl.Value) (apl.Value, apl.Value, bool) {
	if L == nil {
		return L, R, false
	}
	isScalar := domain.IsScalar(nil)
	_, scalarL := isScalar.To(a, L)
	_, scalarR := isScalar.To(a, R)
	if scalarL && scalarR {
		return L, R, false
	}

	al, isLarray := L.(apl.Array)
	ar, isRarray := R.(apl.Array)

	// Both are arrays and must have the same shape or at least one is empty.
	// They must also contain only numbers.
	if isLarray && isRarray {
		// 0-Size arrays are converted to empty arrays.
		// TODO: is this correct?
		if apl.ArraySize(al) == 0 {
			return apl.EmptyArray{}, ar, true
		} else if apl.ArraySize(ar) == 0 {
			return al, apl.EmptyArray{}, true
		}
		lshape := al.Shape()
		rshape := ar.Shape()
		if len(lshape) == len(rshape) {
			for i := range lshape {
				if lshape[i] != rshape[i] {
					break
				}
				if i == len(lshape)-1 {
					// Both arrays have the same shape
					return al, ar, true
				}
			}
		}
		// Convert single element arrays to scalars.
		if apl.ArraySize(al) == 1 {
			if v, err := al.At(0); err == nil {
				return v, ar, true
			}
		}
		if apl.ArraySize(ar) == 1 {
			if v, err := ar.At(0); err == nil {
				return al, v, true
			}
		}
		return L, R, false
	}
	if isLarray && scalarR {
		return L, R, true
	} else if scalarL && isRarray {
		return L, R, true
	}
	return L, R, false
}
func (ars arrays) String(a *apl.Apl) string { return "arithmetic arrays" }

// ArraysWithAxis is the domain for binary arithmetic functions
// with an axis specification.
// The axis specification is bound to the right argument in an Axis value.
// It converts L, R and the axis specification to Arrays.
type arraysWithAxis struct{}

func (ars arraysWithAxis) To(a *apl.Apl, L, R apl.Value) (apl.Value, apl.Value, bool) {
	if L == nil {
		return L, R, false
	}
	ax, ok := R.(apl.Axis)
	if ok == false {
		return L, R, false
	}

	toArray := domain.ToArray(nil)
	al, ok := toArray.To(a, L)
	if ok == false {
		return L, R, false
	}
	ar, ok := toArray.To(a, ax.R)
	if ok == false {
		return L, R, false
	}

	toIdxArray := domain.ToIndexArray(nil)
	x, ok := toIdxArray.To(a, ax.A)
	if ok == false {
		return L, R, false
	}

	return al, apl.Axis{A: x, R: ar}, true
}
func (ars arraysWithAxis) String(a *apl.Apl) string { return "arithmetic arrays with axis" }

// array1 tries to apply the elementary function returned by arith1(fn)
// monadically to each element of the array R
func array1(symbol string, fn func(*apl.Apl, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	efn := arith1(symbol, fn)
	return func(a *apl.Apl, _ apl.Value, R apl.Value) (apl.Value, error) {
		ar := R.(apl.Array)
		res := apl.MixedArray{
			Values: make([]apl.Value, apl.ArraySize(ar)),
			Dims:   apl.CopyShape(ar),
		}
		for i := range res.Values {
			e, err := ar.At(i)
			if err != nil {
				return nil, err
			}
			val, err := efn(a, nil, e)
			if err != nil {
				return nil, err
			}
			res.Values[i] = val
		}
		return res, nil
	}
}

// array2 tries to apply the elementary function returned by arith2(fn)
// dyadically to the elements of the arrays L and R.
// L and R have been tested and converted by arrays.
func array2(symbol string, fn func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	efn := arith2(symbol, fn)
	return func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		// If one or both are empty, return an EmptyArray{}
		// TODO: or should we test against the array size (0 dimsions anywhere)
		// and return what? EmptyArray{} or array of the same shape?
		_, emptyL := L.(apl.EmptyArray)
		_, emptyR := R.(apl.EmptyArray)
		if emptyL || emptyR {
			return apl.EmptyArray{}, nil
		}

		al, isLarray := L.(apl.Array)
		ar, isRarray := R.(apl.Array)

		var shape []int
		if isLarray == false {
			shape = apl.CopyShape(ar)
		} else {
			shape = apl.CopyShape(al)
		}
		res := apl.MixedArray{Dims: shape}
		res.Values = make([]apl.Value, apl.ArraySize(res))
		var err error
		for i := range res.Values {
			lv := L
			if isLarray {
				lv, err = al.At(i)
				if err != nil {
					return nil, err
				}
			}
			rv := R
			if isRarray {
				rv, err = ar.At(i)
				if err != nil {
					return nil, err
				}
			}
			if val, err := efn(a, lv, rv); err != nil {
				return nil, err
			} else {
				res.Values[i] = val
			}
		}
		return res, nil
	}
}

// ArrayAxis is like array2 but with R bound in an axis specification.
func arrayAxis(symbol string, fn func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	efn := arith2(symbol, fn)
	return func(a *apl.Apl, L, ax apl.Value) (apl.Value, error) {
		axis := ax.(apl.Axis)
		R := axis.R
		X := axis.A.(apl.IndexArray)

		_, emptyL := L.(apl.EmptyArray)
		_, emptyR := R.(apl.EmptyArray)
		if emptyL || emptyR {
			return apl.EmptyArray{}, nil
		}

		al := L.(apl.Array)
		ar := R.(apl.Array)
		ls := al.Shape()
		rs := ar.Shape()

		// We assume L has higher rank, otherwise flip L and R.
		flip := false
		if len(ls) < len(rs) {
			flip = true
			al, ar = ar, al
			ls, rs = rs, ls
		}

		// See APL2 p.55 for conformance:
		// 	L f[X] R
		// (⍴,X) ←→ (⍴⍴L)⌊⍴⍴R
		// (⍴,X) ←→ ∧/X∊⍳(⍴⍴L)⌈⍴⍴R
		if len(X.Dims) != 1 {
			return nil, fmt.Errorf("axis specification must have rank 1: %T", len(X.Dims))
		}

		// X≡X[⍋X]
		x := make([]int, len(X.Ints))
		copy(x, X.Ints)
		sort.Ints(x)
		for i := range x {
			x[i] -= a.Origin
		}

		// (⍴L)[X] ←→ (⍴R).
		if len(rs) != len(x) {
			return nil, fmt.Errorf("axis rank must match lower argument rank")
		}
		for i, n := range x {
			if n < 0 || n >= len(ls) {
				return nil, fmt.Errorf("axis exceeds higher argument rank")
			}
			if i > 0 && n == x[i-1] {
				return nil, fmt.Errorf("axis values are not unique")
			}
			if ls[n] != rs[i] {
				return nil, fmt.Errorf("arguments with axis do not conform")
			}
		}

		// There is no explicit algorithm description in APL2, DyaRef or ISO. We do:
		// Extend the rank of argument with lower rank (already flipped to R),
		// to the higher rank by filling missing axes with 1s.
		// Apply the function elementwise, but use index 1 if an axis has only one element.
		rightShape := make([]int, len(ls))
		for i := range rightShape {
			rightShape[i] = 1
		}
		for i, n := range x {
			rightShape[n] = rs[i]
		}

		var err error
		var lv, rv, v apl.Value
		res := apl.MixedArray{Dims: apl.CopyShape(al)}
		res.Values = make([]apl.Value, apl.ArraySize(res))
		idx := make([]int, len(res.Dims))
		ic, rdx := apl.NewIdxConverter(rightShape)
		for i := range res.Values {
			copy(rdx, idx)
			for k := range rdx {
				if rdx[k] >= rightShape[k] {
					rdx[k] = 0
				}
			}

			lv, err = al.At(i)
			if err != nil {
				return nil, err
			}

			rv, err = ar.At(ic.Index(rdx))
			if err != nil {
				return nil, err
			}

			if flip {
				lv, rv = rv, lv
			}
			v, err = efn(a, lv, rv)
			if err != nil {
				return nil, err
			}
			res.Values[i] = v

			apl.IncArrayIndex(idx, res.Dims)
		}
		return res, nil
	}
}
