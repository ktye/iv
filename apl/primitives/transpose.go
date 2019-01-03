package primitives

import (
	"fmt"
	"reflect"

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
		doc:    "table from object",
		Domain: Monadic(IsObject(nil)),
		fn:     transposeObject,
	})
	register(primitive{
		symbol: "⍉",
		doc:    "dict from table",
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
	res := apl.MixedArray{
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
	flat := make([]int, apl.ArraySize(apl.MixedArray{Dims: shape}))
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
	uniform := func(v apl.Value) (bool, int) {
		ar, ok := v.(apl.Array)
		if ok == false {
			return true, 1
		}
		shape := ar.Shape()
		if len(shape) != 1 {
			return false, 0
		}
		if _, ok := v.(apl.Uniform); ok {
			return true, shape[0]
		}
		if shape[0] < 2 {
			return true, shape[0]
		}
		v, _ = ar.At(0)
		t := reflect.TypeOf(v)
		for i := 1; i < shape[0]; i++ {
			if v, err := ar.At(i); err != nil {
				return false, 0
			} else if reflect.TypeOf(v) != t {
				return false, 0
			}
		}
		return true, shape[0]
	}
	tab := apl.Table{Dict: &apl.Dict{}}
	keys := o.Keys()
	n := 0
	for i, k := range keys {
		col := o.At(a, k)
		if col == nil {
			return nil, fmt.Errorf("table: column %s does not exist", k.String(a))
		}
		ok, size := uniform(col)
		if ok == false {
			return nil, fmt.Errorf("table: column %s has mixed types", k.String(a))
		}
		if size == -1 {
			size = 1
			col = apl.List{col}
		}
		if i == 0 {
			n = size
		} else if size != n {
			return nil, fmt.Errorf("table: columns have different sizes")
		}
		tab.K = append(tab.K, k) // TODO: copy k
		if tab.Dict.M == nil {
			tab.Dict.M = make(map[apl.Value]apl.Value)
		}
		tab.M[k] = col // TODO: copy
	}
	tab.Rows = n
	return tab, nil
}

// transposeTable returns a Dict by transposing a table.
func transposeTable(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return R.(apl.Table).Dict, nil // TODO: copy
}
