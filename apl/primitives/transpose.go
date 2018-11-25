package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⍉",
		doc:    "cant, transpose, reverse axes",
		Domain: Monadic(IsArray(nil)),
		fn:     transpose,
	})
	register(primitive{
		symbol: "⍉",
		doc:    "cant, transpose, general transpose",
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

	return nil, nil, fmt.Errorf("TODO: duplicate axis")

	// Case 2: There are duplicate axes.
	// The rank of the result is the rank of R minus the duplicated axes.
	rank := max - a.Origin + 1
	shape := make([]int, rank)

	// Build the shape of the result.
	dup := make([][]int, rank) // dup holds the axis of R with origin 0
	for i := range shape {
		// 2 2 1 => [[1 2], [3]]
		// 2 1 2 => [[1 3], [2]]
		// 1 1 2 => [[1 2], [3]]
		// TODO: this is wrong...
		dup[i] = append(dup[i], i)
		for k := i + 1; k < len(al.Ints); k++ { // range al.Ints {
			//if n-a.Origin == i {
			if al.Ints[i] == al.Ints[k] {
				dup[i] = append(dup[i], k)
			}
		}
		if len(dup[i]) == 0 {
			return nil, nil, fmt.Errorf("transpose: len(dup_i)=0: this should not happen")
		} else if len(dup[i]) == 1 {
			shape[i] = rs[dup[i][0]]
		} else {
			min := rs[al.Ints[dup[i][0]]]
			for _, d := range dup[i] {
				if m := rs[al.Ints[d]-a.Origin]; m < min {
					min = m
				}
			}
			shape[i] = min
		}
	}
	fmt.Println("dup", dup)

	// Iterate over the result indexes and map backwards.
	flat := make([]int, apl.ArraySize(apl.GeneralArray{Dims: shape}))
	fidx := make([]int, len(shape))
	ic, ridx := apl.NewIdxConverter(rs)
	for i := range flat {
		for k, f := range fidx {
			for _, g := range dup[k] {
				ridx[g] = f
			}
		}
		flat[i] = ic.Index(ridx)
		apl.IncArrayIndex(fidx, shape)
	}
	fmt.Println("flat", flat)
	fmt.Println("shape", shape)
	return flat, shape, nil

	// TODO: the result needs to be transposed
	// by: ((L⍳L)=⍳⍴L)/L

}
