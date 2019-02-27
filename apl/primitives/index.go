package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	// An expression such as A[1;2;] is translated by the parser to
	//	[1;2;] ⌷ A
	// ⌷ cannot be used directly, as an index specification is converted by the parser.
	register(primitive{
		symbol: "⌷",
		doc:    "index, []",
		Domain: Dyadic(Split(indexSpec{}, ToArray(nil))),
		fn:     index,
		sel:    indexSelection,
	})
	register(primitive{
		symbol: "⌷",
		doc:    "index list, []",
		Domain: Dyadic(Split(indexSpec{}, IsList(nil))),
		fn:     listIndex,
		sel:    listSelection,
	})
	register(primitive{
		symbol: "⌷",
		doc:    "index object, []",
		Domain: Dyadic(Split(indexSpec{}, IsObject(nil))),
		fn:     objIndex,
		sel:    objSelection,
	})
	register(primitive{
		symbol: "⌷",
		doc:    "index table, []",
		Domain: Dyadic(Split(indexSpec{}, IsTable(nil))),
		fn:     tableIndex,
		sel:    tableSelection,
	})
}

// indexSpec is the domain type for an index specification.
type indexSpec struct{}

func (i indexSpec) To(a *apl.Apl, v apl.Value) (apl.Value, bool) {
	if _, ok := v.(apl.IdxSpec); ok {
		return v, true
	}
	return v, false
}
func (i indexSpec) String(a *apl.Apl) string {
	return "[index specification]"
}

func index(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	spec := L.(apl.IdxSpec)
	ar := R.(apl.Array)

	// Special case for empty brackets.
	if len(spec) == 0 {
		return R, nil
	}

	idx, err := indexArray(a, spec, ar.Shape())
	if err != nil {
		return nil, err
	}

	// Special case, if the result is a scalar.
	if len(idx.Ints) == 1 && len(idx.Dims) == 0 {
		if err := apl.ArrayBounds(ar, int(idx.Ints[0])); err != nil {
			return nil, err
		}
		return ar.At(int(idx.Ints[0])), nil
	}

	res := apl.MixedArray{
		Dims:   apl.CopyShape(idx),
		Values: make([]apl.Value, apl.ArraySize(idx)),
	}
	for i, n := range idx.Ints {
		if err := apl.ArrayBounds(ar, n); err != nil {
			return nil, err
		}
		res.Values[i] = ar.At(n) // TODO copy?
	}
	return res, nil
}

func indexSelection(a *apl.Apl, L, R apl.Value) (apl.IntArray, error) {
	spec := L.(apl.IdxSpec)
	ar := R.(apl.Array)

	// Special case for empty brackets.
	if len(spec) == 0 {
		ai := apl.IntArray{Dims: apl.CopyShape(ar), Ints: make([]int, apl.ArraySize(ar))}
		for i := range ai.Ints {
			ai.Ints[i] = i
		}
		return ai, nil
	}

	if len(spec) == 1 {
		if f, ok := spec[0].(apl.Function); ok {
			return funcArrayIndex(a, f, ar)
		}
	}

	return indexArray(a, spec, ar.Shape())
}

func objSelection(a *apl.Apl, L, R apl.Value) (apl.IntArray, error) {
	obj := R.(apl.Object)
	d, isd := R.(*apl.Dict)
	spec := L.(apl.IdxSpec)
	if len(spec) != 1 {
		return objDepthSelection(a, obj, spec, apl.IntArray{})
	}

	keys := make(map[apl.Value]int)
	for i, k := range obj.Keys() {
		keys[k] = i + a.Origin
	}

	as, ok := spec[0].(apl.Array)
	if ok == false {
		if idx, ok := keys[spec[0]]; ok == false {
			if isd {
				// Index-assignment into a non-existing key in a dict, creates a new key.
				if err := d.Set(a, spec[0], apl.EmptyArray{}); err != nil {
					return apl.IntArray{}, err
				} else {
					return apl.IntArray{Dims: []int{1}, Ints: []int{len(keys) + a.Origin}}, nil
				}
			} else {
				return apl.IntArray{}, fmt.Errorf("key does not exist: %s", spec[0].String(a))
			}
		} else {
			return apl.IntArray{Dims: []int{1}, Ints: []int{idx}}, nil
		}
	}

	ai := apl.IntArray{Dims: []int{as.Size()}, Ints: make([]int, as.Size())}
	for i := 0; i < as.Size(); i++ {
		key := as.At(i)
		k, ok := keys[key]
		if ok == false {
			if isd {
				if err := d.Set(a, key, apl.EmptyArray{}); err != nil {
					return apl.IntArray{}, err
				} else {
					k = len(keys) + a.Origin
					keys[key] = k
				}
			} else {
				return apl.IntArray{}, fmt.Errorf("key does not exist: %s", key.String(a))
			}
		}
		ai.Ints[i] = k
	}
	return ai, nil
}

