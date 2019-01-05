package primitives

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⍋",
		doc:    "grade up, sort index",
		Domain: Monadic(IsArray(nil)),
		fn:     grade(true),
	})
	register(primitive{
		symbol: "⍒",
		doc:    "grade down, reverse sort index",
		Domain: Monadic(IsArray(nil)),
		fn:     grade(false),
	})
	register(primitive{
		symbol: "⍋",
		doc:    "grade up with collating sequence",
		Domain: Dyadic(Split(IsVector(nil), IsArray(nil))),
		fn:     grade2(true),
	})
	register(primitive{
		symbol: "⍒",
		doc:    "grade down with collating sequence",
		Domain: Dyadic(Split(IsVector(nil), IsArray(nil))),
		fn:     grade2(false),
	})
}

func grade(up bool) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	return func(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
		si, err := gradeSetup(a, R)
		if err != nil {
			return nil, err
		}
		if up {
			sort.Sort(si)
		} else {
			sort.Sort(sort.Reverse(si))
		}
		return apl.IndexArray{
			Ints: si.idx,
			Dims: []int{len(si.idx)},
		}, nil
	}
}

// gradeSetup preparse grading.
func gradeSetup(a *apl.Apl, R apl.Value) (sortIndexes, error) {
	ar := R.(apl.Array)
	shape := ar.Shape()
	if len(shape) == 0 {
		return sortIndexes{}, fmt.Errorf("gradeup: not an array") // this should not happen
	}

	// We store a copy of all values in b.
	// The subarrays of the higher axis are flattened to b[i].
	// Is this ok for comparison?
	b := make([][]apl.Value, shape[0])
	if len(shape) == 1 {
		// In the vector case, wrap the elements to a single element slice.
		for i := range b {
			b[i] = []apl.Value{ar.At(i)} // TODO: copy?
		}
	} else {
		subsize := apl.ArraySize(apl.MixedArray{Dims: shape[1:]})
		off := 0
		for i := range b {
			b[i] = make([]apl.Value, subsize)
			for k := range b[i] {
				b[i][k] = ar.At(off + k) // TODO: copy?
			}
			off += subsize
		}
	}

	// All values must be numeric, or of the same type.
	// The type must implement a Less method.
	sametype := func() bool {
		var t reflect.Type
		for i := range b {
			for k := range b[i] {
				v := b[i][k]
				if i == 0 && k == 0 {
					t = reflect.TypeOf(v)
				} else {
					if reflect.TypeOf(v) != t {
						return false
					}
				}
			}
		}
		return true
	}

	// Convert all values to numbers of the highest class.
	sameNumberTypes := func() bool {
		class := -1
		for i := range b {
			for k := range b[i] {
				v := b[i][k]
				var num apl.Number
				if n, ok := v.(apl.Index); ok {
					num = a.Tower.FromIndex(int(n))
				} else if n, ok := v.(apl.Bool); ok {
					num = a.Tower.FromBool(n)
				} else if n, ok := v.(apl.Number); ok {
					num = n
				} else {
					return false
				}
				n, ok := a.Tower.Numbers[reflect.TypeOf(num)]
				if ok == false {
					return false
				}
				b[i][k] = num
				if n.Class > class {
					class = n.Class
				}
			}
		}
		for i := range b {
			for k := range b[i] {
				n := b[i][k].(apl.Number)
				num := a.Tower.Numbers[reflect.TypeOf(n)]
				for c := num.Class; c < class; c++ {
					u, ok := num.Uptype(n)
					if ok == false {
						return false
					}
					b[i][k] = u
				}
			}
		}
		return true
	}

	issame := sametype()
	if issame == true {
		if _, ok := b[0][0].(lesser); ok == false {
			return sortIndexes{}, fmt.Errorf("grade up: types are not comparable")
		}
	} else {
		if sameNumberTypes() == false {
			return sortIndexes{}, fmt.Errorf("grade up: cannot convert to numbers")
		}
		if _, ok := b[0][0].(lesser); ok == false {
			return sortIndexes{}, fmt.Errorf("grade up: cannot compare number type %T", b[0][0])
		}
	}

	si := sortIndexes{
		b:   b,
		idx: make([]int, len(b)),
	}
	for i := range si.idx {
		si.idx[i] = i + a.Origin
	}
	return si, nil
}

// grade2 is the dyadic grade up/down.
// It is only implemented for vector left arguments.
// If L is a vector: L⍋R ←→ ⍋L⍳R
func grade2(up bool) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	return func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		LiotaR, err := indexof(a, L, R)
		if err != nil {
			return nil, err
		}
		g := grade(up)
		return g(a, nil, LiotaR)
	}
}

type sortIndexes struct {
	b   [][]apl.Value
	idx []int
}

func (s sortIndexes) Len() int { return len(s.b) }
func (s sortIndexes) Less(i, j int) bool {
	x := s.b[i]
	y := s.b[j]
	for n := range x {
		xl := x[n].(lesser)
		yl := y[n].(lesser)
		if isless, _ := xl.Less(y[n]); isless {
			return true
		} else if isless, _ := yl.Less(x[n]); isless {
			return false
		}
		// On equality the next element is compared.
	}
	return false
}
func (s sortIndexes) Swap(i, j int) {
	s.b[i], s.b[j] = s.b[j], s.b[i]
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}
