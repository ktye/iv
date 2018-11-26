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
		Domain: Dyadic(Split(IsArray(nil), IsNumber(nil))), // This matches (⍳0)⍉5
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
	// Special case: L is the empty array and R is scalar: return R.
	if _, ok := L.(apl.EmptyArray); ok {
		if _, ok := R.(apl.Array); ok == false {
			return R, nil
		} else {
			return nil, fmt.Errorf("transpose: L is empty, R not scalar")
		}
	}

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

	// Add 1 to L, if Origin is 0.
	if a.Origin == 0 {
		for i := range al.Ints {
			al.Ints[i] += 1
		}
	}

	// All values of ⍳⌈/L must be included in L.
	// Iso requires both: ^/L∊⍳⌈/0,L and ^/(⍳⌈/0,L)∊L to evaluate to 1.
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

	maxRS := 0
	for _, i := range rs {
		if i > maxRS {
			maxRS = i
		}
	}

	// Element i of shape is ⌊/(L=i)/⍴R.
	shape := make([]int, max)
	for i := range shape {
		min := maxRS
		for k := range rs {
			if al.Ints[k] == i+1 {
				if rs[k] < min {
					min = rs[k]
				}
			}
		}
		shape[i] = min
	}

	// The index list of the result is for item i is: 1+(⍴R)⊥((shape)⊤i)[L]
	flat := make([]int, apl.ArraySize(apl.GeneralArray{Dims: shape}))
	ics, sidx := apl.NewIdxConverter(shape)
	icr, ridx := apl.NewIdxConverter(rs)
	for i := range flat {
		ics.Indexes(i, sidx) // sidx ← (shape)⊤i
		for k, n := range al.Ints {
			ridx[k] = sidx[n-1] // ridx ← ((shape)⊤i)[L]
		}
		flat[i] = icr.Index(ridx) // 1+(⍴R)⊥((shape)⊤i)[L]
	}
	return flat, shape, nil
}
