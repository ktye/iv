package primitives

import (
	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⊣",
		doc:    "left tack, same",
		Domain: Monadic(nil),
		fn:     same,
	})
	register(primitive{
		symbol: "⊢",
		doc:    "right tack, same",
		Domain: Monadic(nil),
		fn:     same,
	})
	register(primitive{
		symbol: "⊣",
		doc:    "left tack, left argument",
		Domain: Dyadic(nil),
		fn:     left,
	})
	register(primitive{
		symbol: "⊢",
		doc:    "right tack, right argument",
		Domain: Dyadic(nil),
		fn:     right,
	})
}

func same(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return R, nil
}
func left(a *apl.Apl, L, _ apl.Value) (apl.Value, error) {
	return L, nil
}
func right(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return R, nil
}
