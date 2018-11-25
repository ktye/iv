package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⍉",
		doc:    "transpose, reverse axes",
		Domain: Monadic(IsArray(nil)),
		fn:     transpose,
	})
	register(primitive{
		symbol: "⍉",
		doc:    "transpose, general transpose",
		Domain: Dyadic(Split(ToIndexArray(nil), IsArray(nil))),
		fn:     transpose,
	})
}

func transpose(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	idx, shape, err := transposeIndexes(a, L, R)
	if err != nil {
		return nil, err
	}
	res := apl.GeneralArray{
		Values: make([]apl.Value, len(idx)),
		Dims:   shape,
	}
	ar := R.(apl.Array)
	for i, k := range idx {
		v, err := ar.At(k)
		if err != nil {
			return nil, err
		}
		res.Values[i] = v
	}
	return res, nil
}

func transposeIndexes(a *apl.Apl, L, R apl.Value) ([]int, []int, error) {
	ar := R.(apl.Array)
	rs := ar.Shape()

	// Monadic transpose: reverse axis.
	if L == nil {
		l := apl.IndexArray{
			Dims: []int{len(rs)},
			Ints: make([]int, len(rs)),
		}
		n := len(l.Ints)
		for i := range l.Ints {
			l.Ints[i] = n - i - 1 + a.Origin
		}
		L = l
	}
	al := L.(apl.IndexArray)
	ls := al.Shape()

	if len(ls) != 1 {
		return nil, nil, fmt.Errorf("transpose: L must be a vector or a scalar")
	}
	if ls[0] != len(rs) {
		return nil, nil, fmt.Errorf("transpose: length of L must be the rank of R")
	}

	// All values of ⍳⌈/L must be included in L
	max := -1
	m := make(map[int]bool)
	for _, v := range al.Ints {
		if v < a.Origin {
			return nil, nil, fmt.Errorf("transpose: value in L out of range: %d", v)
		}
		if v > max {
			max = v
		}
		m[v] = true
	}
	for i := a.Origin; i <= max; i++ {
		if m[i] == false {
			return nil, nil, fmt.Errorf("transpose: all of ⍳⌈/L must be included in L: %d is missing", i)
		}
	}

	// Case 1: All axes are included.
	if max-a.Origin == len(rs)-1 {
		flat := make([]int, apl.ArraySize(ar))
		shape := make([]int, len(rs))
		for i := range shape {
			shape[i] = rs[al.Ints[i]-a.Origin]
		}
		ic, idx := apl.NewIdxConverter(shape)
		ridx := make([]int, len(rs))
		for i := range flat {
			for k := range idx {
				idx[k] = ridx[al.Ints[k]-a.Origin]
			}
			flat[ic.Index(idx)] = i
			apl.IncArrayIndex(ridx, rs)
		}
		return flat, shape, nil
	}

	// Case 2: There are duplicate axes.

	return nil, nil, fmt.Errorf("transpose: TODO: repeated axes")
}
