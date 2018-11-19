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
		doc:    "roll",
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
	if apl.ArraySize(ar) == 0 {
		return apl.EmptyArray{}, nil
	}
	res := apl.GeneralArray{
		Values: make([]apl.Value, apl.ArraySize(ar)),
		Dims:   apl.CopyShape(ar),
	}
	for i := range res.Values {
		if v, err := ar.At(i); err != nil {
			return nil, err
		} else if z, ok := v.(apl.Number); ok == false {
			return nil, fmt.Errorf("roll: array value is not numeric")
		} else {
			if n, err := rollNumber(a, z); err != nil {
				return nil, err
			} else {
				res.Values[i] = n
			}
		}
	}
	return res, nil
}

// rollNumber returns a random integer upto n, which must be integer.
// If n is 0, it returns a random float between 0 and 1.
func rollNumber(a *apl.Apl, n apl.Number) (apl.Number, error) {
	m, ok := n.ToIndex()
	if ok == false || m < 0 {
		return nil, fmt.Errorf("roll: values of R must be integer > 0: %T", n)
	}
	if m == 0 {
		// TODO: should we exclude 0?
		f := rand.Float64()
		return numbers.Float(f), nil // This only works with the default tower.
	} else {
		return a.Tower.FromIndex(rand.Intn(m) + a.Origin), nil
	}
}

// deal selects L random numbers from ⍳R without repetition.
func deal(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	// TODO []RL (random link)
	n := int(L.(apl.Index))
	m := int(R.(apl.Index))
	if n <= 0 || m < n {
		return nil, fmt.Errorf("deal: L must be > 0 and R >= L")
	}
	p := rand.Perm(m)
	p = p[:n]
	for i := range p {
		p[i] += a.Origin
	}
	return apl.IndexArray{
		Ints: p,
		Dims: []int{n},
	}, nil
}
