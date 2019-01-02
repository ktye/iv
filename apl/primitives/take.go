package primitives

import (
	"fmt"
	"sort"

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
		symbol: "↑",
		doc:    "take from channel",
		Domain: Dyadic(Split(ToIndexArray(nil), IsChannel(nil))),
		fn:     takeChannel2,
	})
	register(primitive{
		symbol: "↑",
		doc:    "take one from channel",
		Domain: Monadic(IsChannel(nil)),
		fn:     takeChannel1,
	})
	register(primitive{
		symbol: "↓",
		doc:    "drop",
		Domain: Dyadic(Split(ToIndexArray(nil), nil)),
		fn:     drop,
		sel:    dropSelection,
	})
	register(primitive{
		symbol: "↓",
		doc:    "drop to channel",
		Domain: Dyadic(Split(IsChannel(nil), nil)),
		fn:     sendChannel,
	})
	register(primitive{
		symbol: "↓",
		doc:    "close channel",
		Domain: Monadic(IsChannel(nil)),
		fn:     closeChannel,
	})
	register(primitive{
		symbol: "↓",
		doc:    "cut",
		Domain: Dyadic(Split(ToIndexArray(nil), IsList(nil))),
		fn:     cut,
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
	if len(ai.Dims) > 1 {
		return nil, fmt.Errorf("take/drop: L must be a vector")
	}

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

	// If R is a scalar, set it's shape to (⍴,L)⍴1.
	ar, ok := R.(apl.Array)
	if ok == false {
		r := apl.MixedArray{Values: []apl.Value{R}} // TODO copy?
		r.Dims = make([]int, len(ai.Ints))
		for i := range r.Dims {
			r.Dims[i] = 1
		}
		ar = r
	}
	rs := ar.Shape()

	// The default axis is the shape list of R: L↑R ←→ L↑[⍳⍴⍴R]R, same for drop.
	// Shorter axis are filled with the missing default items.
	// Elements may not repeat.
	if len(x) > len(rs) {
		return nil, fmt.Errorf("axis is too long")
	}
	m := make(map[int]int)
	for i := range rs {
		m[i] = i
	}
	axis := make([]int, len(rs))
	for i, n := range x {
		if k, ok := m[n]; ok == false {
			return nil, fmt.Errorf("axis does not conform")
		} else {
			axis[i] = k
			delete(m, k)
		}
	}
	tail := make([]int, len(m))
	i := 0
	for _, n := range m {
		tail[i] = n
		i++
	}
	sort.Ints(tail)
	copy(axis[len(x):], tail)
	x = axis

	// Missing items in L default to values of ⍴R[x] for take and 0 for drop.
	if len(ai.Ints) > len(rs) {
		return nil, fmt.Errorf("take/drop: length of L is too large")
	} else if len(ai.Ints) < len(rs) {
		n := make([]int, len(rs))
		copy(n, ai.Ints)
		if take {
			for i := len(ai.Ints); i < len(rs); i++ {
				n[i] = rs[x[i]]
			}
		}
		ai.Ints = n
		ai.Dims[0] = len(n)
	}

	if take == false {
		// Convert L to the left argument of an equivalent call to take.
		for i, n := range ai.Ints {
			m := rs[x[i]]
			if n >= m || -n >= m { // over drop
				ai.Ints[i] = 0
			} else if n > 0 {
				ai.Ints[i] = n - m
			} else {
				ai.Ints[i] = n + m
			}
		}
	}
	// Take is defined in opearators/rank.go
	return operators.Take(a, ai, ar, x)
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

// Cut list R at indexes L.
// This is similar to _ in q.
// Indexes may be negative.
func cut(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	ai := L.(apl.IndexArray)
	r := R.(apl.List)
	if len(ai.Shape()) != 1 {
		return nil, fmt.Errorf("cut: left argument must be an index vector")
	}
	idx := make([]int, len(ai.Ints))
	for i := range idx {
		idx[i] = ai.Ints[i] - a.Origin
		if idx[i] < 0 {
			idx[i] = len(r) + idx[i]
		}
		if i > 0 && idx[i] <= idx[i-1] {
			return nil, fmt.Errorf("cut: indexes may not decrease")
		}
		if idx[i] < 0 || idx[i] >= len(r) {
			return nil, fmt.Errorf("cut: indexes out of range")
		}
	}
	if len(idx) == 1 {
		return r[idx[0]:], nil // TODO: copy
	}
	res := make(apl.List, len(idx))
	for i := range res {
		stp := len(r)
		if i < len(idx)-1 {
			stp = idx[i+1]
		}
		res[i] = r[idx[i]:stp] // TODO: copy
	}
	return res, nil
}

func takeChannel1(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	c := R.(apl.Channel)
	v, ok := <-c[0]
	if ok == false {
		return nil, fmt.Errorf("channel is closed")
	}
	return v, nil
}

// takeChannel2 takes multiple values from channel R and reshapes according to L.
func takeChannel2(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	ai := L.(apl.IndexArray)
	if len(ai.Shape()) != 1 {
		return nil, fmt.Errorf("take channel: L must be an index vector")
	}
	for _, n := range ai.Ints {
		if n <= 0 {
			return nil, fmt.Errorf("take channel: values in L must be positive")
		}
	}
	shape := make([]int, len(ai.Ints))
	copy(shape, ai.Ints)
	res := apl.MixedArray{Dims: shape}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	c := R.(apl.Channel)
	for i := range res.Values {
		v, ok := <-c[0]
		if ok == false {
			return nil, fmt.Errorf("not enough data in channel")
		}
		res.Values[i] = v
	}
	return res, nil
}

// sendChannel sends the value R to the channel L and returns R.
func sendChannel(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	c := L.(apl.Channel)
	c[1] <- R // TODO: copy?
	return R, nil
}

// closeChannel closes the channel R and returns 1.
func closeChannel(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	R.(apl.Channel).Close()
	return apl.Index(1), nil
}
