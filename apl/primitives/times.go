package primitives

import (
	"fmt"
	"math/cmplx"
	"strings"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(signum)
	register(direction)
	register(multiply)
	register(strrepeat1)
	register(strrepeat2)
}

var signum = primitive{
	symbol: "×",
	doc:    "signum, sign of",
	Domain: Monadic(ToFloat(nil)),
	fn:     sign,
}

var direction = primitive{
	symbol: "×",
	doc:    "direction",
	Domain: Monadic(ToNumber(IsComplex(nil))),
	fn: func(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
		c := complex128(R.(apl.Complex))
		if c == 0 {
			return apl.Complex(c), nil
		}
		r := cmplx.Abs(c)
		return apl.Complex(complex(real(c)/r, imag(c)/r)), nil
	},
}

var multiply = primitive{
	symbol: "×",
	doc:    "multiply",
	Domain: arithmetic{},
	fn:     arith(elementaryTimes),
}

var strrepeat1 = primitive{
	symbol: "×",
	doc:    "repeat strings n times",
	Domain: Dyadic(Split(IsString(nil), ToInt(nil))),
	fn:     stringTimesInt,
}
var strrepeat2 = primitive{
	symbol: "×",
	doc:    "repeat strings n times",
	Domain: Dyadic(Split(ToInt(nil), IsString(nil))),
	fn:     intTimesString,
}

func sign(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	f := R.(apl.Float)
	if f > 0 {
		return apl.Int(1), nil
	} else if f < 0 {
		return apl.Int(-1), nil
	}
	return apl.Int(0), nil
}

func elementaryTimes(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	switch L.(type) {
	case apl.Bool:
		// We always uptype to Int for ×.
		l, r := bools(L, R)
		if l && r {
			return apl.Int(1), nil
		}
		return apl.Int(0), nil
	case apl.Int:
		l, r := ints(L, R)
		return apl.Int(l * r), nil
	case apl.Float:
		l, r := floats(L, R)
		return apl.Float(l * r), nil
	case apl.Complex:
		l, r := complexs(L, R)
		return apl.Complex(l * r), nil
	}
	return nil, fmt.Errorf("impossible: elementary + received: %T %T", L, R)
}

func stringTimesInt(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	s := string(L.(apl.String))
	n := int(R.(apl.Int))
	return strrepeat(s, n)
}
func intTimesString(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	n := int(L.(apl.Int))
	s := string(R.(apl.String))
	return strrepeat(s, n)
}

// strrepeat catenates the string n times.
func strrepeat(s string, n int) (apl.Value, error) {
	return apl.String(strings.Repeat(s, n)), nil
}
