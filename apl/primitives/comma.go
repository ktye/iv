package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: ",",
		doc:    "ravel, create row vector",
		Domain: Monadic(ToArray(nil)),
		fn:     ravel,
		sel:    ravelSelection,
	})
	register(primitive{
		symbol: "∊",
		doc:    "enlist, create simple vector",
		Domain: Monadic(ToArray(nil)),
		fn:     ravel, // for simple arrays ravel and enlist is the same
	})
	register(primitive{
		symbol: ",",
		doc:    "catenate, join along last axis",
		Domain: Dyadic(nil),
		fn:     catenate,
	})
	// TODO ravel with axis
	// TODO catenate with axis
	// TODO laminate
}

// ravel returns a vector from all elements of R.
// R is already converted to an array.
func ravel(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	ar, _ := R.(apl.Array)
	n := apl.ArraySize(ar)
	res := apl.GeneralArray{
		Dims:   []int{n},
		Values: make([]apl.Value, n),
	}
	for i := range res.Values {
		v, err := ar.At(i)
		if err != nil {
			return nil, err
		}
		res.Values[i] = v
	}
	return res, nil
}

func ravelSelection(a *apl.Apl, L, R apl.Value) (apl.IndexArray, error) {
	ar, ok := R.(apl.Array)
	if ok == false {
		return apl.IndexArray{}, fmt.Errorf("ravel: cannot select from non-array: %T", R)
	}
	ai := apl.IndexArray{Dims: []int{apl.ArraySize(ar)}}
	ai.Ints = make([]int, ai.Dims[0])
	for i := range ai.Ints {
		ai.Ints[i] = i
	}
	return ai, nil
}

// L and R are conformable if
//	they have the same rank, or
//	at least one argument is scalar
//	they differ in rank by 1
// For arrays the length of all axis but the last must be the same.
func catenate(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	al, isLarray := L.(apl.Array)
	ar, isRarray := R.(apl.Array)

	// Left or right is an empty array
	if isLarray && apl.ArraySize(al) == 0 {
		return R, nil
	} else if isRarray && apl.ArraySize(ar) == 0 {
		return L, nil
	}

	if isLarray == false && isRarray == false {
		return apl.GeneralArray{
			Values: []apl.Value{L, R}, // TODO: do we need a copy?
			Dims:   []int{2},
		}, nil
	}

	reshapeScalar := func(scalar apl.Value, othershape []int) apl.Array {
		othershape[len(othershape)-1] = 1
		ary := apl.GeneralArray{
			Dims: othershape,
		}
		ary.Values = make([]apl.Value, apl.ArraySize(ary))
		for i := range ary.Values {
			ary.Values[i] = scalar // TODO: copy?
		}
		return ary
	}

	// If one is scalar, reshape to match the other's shape, with
	// the last axis length to 1.
	if isLarray == false {
		al = reshapeScalar(L, apl.CopyShape(ar))
	} else if isRarray == false {
		ar = reshapeScalar(R, apl.CopyShape(al))
	}

	reshape := func(ary apl.Array, shape []int) (apl.Array, error) {
		if rs, ok := ary.(apl.Reshaper); ok {
			return rs.Reshape(shape).(apl.Array), nil
		} else {
			return nil, fmt.Errorf("cannot reshape %T", ary)
		}
	}

	// Catenate arrays.
	sl := al.Shape()
	sr := ar.Shape()
	var err error
	// If ranks differ by 1: reshape.
	if d := len(sl) - len(sr); d != 0 {
		if d < -1 || d > 2 {
			return nil, fmt.Errorf("catenate: ranks differ more that 1")
		}
		if d == -1 {
			sl = append(apl.CopyShape(al), 1)
			al, err = reshape(al, sl)
		} else if d == 1 {
			sr = append(apl.CopyShape(ar), 1)
			ar, err = reshape(ar, sr)
		}
	}
	if err != nil {
		return nil, err
	}

	// All axis lengths but the last must match.
	newshape := make([]int, len(sl))
	for i := range sl {
		newshape[i] = sl[i]
		if i == len(sl)-1 {
			newshape[i] = sl[i] + sr[i]
		} else if sl[i] != sr[i] {
			return nil, fmt.Errorf("catenate: all axis lengths but the last must match")
		}
	}
	ret := apl.GeneralArray{
		Dims: newshape,
	}
	ret.Values = make([]apl.Value, apl.ArraySize(ret))

	// Iterate over combined elements, taking from L or R.
	lidx, ridx := 0, 0
	nl, nr := sl[len(sl)-1], sr[len(sr)-1] // inner length
	kl, kr := 0, 0
	var v apl.Value
	for i := range ret.Values {
		if kl < nl {
			v, err = al.At(lidx)
			kl++
			lidx++
		} else if kr < nr {
			v, err = ar.At(ridx)
			kr++
			ridx++
			if kr == nr {
				kl = 0
				kr = 0
			}
		} else {
			panic("catenate: illegal state: this should not happen")
		}
		if err != nil {
			return nil, err
		}
		ret.Values[i] = v
	}
	return ret, nil
}
