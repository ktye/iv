package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/operators"
)

func init() {
	register(primitive{
		symbol: "⊥",
		doc:    "decode, polynom, base value",
		Domain: Dyadic(Split(ToArray(nil), ToArray(nil))),
		fn:     decode,
	})
}

func decode(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	al := L.(apl.Array)
	ar := R.(apl.Array)
	ls := al.Shape()
	rs := ar.Shape()

	// The last axis of L must match the first axis of R.
	// Single element axis are extended.
	if n := ls[len(ls)-1]; n != rs[0] {
		if n == 1 {
			al, ls = extendAxis(al, len(ls)-1, rs[0])
		} else if rs[0] == 1 {
			ar, rs = extendAxis(ar, 0, n)
		} else {
			return nil, fmt.Errorf("decode: last axis of L must match first axis of R: %v %v", ls, rs)
		}
	}

	// The result of decode is a scalar product between a power matrix and R.
	// The power matrix multiplies L along the last axis recursively from right to left,
	// similar as the Index method of apl.IdxConverter.
	p := apl.GeneralArray{
		Values: make([]apl.Value, apl.ArraySize(al)),
		Dims:   apl.CopyShape(al),
	}
	for i := range p.Values {
		v, err := al.At(i)
		if err != nil {
			return nil, err
		}
		p.Values[i] = v // TODO: copy?
	}
	N := ls[len(ls)-1]

	// Shift last axis by 1 to the left than multiply scan from the right.
	fmul := arith2("×", mul2)
	for off := 0; off < len(p.Values); off += N {
		for k := 0; k < N-1; k++ {
			p.Values[off+k] = p.Values[off+k+1]
		}
		p.Values[off+N-1] = apl.Index(1)
		for k := off + N - 1; k >= off+1; k-- {
			v, err := fmul(a, p.Values[k], p.Values[k-1])
			if err != nil {
				return nil, err
			}
			p.Values[k-1] = v
		}
	}

	dot := operators.Scalarproduct(a, apl.Primitive("+"), apl.Primitive("×"))
	return dot.Call(a, p, ar)
}

// extendAxis extends the axis of length 1 to n
func extendAxis(ar apl.Array, axis, n int) (apl.Array, []int) {
	res := apl.GeneralArray{Dims: apl.CopyShape(ar)}
	res.Dims[axis] = n
	res.Values = make([]apl.Value, apl.ArraySize(res))
	ridx := make([]int, len(res.Dims))
	ic, idx := apl.NewIdxConverter(ar.Shape())
	for i := range res.Values {
		copy(idx, ridx)
		idx[axis] = 0
		v, _ := ar.At(ic.Index(idx))
		res.Values[i] = v // TODO copy?
		apl.IncArrayIndex(ridx, res.Dims)
	}
	return res, res.Dims
}
