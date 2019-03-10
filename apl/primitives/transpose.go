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
		sel:    selection(transpose),
	})
	register(primitive{
		symbol: "⍉",
		doc:    "table from object, transpose, flip",
		Domain: Monadic(IsObject(nil)),
		fn:     transposeObject,
	})
	register(primitive{
		symbol: "⍉",
		doc:    "dict from table, transpose, flip",
		Domain: Monadic(IsTable(nil)),
		fn:     transposeTable,
	})
	register(primitive{
		symbol: "⍉",
		doc:    "cant, transpose, general transpose",
		Domain: Dyadic(Split(IsArray(nil), IsNumber(nil))), // This matches (⍳0)⍉5
		fn:     transpose,
		sel:    selection(transpose),
	})
	register(primitive{
		symbol: "⍉",
		doc:    "cant, transpose, general transpose",
		Domain: Dyadic(Split(ToIndexArray(nil), IsArray(nil))),
		fn:     transpose,
		sel:    selection(transpose),
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

	ar := R.(apl.Array)
	res := apl.MakeArray(ar, shape)
	for i, k := range idx {
		res.Set(i, ar.At(k).Copy())
	}
	return res, nil
}

func transposeIndexes(a *apl.Apl, L, R apl.Value) ([]int, []int, error) {
	ar := R.(apl.Array)
	rs := ar.Shape()

	// Monadic transpose: reverse axis.
	if L == nil {
		l := apl.IntArray{
			Dims: []int{len(rs)},
			Ints: make([]int, len(rs)),
		}
		n := len(l.Ints)
		for i := range l.Ints {
			l.Ints[i] = n - i - 1 + a.Origin
		}
		L = l
	}
	al := L.(apl.IntArray)
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
		if v < 1 {
			return nil, nil, fmt.Errorf("transpose: value in L out of range: %d", v)
		}
		if v > max {
			max = v
		}
		m[v] = true
	}
	for i := 1; i <= max; i++ {
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
	flat := make([]int, apl.Prod(shape))
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

// transposeObject returns a Table by transposing an object.
func transposeObject(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	o := R.(apl.Object)
	tab := apl.Table{Dict: &apl.Dict{}}
	keys := o.Keys()
	n := 0
	for i, k := range keys {
		col := o.At(k).Copy()
		if col == nil {
			return nil, fmt.Errorf("table: column %s does not exist", k.String(apl.Format{}))
		}
		if _, ok := col.(apl.Object); ok {
			return nil, fmt.Errorf("table: contains an object: %s", k.String(apl.Format{}))
		}
		if _, ok := col.(apl.Array); ok == false {
			col = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{col}}
		}

		size := 1
		ar := col.(apl.Array)
		if shape := ar.Shape(); len(shape) != 1 {
			return nil, fmt.Errorf("table: column %s has rank != 1", k.String(apl.Format{}))
		} else {
			size = shape[0]
		}
		u, ok := a.Unify(ar, true)
		if ok == false {
			return nil, fmt.Errorf("table: cannot unify column %s (mixed types)", k.String(apl.Format{}))
		}

		if i == 0 {
			n = size
		} else if size != n {
			return nil, fmt.Errorf("table: columns have different sizes")
		}
		tab.K = append(tab.K, k.Copy())
		if tab.Dict.M == nil {
			tab.Dict.M = make(map[apl.Value]apl.Value)
		}
		tab.M[k] = u // already copied
	}
	tab.Rows = n
	return tab, nil
}

// transposeTable returns a Dict by transposing a table.
func transposeTable(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return R.(apl.Table).Dict.Copy(), nil
}
