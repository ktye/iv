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
	register(primitive{
		symbol: "⊤",
		doc:    "encode, representation",
		Domain: Dyadic(nil),
		fn:     encode,
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

func encode(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	al, lok := L.(apl.Array)
	ar, rok := R.(apl.Array)

	// If L is a scalar, return L|R.
	if lok == false {
		if rok == false {
			mod := arith2("|", abs2)
			return mod(a, L, R)
		} else {
			mod := array2("|", abs2)
			return mod(a, L, R)
		}
	}

	ls := al.Shape()
	rs := []int{}
	if rok {
		rs = ar.Shape()
	}

	// The shape of the result is the catenation of the shapes of L and R.
	shape := make([]int, len(ls)+len(rs))
	copy(shape[:len(ls)], ls)
	copy(shape[len(ls):], rs)
	res := apl.GeneralArray{Dims: shape}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	// enc represents r in the given radix power vector and sets the result to vec.
	fdiv := arith2("/", div2)
	flor := arith1("⌊", min)
	fmul := arith2("×", mul2)
	fsub := arith2("-", sub2)
	mod := arith2("|", abs2)
	eq := arith2("=", compare("="))
	enc := func(rad []apl.Value, r apl.Value, vec []apl.Value) error {
		var p apl.Value
		for i := range rad {
			p = apl.Index(1)
			if i < len(rad)-1 {
				p = rad[i+1]
			}
			v, err := fdiv(a, r, p)
			if err != nil {
				return err
			}
			// Dont take the residue for the last value.
			if i < len(rad)-1 {
				v, err = flor(a, nil, v)
				if err != nil {
					return err
				}
			}
			vec[i] = v
			u, err := fmul(a, v, p)
			if err != nil {
				return err
			}
			r, err = fsub(a, r, u)
			if err != nil {
				return err
			}
		}
		// If L has no zero lead, numbers exceeding the representation is incomplete.
		zerold, err := eq(a, rad[0], apl.Index(0))
		if err != nil {
			return err
		}
		if zerold.(apl.Bool) == false {
			vec[0], err = mod(a, rad[0], vec[0])
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Apply the powerradix vector to r and set the result
	vec := make([]apl.Value, shape[0])
	apply := func(rad []apl.Value, r apl.Value, n, off int) error {
		if err := enc(rad, r, vec); err != nil {
			return err
		}
		for i := range vec {
			res.Values[i*n+off] = vec[i] // TODO copy?
		}
		return nil
	}

	// Powerradix recursively multiplies from the right.
	// The index 0 value is preserved to determine
	// underrepesented values.
	powerradix := func(rad []apl.Value) error {
		for i := len(rad) - 2; i > 0; i-- {
			v, err := fmul(a, rad[i], rad[i+1])
			if err != nil {
				return err
			}
			rad[i] = v
		}
		return nil
	}

	// Number of iterations over L omitting the first axis
	NL := 1
	if len(ls) > 1 {
		NL = apl.ArraySize(apl.GeneralArray{Dims: ls[1:]})
	}
	// Number of iterations over R
	NR := 0
	if rok {
		NR = apl.ArraySize(ar)
	}
	// Number of result elements divided by length of first axis
	NN := 1
	if len(shape) > 1 {
		NN = apl.ArraySize(apl.GeneralArray{Dims: shape[1:]})
	}
	rad := make([]apl.Value, shape[0])
	off := 0
	for i := 0; i < NL; i++ {
		// Build radix vec from the first axis of L
		for k := 0; k < len(rad); k++ {
			v, err := al.At(k*NL + i)
			if err != nil {
				return nil, err
			}
			rad[k] = v
		}
		powerradix(rad)

		if rok == false {
			if err := apply(rad, R, NN, off); err != nil {
				return nil, err
			}
			off++
		}
		for k := 0; k < NR; k++ {
			r, err := ar.At(k)
			if err != nil {
				return nil, err
			}
			if err := apply(rad, r, NN, off); err != nil {
				return nil, err
			}
			off++
		}
	}
	return res, nil
}
