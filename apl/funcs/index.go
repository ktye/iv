package funcs

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

func init() {
	register("⌷", handle(index))

	addDoc("[", `[] primitive function: indexing, bracket index, ⌷

Z←A[I] A array, I index specification
	index specification is a list of integer arrays separated by semicolon
	Example: A[3], A[1 2 3], A[1;2], A[;2]
	
Z←A[I] A empty, Z empty
`)
}

// the parser changes the expression A[spec] to (spec ⌷ A)
func index(a *apl.Apl, specVal, Aval apl.Value) (bool, apl.Value, error) {

	idxSpec, ok := specVal.(apl.IdxSpec)
	if ok == false {
		return false, nil, nil
	}

	ar, ok := Aval.(apl.Array)
	if ok == false {
		return false, nil, nil
	}

	idx, err := bracketIndex{idxSpec}.Indexes(ar)
	if err != nil {
		return true, nil, err
	} else if len(idx) == 0 {
		return true, apl.EmptyArray{}, nil
	}

	shape, err := idxSpec.Shape(ar.Shape())
	if err != nil {
		return true, nil, err
	}

	rr := apl.MakeArray(ar, shape)
	for i, k := range idx {
		v, err := ar.At(k)
		if err != nil {
			return true, nil, err
		}
		if err := rr.Set(i, v); err != nil {
			return true, nil, err
		}
	}
	return true, rr, nil
}

// bracketIndex implements the ?Indexer? interface for selective assignment.
type bracketIndex struct {
	apl.IdxSpec
}

// Indexes returns an index slice with a length that fits into of IdxSpec.Shape().
// If it cannot, it returns nil and/or an error.
func (spec bracketIndex) Indexes(ar apl.Array) ([]int, error) {
	shape := ar.Shape()

	n := 1
	for _, i := range shape {
		n *= i
	}
	if n == 0 {
		return nil, nil
	}

	return nil, fmt.Errorf("TODO funcs.index.brackedIndex.Indexes")
}

/*


	dims := ar.Shape()
	if len(dims) == 0 {
		return apl.EmptyArray{}
	}

	if len(idxSpec) != len(dims) {
		return true, nil, fmt.Errorf("index specification has %d values but array has %d dimensions", len(idxSpec), len(dims))
	}

	// Special case, array is a vector and idx a scalar.
	if s, ok := apl.ScalarValue(idxSpec[0]); ok && len(dims) == 1 {
		if n, ok := apl.IntValue(s); ok == false {
			return true, nil, fmt.Errorf("index specification is not an integer value: %T", s)
		} else {
			v, err := ar.At(int(n) - a.Origin)
			return true, v, err
		}
	}

	// Special case: array is a vector and idx is an array.
	if idx, ok := idxSpec[0].(Array); ok && len(dims) == 1 {
		idxshape := idx.Shape()
		shape := make([]int, len(idxshape))
		copy(shape, idxshape)
		rr := apl.MakeArray(a, shape)
		for i := range apl.ArraySize(rr) {
			v, err := idx.At(i)
			if err != nil {
				return true, nil, err
			}
			if n, ok := apl.IntValue(v); ok == false {
				return true, nil, fmt.Errorf("index specification is not an integer value: %T", v)
			} else {
				v, err := idx.At(int(n) + a.Origin)
				if err != nil {
					return true, nil, err
				}
				if s, ok := rr.(apl.Setter); ok {
					if err := s.Set(i, v); err != nil {
						return true, nil, err
					}
				} else {
					return true, nil, fmt.Errorf("array is not settable") // Should not happen.
				}
			}
		}
		return true, rr, nil
	}

	// Convert idxSpec to [][]int.
	idx := make([][]int, len(dims))
	for i := range idx {
		is := idxSpec[i]

		// Fill empty indexes.
		if _, ok := s.(apl.EmptyArray); ok {
			idx[i] = make([]int, len(dims[i]))
			for k := range idx[i] {
				idx[i][k] = k
			}
			continue
		}

		if ia, ok := s.(apl.Array); ok == false {
			if n, ok := apl.IntValue(s); ok == false {
				return true, nil, fmt.Errorf("index specification is not an integer value: %T", s)
			} else {
				idx[i] = []int{int(n) - a.Origin}
			}
		} else {
			shape := ia.Shape()
			if len(shape) != 1 {
				return true, nil, fmt.Errorf("index specification must be a vector, but shape is %v", shape)
			}
			idx[i] = make([]int, len(shape[0]))
			for k := range idx[i] {
				s, err := ia.At(k)
				if err != nil {
					return true, nil, err
				}
				if n, ok := apl.IntValue(s); ok == false {
					return true, nil, fmt.Errorf("index specification is not an integer value: %T", s)
				} else {
					idx[i][k] = int(n) - a.Origin
				}
			}
		}
	}

	panic("TODO")

	newshape := make([]int, len(idx))
	var err error
	for i, v := range idx {
		newshape[i] = len(v)
	}

	rr := apl.MakeArray(newshape)
	pos := make([]int, len(idx))
	for i := range apl.ArraySize(rr) {
		apl.ArrayAt(newshape)
	}
*/
