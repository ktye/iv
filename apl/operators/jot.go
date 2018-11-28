package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "∘",
		Domain:  DyadicOp(Split(nil, nil)),
		doc:     "compose",
		derived: compose,
	})
}

func compose(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		fn, isfunc := f.(apl.Function)
		gn, isgunc := g.(apl.Function)
		if L == nil {
			if isfunc && isgunc {
				// Form 1: f∘g R
				v, err := gn.Call(a, nil, R)
				if err != nil {
					return nil, err
				}
				return fn.Call(a, nil, v)
			} else if isgunc {
				// Form II: A∘g R
				return gn.Call(a, f, R)
			} else if isfunc {
				// Form III: (f∘X) R
				return fn.Call(a, R, g)
			}
		} else {
			if isfunc && isgunc {
				// Form IV: L f∘g R
				v, err := gn.Call(a, nil, R)
				if err != nil {
					return nil, err
				}
				return fn.Call(a, L, v)
			}
		}
		return nil, fmt.Errorf("compose: cannot handle %T %T ∘ %T %T", L, f, g, R)
	}
	return function(derived)
}
