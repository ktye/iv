package primitives

import (
	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⊥",
		doc:    "decode, polynom, base value",
		Domain: Dyadic(Split(ToNumber(nil), ToVector(nil))),
		fn:     poly,
	})
	// TODO: other cases, APL2 p90
}

// poly evaluates a polynom. L is a number and R a vector or empty.
func poly(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.EmptyArray); ok {
		return apl.EmptyArray{}, nil
	}

	ar := R.(apl.Array)
	fpow := arith2("*", pow2)
	fmul := arith2("×", mul2)
	fadd := arith2("+", add2)
	sum := apl.Value(a.Tower.FromIndex(0))
	n := apl.ArraySize(ar)
	var val apl.Value
	for i := 0; i < n; i++ {
		coeff, err := ar.At(i)
		if err != nil {
			return nil, err
		}
		exponent := n - 1 - i
		val, err = fpow(a, L, a.Tower.FromIndex(exponent))
		if err != nil {
			return nil, err
		}
		val, err = fmul(a, coeff, val)
		if err != nil {
			return nil, err
		}
		sum, err = fadd(a, sum, val)
		if err != nil {
			return nil, err
		}
	}
	return sum, nil
}
