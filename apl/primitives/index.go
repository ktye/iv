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
		doc:    "object index, []",
		Domain: Dyadic(Split(indexSpec{}, IsObject(nil))),
		fn:     objIndex,
		sel:    objSelection,
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
		if v, err := ar.At(idx.Ints[0]); err != nil {
			return nil, err
		} else {
			return v, err
		}
	}

	res := apl.MixedArray{
		Dims:   apl.CopyShape(idx),
		Values: make([]apl.Value, apl.ArraySize(idx)),
	}
	for i, n := range idx.Ints {
		v, err := ar.At(n)
		if err != nil {
			return nil, err
		}
		res.Values[i] = v // TODO copy?
	}
	return res, nil
}

func indexSelection(a *apl.Apl, L, R apl.Value) (apl.IndexArray, error) {
	spec := L.(apl.IdxSpec)
	ar := R.(apl.Array)

	// Special case for empty brackets.
	if len(spec) == 0 {
		ai := apl.IndexArray{Dims: apl.CopyShape(ar), Ints: make([]int, apl.ArraySize(ar))}
		for i := range ai.Ints {
			ai.Ints[i] = i
		}
		return ai, nil
	}

	return indexArray(a, spec, ar.Shape())
}

func objSelection(a *apl.Apl, L, R apl.Value) (apl.IndexArray, error) {
	obj := R.(apl.Object)
	d, isd := R.(*apl.Dict)
	spec := L.(apl.IdxSpec)
	if len(spec) != 1 {
		return apl.IndexArray{}, fmt.Errorf("object index must be a vector")
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
					return apl.IndexArray{}, err
				} else {
					return apl.IndexArray{Dims: []int{1}, Ints: []int{len(keys) + a.Origin}}, nil
				}
			} else {
				return apl.IndexArray{}, fmt.Errorf("key does not exist: %s", spec[0].String(a))
			}
		} else {
			return apl.IndexArray{Dims: []int{1}, Ints: []int{idx}}, nil
		}
	}

	ai := apl.IndexArray{Dims: []int{as.Size()}, Ints: make([]int, as.Size())}
	for i := 0; i < as.Size(); i++ {
		key, err := as.At(i)
		if err != nil {
			return apl.IndexArray{}, err
		}
		k, ok := keys[key]
		if ok == false {
			if isd {
				if err := d.Set(a, key, apl.EmptyArray{}); err != nil {
					return apl.IndexArray{}, err
				} else {
					k = len(keys) + a.Origin
					keys[key] = k
				}
			} else {
				return apl.IndexArray{}, fmt.Errorf("key does not exist: %s", key.String(a))
			}
		}
		ai.Ints[i] = k
	}
	return ai, nil
}

// indexArray returns the indexes within the array A for the given index specification.
// The result may have a smaller size and shape as the input array.
// The indexes in the spec are origin dependend, but in IndexArray are always origin 0.
func indexArray(a *apl.Apl, spec apl.IdxSpec, shape []int) (apl.IndexArray, error) {
	intspec, err := spec2ints(a, spec, shape)
	if err != nil {
		return apl.IndexArray{}, nil
	}

	// Initially the rank is the same as spec.
	// Single element axis will be reduced later.
	res := apl.IndexArray{Dims: make([]int, len(intspec))}
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
		ia := v.(apl.IndexArray)
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

// listIndex returns a dictionary with only the given keys.
// Keys may be indexed by integers, or strings.
func objIndex(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	obj := R.(apl.Object)
	spec := L.(apl.IdxSpec)
	if len(spec) != 1 {
		// TODO: this could be extended to index into an array value.
		return nil, fmt.Errorf("object index: index spec must be a single scalar or vector")
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
		key, err := sv.At(i)
		if err != nil {
			return nil, err
		}
		v := obj.At(a, key)
		if v == nil {
			return nil, fmt.Errorf("key does not exist: %s", key.String(a))
		}
		k[i] = key // TODO: copy?
		m[key] = v // TODO: copy?
	}
	return &apl.Dict{K: k, M: m}, nil
}

// listIndexing indexes a list at depth.
// indexes may be negative.
func listIndex(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	lst := R.(apl.List)

	ai, err := listSelection(a, L, R)
	if err != nil {
		return nil, err
	}
	idx := ai.Ints

	// Index at depth.
	// Indexes may be negative (count from the end).
	for i, k := range idx {
		v := lst[k]
		if i == len(idx)-1 {
			return v, nil // TODO: copy?
		}
		lst = v.(apl.List)
	}
	return lst, nil // TODO: copy?
}

func listSelection(a *apl.Apl, L, R apl.Value) (apl.IndexArray, error) {
	lst := R.(apl.List)
	spec := L.(apl.IdxSpec)

	// Convert spec to ints.
	var ai apl.IndexArray
	to := ToIndexArray(nil)
	idx := make([]int, len(spec))
	for i := range spec {
		v, ok := to.To(a, spec[i])
		if ok == false {
			return ai, fmt.Errorf("list index is no integer")
		}
		ai = v.(apl.IndexArray)
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
			idx[i] = k
		}
		if k < 0 || k >= len(lst) {
			return ai, fmt.Errorf("list index out of range")
		}
		v := lst[k]
		if i < len(idx)-1 {
			if l, ok := v.(apl.List); ok {
				lst = l
			} else {
				return ai, fmt.Errorf("list index is too deep")
			}
		}
	}
	return apl.IndexArray{Dims: []int{len(idx)}, Ints: idx}, nil
}
