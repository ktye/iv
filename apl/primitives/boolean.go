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
}

func logical(logical string) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool) {
	return func(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
		boolean := func(v apl.Value) (bool, bool) {
			if n, ok := v.(apl.Number); ok == false {
				return false, false
			} else if m, ok := n.ToIndex(); ok == false {
				return false, false
			} else if m < 0 || m > 1 {
				return false, false
			} else {
				return m == 1, true
			}
		}
		l, ok := boolean(L)
		if ok == false {
			return nil, false
		}
		r, ok := boolean(R)
		if ok == false {
			return nil, false
		}
		var t bool
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
		return apl.Bool(t), true
	}
}
