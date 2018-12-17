package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/operators"
)

func init() {
	register(primitive{
		symbol: "↑",
		doc:    "take",
		Domain: Dyadic(Split(ToIndexArray(nil), nil)),
		fn:     take,
		sel:    takeSelection,
	})
	register(primitive{
		symbol: "↓",
		doc:    "drop",
		Domain: Dyadic(Split(ToIndexArray(nil), nil)),
		fn:     drop,
		sel:    dropSelection,
	})
}

func take(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	return takedrop(a, L, R, true)
}
func drop(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	return takedrop(a, L, R, false)
}

func takeSelection(a *apl.Apl, L, R apl.Value) (apl.IndexArray, error) {
	v, err := takeDropSelection(a, L, R, true)
	return v, err
}
func dropSelection(a *apl.Apl, L, R apl.Value) (apl.IndexArray, error) {
	return takeDropSelection(a, L, R, false)
}

// takedrop does the preprocessing, that is common to both take and drop.
func takedrop(a *apl.Apl, L, R apl.Value, take bool) (apl.Value, error) {
	// Special case, L is the empty array, return R.
	if _, ok := L.(apl.EmptyArray); ok {
		return R, nil
	}

	var x []int
	var err error
	R, x, err = splitAxis(a, R)
	if err != nil {
		return nil, err
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
	rs := ar.Shape()

	if take == false {
		// Missing items in L default to 0.
		if n := len(rs) - ai.Dims[0]; n > 0 {
			zeros := make([]int, n)
			ai.Ints = append(ai.Ints, zeros...)
			ai.Dims[0] = len(ai.Ints)
		}
	}

	// APL2: Default axis: L↑R ←→ L↑[⍳⍴⍴R]R, same for drop.
	// Is this true? or should that be L↑[⍳⍴⍴L]R ?
	if x == nil {
		x = make([]int, len(ai.Ints))
		for i := range x {
			x[i] = i
		}
	}

	// X∊⍳⍴⍴R, the order is not important, but no repetitions.
	m := make(map[int]bool)
	for _, n := range x {
		if m[n] {
			return nil, fmt.Errorf("axis has repeated elements")
		}
		m[n] = true
		if n < 0 || n >= len(rs) {
			return nil, fmt.Errorf("axis is out of range")
		}
	}

	// ⍴L must match the length of x
	if len(ai.Ints) != len(x) {
		return nil, fmt.Errorf("take/drop: length of L must match length of axis vector")
	}

	if take {
		// Take is defined in opearators/rank.go
		return operators.Take(a, ai, ar, x)
	} else {
		return dodrop(a, ai, ar, x)
	}
}

func dodrop(a *apl.Apl, L apl.IndexArray, R apl.Array, x []int) (apl.Value, error) {
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
	return operators.Take(a, b, R, x)
}

func takeDropSelection(a *apl.Apl, L, R apl.Value, take bool) (apl.IndexArray, error) {
	var x []int
	var err error
	R, x, err = splitAxis(a, R)
	if err != nil {
		return apl.IndexArray{}, err
	}

	ar, ok := R.(apl.Array)
	if ok == false {
		return apl.IndexArray{}, fmt.Errorf("cannot select from non-array: %T", R)
	}

	// Take/drop from an index array instead of R of the same shape.
	// Take/drop fills with zeros, so count with origin 1 temporarily.
	r := apl.IndexArray{Dims: apl.CopyShape(ar)}
	r.Ints = make([]int, apl.ArraySize(r))
	for i := range r.Ints {
		r.Ints[i] = i + 1
	}

	R = r
	if x != nil {
		for i := range x {
			x[i] += a.Origin
		}
		R = apl.Axis{R: r, A: apl.IndexArray{Dims: []int{len(x)}, Ints: x}}
	}

	var ai apl.IndexArray
	res, err := takedrop(a, L, R, take)
	if err != nil {
		return ai, err
	}

	to := ToIndexArray(nil)
	if v, ok := to.To(a, res); ok == false {
		return ai, fmt.Errorf("could not convert selection to index array: %T", res)
	} else {
		ai = v.(apl.IndexArray)
	}

	for i := range ai.Ints {
		ai.Ints[i]--

		// TODO: Elements < 0 are the result of overtake.
		// These elements should be removed.
		if ai.Ints[i] < 0 {
			return ai, fmt.Errorf("TODO: overtake/drop with selection")
		}
	}
	return ai, nil
}
