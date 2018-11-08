package funcs

import (
	"math/cmplx"
	"strings"

	"github.com/ktye/iv/apl"
)

func init() {
	register("×", both(signum, arrayWrap(multiply)))
	register("×", handle(stringsRepeat))

	addDoc("×", `× primitive function: signum, sign of, direction, multiply
Z←×R: R Bool, Int, Float
	Z: -1, 0 or 1 depending on the sign of R.
Z←×R: R Complex
	Z: direction as a complex number, with the same phase as R
	but magnitude 1.
	Z: 0, if R is 0.
Z←L×R: L, R: numeric
	Z: multiplication L times R
Z←L×R: L int, R: string catenate R L times (repeat), same for L int, R string.
`)
}

func signum(a *apl.Apl, ignored, v apl.Value) (bool, apl.Value, error) {
	var fv apl.Float
	switch v := v.(type) {
	case apl.Bool:
		if v {
			fv = 1
		}
	case apl.Int:
		fv = apl.Float(v)
	case apl.Complex:
		return true, direction(complex128(v)), nil
	case apl.Array:
		rv, err := v.ApplyMonadic(a, handle(signum))
		return true, rv, err
	default:
		return false, nil, nil
	}
	if fv > 0 {
		return true, apl.Int(1), nil
	} else if fv < 0 {
		return true, apl.Int(-1), nil
	}
	return true, apl.Int(0), nil
}

// direction returns the direction of the complex value as a complex number
// on the unit circle.
// It returns 0,
func direction(c complex128) apl.Value {
	if c == 0 {
		return apl.Complex(c)
	}
	r := cmplx.Abs(c)
	return apl.Complex(complex(real(c)/r, imag(c)/r))
}

// Multiply two scalar values of the same type.
func multiply(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	switch lv := l.(type) {
	case apl.Bool:
		return true, apl.Bool(lv && r.(apl.Bool)), nil
	case apl.Int:
		return true, lv * r.(apl.Int), nil
	case apl.Float:
		return true, lv * r.(apl.Float), nil
	case apl.Complex:
		return true, lv * r.(apl.Complex), nil
	default:
		return false, nil, nil
	}
}

// StringsRepeat catenates the string n times.
func stringsRepeat(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	if l == nil {
		return false, nil, nil // cannot be called in monadic context.
	}
	if n, ok := apl.ToInt(l); ok {
		if s, ok := r.(apl.String); ok {
			return true, apl.String(strings.Repeat(string(s), n)), nil
		}
	}
	if s, ok := l.(apl.String); ok {
		if n, ok := apl.ToInt(r); ok {
			return true, apl.String(strings.Repeat(string(s), n)), nil
		}
	}
	return false, nil, nil
}
