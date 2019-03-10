package operators

import (
	"fmt"
	"io"

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

// rank is extended for sending subarrays over a channel:
//	<⍤3 A  send rank-3 subarray of A sequentially over the returned channel
//	<⍤3 C  read strings from input channel C, parse rank-3 subarrays and send them over a return channel
func rank(a *apl.Apl, LO, RO apl.Value) apl.Function {
	// Cell, Frame, Conform: ISO: 9.3.3, p123
	// Rank operator: ISO; 9.3.4, p 124
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		f := LO.(apl.Function)
		doSend := false
		if pf, ok := f.(apl.Primitive); ok && string(pf) == "<" {
			doSend = true
		}

		ai := RO.(apl.IntArray)
		if len(ai.Shape()) != 1 {
			return nil, fmt.Errorf("rank: RO must be a vector")
		}

		// p, q and r are the values of RO.
		var p, q, r int
		if n := apl.ArraySize(ai); n < 1 || n > 3 {
			return nil, fmt.Errorf("rank: RO vector has %d elements, must be between 1..3", n)
		} else if n == 1 {
			p, q, r = int(ai.Ints[0]), int(ai.Ints[0]), int(ai.Ints[0])
			// Special case: <⍤R C
			if c, ok := R.(apl.Channel); ok && doSend == true {
				return sendParseSubArray(a, r, c)
			}
		} else if n == 2 {
			p, q, r = int(ai.Ints[1]), int(ai.Ints[0]), int(ai.Ints[1])
		} else if n == 3 {
			p, q, r = int(ai.Ints[0]), int(ai.Ints[1]), int(ai.Ints[2])
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
				return x.At(n), nil
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
			cell := apl.MixedArray{Dims: subshape}
			cell.Values = make([]apl.Value, apl.ArraySize(cell))
			m := len(cell.Values)
			for i := range cell.Values {
				cell.Values[i] = x.At(n*m + i).Copy()
			}
			return a.UnifyArray(cell), nil
		}

		// subcells returns the number of rank-cells of x.
		subcells := func(x apl.Array, rank int) int {
			s := x.Shape()
			// The number is the product of the frame of x with respect to rank.
			return apl.ArraySize(apl.MixedArray{Dims: s[:len(s)-rank]})
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
				results = append(results, v.Copy())
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
				if doSend {
					results = append(results, s)
				} else {
					v, err := f.Call(a, nil, s)
					if err != nil {
						return nil, err
					}
					results = append(results, v.Copy())
				}
			}
			if doSend {
				c := apl.NewChannel()
				go c.SendAll(results)
				return c, nil
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
					ga := apl.MixedArray{Dims: common}
					ga.Values = make([]apl.Value, apl.ArraySize(ga))
					for i := range ga.Values {
						ga.Values[i] = results[n].Copy()
					}
					results[n] = a.UnifyArray(ga)
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
					idx := apl.IntArray{Dims: []int{len(common)}}
					idx.Ints = make([]int, len(common))
					for i := range common {
						idx.Ints[i] = int(common[i])
					}
					var err error
					vr, err = Take(a, idx, vr, nil)
					if err != nil {
						return nil, err
					}
				}
				results[n] = vr.Copy()
			}
		}

		// The result has the shape: frame, conform
		res := apl.MixedArray{}
		res.Dims = append(res.Dims, frame...)
		res.Dims = append(res.Dims, common...)
		res.Values = make([]apl.Value, apl.ArraySize(res))
		if len(common) == 0 {
			if len(results) != len(res.Values) {
				return nil, fmt.Errorf("rank: unexpected number of scalar results %d instead of %d", len(results), len(res.Values)) // Should not happen
			}
			res.Values = results
			return a.UnifyArray(res), nil
		}
		commonsize := apl.ArraySize(apl.MixedArray{Dims: common})
		off := 0
		for i := range results {
			if len(common) == 0 {
				res.Values[i] = results[i].Copy()
			} else {
				vr := results[i].(apl.Array)
				for n := 0; n < commonsize; n++ {
					res.Values[off+n] = vr.At(n).Copy()
				}
				off += commonsize
			}
		}
		return a.UnifyArray(res), nil
	}
	return function(derived)
}

// sendParseSubArray assembles an array of the given rank from strings read on channel c.
// Strings are parsed to numbers or strings.
// It returns a channel and sends arrays of the rank.
func sendParseSubArray(a *apl.Apl, rank int, in apl.Channel) (apl.Value, error) {
	out := apl.NewChannel()
	go func() {
		defer close(out[0])
		scn := apl.RuneScanner{C: in, O: out}
		for {
			v, err := a.ScanRankArray(&scn, rank)
			if err == io.EOF || err == io.ErrClosedPipe {
				return
			} else if err != nil {
				out[0] <- apl.Error{err}
				return
			}
			select {
			case _, ok := <-in[1]:
				if !ok {
					return
				}
			case out[0] <- v:
			}
		}
	}()
	return out, nil
}

// Take is defined and exported here, because it is used by both the rank operator and the take primitive function.

func Take(a *apl.Apl, ai apl.IntArray, ar apl.Array, x []int) (apl.Array, error) {
	rs := ar.Shape()

	if x == nil {
		x = make([]int, len(ai.Ints))
		for i := range x {
			x[i] = i
		}
	}

	if len(ai.Ints) != len(x) {
		return nil, fmt.Errorf("take: length of L and axis must match")
	}

	shape := make([]int, len(rs))
	for i := range shape {
		shape[i] = rs[i]
	}
	for i, k := range x {
		shape[k] = int(ai.Ints[i])
		if shape[k] < 0 {
			shape[k] = -shape[k]
		}
	}

	// Offset for negative arguments.
	off := make([]int, len(ar.Shape()))
	for i := range off {
		if i < len(ai.Ints) {
			if n := ai.Ints[i]; n < 0 {
				k := x[i]
				off[k] = rs[k] + int(n)
			}
		}
	}

	res := apl.MakeArray(ar, shape)
	var z apl.Value
	if u, ok := res.(apl.Uniform); ok {
		z = u.Zero()
	} else {
		z = apl.Int(0)
	}
	idx := make([]int, len(shape))
	ic, src := apl.NewIdxConverter(ar.Shape())
	for i := 0; i < res.Size(); i++ {
		copy(src, idx)
		zero := false
		for k, n := range off {
			src[k] += n
			if src[k] < 0 || src[k] >= rs[k] {
				zero = true
			}
		}
		if zero {
			res.Set(i, z)
		} else {
			res.Set(i, ar.At(ic.Index(src)).Copy())
		}
		apl.IncArrayIndex(idx, shape)
	}
	return res, nil
}
