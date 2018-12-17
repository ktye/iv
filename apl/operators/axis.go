package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:    "[]",
		Domain:    DyadicOp(nil),
		doc:       "axis specification",
		derived:   axis,
		selection: selection(axis),
	})
}

func axis(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		if p, ok := f.(apl.Primitive); ok {
			return p.Call(a, L, apl.Axis{R: R, A: g})
		} else {
			return nil, fmt.Errorf("axis: expected primitive on the left: %T", f)
		}
	}
	return function(derived)
}
