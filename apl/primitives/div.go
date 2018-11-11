package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	/*
		register(primitive{
			symbol: "÷",
			doc:    "reciprocal",
			Domain: Monadic(ToNumber(nil)),
			fn:     reciprocal,
		})
	*/
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
	fn:     arithmonads(reciprocal),
}

var div = primitive{
	symbol: "÷",
	doc:    "div, division, divide",
	Domain: arithmetic{},
	fn:     arith(elementaryDiv),
}

func oneby(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	switch v := R.(type) {
	case apl.Bool:
		return nil, fmt.Errorf("cannot ÷bool")
	case apl.Int:
		return apl.Float(1.0 / float64(v)), nil
	case apl.Float:
		return apl.Float(1.0 / float64(v)), nil
	case apl.Complex:
		return apl.Complex(complex(1, 0) / complex128(v)), nil
	default:
		return nil, fmt.Errorf("reciprocal: unknown %T", R)
	}
}

func elementaryDiv(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	var f1, f2 float64
	switch L.(type) {
	case apl.Bool:
		l, r := bools(L, R)
		if l {
			f1 = 1
		}
		if r {
			f2 = 1
		}
	case apl.Int:
		l, r := ints(L, R)
		f1, f2 = float64(l), float64(r)
	case apl.Float:
		l, r := floats(L, R)
		f1, f2 = l, r
	case apl.Complex:
		l, r := complexs(L, R)
		res := apl.Complex(l / r)
		return res, IsFloatErr(res)
	default:
		return nil, fmt.Errorf("impossible: elementary + received: %T %T", L, R)
	}
	res := apl.Float(f1 / f2)
	return res, IsFloatErr(res)
}
