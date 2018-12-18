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
	var err error
	var x int
	var frac bool
	R, x, frac, err = splitCatAxis(a, L, R)
	if err != nil {
		return nil, err
	}
	if frac {
		return laminate(a, L, R, x)
	}
	_ = x

	al, isLarray := L.(apl.Array)
	ar, isRarray := R.(apl.Array)

	// Left or right is an empty array
	if isLarray && apl.ArraySize(al) == 0 {
		return R, nil
	} else if isRarray && apl.ArraySize(ar) == 0 {
		return L, nil
	}

	// Catenate two scalars.
	if isLarray == false && isRarray == false {
		return apl.GeneralArray{
			Values: []apl.Value{L, R}, // TODO: copy?
			Dims:   []int{2},
		}, nil
	}

	reshapeScalar := func(scalar apl.Value, othershape []int) apl.Array {
		othershape[x] = 1
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
	// the x axis length to 1.
	if isLarray == false {
		al = reshapeScalar(L, apl.CopyShape(ar))
	} else if isRarray == false {
		ar = reshapeScalar(R, apl.CopyShape(al))
	}
	sl := al.Shape()
	sr := ar.Shape()

	reshape := func(ary apl.Array, shape []int) (apl.Array, error) {
		if rs, ok := ary.(apl.Reshaper); ok {
			return rs.Reshape(shape).(apl.Array), nil
		} else {
			return nil, fmt.Errorf("cannot reshape %T", ary)
		}
	}

	// If ranks differ by 1: insert 1 at the axis.
	insert1 := func(s []int, i int) []int {
		s = append(s, 0)
		copy(s[i+1:], s[i:])
		s[x] = 1
		return s
	}
	if d := len(sl) - len(sr); d != 0 {
		if d < -1 || d > 2 {
			return nil, fmt.Errorf("catenate: ranks differ more that 1")
		}
		if d == -1 {
			sl = insert1(apl.CopyShape(al), x)
			al, err = reshape(al, sl)
		} else if d == 1 {
			sr = insert1(apl.CopyShape(ar), x)
			ar, err = reshape(ar, sr)
		}
	}
	if err != nil {
		return nil, err
	}

	// All axis lengths except for x must match.
	newshape := make([]int, len(sl))
	for i := range sl {
		newshape[i] = sl[i]
		if i == x { // i == len(sl)-1 {
			newshape[i] = sl[i] + sr[i]
		} else if sl[i] != sr[i] {
			return nil, fmt.Errorf("catenate: all axis lengths except for the catenation axis must match")
		}
	}
	res := apl.GeneralArray{
		Dims: newshape,
	}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	// Iterate over combined elements, taking from L or R.
	split := sl[x]
	dst := make([]int, len(newshape))
	lc, src := apl.NewIdxConverter(sl)
	rc, _ := apl.NewIdxConverter(sr)
	for i := range res.Values {
		var v apl.Value
		copy(src, dst)
		if n := src[x]; n >= split {
			src[x] -= split
			v, err = ar.At(rc.Index(src))
		} else {
			v, err = al.At(lc.Index(src))
		}
		res.Values[i] = v // TODO: copy?
		apl.IncArrayIndex(dst, newshape)
	}
	return res, nil

	/*
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
	*/
}

func laminate(a *apl.Apl, L, R apl.Value, x int) (apl.Value, error) {
	return nil, fmt.Errorf("TODO: laminate")
}
