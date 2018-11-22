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
		fn:     gradeup,
	})
}

func gradeup(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	ar := R.(apl.Array)
	shape := ar.Shape()
	if len(shape) == 0 {
		return nil, fmt.Errorf("gradeup: not an array") // this should not happen
	}

	// We store a copy of all values in b.
	// The subarrays of the higher axis are flattened to b[i].
	// Is this ok for comparison?
	b := make([][]apl.Value, shape[0])
	if len(shape) == 1 {
		// In the vector case, wrap the elements to a single element slice.
		for i := range b {
			v, err := ar.At(i) // TODO: copy?
			if err != nil {
				return nil, err
			}
			b[i] = []apl.Value{v}
		}
	} else {
		subsize := apl.ArraySize(apl.GeneralArray{Dims: shape[1:]})
		off := 0
		for i := range b {
			b[i] = make([]apl.Value, subsize)
			for k := range b[i] {
				v, err := ar.At(off + k) // TODO: copy?
				if err != nil {
					return nil, err
				}
				b[i][k] = v
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
			return nil, fmt.Errorf("grade up: types are not comparable")
		}
	} else {
		if sameNumberTypes() == false {
			return nil, fmt.Errorf("grade up: cannot convert to numbers")
		}
		if _, ok := b[0][0].(lesser); ok == false {
			return nil, fmt.Errorf("grade up: cannot compare number type %T", b[0][0])
		}
	}

	si := sortIndexes{
		b:   b,
		idx: make([]int, len(b)),
	}
	for i := range si.idx {
		si.idx[i] = i
	}
	sort.Sort(si)
	return apl.IndexArray{
		Ints: si.idx,
		Dims: []int{len(si.idx)},
	}, nil
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
