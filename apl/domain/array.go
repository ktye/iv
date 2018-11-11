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
	ga := apl.GeneralArray{
		Values: []apl.Value{V},
		Dims:   []int{1},
	}
	return propagate(a, ga, v.child)
}
func (v array) String(a *apl.Apl) string {
	name := "array"
	if v.conv {
		name = "toarray"
	}
	if v.child == nil {
		return name
	}
	return name + " " + v.child.String(a)
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
		ga := apl.GeneralArray{
			Values: []apl.Value{V},
			Dims:   []int{1},
		}
		return propagate(a, ga, v.child)
	}

	shape := ar.Shape()
	if len(shape) == 1 {
		return propagate(a, V, v.child)
	}
	as := apl.ArraySize(ar)

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
	ret := apl.GeneralArray{
		Values: make([]apl.Value, as),
		Dims:   []int{as},
	}
	for i := 0; i < as; i++ {
		s, _ := ar.At(i)
		ret.Values[i] = s
	}
	return propagate(a, ret, v.child)
}
func (v vector) String(a *apl.Apl) string {
	name := "vector"
	if v.conv {
		name = "tovector"
	}
	if v.child == nil {
		return name
	}
	return name + " " + v.child.String(a)
}

// ToIntArray accepts arrays that contain only numbers that are convertibel to ints.
// It accepts also scalars.
// Size-0 arrays are returns as empty.
func ToIntArray(child SingleDomain) SingleDomain {
	return intarray{child, true}
}

// IsIntArray accepts only an IntArray.
func IsIntArray(child SingleDomain) SingleDomain {
	return intarray{child, false}
}

type intarray struct {
	child SingleDomain
	conv  bool
}

func (ia intarray) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
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

	// Scalar
	if ok == false {
		if n, ok := toInt(V, true); ok {
			return propagate(a, apl.IntArray{
				Ints: []int{int(n.(apl.Int))},
				Dims: []int{1},
			}, ia.child)
		} else {
			return V, false
		}
	}

	// Empty array.
	if apl.ArraySize(ar) == 0 {
		return propagate(a, apl.EmptyArray{}, ia.child)
	}

	// Make a new array and try to convert all values.
	res := apl.IntArray{
		Ints: make([]int, apl.ArraySize(ar)),
		Dims: apl.CopyShape(ar),
	}

	for i := range res.Ints {
		s, _ := ar.At(i)
		if n, ok := toInt(s, true); ok == false {
			return V, false
		} else {
			res.Ints[i] = int(n.(apl.Int))
		}
	}
	return propagate(a, res, ia.child)
}
func (ia intarray) String(a *apl.Apl) string {
	name := "intarray"
	if ia.conv == true {
		name = "tointarray"
	}
	if ia.child == nil {
		return name
	}
	return name + " " + ia.child.String(a)
}
