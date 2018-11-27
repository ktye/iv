package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "↑",
		doc:    "take",
		Domain: Dyadic(Split(ToIndexArray(nil), nil)),
		fn:     take,
	})
	register(primitive{
		symbol: "↓",
		doc:    "drop",
		Domain: Dyadic(Split(ToIndexArray(nil), nil)),
		fn:     drop,
	})
}

func take(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	return takedrop(a, L, R, true)
}
func drop(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	return takedrop(a, L, R, false)
}

// takedrop does the preprocessing, that is common to both take and drop.
func takedrop(a *apl.Apl, L, R apl.Value, take bool) (apl.Value, error) {
	// Special case, L is the empty array, return R.
	if _, ok := L.(apl.EmptyArray); ok {
		return R, nil
	}

	ai := L.(apl.IndexArray)

	// If R is an empty array, return 0s of the size of |L.
	if _, ok := R.(apl.EmptyArray); ok {
		if len(ai.Ints) == 1 {
			n := ai.Ints[0]
			if n < 0 {
				n = -n
			}
			return apl.IndexArray{
				Ints: make([]int, n),
				Dims: []int{n},
			}, nil
		}
	}

	ar, ok := R.(apl.Array)

	if len(ai.Dims) > 1 {
		return nil, fmt.Errorf("take/drop: L must be a vector")
	}

	// If R is a scalar, set it's shape to (⍴,L)⍴1.
	if ok == false {
		r := apl.GeneralArray{Values: []apl.Value{R}} // TODO copy?
		r.Dims = make([]int, len(ai.Ints))
		for i := range r.Dims {
			r.Dims[i] = 1
		}
		ar = r
	}

	if take == false {
		// Missing items in L default to 0.
		if n := len(ar.Shape()) - ai.Dims[0]; n > 0 {
			zeros := make([]int, n)
			ai.Ints = append(ai.Ints, zeros...)
			ai.Dims[0] = len(ai.Ints)
		}
	}

	// ⍴L must match the rank of ar.
	if ai.Dims[0] != len(ar.Shape()) {
		return nil, fmt.Errorf("take/drop: ⍴,L must match ⍴⍴R")
	}

	if take {
		return dotake(a, ai, ar)
	} else {
		return dodrop(a, ai, ar)
	}
}

func dotake(a *apl.Apl, ai apl.IndexArray, ar apl.Array) (apl.Value, error) {
	rs := ar.Shape()

	// The shape of the result is ,|L
	shape := make([]int, len(ai.Ints))
	for i, n := range ai.Ints {
		if n < 0 {
			shape[i] = -n
		} else {
			shape[i] = n
		}
	}
	res := apl.GeneralArray{Dims: shape}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	ic, J := apl.NewIdxConverter(rs)
	idx := make([]int, len(shape))
	for i := range res.Values {
		for k := range J {
			J[k] = idx[k]
			if n := ai.Ints[k]; n < 0 {
				J[k] += n + rs[k]
			}
		}
		iszero := false
		for k := range J {
			if J[k] < 0 || J[k] >= rs[k] {
				iszero = true
				break
			}
		}
		if iszero {
			res.Values[i] = apl.Index(0) // TODO: typical element of R?
		} else {
			n := ic.Index(J)
			v, err := ar.At(n)
			if err != nil {
				return nil, err
			}
			res.Values[i] = v // TODO: copy?
		}

		apl.IncArrayIndex(idx, shape)
	}
	return res, nil
}

func dodrop(a *apl.Apl, L apl.IndexArray, R apl.Array) (apl.Value, error) {
	// (((L<0)×0⌈L+⍴R)+(L≥0)×x0⌊L-⍴R) ↑R
	b := apl.IndexArray{
		Dims: apl.CopyShape(L),
		Ints: make([]int, len(L.Ints)),
	}
	rs := R.Shape()
	for i := range b.Ints {
		l := L.Ints[i]
		r := rs[i]
		t := l - r // L-⍴R
		if t > 0 {
			t = 0 // 0⌊L-⍴R
		}
		if l < 0 {
			t = 0 // (L≥0)×x0⌊L-⍴R
		}
		s := l + r // L+⍴R
		if s < 0 {
			s = 0 // 0⌈L+⍴R
		}
		if l >= 0 {
			s = 0 // ((L<0)×0⌈L+⍴R)
		}
		b.Ints[i] = s + t // (((L<0)×0⌈L+⍴R)+(L≥0)×x0⌊L-⍴R)
	}
	return dotake(a, b, R)
}
