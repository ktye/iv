package primitives

import (
	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "∪",
		doc:    "unique",
		Domain: Monadic(ToVector(nil)),
		fn:     unique,
	})
	register(primitive{
		symbol: "∪",
		doc:    "union",
		Domain: Dyadic(Split(ToVector(nil), ToVector(nil))),
		fn:     union,
	})
}

// unique: R is a vector.
// DyRef gives an example of an array: ∪3 4 5⍴⍳20
// which fails in tryapl.
// ISO also allows only vectors.
func unique(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	ar := R.(apl.Array)

	var values []apl.Value
	for i := 0; i < apl.ArraySize(ar); i++ {
		v := ar.At(i)
		u := true
		for k := range values {
			if isEqual(a, v, values[k]) {
				u = false
				break
			}
		}
		if u {
			values = append(values, v) // TODO copy?
		}
	}
	return apl.MixedArray{Values: values, Dims: []int{len(values)}}, nil
}

// union of L and R, both are vectors.
func union(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if _, ok := L.(apl.EmptyArray); ok {
		return R, nil
	}
	if _, ok := R.(apl.EmptyArray); ok {
		return L, nil
	}
	al := L.(apl.Array)
	ar := R.(apl.Array)

	var values []apl.Value
	appendvec := func(vec apl.Array) error {
		for i := 0; i < apl.ArraySize(vec); i++ {
			v := vec.At(i)
			u := true
			for k := range values {
				if isEqual(a, v, values[k]) == true {
					u = false
					break
				}
			}
			if u {
				values = append(values, v) // TODO copy?
			}
		}
		return nil
	}
	if err := appendvec(al); err != nil {
		return nil, err
	}
	if err := appendvec(ar); err != nil {
		return nil, err
	}
	return apl.MixedArray{Dims: []int{len(values)}, Values: values}, nil
}
