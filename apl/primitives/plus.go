package primitives

/* TODO remove
import (
	"fmt"
	"math/cmplx"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(identity)
	register(conjugate)
	register(plusarray)
	register(plus)
}

var identity = primitive{
	symbol: "+",
	doc:    "identity",
	Domain: Monadic(ToNumber(Not(IsComplex(nil)))),
	fn: func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		return R, nil
	},
}

var conjugate = primitive{
	symbol: "+",
	doc:    "complex conjugate",
	Domain: Monadic(ToNumber(IsComplex(nil))),
	fn: func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		return apl.Complex(cmplx.Conj(complex128(R.(apl.Complex)))), nil
	},
}

var plusarray = primitive{
	symbol: "+",
	doc:    "plus, addition",
	Domain: Monadic(IsArray(nil)),
	fn:     arithmonads(identity, conjugate),
}

var plus = primitive{
	symbol: "+",
	doc:    "plus, addition",
	Domain: arithmetic{},
	fn:     arith(elementaryPlus),
}

func elementaryPlus(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	switch L.(type) {
	case apl.Bool:
		// We always uptype to Int for +.
		l, r := bools(L, R)
		if l && r {
			return apl.Int(2), nil
		}
		if l || r {
			return apl.Int(1), nil
		}
		return apl.Int(0), nil
	case apl.Int:
		l, r := ints(L, R)
		return apl.Int(l + r), nil
	case apl.Float:
		l, r := floats(L, R)
		return apl.Float(l + r), nil
	case apl.Complex:
		l, r := complexs(L, R)
		return apl.Complex(l + r), nil
	}
	return nil, fmt.Errorf("impossible: elementary + received: %T %T", L, R)
}

// Helper functions for elementary operations.
// The conversion is not checked, that must have been done before.
func bools(L, R apl.Value) (bool, bool) {
	return bool(L.(apl.Bool)), bool(R.(apl.Bool))
}
func ints(L, R apl.Value) (int, int) {
	return int(L.(apl.Int)), int(R.(apl.Int))
}
func floats(L, R apl.Value) (float64, float64) {
	return float64(L.(apl.Float)), float64(R.(apl.Float))
}
func complexs(L, R apl.Value) (complex128, complex128) {
	return complex128(L.(apl.Complex)), complex128(R.(apl.Complex))
}
*/
