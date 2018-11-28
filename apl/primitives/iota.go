package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⍳",
		doc:    "interval, index generater, progression",
		Domain: Monadic(ToScalar(ToIndex(nil))),
		fn:     interval,
	})
	register(primitive{
		symbol: "⍳",
		doc: `index of, first occurrence of L in items of R
If an item is not found, the value is ⍴L+⎕IO.
If an item recurs: the value is the index of the first occurence`,
		Domain: Dyadic(Split(ToVector(nil), ToArray(nil))),
		fn:     indexof,
	})
	register(primitive{
		symbol: "∊",
		doc:    `membership`,
		Domain: Dyadic(nil),
		fn:     membership,
	})
	register(primitive{
		symbol: "⍸",
		doc:    "where",
		Domain: Monadic(ToIndexArray(nil)),
		fn:     where,
	})
}

// interval: R: integer. index generator.
func interval(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	n := int(R.(apl.Index))
	if n < 0 {
		return nil, fmt.Errorf("iota: L is negative")
	}
	if n == 0 {
		return apl.EmptyArray{}, nil
	}
	ar := apl.IndexArray{
		Ints: make([]int, n),
		Dims: []int{n},
	}
	for i := 0; i < n; i++ {
		ar.Ints[i] = a.Origin + i
	}
	return ar, nil
}

// indexof: L: vector, R: array
func indexof(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	al := L.(apl.Array) // vector
	ar := R.(apl.Array)

	nl := apl.ArraySize(al)
	notfound := nl + a.Origin
	vals := make([]apl.Value, nl)
	for i := range vals {
		v, _ := al.At(i)
		vals[i] = v
	}

	index := func(x apl.Value) int {
		for i := 0; i < nl; i++ {
			if ok := isEqual(a, x, vals[i]); ok {
				return i + a.Origin
			}
		}
		return notfound
	}

	ai := apl.IndexArray{
		Ints: make([]int, apl.ArraySize(ar)),
		Dims: apl.CopyShape(ar),
	}
	for i := range ai.Ints {
		v, err := ar.At(i)
		if err != nil {
			return nil, err
		}
		ai.Ints[i] = index(v)
	}
	return ai, nil
}

// membership. L and R may be arrays.
func membership(a *apl.Apl, L, R apl.Value) (apl.Value, error) {

	ar, ok := R.(apl.Array)
	if ok == false {
		ar = apl.GeneralArray{
			Dims:   []int{1},
			Values: []apl.Value{R},
		}
	}
	n := apl.ArraySize(ar)

	al, ok := L.(apl.Array)
	if !ok {
		// Scalar L: return a scalar boolean.
		for i := 0; i < n; i++ {
			v, err := ar.At(i)
			if err != nil {
				return nil, err
			}
			if isEqual(a, v, L) == true {
				return apl.Bool(true), nil
			}
		}
		return apl.Bool(false), nil
	}

	res := apl.IndexArray{
		Dims: apl.CopyShape(al),
		Ints: make([]int, apl.ArraySize(al)),
	}
	for k := range res.Ints {
		l, err := al.At(k)
		if err != nil {
			return nil, err
		}

		ok = false
		for i := 0; i < n; i++ {
			r, err := ar.At(i)
			if err != nil {
				return nil, err
			}
			if isEqual(a, l, r) == true {
				ok = true
				break
			}
		}
		if ok {
			res.Ints[k] = 1
		}
	}
	return res, nil
}

// where R is an IndexArray but only a boolean is allowed.
// In Dyalog where returns a nested index array for higher dimensional arrays.
// Here only vectors are supported.
func where(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	ar := R.(apl.IndexArray)
	shape := ar.Shape()
	if apl.ArraySize(ar) == 0 {
		return apl.EmptyArray{}, nil
	}

	if len(shape) != 1 {
		return nil, fmt.Errorf("where: only vectors are supported")
	}

	count := 0
	for _, n := range ar.Ints {
		if n == 1 {
			count++
		} else if n != 0 {
			return nil, fmt.Errorf("where: right argument must be a boolean array")
		}
	}

	res := apl.IndexArray{Dims: []int{count}, Ints: make([]int, count)}
	n := 0
	for i, v := range ar.Ints {
		if v == 1 {
			res.Ints[n] = i + a.Origin
			n++
		}
	}
	return res, nil
}

// IsEqual compares if the values are equal.
// If they are numbers of different type, they are converted before comparison.
func isEqual(a *apl.Apl, x, y apl.Value) bool {
	// TODO: should we use CT (comparison tolerance)?
	if x == y {
		return true
	}
	xn, isxnum := x.(apl.Number)
	yn, isynum := y.(apl.Number)
	if isxnum == false || isynum == false {
		return false
	}
	if xn, yn, err := a.Tower.SameType(xn, yn); err == nil && xn == yn {
		return true
	}
	return false
}