// objDepthSelection returns a depth index into an object tree.
// Depth indexes for objects are returned as negative indexes starting at -1
// to distinguish them from vector indexes (multiple keys at the same level).
func objDepthSelection(a *apl.Apl, o apl.Object, spec apl.IdxSpec, ia apl.IntArray) (apl.IntArray, error) {
	key := spec[0]
	val := o.At(a, key)
	if val == nil {
		return ia, fmt.Errorf("obj depth sel: key does not exist: %v", key)
	}

	idx := -1
	keys := o.Keys()
	for i, k := range keys {
		if k == key {
			idx = i
		}
	}
	if idx < 0 {
		return ia, fmt.Errorf("obj depth sel: cannot find index")
	}
	if len(ia.Dims) == 0 {
		ia.Dims = []int{1}
		ia.Ints = []int{-1 - idx}
	} else {
		ia.Dims[0]++
		ia.Ints = append(ia.Ints, -1-idx)
	}
	if len(spec) < 2 {
		return ia, nil
	}

	add := func(idx apl.IntArray) {
		if n := len(idx.Ints); n > 0 {
			ia.Dims[0] += n
			ia.Ints = append(ia.Ints, idx.Ints...)
		}
	}

	if _, ok := val.(apl.Table); ok {
		return ia, fmt.Errorf("obj depth sel: index into table is not supported")
	} else if o, ok := val.(apl.Object); ok {
		return objDepthSelection(a, o, spec[1:], ia)
	} else if l, ok := val.(apl.List); ok {
		idx, err := listSelection(a, l, spec[1:])
		if err != nil {
			return ia, err
		}
		add(idx)
		return ia, nil
	} else {
		idx, err := indexSelection(a, spec[1:], val)
		if err != nil {
			return ia, err
		}
		add(idx)
		return ia, nil
	}
}

// indexArray returns the indexes within the array A for the given index specification.
// The result may have a smaller size and shape as the input array.
// The indexes in the spec are origin dependend, but in IndexArray are always origin 0.
func indexArray(a *apl.Apl, spec apl.IdxSpec, shape []int) (apl.IntArray, error) {
	intspec, err := spec2ints(a, spec, shape)
	if err != nil {
		return apl.IntArray{}, nil
	}

	// Initially the rank is the same as spec.
	// Single element axis will be reduced later.
	res := apl.IntArray{Dims: make([]int, len(intspec))}
	for i := range intspec {
		res.Dims[i] = len(intspec[i])
	}

	res.Ints = make([]int, apl.ArraySize(res))
	ic, src := apl.NewIdxConverter(shape)
	dst := make([]int, len(res.Dims))
	for i := range res.Ints {
		for k, n := range dst {
			src[k] = intspec[k][n]
		}
		res.Ints[i] = ic.Index(src)
		apl.IncArrayIndex(dst, res.Dims)
	}

	// Reduce rank by collapsing single element axis.
	rs := make([]int, 0, len(res.Dims))
	for _, v := range res.Dims {
		if v != 1 {
			rs = append(rs, v)
		}
	}
	res.Dims = rs
	return res, nil
}

// Spec2ints converts an index specification to [][]int for the given shape.
// spec is origin dependent, the result has always origin 0.
func spec2ints(a *apl.Apl, spec apl.IdxSpec, shape []int) ([][]int, error) {
	if len(spec) != len(shape) {
		return nil, fmt.Errorf("indexing: array and index specification have different rank")
	}

	to := ToIndexArray(nil)
	idx := make([][]int, len(shape))
	for i := range spec {
		v, ok := to.To(a, spec[i])
		if ok == false {
			return nil, fmt.Errorf("index specification for axis %d is illegal: %T", i+1, spec[i])
		}

		// Empty axis are expanded to all elements of the axis.
		if _, ok := v.(apl.EmptyArray); ok {
			idx[i] = make([]int, shape[i])
			for k := range idx[i] {
				idx[i][k] = k
			}
			continue
		}
		ia := v.(apl.IntArray)
		idx[i] = make([]int, len(ia.Ints))
		for k := range ia.Ints {
			if n := ia.Ints[k] - a.Origin; n < 0 || n >= shape[i] {
				return nil, fmt.Errorf("index specification for axis %d is out of range", i+1)
			} else {
				idx[i][k] = n
			}
		}
	}
	return idx, nil
}

