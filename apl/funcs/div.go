package funcs

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

func init() {
	register("÷", both(reciprocal, wrap(divide)))
	addDoc("÷", `÷ reciprocal, divide
Z←÷R: R numeric
	Z: reciprocal: 1÷R
Z←L÷R: R numeric
	division: L÷R
`)
}

// Reciprocal of the right argument.
func reciprocal(a *apl.Apl, ignored, v apl.Value) (bool, apl.Value, error) {
	var fv apl.Float
	switch v := v.(type) {
	case apl.Bool:
		return true, nil, fmt.Errorf("cannot ÷ bool")
	case apl.Int:
		fv = apl.Float(v)
	case apl.Complex:
		return true, apl.Complex(complex(1, 0) / complex128(v)), nil
	case apl.Array:
		rv, err := v.ApplyMonadic(a, reciprocal)
		return true, rv, err
	default:
		return false, nil, nil
	}
	return true, apl.Float(1.0 / float64(fv)), nil
}

// Divide two scalar values of the same type.
func divide(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	switch lv := l.(type) {
	case apl.Bool:
		return true, nil, fmt.Errorf("cannot ÷ boolean values")
	case apl.Int:
		return true, apl.Float(lv) / apl.Float(r.(apl.Int)), nil
	case apl.Float:
		return true, lv / r.(apl.Float), nil
	case apl.Complex:
		return true, lv / r.(apl.Complex), nil
	default:
		return false, nil, nil
	}
}
