package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: ",",
		doc:    "ravel, ravel with axis",
		Domain: Monadic(nil),
		fn:     ravel,
		sel:    ravelSelection,
	})
	register(primitive{
		symbol: "∊",
		doc:    "enlist",
		Domain: Monadic(nil),
		fn:     enlist,
	})
	register(primitive{
		symbol: ",",
		doc:    "catenate, join along last axis",
		Domain: Dyadic(nil),
		fn:     catenate,
	})
	register(primitive{
		symbol: "⍪",
		doc:    "catenate first",
		Domain: Dyadic(nil),
		fn:     catenateFirst,
	})
	register(primitive{
		symbol: "⍪",
		doc:    "table",
		Domain: Monadic(nil),
		fn:     table,
	})
}

// ravel returns a vector from all elements of R.
// R is already converted to an array.
func ravel(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.Axis); ok {
		return ravelWithAxis(a, R)
	}

	to := ToArray(nil)
	r, ok := to.To(a, R)
	if ok == false {
		return nil, fmt.Errorf("ravel: cannot convert to array: %T", R)
	}

	ar, _ := r.(apl.Array)
	n := apl.ArraySize(ar)
	res := apl.MixedArray{
		Dims:   []int{n},
		Values: make([]apl.Value, n),
	}
	for i := range res.Values {
		res.Values[i] = ar.At(i)
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

func ravelWithAxis(a *apl.Apl, R apl.Value) (apl.Value, error) {
	var x []int
	if r, vec, err := splitAxis(a, R); err == nil {
		R = r
		x = vec
	} else if r, n, frac, err := splitCatAxis(a, apl.Index(0), R); err != nil {
		return nil, fmt.Errorf("ravel with axis: %s", err)
	} else {
		// The result has rank ⍴⍴R+1 with the same shape as R,
		// but a new axis 1 at position x.
		if frac == false {
			return nil, fmt.Errorf("ravel with axis: expected fractional axis")
		}
		x := n + 1
		R = r
		ar, ok := R.(apl.Array)
		if ok == false {
			return apl.MixedArray{Dims: []int{1}, Values: []apl.Value{R}}, nil
		}

		rs := ar.Shape()
		shape := make([]int, len(rs)+1)
		off := 0
		for i := range shape {
			if i == x {
				shape[i] = 1
				off = -1
			} else {
				shape[i] = rs[i+off]
			}
		}
		if rs, ok := ar.(apl.Reshaper); ok {
			return rs.Reshape(shape), nil
		} else {
			return nil, fmt.Errorf("cannot reshape %T", R)
		}
	}

	// The axis is an integer vector.
	// It must be continuous and in ascending order.
	for i := range x {
		if i > 0 && x[i-1] != x[i]-1 {
			return nil, fmt.Errorf("ravel with axis: axis must be ascending and continuous")
		}
	}
	if len(x) == 1 {
		return R, nil
	}

	ar, Rarray := R.(apl.Array)

	// APL2: if the axis is empty, a new last axis of length 1 is appended.
	if len(x) == 0 {
		if Rarray == false {
			return apl.MixedArray{Dims: []int{1}, Values: []apl.Value{R}}, nil
		}
		shape := apl.CopyShape(ar)
		shape = append(shape, 1)
		if rs, ok := ar.(apl.Reshaper); ok {
			return rs.Reshape(shape), nil
		} else {
			return nil, fmt.Errorf("cannot reshape %T", R)
		}
	}

	if Rarray == false {
		return nil, fmt.Errorf("ravel with axis: R must be an array: %T", R)
	}
	rs := ar.Shape()

	// The axis in x are combined.
	prod := 1
	for _, n := range x {
		prod *= rs[n]
	}
	shape := make([]int, 1+len(rs)-len(x))
	off := 0
	for i := range shape {
		if i == x[0] {
			shape[i] = prod
			off = len(x) - 1
		} else {
			shape[i] = rs[i+off]
		}
	}
	if rs, ok := ar.(apl.Reshaper); ok {
		return rs.Reshape(shape), nil
	} else {
		return nil, fmt.Errorf("cannot reshape %T", R)
	}
}

// L and R are conformable if
//	they have the same rank, or
//	at least one argument is scalar
//	they differ in rank by 1
// For arrays the length of all axis but the last must be the same.
func catenate(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if l, r, first, ok := isTableCat(a, L, R); ok {
		return catenateTables(a, l, r, first)
	}

	_, lst := L.(apl.List)
	_, rst := R.(apl.List)
	if lst || rst {
		return catenateLists(a, L, R)
	}

	var err error
	var x int
	var frac bool
	R, x, frac, err = splitCatAxis(a, L, R)
	if err != nil {
		return nil, err
	}
	if frac {
		return laminate(a, L, R, x+1)
	}

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
		v, _ := a.Unify(apl.MixedArray{
			Values: []apl.Value{L, R}, // TODO: copy?
			Dims:   []int{2},
		}, false)
		return v, nil
	}

	reshapeScalar := func(scalar apl.Value, othershape []int) apl.Array {
		othershape[x] = 1
		ary := apl.MixedArray{
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
	res := apl.MixedArray{
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
			v = ar.At(rc.Index(src))
		} else {
			v = al.At(lc.Index(src))
		}
		res.Values[i] = v // TODO: copy?
		apl.IncArrayIndex(dst, newshape)
	}
	return res, nil
}

func catenateLists(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	l, lok := L.(apl.List)
	r, rok := R.(apl.List)
	if lok == false {
		return append(apl.List{L}, r...), nil // TODO: copy
	} else if rok == false {
		return append(l, R), nil // TODO: copy
	}
	res := make(apl.List, len(l)+len(r))
	copy(res, l)          // TODO: copy
	copy(res[len(l):], r) // TODO: copy
	return res, nil
}

func laminate(a *apl.Apl, L, R apl.Value, x int) (apl.Value, error) {
	al, lok := L.(apl.Array)
	ar, rok := R.(apl.Array)
	if lok == false && rok == false {
		if x != 0 {
			return nil, fmt.Errorf("cannot laminate two scalar for given axis")
		}
		return apl.MixedArray{Dims: []int{2}, Values: []apl.Value{L, R}}, nil
	}

	reshape := func(scalar apl.Value, shape []int) apl.Array {
		ary := apl.MixedArray{
			Dims: shape,
		}
		ary.Values = make([]apl.Value, apl.ArraySize(ary))
		for i := range ary.Values {
			ary.Values[i] = scalar // TODO: copy?
		}
		return ary
	}

	if lok == false {
		al = reshape(L, apl.CopyShape(ar))
	} else if rok == false {
		ar = reshape(R, apl.CopyShape(al))
	}
	ls := al.Shape()
	rs := ar.Shape()

	if len(ls) != len(rs) {
		return nil, fmt.Errorf("laminate: arguments must have the same rank")
	}
	for i := range ls {
		if ls[i] != rs[i] {
			return nil, fmt.Errorf("laminate: arguments must have the same shape")
		}
	}

	// The new array has one more dimension with length 2 at axis x,
	// otherwise the shape is the same as for L and R.
	shape := make([]int, len(ls)+1)
	off := 0
	for i := range shape {
		if i == x {
			shape[i] = 2
			off = -1
		} else {
			shape[i] = ls[i+off]
		}
	}

	// Iterate over the result and copy values from L or R depending,
	// if the the index at axis x is 0 or 1.
	res := apl.MixedArray{Dims: shape}
	res.Values = make([]apl.Value, apl.ArraySize(res))
	dst := make([]int, len(shape))
	ic, src := apl.NewIdxConverter(ls)
	for i := range res.Values {
		var v apl.Value
		copy(src[:x], dst[:x])
		copy(src[x:], dst[x+1:])
		if dst[x] == 0 {
			v = al.At(ic.Index(src))
		} else {
			v = ar.At(ic.Index(src))
		}
		res.Values[i] = v // TODO: copy?
		apl.IncArrayIndex(dst, shape)
	}
	return res, nil
}

func catenateFirst(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.Axis); ok == true {
		return catenate(a, L, R)
	}
	return catenate(a, L, apl.Axis{A: apl.Index(a.Origin), R: R})
}

func table(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	ar, ok := R.(apl.Array)
	if ok == false {
		return apl.MixedArray{Dims: []int{1, 1}, Values: []apl.Value{R}}, nil
	}
	rs := ar.Shape()

	prod := 1
	for _, n := range rs[1:] {
		prod *= n
	}
	shape := []int{rs[0], prod}

	if rs, ok := ar.(apl.Reshaper); ok == false {
		return nil, fmt.Errorf("cannot reshape %T", R)
	} else {
		return rs.Reshape(shape), nil
	}
}

// enlist creates a flat list from a nested list catenating all elements by depth first.
func enlist(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	r, ok := R.(apl.List)
	if ok == false {
		return apl.List{R}, nil // TODO: copy
	}

	var f func(l apl.List) apl.List
	f = func(l apl.List) apl.List {
		var res apl.List
		for _, e := range l {
			if v, ok := e.(apl.List); ok {
				v = f(v)
				res = append(res, v...) // TODO: copy
			} else {
				res = append(res, e) // TODO: copy
			}
		}
		return res
	}
	return f(r), nil
}
