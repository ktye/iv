package primitives

import (
	"reflect"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "≡",
		doc:    "depth, level of nesting",
		Domain: Monadic(nil),
		fn:     depth,
	})
	register(primitive{
		symbol: "≢",
		doc:    "tally, number of major cells",
		Domain: Monadic(nil),
		fn:     tally,
	})

	register(primitive{
		symbol: "≡",
		doc:    "match",
		Domain: Dyadic(nil),
		fn:     match,
	})
	register(primitive{
		symbol: "≢",
		doc:    "not match",
		Domain: Dyadic(nil),
		fn:     notmatch,
	})
}

// depth reports the level of nesting.
// Nested arrays are not supported, so depth is always 1 for arrays and 0 for scalars.
func depth(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.Array); ok {
		return apl.Index(1), nil
	}
	return apl.Index(0), nil
}

// tally returns the number of major cells of R.
// It is equlivalent to {⍬⍴(⍴⍵),1}.
func tally(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	ar, ok := R.(apl.Array)
	if ok == false {
		return apl.Index(1), nil
	}
	shape := ar.Shape()
	if len(shape) == 0 {
		return apl.Index(0), nil
	}
	return apl.Index(shape[0]), nil
}

func match(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	al, isal := L.(apl.Array)
	ar, isar := R.(apl.Array)
	if isal != isar {
		return apl.Bool(false), nil
	}
	if isal == false {
		// Compare scalars
		// arith2 converts numbers to the same type.
		feq := arith2("=", compare("="))
		return feq(a, L, R)
	} else {
		sl := al.Shape()
		sr := ar.Shape()
		if len(sr) != len(sl) {
			return apl.Bool(false), nil
		} else if len(sr) == 0 {
			// Empty arrays must have the same type.
			if reflect.TypeOf(ar) == reflect.TypeOf(al) {
				return apl.Bool(true), nil
			} else {
				return apl.Bool(false), nil
			}
		}
		for i := range sl {
			if sl[i] != sr[i] {
				return apl.Bool(false), nil
			}
		}
		feq := arith2("=", compare("="))
		for i := 0; i < apl.ArraySize(ar); i++ {
			rv, err := ar.At(i)
			if err != nil {
				return nil, err
			}
			lv, err := al.At(i)
			if err != nil {
				return nil, err
			}
			if iseq, err := feq(a, rv, lv); err != nil {
				return nil, err
			} else if iseq.(apl.Bool) == false {
				return apl.Bool(false), nil
			}
		}
		return apl.Bool(true), nil
	}
}

func notmatch(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if eq, err := match(a, L, R); err != nil {
		return nil, err
	} else {
		return !(eq.(apl.Bool)), nil
	}
}
