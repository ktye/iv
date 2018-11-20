package primitives

import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	tab := []struct {
		symbol    string
		doc, doc2 string
		monadic   func(*apl.Apl, apl.Value) (apl.Value, bool)
		dyadic    func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool)
	}{
		{"+", "identity, complex conjugate", "plus, addition", add, add2},
		{"-", "reverse sign", "substract, substraction", sub, sub2},
		{"×", "signum, sign of, direction", "multiply", mul, mul2},
		{"÷", "reciprocal", "div, division, divide", div, div2},
		{"*", "exponential", "power", pow, pow2},
		{"⍟", "natural logarithm", "log, logarithm", log, log2},
		{"|", "magnitude, absolute value", "magnitude, absolute value", abs, abs2},
		{"⌊", "floor", "min, minumum", min, min2},
		{"⌈", "ceil", "max, maximum", max, max2},
	}

	for _, e := range tab {
		register(primitive{
			symbol: e.symbol,
			doc:    e.doc,
			Domain: Monadic(IsScalar(nil)),
			fn:     arith1(e.symbol, e.monadic),
		})
		register(primitive{
			symbol: e.symbol,
			doc:    e.doc,
			Domain: Monadic(IsArray(nil)),
			fn:     array1(e.symbol, e.monadic),
		})
		register(primitive{
			symbol: e.symbol,
			doc:    e.doc2,
			Domain: Dyadic(Split(IsScalar(nil), IsScalar(nil))),
			fn:     arith2(e.symbol, e.dyadic),
		})
		register(primitive{
			symbol: e.symbol,
			doc:    e.doc2,
			Domain: arrays{},
			fn:     array2(e.symbol, e.dyadic),
		})
	}
}

// arith1 tries to apply fn to the right argument.
// If it does not succeed directly, it tests if the argument is a number and uptypes until
// the function application succeeds.
func arith1(symbol string, fn func(*apl.Apl, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {

	return func(a *apl.Apl, _ apl.Value, R apl.Value) (apl.Value, error) {
		// Try to call the function directly.
		if res, ok := fn(a, R); ok {
			return res, nil
		}

		if v, ok := R.(apl.Index); ok {
			R = a.Tower.FromIndex(int(v))
		}

		// Uptype and call the function.
		n, ok := R.(apl.Number)
		if ok == false {
			return nil, fmt.Errorf("%s: not a numeric type %T", symbol, R)
		}

		num, ok := a.Tower.Numbers[reflect.TypeOf(n)]
		if ok == false {
			return nil, fmt.Errorf("%s: unknown numeric type %T", symbol, n)
		}

		for i := num.Class; i < len(a.Tower.Numbers); i++ {
			if res, ok := fn(a, n); ok {
				return res, nil
			}
			n, ok = num.Uptype(n)
			if ok == false {
				break
			}
			num = a.Tower.Numbers[reflect.TypeOf(n)]
		}
		return nil, fmt.Errorf("%s: not supported for %T", symbol, R)
	}
}

// arith2 tries to apply fn dyadically to the left and right argument.
// If they are not of the same type, it tests if the aruguments are numeric and
// uptypes to the same numeric type.
func arith2(symbol string, fn func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {

	return func(a *apl.Apl, L apl.Value, R apl.Value) (apl.Value, error) {
		// Try to call the function directly.
		if reflect.TypeOf(L) == reflect.TypeOf(R) {
			if res, ok := fn(a, L, R); ok {
				return res, nil
			}
		}

		if v, ok := L.(apl.Index); ok {
			L = a.Tower.FromIndex(int(v))
		}
		if v, ok := R.(apl.Index); ok {
			R = a.Tower.FromIndex(int(v))
		}

		ln, ok := L.(apl.Number)
		if ok == false {
			return nil, fmt.Errorf("%s: left argument is not a numeric type %T", symbol, L)
		}
		rn, ok := R.(apl.Number)
		if ok == false {
			return nil, fmt.Errorf("%s: right argument is not a numeric type %T", symbol, R)
		}
		var err error
		ln, rn, err = a.Tower.SameType(ln, rn)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", symbol, err)
		}

		num := a.Tower.Numbers[reflect.TypeOf(ln)]
		for i := num.Class; i < len(a.Tower.Numbers); i++ {
			if res, ok := fn(a, ln, rn); ok {
				return res, nil
			}
			ln, ok = num.Uptype(ln)
			if ok == false {
				break
			}
			rn, ok = num.Uptype(rn)
			num = a.Tower.Numbers[reflect.TypeOf(ln)]
		}
		return nil, fmt.Errorf("%s: not supported for types %T", symbol, L)
	}
}

// + add, add2
type adder interface {
	Add() (apl.Value, bool)
}

type adder2 interface {
	Add2(apl.Value) (apl.Value, bool)
}

func add(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	if d, ok := R.(adder); ok {
		return d.Add()
	}
	return nil, false
}
func add2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	if d, ok := L.(adder2); ok {
		return d.Add2(R)
	}
	return nil, false
}

// - sub, sub2
type substracter interface {
	Sub() (apl.Value, bool)
}

type substracter2 interface {
	Sub2(apl.Value) (apl.Value, bool)
}

func sub(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	if d, ok := R.(substracter); ok {
		return d.Sub()
	}
	return nil, false
}
func sub2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	if d, ok := L.(substracter2); ok {
		return d.Sub2(R)
	}
	return nil, false
}

