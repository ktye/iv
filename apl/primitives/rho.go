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

var reshape = primitive{
	symbol: "⍴",
	doc:    "reshape",
	Domain: Dyadic(Split(ToVector(ToIntArray(nil)), ToArray(nil))),
	fn:     rho,
}

var shape = primitive{
	symbol: "⍴",
	doc:    "shape",
	Domain: Monadic(nil),
	fn: func(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
		if _, ok := R.(apl.Array); ok == false {
			return apl.EmptyArray{}, nil
		}
		ar := R.(apl.Array)
		shape := ar.Shape()
		ret := apl.GeneralArray{
			Values: make([]apl.Value, len(shape)),
			Dims:   []int{len(shape)},
		}
		for i, n := range shape {
			ret.Values[i] = apl.Int(n)
		}
		return ret, nil
	},
}

// rho is dyadic reshape, L is empty or int array, R is array.
func rho(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	// L is empty, returns empty.
	if apl.ArraySize(L.(apl.Array)) == 0 {
		return apl.EmptyArray{}, nil
	}

	l := L.(apl.IntArray)
	shape := make([]int, len(l.Ints))
	copy(shape, l.Ints)

	if rs, ok := R.(apl.Reshaper); ok {
		return rs.Reshape(shape), nil
	}
	return nil, fmt.Errorf("cannot reshape %T", R)
}
