package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

// splitAxis returns ax.R and converts ax.A to []int taking account of index origin.
// It R is not an axis it returns R and nil.
func splitAxis(a *apl.Apl, R apl.Value) (apl.Value, []int, error) {
	ax, ok := R.(apl.Axis)
	if ok == false {
		return R, nil, nil
	}
	to := domain.ToIndexArray(nil)
	X, ok := to.To(a, ax.A)
	if ok == false {
		return nil, nil, fmt.Errorf("axis is not an index array")
	}
	ar := X.(apl.IndexArray)
	shape := ar.Shape()
	if len(shape) != 1 {
		return nil, nil, fmt.Errorf("axis has wrong shape: %d", len(shape))
	}
	x := make([]int, len(ar.Ints))
	for i, n := range ar.Ints {
		x[i] = n - a.Origin
	}
	return ax.R, x, nil
}

// SplitCatAxis splits the right argument, if it contains an axis.
// The axis must be a numeric scalar value or a single element array.
// It may contain a fractional part.
// If it does not exist, it is set to the last axis.
// It returns the integer part of the axis and indicates if it is fractional.
// The index origin is substracted.
func splitCatAxis(a *apl.Apl, L, R apl.Value) (apl.Value, int, bool, error) {
	ax, ok := R.(apl.Axis)
	if ok == false {
		if ar, ok := R.(apl.Array); ok == false {
			return R, 0, false, nil
		} else {
			return R, len(ar.Shape()) - 1, false, nil
		}
	}
	R = ax.R

	// X∊⍳(⍴⍴L)⌈⍴⍴R
	rkL := 0
	al, ok := L.(apl.Array)
	if ok {
		rkL = len(al.Shape())
	}

	rkR := 0
	ar, ok := R.(apl.Array)
	if ok {
		rkR = len(ar.Shape())
	}

	max := rkL
	if rkR > max {
		max = rkR
	}

	var x apl.Value
	if xr, ok := ax.A.(apl.Array); ok {
		if apl.ArraySize(xr) != 1 {
			return nil, 0, false, fmt.Errorf(",: axis must be a scalar or single element array")
		} else {
			x, _ = xr.At(0)
		}
	} else {
		x = ax.A
	}
	num, ok := x.(apl.Number)
	if ok == false {
		return nil, 0, false, fmt.Errorf("axis is not numeric")
	}
	if n, ok := num.ToIndex(); ok {
		n -= a.Origin
		if n < 0 || n >= max {
			return nil, 0, false, fmt.Errorf("axis is out of range")
		}
		return R, n, false, nil
	}

	// The axis is fractional, depending on the numerical tower.
	// Substract index origin from the axis.
	if s, ok := num.(substracter2); ok == false {
		return nil, 0, false, fmt.Errorf("cannot substract from axis: %T", num)
	} else {
		if v, ok := s.Sub2(apl.Index(a.Origin)); ok == false {
			return nil, 0, false, fmt.Errorf("cannot substract index origin from axis")
		} else {
			num = v.(apl.Number)
		}
	}

	less, ok := num.(lesser)
	if ok == false {
		return nil, 0, false, fmt.Errorf("axis is not comparable: %T", num)
	}

	// num must be between -1 and max.
	if l, ok := less.Less(apl.Index(-1)); ok == false {
		return nil, 0, false, fmt.Errorf("axis cannot be compared")
	} else if l == true {
		return nil, 0, false, fmt.Errorf("axis is < ¯1+⎕IO")
	}
	for i := 0; i <= max; i++ {
		if l, ok := less.Less(apl.Index(i)); ok == false {
			return nil, 0, false, fmt.Errorf("axis cannot be compared")
		} else if l {
			return R, i, true, nil
		}
	}
	return nil, 0, false, fmt.Errorf("axis is too large")
}
