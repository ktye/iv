package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "/",
		Domain:  MonadicOp(Function(nil)),
		doc:     "reduce, n-wise reduction",
		derived: reduceLast,
	})
	register(operator{
		symbol:  "⌿",
		Domain:  MonadicOp(Function(nil)),
		doc:     "reduce first, n-wise reduction",
		derived: reduceFirst,
	})
	register(operator{
		symbol:  `\`,
		Domain:  MonadicOp(Function(nil)),
		doc:     "scan",
		derived: scanLast,
	})
	register(operator{
		symbol:  `⍀`,
		Domain:  MonadicOp(Function(nil)),
		doc:     "scan first axis",
		derived: scanFirst,
	})
	register(operator{
		symbol:  "/",
		Domain:  MonadicOp(ToIndexArray(nil)),
		doc:     "replicate, compress",
		derived: replicateLast,
		selection: selectSimple(func(a *apl.Apl, LO, R apl.Value) (apl.Value, error) {
			return Replicate(a, LO, R, -1)
		}),
	})
	register(operator{
		symbol:  "⌿",
		Domain:  MonadicOp(ToIndexArray(nil)),
		doc:     "replicate, compress first axis",
		derived: replicateFirst,
	})
	register(operator{
		symbol:  `\`,
		Domain:  MonadicOp(ToIndexArray(nil)),
		doc:     "expand",
		derived: expandLast,
		selection: selectSimple(func(a *apl.Apl, LO, R apl.Value) (apl.Value, error) {
			return Expand(a, LO, R, -1)
		}),
	})
	register(operator{
		symbol:  `⍀`,
		Domain:  MonadicOp(ToIndexArray(nil)),
		doc:     "expand first axis",
		derived: expandFirst,
	})
}

func reduceLast(a *apl.Apl, f, _ apl.Value) apl.Function {
	return reduction(a, f, -1)
}
func reduceFirst(a *apl.Apl, f, _ apl.Value) apl.Function {
	return reduction(a, f, 0)
}

func scanLast(a *apl.Apl, f, _ apl.Value) apl.Function {
	return scanArray(a, f, -1)
}
func scanFirst(a *apl.Apl, f, _ apl.Value) apl.Function {
	return scanArray(a, f, 0)
}

func replicateLast(a *apl.Apl, f, _ apl.Value) apl.Function {
	return replicate(a, f, -1)
}
func replicateFirst(a *apl.Apl, f, _ apl.Value) apl.Function {
	return replicate(a, f, 0)
}

func expandLast(a *apl.Apl, f, _ apl.Value) apl.Function {
	return expand(a, f, -1)
}
func expandFirst(a *apl.Apl, f, _ apl.Value) apl.Function {
	return expand(a, f, 0)
}

// Reduction returns the derived function f/ .
func reduction(a *apl.Apl, f apl.Value, axis int) apl.Function {

	// Special cases: left, right tack.
	if p, ok := f.(apl.Primitive); ok {
		if p == "⊣" {
			return reduceTack(true)
		} else if p == "⊢" {
			return reduceTack(false)
		}
	}
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		return reduct(a, f.(apl.Function), L, R, axis)
	}
	return function(derived)
}

func reduct(a *apl.Apl, f apl.Function, l, r apl.Value, axis int) (apl.Value, error) {

	if c, ok := r.(apl.Channel); ok {
		if axis == 0 {
			// l f⌿ c applies f and filters empty values.
			return c.Apply(a, f, l, true), nil
		}
		return reduceChannel(a, l, f, c)
	}

	if _, ok := r.(apl.Axis); ok {
		if rr, n, err := splitAxis(a, r); err != nil {
			return nil, err
		} else {
			r = rr
			if len(n) != 1 {
				return nil, fmt.Errorf("reduce with axis: axis must be a scalar")
			}
			axis = n[0]
		}
	}

	if _, ok := r.(apl.Table); ok {
		return reduceTable(a, f, l, r, false) // axis is ignored.
	} else if _, ok := r.(apl.Object); ok {
		return reduceTable(a, f, l, r, false)
	}

	if l != nil {
		return nwise(a, f, l, r, axis)
	}

	// If R is a scalar, the operation is not applied and Z←R
	// Single element arrays return the element.
	ar, ok := r.(apl.Array)
	if ok == false {
		return r, nil
	}
	shape := ar.Shape()
	if len(shape) == 1 && shape[0] == 1 {
		return ar.At(0), nil // higher rank single element arrays are reduced.
	}

	if len(shape) == 0 {
		if id := identityItem(f.(apl.Value)); id == nil {
			return nil, fmt.Errorf("no identity item for reduce over empty array")
		} else {
			return id, nil
		}
	}

	if axis < 0 {
		axis = len(shape) + axis
	}
	if axis < 0 || axis >= len(shape) {
		return nil, fmt.Errorf("reduce: axis rank is %d but axis %d", len(shape), axis)
	}

	// n is the number of values being reduced (the length if the reduction axis).
	n := shape[axis]

	// Dims is the shape of the result.
	dims := make([]int, len(shape)-1)
	k := 0
	for i := range shape {
		if i != axis {
			dims[k] = shape[i]
			k++
		}
	}

	// If the length of the axis is 1, the result is a reshape.
	// TODO: if the length of any other axis is 0, this should be triggered as well.
	if n == 1 {
		if rs, ok := ar.(apl.Reshaper); ok {
			return rs.Reshape(dims), nil
		} else {
			return nil, fmt.Errorf("reduce with axis length 1: cannot reshape %T", ar)
		}
	}
	if n == 0 {
		// If the axis is 0, apply an identity function, DyaRef p 169
		if id := identityItem(f.(apl.Value)); id == nil {
			return nil, fmt.Errorf("reduce empty axis: cannot get identify item for %T", f)
		} else {
			ida := apl.MixedArray{Dims: []int{1}, Values: []apl.Value{id}}
			return ida.Reshape(dims), nil
		}
	}

	// Reduce directly, if R is a vector.
	if len(shape) == 1 {
		vec := make([]apl.Value, shape[0])
		for i := range vec {
			vec[i] = ar.At(i)
		}
		v, err := reduce(a, vec, f)
		return v, err
	}

	// Create a new array with the given axis removed.
	values := make([]apl.Value, apl.ArraySize(apl.MixedArray{Dims: dims}))
	v := apl.MixedArray{
		Dims:   dims,
		Values: values,
	}

	vec := make([]apl.Value, n)
	ic, sidx := apl.NewIdxConverter(shape)
	tidx := make([]int, len(dims))
	for k := range v.Values {
		// Copy target index over the source index,
		// leaving the reduced axis unset.
		copy(sidx, tidx[:axis])
		copy(sidx[axis+1:], tidx[axis:])
		// Iterate over the reduced axis
		for i := range vec {
			sidx[axis] = i
			// TODO: maybe this could be done more efficiently
			// e.g. by iteration with a fixed increase.
			vec[i] = ar.At(ic.Index(sidx))
		}
		apl.IncArrayIndex(tidx, dims)

		if res, err := reduce(a, vec, f); err != nil {
			return nil, fmt.Errorf("cannot reduce: %s", err)
		} else {
			v.Values[k] = res
		}
	}
	return v, nil
}

// ScanArray is the derived function f\ .
func scanArray(a *apl.Apl, f apl.Value, axis int) apl.Function {
	return function(func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		return scanfunc(a, f.(apl.Function), L, R, axis)
	})
}

func scanfunc(a *apl.Apl, f apl.Function, L, R apl.Value, axis int) (apl.Value, error) {
	if L != nil {
		return nil, fmt.Errorf("scan: derived function is not defined for dyadic context")
	}

	if _, ok := R.(apl.Axis); ok {
		if r, n, err := splitAxis(a, R); err != nil {
			return nil, err
		} else {
			R = r
			if len(n) != 1 {
				return nil, fmt.Errorf("scan with axis: axis must be a scalar")
			}
			axis = n[0]
		}
	}

	if _, ok := R.(apl.Table); ok {
		return reduceTable(a, f, L, R, true) // axis is ignored.
	} else if _, ok := R.(apl.Object); ok {
		return reduceTable(a, f, L, R, true)
	}
	if c, ok := R.(apl.Channel); ok {
		return scanChannel(a, f, c)
	}

	ar, ok := R.(apl.Array)
	if ok == false {
		// If R is scalar, return unchanged.
		return R, nil
	}

	// The result has the same shape as R.
	dims := apl.CopyShape(ar)
	res := apl.MixedArray{
		Values: make([]apl.Value, apl.ArraySize(ar)),
		Dims:   dims,
	}

	if len(dims) == 0 {
		return apl.EmptyArray{}, nil
	}

	if axis < 0 {
		axis = len(dims) + axis
	}
	if axis < 0 || axis >= len(dims) {
		return nil, fmt.Errorf("scan: axis rank is %d but axis %d", len(dims), axis)
	}

	// Shortcut, if R is a vector
	if len(dims) == 1 {
		vec := make([]apl.Value, dims[0])
		for i := range vec {
			vec[i] = ar.At(i)
		}
		vec, err := scan(a, vec, f)
		if err != nil {
			return nil, err
		}
		return apl.MixedArray{
			Values: vec,
			Dims:   []int{len(vec)},
		}, nil
	}

	// Loop over the indexes, with the scan axis length set to 1.
	lidx := apl.CopyShape(ar)
	lidx[axis] = 1
	ic, idx := apl.NewIdxConverter(dims)
	vec := make([]apl.Value, dims[axis])
	for i := 0; i < apl.ArraySize(apl.MixedArray{Dims: lidx}); i++ {
		// Build the scan vector, by iterating over the axis.
		for k := range vec {
			idx[axis] = k
			vec[k] = ar.At(ic.Index(idx))
		}
		vals, err := scan(a, vec, f)
		if err != nil {
			return nil, err
		}

		// Assign the values to the destination indexes.
		for k := range vals {
			idx[axis] = k
			res.Values[ic.Index(idx)] = vals[k]
		}

		// Reset the index vector and increment.
		idx[axis] = 0
		apl.IncArrayIndex(idx, lidx)
	}
	return res, nil
}

func scanChannel(a *apl.Apl, f apl.Function, c apl.Channel) (apl.Value, error) {
	var vec []apl.Value
	var err error
	var s apl.Value
	for v := range c[0] {
		if vec == nil {
			vec = append(vec, v)
		} else {
			s, err = f.Call(a, vec[len(vec)-1], v)
			if err != nil {
				break
			}
			vec = append(vec, s)
		}
	}
	c.Close()
	if err != nil {
		return nil, err
	} else if vec == nil {
		return apl.EmptyArray{}, nil
	} else {
		return apl.MixedArray{Dims: []int{len(vec)}, Values: vec}, nil
	}
}

// replicate, compress
// L is an index array. Only vectors are allowed.
func replicate(a *apl.Apl, L apl.Value, axis int) apl.Function {
	return function(func(a *apl.Apl, dummyL, R apl.Value) (apl.Value, error) {
		// Replicate should not be an operator, but a dyadic function instead.
		// We use LO as the left argument instead.
		if dummyL != nil {
			return nil, fmt.Errorf("replicate: derived function cannot be called dyadically")
		}
		return Replicate(a, L, R, axis)
	})
}

// Replicate is the function L over R (L/R) where L and R are arrays.
func Replicate(a *apl.Apl, L, R apl.Value, axis int) (apl.Value, error) {
	ai, ar, ax, err := commonReplExp(a, L, R, axis)
	if err != nil {
		return nil, fmt.Errorf("replicate: %s", err)
	}
	axis = ax

	rs := ar.Shape()

	// If L is a 1-element vector (or was a scalar), extend it to match the axis of R.
	if ai.Dims[0] == 1 && rs[axis] > 1 {
		n := ai.Ints[0]
		ai = apl.IntArray{
			Dims: []int{rs[axis]},
		}
		ai.Ints = make([]int, apl.ArraySize(ai))
		for i := range ai.Ints {
			ai.Ints[i] = n
		}
	}
	if ai.Dims[0] != rs[axis] {
		return nil, fmt.Errorf("replicate: length of L must conform to length of R[axis]")
	}

	iscompress := true
	for i := range ai.Ints {
		if ai.Ints[i] < 0 || ai.Ints[i] > 1 {
			iscompress = false
			break
		}
	}
	if iscompress {
		return compress(a, ai, ar, axis)
	}

	// Replicate along axis.
	shape := apl.CopyShape(ar)
	count := 0
	var axismap []int
	for k, n := range ai.Ints {
		if n > 0 {
			count += n
			for i := 0; i < n; i++ {
				axismap = append(axismap, k)
			}
		} else if n < 0 {
			count += -n
			for i := 0; i < -n; i++ {
				axismap = append(axismap, -1)
			}
		}
	}
	shape[axis] = int(count)
	res := apl.MixedArray{Dims: shape}
	res.Values = make([]apl.Value, apl.ArraySize(res))
	ic, idx := apl.NewIdxConverter(rs)
	dst := make([]int, len(shape))
	for i := range res.Values {
		k := dst[axis]
		if n := axismap[k]; n == -1 {
			res.Values[i] = apl.Int(0) // TODO: When is a Fill value different from 0?
		} else {
			copy(idx, dst)
			idx[axis] = int(n)
			res.Values[i] = ar.At(ic.Index(idx)) // TODO copy
		}
		apl.IncArrayIndex(dst, shape)
	}
	return res, nil
}

// expand.
// L is an index array. Only vectors are allowed.
func expand(a *apl.Apl, L apl.Value, axis int) apl.Function {
	return function(func(a *apl.Apl, dummyL, R apl.Value) (apl.Value, error) {
		// Expand should not be an operator, but a dyadic function instead.
		// We use LO as the left argument instead.
		if dummyL != nil {
			return nil, fmt.Errorf("expand: derived function cannot be called dyadically")
		}
		return Expand(a, L, R, axis)
	})
}

// Expand is the function L\R where L and R are arrays.
func Expand(a *apl.Apl, L, R apl.Value, axis int) (apl.Value, error) {
	// Special case: L is empty.
	if _, ok := L.(apl.EmptyArray); ok {
		if ar, ok := R.(apl.Array); ok == false {
			return apl.EmptyArray{}, nil
		} else {
			rs := ar.Shape()
			if ir, ok := ar.(apl.IntArray); ok {
				fmt.Println("ar is index array: dims:", ir.Dims)
			}

			ax := int(axis)
			if ax < 0 {
				ax = len(rs) + ax
			}
			if ax >= 0 && len(rs) >= ax && rs[ax] == 0 {
				return apl.IntArray{
					Dims: apl.CopyShape(ar),
				}, nil
			}
			return nil, fmt.Errorf("expand: L is empty, R must be scalar")
		}
	}

	ai, ar, ax, err := commonReplExp(a, L, R, axis)
	if err != nil {
		return nil, fmt.Errorf("expand: %s", err)
	}
	axis = ax

	// Special case: R is empty. L may be 0 and is returned.
	if _, ok := R.(apl.EmptyArray); ok {
		if len(ai.Ints) == 1 && ai.Ints[0] == 0 {
			return ai, nil
		} else {
			return nil, fmt.Errorf("exand: R is empty, but L is not 0")
		}
	}

	// The shape of the result is the shape of R,
	// with the length of the axis set to +/1⌈|L.
	shape := apl.CopyShape(ar)
	sum := 0
	for _, k := range ai.Ints {
		if k < 0 {
			k = -k
		}
		if k > 1 {
			sum += k
		} else {
			sum++
		}
	}
	shape[axis] = int(sum)

	res := apl.MixedArray{Dims: shape}
	n := apl.ArraySize(res)
	res.Values = make([]apl.Value, n)

	short := apl.CopyShape(res)
	short[axis] = 1

	ic, idx := apl.NewIdxConverter(ar.Shape())
	dic, dst := apl.NewIdxConverter(shape)
	for i := 0; i < n/shape[axis]; i++ {
		copy(idx, dst)
		d := 0
		j := 0 // Count positive indexes in L.
		for _, k := range ai.Ints {
			if k > 0 {
				idx[axis] = j
				j++
				for m := 0; m < int(k); m++ {
					dst[axis] = d
					d++
					res.Values[dic.Index(dst)] = ar.At(ic.Index(idx)) // TODO copy
				}
			} else if k == 0 {
				dst[axis] = d
				d++
				res.Values[dic.Index(dst)] = apl.Int(0)
			} else if k < 0 {
				for m := 0; m < int(-k); m++ {
					dst[axis] = d
					d++
					res.Values[dic.Index(dst)] = apl.Int(0)
				}
			}
		}
		dst[axis] = 0
		apl.IncArrayIndex(dst, short)
	}
	return res, nil
}

// commonReplExp is the common input preprocessing for replicate and expand.
func commonReplExp(a *apl.Apl, L, R apl.Value, axis int) (apl.IntArray, apl.Array, int, error) {
	ai := L.(apl.IntArray)
	if len(ai.Dims) != 1 {
		return ai, nil, axis, fmt.Errorf("LO must be a vector")
	}

	// R may contain an axis from a bracket expression, which overwrites axis.
	if r, n, err := splitAxis(a, R); err != nil {
		return ai, nil, axis, err
	} else {
		R = r
		if len(n) == 1 {
			axis = n[0]
		} else if len(n) > 1 {
			return ai, nil, axis, fmt.Errorf("compress/replicate: axis must be a scalar")
		}
	}

	// If R is scalar or a single-element array, convert it to (⍴L)⍴B
	// If R is a scalar, convert it to a single element array.
	ar, ok := R.(apl.Array)
	if ok == false {
		r := apl.MixedArray{
			Dims:   []int{1},
			Values: []apl.Value{R},
		}
		ar = r
	}
	rs := ar.Shape()

	// Special case for empty R.
	if len(rs) == 0 {
		return ai, ar, 0, nil
	}
	if axis < 0 {
		axis = len(rs) + axis
	}
	if axis < 0 {
		return ai, nil, axis, fmt.Errorf("axis is negative")
	}

	// If R has size 1 along the selected axis and L is larger, extend R.
	if rs[axis] == 1 && len(ai.Ints) > 1 {
		shape := apl.CopyShape(ar)
		shape[axis] = len(ai.Ints)
		r := apl.MixedArray{
			Dims: shape,
		}
		r.Values = make([]apl.Value, apl.ArraySize(r))
		ic, idx := apl.NewIdxConverter(rs)
		dst := make([]int, len(shape))
		for i := range r.Values {
			copy(idx, dst)
			idx[axis] = 0
			r.Values[i] = ar.At(ic.Index(idx)) // TODO copy
			apl.IncArrayIndex(dst, shape)
		}
		ar = r
		rs = ar.Shape()
	}

	if axis < 0 || axis >= len(rs) {
		return ai, nil, axis, fmt.Errorf("replicate: axis out of range: %d", axis)
	}
	return ai, ar, axis, nil
}

// L is an index vector with boolean values.
// R is an array.
func compress(a *apl.Apl, L, R apl.Value, axis int) (apl.Value, error) {
	ai := L.(apl.IntArray)
	ar := R.(apl.Array)
	rs := ar.Shape()

	shape := apl.CopyShape(ar)
	count := 0
	for _, b := range ai.Ints {
		count += int(b)
	}
	shape[axis] = count

	res := apl.MixedArray{
		Dims: shape,
	}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	ridx := make([]int, len(rs))
	n := 0
	for i := 0; i < apl.ArraySize(ar); i++ {
		if b := ridx[axis]; ai.Ints[b] == 1 {
			res.Values[n] = ar.At(i)
			n++
		}
		apl.IncArrayIndex(ridx, rs)
	}
	return res, nil
}

func reduce(a *apl.Apl, vec []apl.Value, d apl.Function) (apl.Value, error) {
	var err error
	v := vec[len(vec)-1] // TODO: copy?
	for i := len(vec) - 2; i >= 0; i-- {
		v, err = d.Call(a, vec[i], v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func reduceChannel(a *apl.Apl, L apl.Value, f apl.Function, c apl.Channel) (apl.Value, error) {
	if L != nil {
		return nil, fmt.Errorf("TODO: n-wise channel reduction")
	}
	var res apl.Value
	var err error
	for v := range c[0] {
		if res == nil {
			res = v
		} else {
			res, err = f.Call(a, res, v)
			if err != nil {
				break
			}
		}
	}
	c.Close()
	if err != nil {
		return nil, err
	} else if res == nil {
		return apl.EmptyArray{}, nil
	} else {
		return res, nil
	}
}

// reduceTable applies to tables or objects/dicts.
// It always uses the same axis.
// It is used from both reduce and scan.
func reduceTable(a *apl.Apl, f apl.Function, L, R apl.Value, scan bool) (apl.Value, error) {
	var o apl.Object
	istable := false
	if t, ok := R.(apl.Table); ok {
		istable = true
		o = t.Dict
	} else {
		o = R.(apl.Object)
	}
	keys := o.Keys()
	d := apl.Dict{K: make([]apl.Value, len(keys)), M: make(map[apl.Value]apl.Value)}
	var err error
	for i, k := range keys {
		col := o.At(a, k)
		var ar apl.Array
		if v, ok := col.(apl.Array); ok {
			ar = v
		} else {
			ar = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{col}} // TODO: copy?
		}
		var v apl.Value
		if scan {
			v, err = scanfunc(a, f, L, ar, -1)
		} else {
			v, err = reduct(a, f, L, ar, -1)
		}
		if err != nil {
			return nil, err
		}
		if _, ok := v.(apl.Array); ok == false && istable {
			// Tables never contain scalars. Values must be enlisted.
			v = apl.List{v}
		}
		d.K[i] = k
		d.M[k] = v
	}
	if istable {
		return dict2table(a, &d)
	}
	return &d, nil
}

// copy from primitives/table.go
func dict2table(a *apl.Apl, d *apl.Dict) (apl.Table, error) {
	rows := 0
	if len(d.K) > 0 {
		v := d.At(a, d.K[0])
		if ar, ok := v.(apl.Array); ok {
			shape := ar.Shape()
			if len(shape) == 1 {
				rows = shape[0]
			}
		}
	}
	return apl.Table{Dict: d, Rows: rows}, nil
}

func scan(a *apl.Apl, vec []apl.Value, d apl.Function) ([]apl.Value, error) {
	// The ith element of the result is: d/I↑V
	res := make([]apl.Value, len(vec))
	res[0] = vec[0] // TODO: copy?
	for i := 1; i < len(res); i++ {
		if v, err := reduce(a, vec[:i+1], d); err != nil {
			return nil, err
		} else {
			res[i] = v
		}
	}
	return res, nil
}

// Nwise is the function handle for n-wise recution.
// l must be a scalar (integer) or a 1 element vector.
func nwise(a *apl.Apl, f apl.Function, L, R apl.Value, axis int) (apl.Value, error) {

	var n int
	neg := false
	to := ToScalar(ToIndex(nil))
	if idx, ok := to.To(a, L); ok == false {
		return nil, fmt.Errorf("nwise reduction: left argument must be a scalar integer: %T", L)
	} else {
		n = int(idx.(apl.Int))
		if n < 0 {
			n = -n
			neg = true
		}
	}

	if _, ok := R.(apl.Table); ok {
		return nil, fmt.Errorf("n-wise reduction over table should not happen")
	}
	if _, ok := R.(apl.Object); ok {
		return nil, fmt.Errorf("n-wise reduction over object should not happen")
	}

	ar, ok := R.(apl.Array)
	if ok == false {
		return nil, fmt.Errorf("n-wise reduction: right argument must be an array")
	}
	rs := ar.Shape()

	if _, ok := R.(apl.EmptyArray); ok {
		if n == 0 {
			return apl.IntArray{Dims: []int{1}, Ints: []int{0}}, nil
		} else if n == 1 {
			return apl.EmptyArray{}, nil
		} else {
			return nil, fmt.Errorf("n-wise reduction: length error")
		}
	}

	if axis == -1 {
		axis = len(rs) + axis
	}
	if axis < 0 || axis >= len(rs) {
		return nil, fmt.Errorf("n-wise reduction: axis out of range")
	}

	shape := apl.CopyShape(ar)
	shape[axis] -= n - 1
	if n-rs[axis] > 2 {
		return nil, fmt.Errorf("n-wise reduction: length error")
	}

	res := apl.MixedArray{Dims: shape}
	res.Values = make([]apl.Value, apl.ArraySize(res))
	if len(res.Values) == 0 {
		return res, nil
	}

	if n == 0 {
		var id apl.Value
		if p, ok := f.(apl.Primitive); ok {
			id = identityItem(p)
		}
		if id == nil {
			return nil, fmt.Errorf("n-wise reduction: unknown identify function")
		}
		for i := range res.Values {
			res.Values[i] = id
		}
		return res, nil
	}

	// Fast accumulative algorithm for +/ and ×/
	var inv apl.Function
	if p, ok := f.(apl.Primitive); ok {
		if p == "+" {
			inv = apl.Primitive("-")
		} else if p == "×" {
			inv = apl.Primitive("÷")
		}
	}

	// Iterate over all items, except for the reduction axis.
	axlen := rs[axis]
	vec := make([]apl.Value, axlen)
	xs := apl.CopyShape(ar)
	xs[axis] = 1
	outer := apl.ArraySize(apl.MixedArray{Dims: xs})
	ic, idx := apl.NewIdxConverter(rs)
	dc, dst := apl.NewIdxConverter(res.Dims)
	for i := 0; i < outer; i++ {
		for k := range vec {
			j := k
			if neg {
				j = axlen - 1 - k
			}
			idx[axis] = j
			vec[k] = ar.At(ic.Index(idx))
		}
		if err := applyNwise(a, vec, n, f, inv); err != nil {
			return nil, err
		}
		copy(dst, idx)
		for k := 0; k < axlen-n+1; k++ {
			j := k
			if neg {
				j = axlen - n - k
			}
			dst[axis] = j
			res.Values[dc.Index(dst)] = vec[k]
		}

		idx[axis] = 0
		apl.IncArrayIndex(idx, xs)
	}

	return res, nil
}

func applyNwise(a *apl.Apl, vec []apl.Value, n int, f, g apl.Function) error {
	var err error
	reduce := func(x []apl.Value) apl.Value {
		r := x[len(x)-1]
		for i := len(x) - 2; i >= 0; i-- {
			r, err = f.Call(a, x[i], r)
			if err != nil {
				return nil
			}
		}
		return r
	}

	// Fast path: Moving window with accumulator.
	if g != nil && n > 3 {
		var acc apl.Value
		window := make([]apl.Value, n)
		p := 0
		reduce = func(x []apl.Value) apl.Value {
			// Initial call: fill the window.
			if acc == nil {
				for i, v := range x {
					window[i] = v
				}
				acc = x[0]
				for _, v := range x[1:] {
					acc, err = f.Call(a, acc, v)
					if err != nil {
						return nil
					}
				}
			} else {
				xnew := x[len(x)-1]
				acc, err = g.Call(a, acc, window[p])
				if err != nil {
					return nil
				}
				window[p] = xnew
				acc, err = f.Call(a, acc, xnew)
				if err != nil {
					return nil
				}
				p++
				if p == len(window) {
					p = 0
				}
			}
			return acc
		}
	}

	for i := 0; i < len(vec)-n+1; i++ {
		vec[i] = reduce(vec[i : i+n])
		if err != nil {
			return err
		}
	}
	return nil
}

// reduceTack is the derived function from ⊣/ or ⊢/ .
type reduceTack bool

func (first reduceTack) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if L != nil {
		return nil, fmt.Errorf("tack-reduce can only be applied monadically")
	}
	ar, ok := R.(apl.Array)
	if ok == false {
		return R, nil
	}
	shape := ar.Shape()
	if len(shape) == 0 {
		return R, nil
	}

	// Reduce a vector to a scalar.
	if len(shape) == 1 {
		if shape[0] <= 0 {
			return apl.EmptyArray{}, nil
		}
		var v apl.Value
		if first {
			v = ar.At(0)
		} else {
			v = ar.At(shape[0] - 1)
		}
		return v, nil
	}

	// Create a new array
	inner := shape[len(shape)-1]
	newshape := apl.CopyShape(ar)
	newshape = newshape[:len(newshape)-1]
	ret := apl.MixedArray{Dims: newshape}
	ret.Values = make([]apl.Value, apl.ArraySize(ret))
	i := 0
	n := 0 // index over inner axis.
	for k := 0; k < apl.ArraySize(ar); k++ {
		if first && n == 0 {
			ret.Values[i] = ar.At(k)
			i++
		} else if first == false && n == inner-1 {
			ret.Values[i] = ar.At(k)
			i++
		}
		n++
		if n == inner {
			n = 0
		}
	}
	return ret, nil
}
