package funcs

import (
	"math/cmplx"

	"github.com/ktye/iv/apl"
)

func init() {
	register("+", both(conjugate, arrayWrap(add)))
	addDoc("+", `+ primitive function: conjugate, add, plus
Z←+R: R Complex
	Z: conjugate complex of R
Z←+R: R Numeric
	Z: Identity
Z←L+R: L, R numeric: addition of L+R
Z←L+R: L, R strings: catenation of L and R
`)
}

func conjugate(a *apl.Apl, ignored, v apl.Value) (bool, apl.Value, error) {
	switch v := v.(type) {
	case apl.Bool, apl.Int, apl.Float:
		return true, v, nil
	case apl.Complex:
		return true, apl.Complex(cmplx.Conj(complex128(v))), nil
	case apl.Array:
		rv, err := v.ApplyMonadic(a, handle(conjugate))
		return true, rv, err
	default:
		return false, nil, nil
	}
}

// Add adds two scalar values of the same type.
func add(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	switch lv := l.(type) {
	case apl.Bool:
		rv := r.(apl.Bool)
		if lv && rv {
			return true, apl.Int(2), nil
		} else if lv || rv {
			return true, apl.Int(1), nil
		}
		return true, apl.Int(0), nil
	case apl.Int:
		return true, lv + r.(apl.Int), nil
	case apl.Float:
		return true, lv + r.(apl.Float), nil
	case apl.Complex:
		return true, lv + r.(apl.Complex), nil
	case apl.String: // String catenation.
		return true, apl.String(string(lv) + string(r.(apl.String))), nil
	default:
		return false, nil, nil
	}
}
