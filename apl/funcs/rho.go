package funcs

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

func init() {
	register("⍴", both(shape, reshape))
	addDoc("⍴", `⍴ shape, reshape
Z←⍴R: R Scalar
	Z: empty array
Z←⍴R: R array
	Z: shape of R
Z←L⍴R: L integer, R: number or array
	Z: reshape of R to vector of with size of L
Z←L⍴R: L vector, R: number or array
	Z: reshape of R to shape L
	If Z has more elements than R, the elements are cycled through.
	If Z has more elements, values of R are truncated
`)
}

func shape(a *apl.Apl, ignored, v apl.Value) (bool, apl.Value, error) {
	switch v := v.(type) {
	case apl.Bool, apl.Int, apl.Float, apl.Complex:
		return true, apl.EmptyArray{}, nil
	case apl.GeneralArray:
		var ar apl.GeneralArray
		ar.Values = make([]apl.Value, len(v.Dims))
		for i, x := range v.Dims {
			ar.Values[i] = apl.Int(x)
		}
		ar.Dims = []int{len(ar.Values)}
		return true, ar, nil
	default:
		return false, nil, nil
	}
}

func reshape(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	var ar apl.Array

	switch r.(type) {
	case apl.Array:
		ar = r.(apl.Array)
	default:
		ar = apl.GeneralArray{
			Values: []apl.Value{r},
			Dims:   []int{1},
		}
	}

	var shape []int
	switch lv := l.(type) {
	case apl.Bool:
		if lv {
			shape = []int{1}
		} else {
			shape = []int{0}
		}
	case apl.Int:
		shape = []int{int(lv)}
	case apl.Float, apl.Complex:
		return true, nil, fmt.Errorf("cannot reshape to %T", l)
	case apl.Array:
		n := 0
		if s := lv.Shape(); len(s) != 1 {
			return true, nil, fmt.Errorf("reshape: shape is not a vector")
		} else {
			n = s[0]
			shape = make([]int, n)
		}
		for i := range shape {
			e, err := lv.At(i)
			if err != nil {
				return true, nil, err
			}
			if b, ok := e.(apl.Bool); ok {
				if b {
					shape[i] = 1
				} else {
					shape[i] = 0
				}
			} else if k, ok := e.(apl.Int); ok {
				shape[i] = int(k)
			} else {
				return true, nil, fmt.Errorf("element shape vector has wrong type: %T", e)
			}
		}
	default:
		return false, nil, nil
	}

	if rs, ok := ar.(apl.Reshaper); ok {
		return true, rs.Reshape(shape), nil
	}
	return false, nil, nil
}
