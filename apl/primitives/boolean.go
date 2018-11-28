package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/operators"
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
	register(primitive{
		symbol: "~",
		doc:    "without",
		Domain: Dyadic(Split(ToVector(nil), ToVector(nil))),
		fn:     without,
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
		l, lok := a.Tower.ToBool(L.(apl.Number))
		r, rok := a.Tower.ToBool(R.(apl.Number))
		if lok == false || rok == false {
			if logical == "and" {
				return lcm(a, L, R)
			} else if logical == "or" {
				return gcd(a, L, R)
			}
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

// without: L and R are vectors.
// L~R is equivalent to (~L∊R)/L.
func without(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.EmptyArray); ok {
		return L, nil
	}
	if _, ok := L.(apl.EmptyArray); ok {
		return apl.EmptyArray{}, nil
	}

	lr, err := membership(a, L, R)
	if err != nil {
		return nil, err
	}

	not := arith1("~", logicalNot)
	if _, ok := L.(apl.Array); ok {
		not = array1("~", logicalNot)
	}
	nlr, err := not(a, nil, lr)
	if err != nil {
		return nil, err
	}

	to := ToIndexArray(nil)
	ia, ok := to.To(a, nlr)
	if ok == false {
		return nil, fmt.Errorf("without: cannot convert (~L∊R) to index array")
	}

	return operators.Replicate(a, ia, L, 0)
}
