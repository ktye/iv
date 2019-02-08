package primitives

import (
	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⍷",
		doc:    "find",
		Domain: Dyadic(Split(ToArray(nil), ToArray(nil))),
		fn:     find,
	})
}

func find(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.EmptyArray); ok {
		return apl.EmptyArray{}, nil
	}

	var al apl.Array
	el, Lempty := L.(apl.EmptyArray)
	if Lempty == false {
		al = L.(apl.Array)
	} else {
		al = el
	}
	ls := al.Shape()

	ar := R.(apl.Array)
	rs := ar.Shape()

	res := apl.IntArray{Dims: apl.CopyShape(ar)}
	res.Ints = make([]int, apl.ArraySize(res))

	// If the rank of L is arger than the rank of R, nothing is found.
	if len(ls) > len(rs) {
		return res, nil
	}

	// If the rank of L is smaller than the rank of R, fill is with ones
	// at the beginning.
	if d := len(rs) - len(ls); d > 0 {
		shape := apl.CopyShape(ar)
		for i := range shape {
			if i < d {
				shape[i] = 1
			} else {
				shape[i] = ls[i-d]
			}
		}
		l := apl.MixedArray{Dims: shape, Values: make([]apl.Value, apl.ArraySize(al))}
		for i := range l.Values {
			l.Values[i] = al.At(i)
		}
		al = l
		ls = shape
	}
	nl := apl.ArraySize(al)

	feq := arith2("=", compare("="))
	ic, idx := apl.NewIdxConverter(rs)
	for i := range res.Ints {
		if nl > len(res.Ints)-i {
			res.Ints[i] = 0
		} else {
			iseq := true
			for k := 0; k < len(idx); k++ {
				idx[k] = 0
			}
			for k := 0; k < nl; k++ {
				eq, err := feq(a, al.At(k), ar.At(i+ic.Index(idx)))
				if err != nil {
					return nil, err
				}
				if eq.(apl.Bool) == false {
					iseq = false
					break
				}
				apl.IncArrayIndex(idx, ls)
			}
			if iseq {
				res.Ints[i] = 1
			}
		}
	}
	return res, nil
}
