package operators

import (
	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "⍨",
		Domain:  MonadicOp(Function(nil)),
		doc:     "commute, duplicate",
		derived: commute,
	})
}

func commute(a *apl.Apl, f, _ apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		f := f.(apl.Function)
		if L == nil {
			L = R.Copy()
		}
		return f.Call(a, R, L)
	}
	return function(derived)
}
