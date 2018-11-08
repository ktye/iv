package funcs

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

func init() {
	register("⍳", both(interval, indexof))
	addDoc("⍳", `⍳ primitive function: iota, interval, progression, index of
Z←⍳R: R nonnegative integer
	Z: vector of integer sequence 1..R
Z←⍳R: R = 0
	Z: empty array
Z←L⍳R: L: vector
	Z: first occurencein L of items in R
	Z ←→ IO++/∧\∼R∘.≡L
	If an item is not found, the item in Z is IO+⍴L
	If an item occurs several times, the index of the first
	occurance is used.
`)
}

// Interval generates a sequence of numbers up to v returnd in the array.
func interval(a *apl.Apl, ignored, v apl.Value) (bool, apl.Value, error) {
	n := 0
	switch v := v.(type) {
	case apl.Bool:
		if v {
			n = 1
		}
	case apl.Int:
		n = int(v)
	default:
		return true, nil, fmt.Errorf("left value of iota must be an integer")
	}
	if n < 0 {
		return true, nil, fmt.Errorf("left value of iota is negative")
	}
	ar := apl.GeneralArray{
		Values: make([]apl.Value, n),
		Dims:   []int{n},
	}
	for i := 0; i < n; i++ {
		ar.Values[i] = apl.Int(a.Origin + i)
	}
	return true, ar, nil
}

// Indexof returns the first occurance of l in the items of r.
func indexof(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	la, ok := l.(apl.Array)
	leftShape := la.Shape()
	if ok == false || len(leftShape) != 1 {
		return false, nil, fmt.Errorf("left argument to iota must be a vector")
	}

	// Convert l to []Int.
	left := make([]apl.Int, leftShape[0])
	for i := range left {
		if v, err := la.At(i); err != nil {
			return true, nil, err
		} else if n, ok := apl.ToInt(v); ok == false {
			return true, nil, fmt.Errorf("left value of iota must contain only integers")
		} else {
			left[i] = apl.Int(n)
		}
	}

	// TODO: Index does a direct comparison, no type conversion.
	// Is that ok? Probably not.
	index := func(x apl.Value) apl.Int {
		for i, lv := range left {
			if eq, _, _ := apl.CompareScalars(x, lv); eq {
				return apl.Int(i + a.Origin)
			}
		}
		return apl.Int(len(left) + a.Origin)
	}

	if _, ok := r.(apl.Array); ok == false {
		r = apl.GeneralArray{
			Values: []apl.Value{r},
			Dims:   []int{1},
		}
	}
	ra := r.(apl.Array)

	rv := apl.GeneralArray{
		Values: make([]apl.Value, apl.ArraySize(ra)),
		Dims:   make([]int, len(ra.Shape())),
	}
	copy(rv.Dims, ra.Shape())

	for i := range rv.Values {
		v, err := ra.At(i)
		if err != nil {
			return true, nil, err
		}
		rv.Values[i] = index(v)
	}
	return true, rv, nil
}
