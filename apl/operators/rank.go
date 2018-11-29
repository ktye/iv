package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "⍤",
		Domain:  DyadicOp(Split(Function(nil), ToIndexArray(nil))),
		doc:     "rank",
		derived: rank,
	})
}

func rank(a *apl.Apl, LO, RO apl.Value) apl.Function {
	// Cell, Frame, Conform: ISO: 9.3.3, p123
	// Rank operator: ISO; 9.3.4, p 124
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		f := LO.(apl.Function)
		ai := RO.(apl.IndexArray)
		if len(ai.Shape()) != 1 {
			return nil, fmt.Errorf("rank: RO must be a vector")
		}

		// p, q and r are the values of RO.
		var p, q, r int
		if n := apl.ArraySize(ai); n < 1 || n > 3 {
			return nil, fmt.Errorf("rank: RO vector has %d elements, must be between 1..3", n)
		} else if n == 1 {
			p, q, r = ai.Ints[0], ai.Ints[0], ai.Ints[0]
		} else if n == 2 {
			p, q, r = ai.Ints[1], ai.Ints[0], ai.Ints[1]
		} else if n == 3 {
			p, q, r = ai.Ints[0], ai.Ints[1], ai.Ints[2]
		}

		ar, ok := R.(apl.Array)
		if ok == false {
			return nil, fmt.Errorf("rank: right argument must be an array: %T", R)
		}
		rs := ar.Shape()

		var al apl.Array
		var ls []int
		if L != nil {
			al, ok = L.(apl.Array)
			ls = al.Shape()
		}

		// subcell returns the rank-subcell number i of the array x.
		subcell := func(x apl.Array, rank int, n int) (apl.Value, error) {
			if rank == 0 {
				return x.At(n)
			}
			shape := x.Shape()
			if rank < 0 || rank > len(shape) {
				return nil, fmt.Errorf("cannot get %d-subcell of array with rank %d", rank, len(shape))
			}
			if n < 0 || n > shape[len(shape)-1-rank] {
				return nil, fmt.Errorf("cannot get %d-subcell number %d: length of axis is %d", rank, n, shape[rank])
			}

			subshape := apl.CopyShape(x)
			subshape = subshape[len(subshape)-rank:]
			cell := apl.GeneralArray{Dims: subshape}
			cell.Values = make([]apl.Value, apl.ArraySize(cell))
			m := len(cell.Values)
			for i := range cell.Values {
				v, err := x.At(n*m + i)
				if err != nil {
					return nil, err
				}
				cell.Values[i] = v
			}
			return cell, nil
		}

		// subcells returns the number of rank-cells of x.
		subcells := func(x apl.Array, rank int) int {
			s := x.Shape()
			// The number is the product of the frame of x with respect to rank.
			return apl.ArraySize(apl.GeneralArray{Dims: s[:len(s)-rank]})
		}

		var err error
		var results []apl.Value
		frame := apl.CopyShape(ar)
		if al != nil {
			// Dyadic context: q specifies rank of L.
			if q < 0 {
				q += len(ls)
			}
			if q < 0 || q > len(ls) {
				return nil, fmt.Errorf("rank: q (%d) exceeds rank of left argument: %d", q, len(ls))
			}
			// r specifies rank of R.
			if r < 0 {
				r += len(rs)
			}
			if r < 0 || r > len(rs) {
				return nil, fmt.Errorf("rank: r (%d) exceeds rank of right argument: %d", r, len(rs))
			}

			// The number of subcells must match or any must be 0 and is repeated.
			ml := subcells(al, q)
			mr := subcells(ar, r)
			if ml != mr && ml > 0 && mr > 0 {
				return nil, fmt.Errorf("rank: L and R have different number of subcells")
			}
			m := ml
			if ml < mr {
				m = mr
			}

			// Apply f successsively between sub arrays of L and R specified by q and r.
			frame = frame[:len(frame)-r]
			if mr == 0 {
				frame = apl.CopyShape(al)
				frame = frame[:len(frame)-q]
			}
			var subl, subr apl.Value
			for i := 0; i < m; i++ {
				if ml == 0 {
					subl = al
				} else {
					subl, err = subcell(al, q, i)
					if err != nil {
						return nil, err
					}
				}
				if mr == 0 {
					subr = ar
				} else {
					subr, err = subcell(ar, r, i)
					if err != nil {
						return nil, err
					}
				}
				v, err := f.Call(a, subl, subr)
				if err != nil {
					return nil, err
				}
				results = append(results, v)
			}

		} else {
			// Monadic context: p specifies rank of R.
			// r specifies rank of R.
			if p < 0 {
				p += len(rs)
			}
			if p < 0 || r > len(rs) {
				return nil, fmt.Errorf("rank: p (%d) exceeds rank of right argument: %d", p, len(rs))
			}

			// Apply f successsively to sub arrays of R specified by p.
			frame = frame[:len(frame)-p]
			for i := 0; i < subcells(ar, p); i++ {
				s, err := subcell(ar, p, i)
				if err != nil {
					return nil, err
				}
				v, err := f.Call(a, nil, s)
				if err != nil {
					return nil, err
				}
				results = append(results, v)
			}
		}

		// Bring all individual results to conforming shape.
		var common []int
		for i := range results {
			if vr, ok := results[i].(apl.Array); ok {
				s := vr.Shape()
				if d := len(s) - len(common); d > 0 {
					common = append(make([]int, d), common...)
				}
				for n := 0; n < len(s); n++ {
					k := len(common) - len(s) + n
					if s[n] > common[k] {
						common[k] = s[n]
					}
				}
			}
		}
		for n := range results {
			if vr, ok := results[n].(apl.Array); ok == false {
				if len(common) > 0 {
					// Reshape scalar to common shape.
					ga := apl.GeneralArray{Dims: common}
					ga.Values = make([]apl.Value, apl.ArraySize(ga))
					for i := range ga.Values {
						ga.Values[i] = results[n] // TODO copy?
					}
					results[n] = ga
				}
			} else {
				// If rank is smaller than common rank,
				// fill ones at the start and reshape.
				shape := apl.CopyShape(vr)
				if d := len(common) - len(shape); d > 0 {
					shape = append(make([]int, d), shape...)
					for i := 0; i < d; i++ {
						shape[i] = 1
					}
				}
				if rs, ok := vr.(apl.Reshaper); ok {
					vr = rs.Reshape(shape).(apl.Array)
				}

				// If the shape is different from common, make a conforming
				// array by: common↑vr
				diffshape := false
				for i := range common {
					if common[i] != shape[i] {
						diffshape = true
						break
					}
				}
				if diffshape {
					idx := apl.IndexArray{Dims: []int{len(common)}}
					idx.Ints = make([]int, len(common))
					copy(idx.Ints, common)
					var err error
					vr, err = Take(a, idx, vr)
					if err != nil {
						return nil, err
					}
				}
				results[n] = vr
			}
		}

		// The result has the shape: frame, conform
		res := apl.GeneralArray{}
		res.Dims = append(res.Dims, frame...)
		res.Dims = append(res.Dims, common...)
		res.Values = make([]apl.Value, apl.ArraySize(res))
		if len(common) == 0 {
			if len(results) != len(res.Values) {
				return nil, fmt.Errorf("rank: unexpected number of scalar results %d instead of %d", len(results), len(res.Values)) // Should not happen
			}
			res.Values = results
			return res, nil
		}
		commonsize := apl.ArraySize(apl.GeneralArray{Dims: common})
		off := 0
		for i := range results {
			if len(common) == 0 {
				res.Values[i] = results[i]
			} else {
				vr := results[i].(apl.Array)
				for n := 0; n < commonsize; n++ {
					v, err := vr.At(n)
					if err != nil {
						return nil, err
					}
					res.Values[off+n] = v // TODO copy?
				}
				off += commonsize
			}
		}
		return res, nil
	}
	return function(derived)
}

// Take is defined and exported here, because it is used by both the rank operator and the take primitive function.

func Take(a *apl.Apl, ai apl.IndexArray, ar apl.Array) (apl.Array, error) {
	rs := ar.Shape()

	// The shape of the result is ,|L
	shape := make([]int, len(ai.Ints))
	for i, n := range ai.Ints {
		if n < 0 {
			shape[i] = -n
		} else {
			shape[i] = n
		}
	}
	res := apl.GeneralArray{Dims: shape}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	ic, J := apl.NewIdxConverter(rs)
	idx := make([]int, len(shape))
	for i := range res.Values {
		for k := range J {
			J[k] = idx[k]
			if n := ai.Ints[k]; n < 0 {
				J[k] += n + rs[k]
			}
		}
		iszero := false
		for k := range J {
			if J[k] < 0 || J[k] >= rs[k] {
				iszero = true
				break
			}
		}
		if iszero {
			res.Values[i] = apl.Index(0) // TODO: typical element of R?
		} else {
			n := ic.Index(J)
			v, err := ar.At(n)
			if err != nil {
				return nil, err
			}
			res.Values[i] = v // TODO: copy?
		}

		apl.IncArrayIndex(idx, shape)
	}
	return res, nil
}
