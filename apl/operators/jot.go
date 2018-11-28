package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "∘",
		Domain:  DyadicOp(Split(Function(nil), Function(nil))),
		doc:     "compose",
		derived: compose,
	})
}

// TODO: compose is added to register the jot symbol for outer product.
func compose(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
		//f := f.(apl.Function)
		//g := g.(apl.Function)
		return nil, fmt.Errorf("TODO compose (jot)")
	}
	return function(derived)
}
