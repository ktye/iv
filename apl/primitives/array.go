package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

// arithmetic is the domain for binary arithmetic functions which
// may be scalars or arrays on both sides.
// Arrays are tested if they contain only basic numeric types.
// If this function suceeds, only these cases are possible:
//	- one or both are empty (apl.EmptyArray)
//	- both are scalar numbers
//	- one is a scalar number and the other an array of numbers
//	- both are arrays of number of the same shape
// Single element arryas are converted to scalars.
type arithmetic struct{}

func (am arithmetic) To(a *apl.Apl, L, R apl.Value) (apl.Value, apl.Value, bool) {
	isNumber := domain.ToNumber(nil)
	l, numL := isNumber.To(a, L)
	r, numR := isNumber.To(a, R)

	// Both are numbers.
	if numL && numR {
		return l, r, true
	}

	al, isLarray := L.(apl.Array)
	ar, isRarray := R.(apl.Array)

	// Both are arrays and must have the same shape or at least one is empty.
	// They must also contain only numbers.
	if isLarray && isRarray {
		// 0-Size arrays are converted to empty arrays.
		// The other array is not tested for numeric values.
		// TODO: is this correct?
		if apl.ArraySize(al) == 0 {
			return apl.EmptyArray{}, ar, true
		} else if apl.ArraySize(ar) == 0 {
			return al, apl.EmptyArray{}, true
		}
		lshape := al.Shape()
		rshape := ar.Shape()
		if len(lshape) != len(rshape) {
			return L, R, false
		}
		for i := range lshape {
			if lshape[i] != rshape[i] {
				return L, R, false
			}
		}
		if isNumericArray(al) == false || isNumericArray(ar) == false {
			return L, R, false
		}
		return al, ar, true
	} else if isLarray && numR {
		if isNumericArray(al) == false {
			return L, R, false
		} else if apl.ArraySize(al) == 0 {
			return apl.EmptyArray{}, r, true
		}
		return al, r, true
	} else if numL && isRarray {
		if isNumericArray(ar) == false {
			return L, R, false
		} else if apl.ArraySize(ar) == 0 {
			return l, apl.EmptyArray{}, true
		}
		return l, ar, true
	}
	return L, R, false
}
func (am arithmetic) String(a *apl.Apl) string { return "arithmetic" }

// isNumericArray tests if the array contains only numbers.
// An array of size 0 passes the test.
func isNumericArray(a apl.Array) bool {
	isNumber := domain.IsNumber(nil)
	for i := 0; i < apl.ArraySize(a); i++ {
		if v, err := a.At(i); err != nil {
			return false
		} else if _, ok := isNumber.To(nil, v); ok == false {
			return false
		}
	}
	return true
}

// arith takes the elementry function f and applies it to it's left and right arguments,
// which may be arrays.
// It has already been checked that the array dimensions are ok and the values
// are of basic numeric type.
func arith(f func(a *apl.Apl, L, R apl.Value) (apl.Value, error)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
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

		// If both are scalars, promote to the same type and apply the elementry function.
		if isLarray == false && isRarray == false {
			L, R = domain.MustPromote(L, R)
			return f(a, L, R)
		}

		// If at least one is an array, create a new GeneralArray of the same shape
		// and apply the elementry function.
		var prototype apl.Array
		if isLarray {
			prototype = al
		} else if isRarray {
			prototype = ar
		}
		res := apl.GeneralArray{
			Values: make([]apl.Value, apl.ArraySize(prototype)),
			Dims:   apl.CopyShape(prototype),
		}

		// Both are arrays of the same size.
		if isLarray && isRarray {
			for i := range res.Values {
				el, _ := al.At(i)
				er, _ := ar.At(i)
				if v, err := f(a, el, er); err != nil {
					return nil, err
				} else {
					res.Values[i] = v
				}
			}
			return res, nil
		}

		// One is the array, the other scalar.
		if isLarray {
			for i := range res.Values {
				el, _ := al.At(i)
				if v, err := f(a, el, R); err != nil {
					return nil, err
				} else {
					res.Values[i] = v
				}
			}
		} else {
			for i := range res.Values {
				er, _ := ar.At(i)
				if v, err := f(a, L, er); err != nil {
					return nil, err
				} else {
					res.Values[i] = v
				}
			}
		}
		return res, nil
	}
}

// arithmonads takes monadic primitives and applies it to the array R.
func arithmonads(p ...primitive) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	return func(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
		r := R.(apl.Array)
		ar := apl.GeneralArray{
			Values: make([]apl.Value, apl.ArraySize(r)),
			Dims:   apl.CopyShape(r),
		}
		for i := range ar.Values {
			v, err := r.At(i)
			if err != nil {
				return nil, err
			}
			// Try each primitive.
			for k, pk := range p {
				_, u, ok := pk.Domain.To(a, nil, v)
				if ok {
					if w, err := pk.fn(a, nil, u); err != nil {
						return nil, err
					} else {
						ar.Values[i] = w
						break
					}
				} else if k == len(p)-1 {
					return nil, fmt.Errorf("cannot apply monadic %s to array element %T", pk.symbol, v)
				}
			}
		}
		return ar, nil
	}
}
