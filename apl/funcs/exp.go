package funcs

import (
	"math"
	"math/cmplx"

	"github.com/ktye/iv/apl"
)

func init() {
	register("*", both(exponential, arrayWrap(power)))
	addDoc("*", `* primitive function: exponential, power
Z←*R: R numeric
	Z: e*R exponential with base 2.78...
Z←L*R: L, R numeric
	Z: raises the base L to the Rth power.
	L: negative → Z: complex.
`)
}

func exponential(a *apl.Apl, ignored, v apl.Value) (bool, apl.Value, error) {
	var f float64
	switch v := v.(type) {
	case apl.Bool:
		if v {
			f = 1
		}
	case apl.Int:
		f = float64(v)
	case apl.Float:
		f = float64(v)
	case apl.Complex:
		return true, apl.Complex(cmplx.Exp(complex128(v))), nil
	case apl.Array:
		rv, err := v.ApplyMonadic(a, handle(exponential))
		return true, rv, err
	default:
		return false, nil, nil
	}
	return true, apl.Float(math.Exp(f)), nil
}

// Power raises l to the power of r.
func power(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	var base, exp float64
	switch lv := l.(type) {
	case apl.Bool:
		if r.(apl.Bool) == false {
			return true, apl.Bool(true), nil
		} else {
			return true, lv, nil
		}
	case apl.Int:
		base = float64(lv)
		exp = float64(r.(apl.Int))
	case apl.Float:
		base = float64(lv)
		exp = float64(r.(apl.Float))
	case apl.Complex:
		return true, apl.Complex(cmplx.Pow(complex128(lv), complex128(r.(apl.Complex)))), nil
	default:
		return false, nil, nil
	}
	if exp == 0 {
		// This may downcast floats. Is this ok?
		return true, apl.Int(1), nil
	}
	if base < 0 {
		return true, apl.Complex(cmplx.Pow(complex(base, 0), complex(exp, 0))), nil
	}
	return true, apl.Float(math.Pow(base, exp)), nil
}
