package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "@",
		Domain:  DyadicOp(nil),
		doc:     "at",
		derived: at,
	})
}

func at(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		// g selects values from R.
		ar, ok := R.(apl.Array)
		if ok == false {
			ar = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{R}}
		}
		rs := ar.Shape()

		// If g is a function, it must return a boolean mask.
		mask := make([]bool, apl.ArraySize(ar))
		var replshape []int
		if fg, ok := g.(apl.Function); ok {
			v, err := fg.Call(a, nil, ar)
			if err != nil {
				return nil, err
			}
			av, ok := v.(apl.Array)
			if ok == false {
				return nil, fmt.Errorf("at: function g did not return an array: %T", v)
			}
			size := apl.ArraySize(av)
			if size != len(mask) {
				return nil, fmt.Errorf("at: array returned by function g has wrong size")
			}
			for i := range mask {
				if err := apl.ArrayBounds(av, i); err != nil {
					return nil, fmt.Errorf("at: %s", err)
				}
				v := av.At(i)
				if n, ok := v.(apl.Number); ok == false {
					return nil, fmt.Errorf("at: function g did not return a number: %T", v)
				} else if b, ok := a.Tower.ToBool(n); ok == false {
					return nil, fmt.Errorf("at: number returned by function g is not a boolean: %T", n)
				} else {
					mask[i] = bool(b)
				}
			}
		} else {
			// g is an index array that selects major cells of R.
			ag, ok := g.(apl.Array)
			if ok == false {
				ag = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{g}}
			}
			var gi apl.IndexArray
			if v, ok := ToIndexArray(nil).To(a, ag); ok == false {
				return nil, fmt.Errorf("at: g is not an index array")
			} else {
				gi = v.(apl.IndexArray)
			}
			if len(gi.Dims) != 1 {
				return nil, fmt.Errorf("at: g should have rank 1: %d", len(gi.Dims))
			}
			n := apl.ArraySize(ar) / rs[0]
			for _, major := range gi.Ints {
				major -= a.Origin
				if major < 0 || major >= rs[0] {
					return nil, fmt.Errorf("at: selected major cell is out of range %d: [1, %d]", major+1, rs[0])
				}
				off := n * major
				for i := 0; i < n; i++ {
					mask[off+i] = true
				}
			}
			// Keep shape of selected subarray.
			replshape = apl.CopyShape(ar)
			replshape[0] = gi.Dims[0]
		}

		res := apl.MixedArray{Dims: apl.CopyShape(ar)}
		res.Values = make([]apl.Value, apl.ArraySize(res))

		// Number of replacements.
		n := 0
		for _, v := range mask {
			if v {
				n++
			}
		}

		repl := make([]apl.Value, n)
		var vr apl.Value

		if fn, ok := f.(apl.Function); ok {
			// Apply fn to the sub-array of R as a whole.
			if replshape == nil {
				replshape = []int{n}
			}
			re := apl.MixedArray{Dims: replshape, Values: repl}
			n := 0
			for i, m := range mask {
				if m {
					if err := apl.ArrayBounds(ar, i); err != nil {
						return nil, err
					}
					re.Values[n] = ar.At(i)
					n++
				}
			}
			if v, err := fn.Call(a, L, re); err != nil {
				return nil, err
			} else {
				vr = v
			}
		} else {
			// f is an array of replacements.
			vr = f
		}

		re, ok := vr.(apl.Array)
		if ok == false {
			re = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{vr}}
		}
		if n := apl.ArraySize(re); n == 1 {
			for i := range repl {
				repl[i] = re.At(0) // TODO: copy?
			}
		} else if n != len(repl) {
			return nil, fmt.Errorf("at: number of replacements does not match selection")
		} else {
			for i := range repl {
				if err := apl.ArrayBounds(ar, i); err != nil {
					return nil, err
				}
				repl[i] = re.At(i)
			}
		}

		k := 0
		for i := range res.Values {
			if mask[i] {
				res.Values[i] = repl[k]
				k++
			} else {
				if err := apl.ArrayBounds(ar, i); err != nil {
					return nil, err
				}
				res.Values[i] = ar.At(i)
			}
		}
		return res, nil
	}
	return function(derived)
}
