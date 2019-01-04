package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(shape)
	register(reshape)
}

var shape = primitive{
	symbol: "⍴",
	doc:    "shape",
	Domain: Monadic(nil),
	fn:     rho1,
}

var reshape = primitive{
	symbol: "⍴",
	doc:    "reshape",
	Domain: Dyadic(Split(ToVector(ToIndexArray(nil)), ToArray(nil))),
	fn:     rho2,
	sel:    selection(rho2),
}

// Rho1 returns the shape of R.
func rho1(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	// Report a table as a two dimensional array.
	if t, ok := R.(apl.Table); ok == true {
		return apl.IndexArray{
			Dims: []int{2},
			Ints: []int{t.Rows, len(t.K)},
		}, nil
	}
	// An object returns the number of keys.
	if o, ok := R.(apl.Object); ok == true {
		n := len(o.Keys())
		return apl.IndexArray{Dims: []int{1}, Ints: []int{n}}, nil
	}

	if _, ok := R.(apl.Array); ok == false {
		return apl.EmptyArray{}, nil
	}
	// Shape of an empty array is 0, rank is 1
	if _, ok := R.(apl.EmptyArray); ok {
		return apl.IndexArray{Ints: []int{0}, Dims: []int{1}}, nil
	}
	ar := R.(apl.Array)
	shape := ar.Shape()
	ret := apl.IndexArray{
		Ints: make([]int, len(shape)),
		Dims: []int{len(shape)},
	}
	for i, n := range shape {
		ret.Ints[i] = n
	}
	return ret, nil
}

// Rho2 is dyadic reshape, L is empty or index array, R is array.
func rho2(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	// L is empty, returns empty.
	if apl.ArraySize(L.(apl.Array)) == 0 {
		return apl.EmptyArray{}, nil
	}

	if _, ok := R.(apl.Object); ok {
		return nil, fmt.Errorf("cannot reshape %T", R)
	}

	l := L.(apl.IndexArray)
	shape := make([]int, len(l.Ints))
	copy(shape, l.Ints)

	if rs, ok := R.(apl.Reshaper); ok {
		return rs.Reshape(shape), nil
	}
	return nil, fmt.Errorf("cannot reshape %T", R)
}