// funcArrayIndex uses the dyadic function f to index into array A.
// f is called with the current index vector on the left and the array value on the right
// for each element. Indexes are selected if the function returns 1 otherwise it may return 0 or empty.
// Example: A←2 3⍴6 ⋄ A[{⍺[2]>⍺[1]&&⍵<4:1}].
// It returns an index vector with the selected indexes only (0-based).
func funcArrayIndex(a *apl.Apl, f apl.Function, A apl.Array) (apl.IntArray, error) {
	shape := A.Shape()
	res := apl.IntArray{Dims: []int{0}}
	res.Ints = make([]int, 0, A.Size())
	idx := make([]int, len(shape))
	L := apl.IntArray{Dims: []int{len(idx)}}
	L.Ints = make([]int, len(idx))
	for i := 0; i < A.Size(); i++ {
		for k := range L.Ints {
			L.Ints[k] = a.Origin + idx[k]
		}
		v, err := f.Call(a, L, A.At(i)) // TODO: copy?
		if err != nil {
			return res, fmt.Errorf("func-array-index: %s", err)
		}
		if _, ok := v.(apl.EmptyArray); ok {
			continue
		}
		num, ok := v.(apl.Number)
		if ok == false {
			return res, fmt.Errorf("func-array-index: return value is not numeric: %T", v)
		}
		n, ok := num.ToIndex()
		if ok == false || (i < 0 && i > 1) {
			return res, fmt.Errorf("func-array-index: return value is not boolean: %T %s", v, v.String(a))
		}
		if n == 1 {
			res.Ints = append(res.Ints, i)
			res.Dims[0]++
		}
		apl.IncArrayIndex(idx, shape)
	}
	return res, nil
}

// objIndex returns a dictionary with only the given keys.
// Keys may be indexed by integers, or strings.
func objIndex(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	obj := R.(apl.Object)
	spec := L.(apl.IdxSpec)
	if len(spec) != 1 {
		return objDepthIndex(a, obj, spec)
	}

	// If the spec is a single value, return the value for the key.
	sv, ok := spec[0].(apl.Array)
	if ok == false {
		v := obj.At(a, spec[0])
		if v == nil {
			return nil, fmt.Errorf("key does not exist")
		}
		return v, nil
	}

	// If the spec is a vector, create a dict with these keys.
	ls := sv.Shape()
	if len(ls) != 1 {
		return nil, fmt.Errorf("object index must be a vector")
	}
	k := make([]apl.Value, ls[0])
	m := make(map[apl.Value]apl.Value)
	for i := 0; i < ls[0]; i++ {
		key := sv.At(i)
		v := obj.At(a, key)
		if v == nil {
			return nil, fmt.Errorf("key does not exist: %s", key.String(a))
		}
		k[i] = key // TODO: copy?
		m[key] = v // TODO: copy?
	}
	return &apl.Dict{K: k, M: m}, nil
}

func objDepthIndex(a *apl.Apl, obj apl.Object, spec apl.IdxSpec) (apl.Value, error) {
	key := spec[0]
	v := obj.At(a, key)
	if v == nil {
		return nil, fmt.Errorf("key does not exist: %q", key.String(a))
	}
	if len(spec) == 1 {
		return v, nil
	}
	if _, ok := v.(apl.Table); ok {
		return nil, fmt.Errorf("obj-depth-index: table is not supported")
	} else if o, ok := v.(apl.Object); ok {
		return objIndex(a, spec[1:], o)
	} else if l, ok := v.(apl.List); ok {
		return listIndex(a, spec[1:], l)
	} else if ar, ok := v.(apl.Array); ok {
		return index(a, spec[1:], ar)
	}
	return nil, fmt.Errorf("obj-depth-index: cannot index into %T", v)
}

// listIndexing indexes a list at depth.
// indexes may be negative.
func listIndex(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	lst := R.(apl.List)
	ai, err := listSelection(a, L, R)
	if err != nil {
		return nil, err
	}

	// Index at depth.
	// Indexes may be negative (count from the end).
	idx := ai.Ints
	for i, k := range idx {
		v := lst[k]
		if i == len(idx)-1 {
			return v, nil // TODO: copy?
		} else {
			lst = v.(apl.List)
		}
	}
	return lst, nil // TODO: copy?
}

