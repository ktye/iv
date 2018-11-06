package funcs

import (
	"github.com/ktye/iv/apl"
)

func init() {
	register("-", both(negate, arrayWrap(substract)))
	addDoc("-", `- primitive function: negate, substract, minus
Z←-R: R Numeric
	Z: reverse the sign of R
Z←L-R: L, R numeric: substraction L-R
Z←L-R: one or both Array: elementwise substraction
`)
}

// Negate reverses the sign of v.
func negate(a *apl.Apl, ignored, v apl.Value) (bool, apl.Value, error) {
	switch v := v.(type) {
	case apl.Bool:
		if v {
			return true, apl.Int(-1), nil
		}
		return true, apl.Int(0), nil
	case apl.Int:
		return true, -v, nil
	case apl.Float:
		return true, -v, nil
	case apl.Complex:
		return true, -v, nil
	case apl.Array:
		rv, err := v.ApplyMonadic(a, handle(negate))
		return true, rv, err
	default:
		return false, nil, nil
	}
}

// Substract two scalar values of the same type.
func substract(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	switch lv := l.(type) {
	case apl.Bool:
		rv := r.(apl.Bool)
		if lv && !rv {
			return true, apl.Int(1), nil
		} else if rv && !lv {
			return true, apl.Int(-1), nil
		}
		return true, apl.Int(0), nil
	case apl.Int:
		return true, lv - r.(apl.Int), nil
	case apl.Float:
		return true, lv - r.(apl.Float), nil
	case apl.Complex:
		return true, lv - r.(apl.Complex), nil
	default:
		return false, nil, nil
	}
}
