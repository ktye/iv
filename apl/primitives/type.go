package primitives

import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⌶",
		doc:    "type",
		Domain: Monadic(nil),
		fn:     typeof,
	})
	register(primitive{
		symbol: "⌶",
		doc:    "convert to type",
		Domain: Dyadic(nil),
		fn:     convert,
	})
	register(primitive{
		symbol: "⌶",
		doc:    "convert to named type",
		Domain: Dyadic(Split(IsString(nil), nil)),
		fn:     convert2,
	})
}

func typeof(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return apl.String(reflect.TypeOf(R).String()), nil
}

// I-beam is a nice symbol for conversion, as it represents an encode-decode pair.

// convert is called with a values of the destination type on the left.
func convert(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	// L can be a prototype, similar to the left argument of format.
	return nil, fmt.Errorf("TODO: convert to destination type")
}

// convert is called with the name of the target type given as a string on the left argument.
func convert2(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	// TODO: L could be the string that monadic ibeam prints or
	// a prototype name similar to the left argument of format.
	s := L.(apl.String)
	switch s {
	case "img":
		to := ToImage(nil)
		if m, ok := to.To(a, R); ok {
			return m, nil
		}
		return nil, fmt.Errorf("cannot convert to image: %T", R)
	default:
		return nil, fmt.Errorf("convert: %T to %s is not supported", R, s)
	}
}
