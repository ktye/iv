package domain

import "github.com/ktye/iv/apl"

// TODO: do we have to test against function or identifiers here?

// ToArray converts scalars to arrays.
func ToArray(child SingleDomain) SingleDomain {
	return array{child, true}
}

// IsArray tests if the value is an array.
func IsArray(child SingleDomain) SingleDomain {
	return array{child, false}
}

type array struct {
	child SingleDomain
	conv  bool
}

func (v array) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	_, ok := V.(apl.Array)
	if v.conv == false && ok == false {
		return V, false
	} else if ok == true {
		return propagate(a, V, v.child)
	}

	// Convert scalars to arrays.
	ga := apl.MixedArray{
		Values: []apl.Value{V},
		Dims:   []int{1},
	}
	u, _ := a.Unify(ga, false)
	return propagate(a, u, v.child)
}
func (v array) String(f apl.Format) string {
	name := "array"
	if v.conv {
		name = "toarray"
	}
	if v.child == nil {
		return name
	}
	return name + " " + v.child.String(f)
}

// ToVector converts scalars and arrays with only one dimension > 0 to a vector.
// It returns an empty array if the size is 0.
// It may convert to a general array, if it has to set the shape.
func ToVector(child SingleDomain) SingleDomain {
	return vector{child, true}
}

// IsVector accepts only arrays with rank 1.
func IsVector(child SingleDomain) SingleDomain {
	return vector{child, false}
}

type vector struct {
	child SingleDomain
	conv  bool
}

func (v vector) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	ar, ok := V.(apl.Array)
	if v.conv == false {
		if ok && len(ar.Shape()) == 1 {
			return propagate(a, V, v.child)
		}
		return V, false
	}

	// Convert scalars.
	if ok == false {
		ga := apl.MixedArray{
			Values: []apl.Value{V},
			Dims:   []int{1},
		}
		return propagate(a, ga, v.child)
	}

	shape := ar.Shape()
	if len(shape) == 1 {
		return propagate(a, V, v.child)
	}
	as := ar.Size()

	// Handle empty case.
	if as == 0 {
		return apl.EmptyArray{}, true
	}

	// Check if maxdim is equal to size.
	maxdim := 0
	for _, n := range ar.Shape() {
		if n > maxdim {
			maxdim = n
		}
	}
	if maxdim != as {
		return V, false
	}

	// Create a new general array.
	ret := apl.NewMixed([]int{as})
	for i := 0; i < as; i++ {
		s := ar.At(i)
		ret.Values[i] = s
	}
	return propagate(a, a.UnifyArray(ret), v.child)
}
func (v vector) String(f apl.Format) string {
	name := "vector"
	if v.conv {
		name = "tovector"
	}
	if v.child == nil {
		return name
	}
	return name + " " + v.child.String(f)
}

func ToBoolArray(child SingleDomain) SingleDomain {
	return boolarray{child, true}
}

func IsBoolArray(child SingleDomain) SingleDomain {
	return indexarray{child, false}
}

type boolarray struct {
	child SingleDomain
	conv  bool
}

func (ba boolarray) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	_, ok := V.(apl.BoolArray)
	if ba.conv == false && ok == false {
		return V, false
	} else if ba.conv == false && ok {
		return propagate(a, V, ba.child)
	} else if ok {
		return propagate(a, V, ba.child)
	}

	iav, ok := indexarray{nil, true}.To(a, V)
	if ok == false {
		return V, false
	}
	ia := iav.(apl.IntArray)
	b := apl.BoolArray{Dims: apl.CopyShape(ia), Bools: make([]bool, len(ia.Ints))}
	for i, n := range ia.Ints {
		if n < 0 || n > 1 {
			return V, false
		} else if n == 1 {
			b.Bools[i] = true
		}
	}
	return propagate(a, b, ba.child)
}
func (ba boolarray) String(f apl.Format) string {
	name := "boolarray"
	if ba.conv == true {
		name = "toboolarray"
	}
	if ba.child == nil {
		return name
	}
	return name + " " + ba.child.String(f)
}

// ToIndexArray accepts arrays that contain only numbers that are convertibel to ints.
// It accepts also scalars.
// Size-0 arrays are returns as empty.
func ToIndexArray(child SingleDomain) SingleDomain {
	return indexarray{child, true}
}

// IsIndexArray accepts only an IntArray.
func IsIndexArray(child SingleDomain) SingleDomain {
	return indexarray{child, false}
}

type indexarray struct {
	child SingleDomain
	conv  bool
}

func (ia indexarray) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	_, ok := V.(apl.IntArray)
	if ia.conv == false && ok == false {
		return V, false
	} else if ia.conv == false && ok {
		return propagate(a, V, ia.child)
	} else if ok {
		return propagate(a, V, ia.child)
	}

	// Try to convert.
	ar, ok := V.(apl.Array)

	// Scalar number
	if ok == false {
		if n, ok := V.(apl.Number); ok {
			if i, ok := n.ToIndex(); ok {
				return propagate(a, apl.IntArray{
					Ints: []int{i},
					Dims: []int{1},
				}, ia.child)
			} else {
				return V, false
			}
		} else if b, ok := V.(apl.Bool); ok {
			n := 0
			if b {
				n = 1
			}
			return propagate(a, apl.IntArray{
				Ints: []int{n},
				Dims: []int{1},
			}, ia.child)
		} else {
			return V, false
		}
	}

	// Empty array.
	if ar.Size() == 0 {
		return propagate(a, apl.EmptyArray{}, ia.child)
	}

	// Make a new array and try to convert all values.
	res := apl.IntArray{
		Ints: make([]int, ar.Size()),
		Dims: apl.CopyShape(ar),
	}

	for i := range res.Ints {
		s := ar.At(i)
		if n, ok := s.(apl.Number); ok {
			if d, ok := n.ToIndex(); ok {
				res.Ints[i] = d
			} else {
				return V, false
			}
		} else if b, ok := s.(apl.Bool); ok {
			res.Ints[i] = 0
			if b {
				res.Ints[i] = 1
			}
		} else {
			return V, false
		}
	}
	return propagate(a, res, ia.child)
}
func (ia indexarray) String(f apl.Format) string {
	name := "indexarray"
	if ia.conv == true {
		name = "toindexarray"
	}
	if ia.child == nil {
		return name
	}
	return name + " " + ia.child.String(f)
}
