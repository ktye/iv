package primitives

import (
	"fmt"
	"math/rand"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/numbers"
)

func init() {
	register(primitive{
		symbol: "?",
		doc:    "roll, rand, randn, bi-randn",
		Domain: Monadic(nil),
		fn:     roll,
	})
	register(primitive{
		symbol: "?",
		doc:    "deal",
		Domain: Dyadic(Split(ToScalar(ToIndex(nil)), ToScalar(ToIndex(nil)))),
		fn:     deal,
	})
}

// roll returns a number or an array of the same shape as R.
// Values of R must be numbers
func roll(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	ar, ok := R.(apl.Array)
	if ok == false {
		if z, ok := R.(apl.Number); ok == false {
			return nil, fmt.Errorf("roll expectes a numeric array as the right argument")
		} else {
			if n, err := rollNumber(a, z); err != nil {
				return nil, err
			} else {
				return n, nil
			}
		}
	}
	if ar.Size() == 0 {
		return apl.EmptyArray{}, nil
	}
	res := apl.MixedArray{
		Dims:   apl.CopyShape(ar),
		Values: make([]apl.Value, ar.Size()),
	}
	for i := range res.Values {
		if z, ok := ar.At(i).(apl.Number); ok == false {
			return nil, fmt.Errorf("roll: array value is not numeric")
		} else {
			if n, err := rollNumber(a, z); err != nil {
				return nil, err
			} else {
				res.Values[i] = n
			}
		}
	}
	return a.UnifyArray(res), nil
}

// rollNumber returns a random integer upto n, which must be integer.
// If n is 0, it returns a random float between 0 and 1.
// If n is negative it returns a random number from a normal distribution with std←|n.
// If n is complex, it returns a random number from a bivariate normal distribution
// with normal parameters given by the real an imag part.
//
// TODO: seed. Currently random numbers are not seeded (equivalent to Seed(1))
// and always return the same values.
func rollNumber(a *apl.Apl, n apl.Number) (apl.Number, error) {
	if f, ok := n.(numbers.Float); ok && float64(f) < 0 {
		return numbers.Float(float64(f) * rand.NormFloat64()), nil
	}
	if z, ok := n.(numbers.Complex); ok {
		return numbers.Complex(complex(real(z)*rand.NormFloat64(), imag(z)*rand.NormFloat64())), nil
	}
	m, ok := n.ToIndex()
	if ok == false || m < 0 {
		return numbers.Float(float64(m) * rand.NormFloat64()), nil
	}
	if m == 0 {
		// TODO: should we exclude 0?
		f := rand.Float64()
		return numbers.Float(f), nil // This only works with the default tower.
	} else {
		return a.Tower.Import(apl.Int(rand.Intn(m) + a.Origin)), nil
	}
}

// deal selects L random numbers from ⍳R without repetition.
func deal(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	// TODO []RL (random link)
	n := int(L.(apl.Int))
	m := int(R.(apl.Int))
	if n <= 0 || m < n {
		return nil, fmt.Errorf("deal: L must be > 0 and R >= L")
	}
	p := rand.Perm(m)
	p = p[:n]
	for i := range p {
		p[i] += a.Origin
	}
	return apl.IntArray{
		Ints: p,
		Dims: []int{n},
	}, nil
}
