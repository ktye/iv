package primitives

import (
	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "←",
		doc:    "sink\nreturns empty array to suppress printing",
		Domain: Monadic(nil),
		fn: func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
			return apl.EmptyArray{}, nil
		},
	})
	register(primitive{
		symbol: "←",
		doc:    "assign, assignment",
		Domain: Dyadic(Split(IsIdentifier(nil), nil)),
		fn: func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
			if err := a.Assign(string(L.(apl.Identifier)), R); err != nil {
				return nil, err
			}
			return R, nil
		},
	})
}
