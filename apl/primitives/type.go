package primitives

import (
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
}

func typeof(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return apl.String(reflect.TypeOf(R).String()), nil
}
