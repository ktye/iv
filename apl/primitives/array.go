package primitives

import (
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

// array1 tries to apply the elementary function returned by arith1(fn)
// monadically to each element of the array R
func array1(symbol string, fn func(*apl.Apl, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	efn := arith1(symbol, fn)
	return func(a *apl.Apl, _ apl.Value, R apl.Value) (apl.Value, error) {
		ar := R.(apl.Array)
		res := apl.GeneralArray{
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
		res := apl.GeneralArray{Dims: shape}
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

/* TODO should we try all primitives for arrays?
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
*/