// × mul, mul2
type multiplier interface {
	Mul() (apl.Value, bool)
}

type multiplier2 interface {
	Mul2(apl.Value) (apl.Value, bool)
}

func mul(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	if d, ok := R.(multiplier); ok {
		return d.Mul()
	}
	return nil, false
}
func mul2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	if d, ok := L.(multiplier2); ok {
		return d.Mul2(R)
	}
	return nil, false
}

// ÷ div, div2
type divider interface {
	Div() (apl.Value, bool)
}

type divider2 interface {
	Div2(apl.Value) (apl.Value, bool)
}

func div(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	if d, ok := R.(divider); ok {
		return d.Div()
	}
	return nil, false
}
func div2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	if d, ok := L.(divider2); ok {
		return d.Div2(R)
	}
	return nil, false
}

// * pow, pow2
type power interface {
	Pow() (apl.Value, bool)
}

type power2 interface {
	Pow2(apl.Value) (apl.Value, bool)
}

func pow(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	if d, ok := R.(power); ok {
		return d.Pow()
	}
	return nil, false
}
func pow2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	if d, ok := L.(power2); ok {
		return d.Pow2(R)
	}
	return nil, false
}

// ⍟ log, log2
type loger interface {
	Log() (apl.Value, bool)
}

type loger2 interface {
	Log2(apl.Value) (apl.Value, bool)
}

func log(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	if d, ok := R.(loger); ok {
		return d.Log()
	}
	return nil, false
}
func log2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	if d, ok := L.(loger2); ok {
		return d.Log2(R)
	}
	return nil, false
}

// | abs, abs2
type abser interface {
	Abs() (apl.Value, bool)
}

func abs(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	// Complex numbers should implement their own Abs method.
	if a, ok := R.(abser); ok {
		return a.Abs()
	}
	zero, r, err := a.Tower.SameType(a.Tower.FromIndex(0), R.(apl.Number))
	if err != nil {
		return nil, false
	}
	if isless, ok := less(r, zero); ok == false {
		return nil, false
	} else if isless {
		return sub(a, R)
	} else {
		return R, true
	}
}
func abs2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	// R-L×⌊R÷L+0=L
	zero, _, err := a.Tower.SameType(a.Tower.FromIndex(0), L.(apl.Number))
	if err != nil {
		return nil, false
	}
	L0 := zero // 0=L
	if Lzero, ok := equals(L.(apl.Number), zero); ok == false {
		return nil, false
	} else if Lzero {
		one, _, err := a.Tower.SameType(a.Tower.FromIndex(1), L.(apl.Number))
		if err != nil {
			return nil, false
		}
		L0 = one
	}
	x, ok := add2(a, L, L0)
	if ok == false {
		return nil, false
	}
	x, ok = div2(a, R, x)
	if ok == false {
		return nil, false
	}
	x, ok = min(a, x)
	if ok == false {
		return nil, false
	}
	x, ok = mul2(a, L, x)
	if ok == false {
		return nil, false
	}
	return sub2(a, R, x)
}

// ⌊ min, min2
type floorer interface {
	Floor() (apl.Value, bool)
}

// min returns the largest integer that is less or equal to R
func min(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	if floor, ok := R.(floorer); ok {
		return floor.Floor()
	}
	return nil, false
}
func min2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	if isless, ok := less(L, R); ok == false {
		return nil, false
	} else {
		if isless {
			return L, true
		} else {
			return R, true
		}
	}
}

// ⌈ max, max2
type ceiler interface {
	Ceil() (apl.Value, bool)
}

// max returns the smallest integer that is larger or equal to R
func max(a *apl.Apl, R apl.Value) (apl.Value, bool) {
	if ceil, ok := R.(ceiler); ok {
		return ceil.Ceil()
	}
	return nil, false
}
func max2(a *apl.Apl, L, R apl.Value) (apl.Value, bool) {
	if isless, ok := less(L, R); ok == false {
		return nil, false
	} else {
		if isless {
			return R, true
		} else {
			return L, true
		}
	}
}