// listSelection returns the index for selective assignment.
// The returned array is a single depth-index, not multiple indexes.
// This is different from array indexing.
func listSelection(a *apl.Apl, L, R apl.Value) (apl.IntArray, error) {
	lst := R.(apl.List)
	spec := L.(apl.IdxSpec)

	var ai apl.IntArray
	if len(spec) == 1 {
		if _, ok := spec[0].(apl.List); ok {
			return ai, fmt.Errorf("indexing with a list is not supported")
		}
	}

	// Convert spec to ints.
	to := ToIndexArray(nil)
	idx := make([]int, len(spec))
	for i := range spec {
		v, ok := to.To(a, spec[i])
		if ok == false {
			return ai, fmt.Errorf("list index is no integer")
		}
		ai := v.(apl.IntArray)
		if s := ai.Shape(); len(s) != 1 || s[0] != 1 {
			return ai, fmt.Errorf("list index is no integer: %T", v)
		}
		idx[i] = ai.Ints[0] - a.Origin
	}

	// Index at depth.
	// Indexes may be negative (count from the end).
	for i, k := range idx {
		if k < 0 {
			k = len(lst) + k
		}
		if k < 0 || k >= len(lst) {
			return ai, fmt.Errorf("list index out of range")
		}
		idx[i] = k
		v := lst[k]
		if i < len(idx)-1 {
			if l, ok := v.(apl.List); ok {
				lst = l
			} else {
				return ai, fmt.Errorf("list index is too deep")
			}
		}
	}
	return apl.IntArray{Dims: []int{len(idx)}, Ints: idx}, nil
}

// tableSelection returns the indexes for selective assignment on tables.
// T[rowidx], T[rowidx; colkeys], T[rowfunc], T[rowfunc, colkeys], T[colkeys]
// It returns a flat index vector (0-based) with catenated row and col indexes.
// The shape only counts row indexes.
// A rowfunc can be used as a selection function.
// It is called in an environment with column variables predeclared with the names of the
// keys (if they are strings, and valid variable names).
// The left is the row index vector and right argument is the table itself.
// E.g. T[{(⍺>10)^Time>2015.11.22}]
// The function should returns a boolean selection vector of the size of number of columns.
func tableSelection(a *apl.Apl, L, R apl.Value) (apl.IntArray, error) {
	T := R.(apl.Table)
	spec := L.(apl.IdxSpec)
	var idx apl.IntArray
	if len(spec) < 1 || len(spec) > 2 {
		return idx, fmt.Errorf("table-select: index spec len must be 1 or 2: %d", len(spec))
	}

	var colidx []int
	cols := T.Dict.Keys()
	colmap := make(map[apl.Value]int)
	for i := range cols {
		colmap[cols[i]] = i
	}
	if len(spec) < 2 {
		// special case: spec[0] is a string, or string vector.
		to := ToStringArray(nil)
		if sa, ok := to.To(a, spec[0]); ok {
			strings := sa.(apl.StringArray).Strings
			colidx = make([]int, len(strings))
			for i, s := range strings {
				n, ok := colmap[apl.String(s)]
				if ok == false {
					return idx, fmt.Errorf("table-select: columns does not exist: %s", s)
				}
				colidx[i] = n
			}
			spec[0] = apl.EmptyArray{}
		} else {
			colidx = make([]int, len(cols))
			for i := range colidx {
				colidx[i] = i
			}
		}
	} else {
		ar, ok := spec[1].(apl.Array)
		if ok == false {
			ar = apl.MixedArray{
				Dims:   []int{1},
				Values: []apl.Value{spec[1]},
			}
		}
		colidx = make([]int, ar.Size())
		for i := 0; i < ar.Size(); i++ {
			key := ar.At(i)
			n, ok := colmap[key]
			if ok == false {
				return idx, fmt.Errorf("table-select: column does not exist: %s", key.String(a))
			}
			colidx[i] = n
		}
		if _, ok := spec[1].(apl.EmptyArray); ok {
			colidx = make([]int, len(cols))
			for i := range colidx {
				colidx[i] = i
			}
		}
	}

	toIdx := ToIndexArray(nil)
	iav, ok := toIdx.To(a, spec[0])
	if _, ok := spec[0].(apl.EmptyArray); ok {
		ia := apl.IntArray{Dims: []int{T.Rows}, Ints: make([]int, T.Rows)}
		for i := range ia.Ints {
			ia.Ints[i] = i + a.Origin
		}
		iav = ia
	}

	var ia apl.IntArray
	if ok == false {
		f, ok := spec[0].(apl.Function)
		if ok == false {
			return idx, fmt.Errorf("table-select: first spec must be an index vector or a function: %T", spec[0])
		}

		vars := make(map[string]apl.Value)
		for _, key := range cols {
			if s, ok := key.(apl.String); ok {
				vars[string(s)] = T.Dict.At(a, key)
			}
		}
		iv, err := interval(a, nil, apl.Int(T.Rows))
		if err != nil {
			return idx, fmt.Errorf("table-select-func: %s", err)
		}
		v, err := a.EnvCall(f, iv, T, vars)
		if err != nil {
			return idx, fmt.Errorf("table-select-func: %s", err)
		}

		to := ToIndexArray(nil)
		ints, ok := to.To(a, v)
		if ok == false {
			return idx, fmt.Errorf("table-select-func: does not return a boolean vector: %T", v)
		}
		if v, err := where(a, nil, ints); err != nil {
			return idx, fmt.Errorf("table-select-func: %s", err)
		} else {
			ar := v.(apl.IntArray)
			ia.Ints = ar.Ints
			ia.Dims = []int{len(ia.Ints)}
		}
	} else {
		ia = iav.(apl.IntArray)
	}

	idx = apl.IntArray{Dims: []int{len(ia.Ints)}, Ints: make([]int, len(ia.Ints)+len(colidx))} // sic!
	for i, k := range ia.Ints {
		if n := k - a.Origin; n < 0 || n >= T.Rows {
			return idx, fmt.Errorf("table-select: row index out of range: %d", k)
		} else {
			idx.Ints[i] = n
		}
	}
	for i, k := range colidx {
		idx.Ints[len(ia.Ints)+i] = k
	}
	return idx, nil
}

