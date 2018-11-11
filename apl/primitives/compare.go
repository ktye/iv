package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	tab := []struct {
		symbol, doc string
	}{
		{"=", "equality"},
		{"<", "less that"},
		{">", "greater than"},
		{"≠", "not equal"},
		{"≤", "less or equal"},
		{"≥", "greater or equal"},
	}
	for _, e := range tab {
		register(primitive{
			symbol: e.symbol,
			doc:    "compare " + e.doc,
			Domain: arithmetic{},
			fn:     arith(comparator(e.symbol)),
		})
	}
}

func comparator(symbol string) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	return func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		eq, lt, err := CompareScalars(L, R)
		if err != nil {
			return nil, err
		}
		switch symbol {
		case "=":
			return apl.Bool(eq), nil
		case "<":
			return apl.Bool(lt), nil
		case ">":
			return apl.Bool(!eq && !lt), nil
		case "≠":
			return apl.Bool(!eq), nil
		case "≤":
			return apl.Bool(eq || lt), nil
		case "≥":
			return apl.Bool(!lt), nil
		default:
			return apl.Bool(false), fmt.Errorf("illegal comparision operator: %s", symbol)
		}
	}
}
