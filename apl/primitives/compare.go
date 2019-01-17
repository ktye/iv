package primitives

import (
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
			doc:    e.doc,
			Domain: Dyadic(Split(IsScalar(nil), IsScalar(nil))),
			fn:     arith2(e.symbol, compare(e.symbol)),
		})
		register(primitive{
			symbol: e.symbol,
			doc:    e.doc,
			Domain: arrays{},
			fn:     array2(e.symbol, compare(e.symbol)),
		})
	}

	register(primitive{
		symbol: "<",
		doc:    "channel send, source",
		Domain: Monadic(nil),
		fn:     channelSource, // channel.go
	})
	register(primitive{
		symbol: "<",
		doc:    "channel copy, connect",
		Domain: Dyadic(Split(IsChannel(nil), IsChannel(nil))),
		fn:     channelCopy, // channel.go
	})
}

func compare(symbol string) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool) {
	return func(a *apl.Apl, L apl.Value, R apl.Value) (apl.Value, bool) {
		switch symbol {
		case "=":
			return equals(L, R)
		case "<":
			return less(L, R)
		case ">":
			eq, ls, ok := equalless(L, R)
			if ok == false {
				return nil, false
			}
			return apl.Bool(!eq && !ls), true
		case "≠":
			eq, ok := equals(L, R)
			if ok == false {
				return nil, false
			}
			return apl.Bool(!eq), true
		case "≤":
			eq, ls, ok := equalless(L, R)
			if ok == false {
				return nil, false
			}
			return apl.Bool(eq || ls), true
		case "≥":
			eq, ls, ok := equalless(L, R)
			if ok == false {
				return nil, false
			}
			return apl.Bool(eq || !ls), true
		}
		return nil, false
	}
}

func equalless(L, R apl.Value) (apl.Bool, apl.Bool, bool) {
	eq, ok := equals(L, R)
	if ok == false {
		return false, false, false
	}
	ls, ok := less(L, R)
	if ok == false {
		return false, false, false
	}
	return eq, ls, true
}
func equals(L, R apl.Value) (apl.Bool, bool) {
	if eq, ok := L.(equaler); ok {
		return eq.Equals(R)
	}
	return apl.Bool(L == R), true
}

type equaler interface {
	Equals(apl.Value) (apl.Bool, bool)
}

func less(L, R apl.Value) (apl.Bool, bool) {
	if ls, ok := L.(lesser); ok {
		return ls.Less(R)
	}
	return false, false
}

type lesser interface {
	Less(apl.Value) (apl.Bool, bool)
}