// tableIndex indexes into a table.
//	- empty index
//		T[]      converts to array
//	- single column (only for string keys)
//		T[`a]    returns a vector
//		T[1⍴`a]  returns a single-column table
//	- single row and column
//		T[3;`a]  returns the scalar value
//	- single row
//		T[3]     returns a dict
//		T[1⍴3]   returns a single-row table
//	- multiple rows or columns
//		T[1 3]     return a sub-table with all column
//		T[;`a`b]   return a sub-table with all rows
//		T[⍳3;`a`c] return a sub-table
//	- functional row index
//		T[{⍺>1}]  return a table
func tableIndex(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	t := R.(apl.Table)
	spec := L.(apl.IdxSpec)

	// Empty spec: convert to array.
	if len(spec) == 0 {
		return table2array(a, t) // table.go
	}

	conv := true
	var agg, grp apl.Value
	if len(spec) > 4 {
		return nil, fmt.Errorf("table: index specification is too long: %d", len(spec))
	}
	if len(spec) > 3 {
		grp = spec[3]
	}
	if len(spec) > 2 {
		agg = spec[2]
		spec = spec[:2]
		conv = false
	}
	if conv == true {
		for i := range spec {
			if _, ok := spec[i].(apl.Function); ok {
				conv = false
			}
			if _, ok := spec[i].(apl.EmptyArray); ok {
				continue
			}
			if _, ok := spec[i].(apl.Array); ok {
				conv = false
			}
		}
	}

	idx, err := tableSelection(a, spec, t)
	if err != nil {
		return nil, err
	}

	// Both row and column index are stored on a single axis in an unusual way.
	rows := idx.Ints[:idx.Dims[0]]
	cols := idx.Ints[idx.Dims[0]:]
	keys := t.Keys()

	// Return a single value.
	if conv && len(rows) == 1 && len(cols) == 1 {
		return t.At(a, keys[cols[0]]).(apl.Array).At(rows[0]), nil // TODO copy
	}

	// Return a single column as an array.
	if conv && len(cols) == 1 {
		return t.At(a, keys[cols[0]]), nil // TODO copy
	}

	res := apl.Table{Rows: len(rows)}
	d := apl.Dict{}
	d.K = make([]apl.Value, len(cols))
	d.M = make(map[apl.Value]apl.Value)
	for i, k := range cols {
		key := keys[k] // TODO: copy
		d.K[i] = key
		srccol := t.At(a, key).(apl.Uniform)
		dstcol := srccol.Make([]int{len(rows)})
		col, ok := dstcol.(apl.ArraySetter)
		if ok == false {
			return nil, fmt.Errorf("column is not settable: %T", dstcol)
		}
		for n, m := range rows {
			if err := col.Set(n, srccol.At(m)); err != nil {
				return nil, err
			}
		}
		d.M[key] = col
	}

	// Return a dict for a single row.
	if conv && len(rows) == 1 {
		return &d, nil
	}
	res.Dict = &d

	if agg != nil {
		return tableQuery(a, res, agg, grp)
	}
	return res, nil
}
