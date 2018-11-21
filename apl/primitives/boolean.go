package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	tab := []struct {
		symbol  string
		doc     string
		logical string
	}{
		{"^", "logical and", "and"},
		{"∧", "logical and", "and"},
		{"∨", "logical or", "or"},
		{"⍲", "logical nand", "nand"},
		{"⍱", "logical nor", "nor"},
	}

	// TODO least common multiply: dyadic ^, if L or R are not bool
	// TODO greatest common divisor: dyadic ∨, if L or R are not bool

	for _, e := range tab {
		register(primitive{
			symbol: e.symbol,
			doc:    e.doc,
			Domain: Dyadic(Split(IsScalar(nil), IsScalar(nil))),
			fn:     arith2(e.symbol, logical(e.logical)),
		})
		register(primitive{
			symbol: e.symbol,
			doc:    e.doc,
			Domain: arrays{},
			fn:     array2(e.symbol, logical(e.logical)),
		})
	}
	register(primitive{
		symbol: "~",
		doc:    "logical not",
		Domain: Monadic(IsScalar(nil)),
		fn:     arith1("~", logicalNot),
	})
	register(primitive{
		symbol: "~",
		doc:    "logical not",
		Domain: Monadic(IsArray(nil)),
		fn:     array1("~", logicalNot),
	})
}

// logical not, R is a Number.
func logicalNot(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	b, ok := a.Tower.ToBool(R.(apl.Number))
	if ok == false {
		return nil, false
	}
	return apl.Bool(!b), true
}

func logical(logical string) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool) {
	return func(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
		l, ok := a.Tower.ToBool(L.(apl.Number))
		if ok == false {
			return nil, false
		}
		r, ok := a.Tower.ToBool(R.(apl.Number))
		if ok == false {
			return nil, false
		}
		var t apl.Bool
		switch logical {
		case "and":
			t = l && r
		case "or":
			t = l || r
		case "nand":
			t = !(l && r)
		case "nor":
			t = !(l || r)
		default:
			panic(fmt.Sprintf("unknown logical: %s", logical))
		}
		return t, true
	}
}
