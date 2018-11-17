package primitives

/* TODO remove

import (
	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(reciprocal)
	register(reciprocalarray)
}

var reciprocal = primitive{
	symbol: "÷",
	doc:    "reciprocal",
	Domain: Monadic(ToNumber(nil)),
	fn:     oneby,
}

var reciprocalarray = primitive{
	symbol: "÷",
	doc:    "reciprocal",
	Domain: Monadic(IsArray(nil)),
	fn:     monadic(oneby),
}

var div = primitive{
	symbol: "÷",
	doc:    "div, division, divide",
	Domain: arithmetic{},
	fn:     arith(divideby),
}

type divider interface {
	Div() (apl.Value, bool)
}

type divider2 interface {
	Div2(b apl.Value) (apl.Value, bool)
}

func oneby(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return uptype(a, R, "÷", func(n apl.Value) (apl.Value, bool) {
		if d, ok := n.(divider); ok {
			return d.Div()
		}
		return nil, false
	})
}

func divideby(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	return uptype2(a, L, R, func(n1, n2 apl.Value) (apl.Value, bool) {
		if d, ok := n1.(divider2); ok {
			return d.Div2(n2)
		}
		return nil, false
	})
}
*/
