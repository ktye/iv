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
	p := apl.MixedArray{
		Dims:   apl.CopyShape(al),
		Values: make([]apl.Value, al.Size()),
	}
	for i := range p.Values {
		p.Values[i] = al.At(i).Copy()
	}
	N := ls[len(ls)-1]

	// Shift last axis by 1 to the left than multiply scan from the right.
	fmul := arith2("×", mul2)
	for off := 0; off < len(p.Values); off += N {
		for k := 0; k < N-1; k++ {
			p.Values[off+k] = p.Values[off+k+1]
		}
		p.Values[off+N-1] = apl.Int(1)
		for k := off + N - 1; k >= off+1; k-- {
			v, err := fmul(a, p.Values[k], p.Values[k-1])
			if err != nil {
				return nil, err
			}
			p.Values[k-1] = v
		}
	}

	dot := operators.Scalarproduct(a, apl.Primitive("+"), apl.Primitive("×"))
	return dot.Call(a, a.UnifyArray(p), ar)
}

// extendAxis extends the axis of length 1 to n
func extendAxis(ar apl.Array, axis, n int) (apl.Array, []int) {
	res := apl.MixedArray{Dims: apl.CopyShape(ar)}
	res.Dims[axis] = n
	res.Values = make([]apl.Value, apl.Prod(res.Dims))
	ridx := make([]int, len(res.Dims))
	ic, idx := apl.NewIdxConverter(ar.Shape())
	for i := range res.Values {
		copy(idx, ridx)
		idx[axis] = 0
		res.Values[i] = ar.At(ic.Index(idx)).Copy()
		apl.IncArrayIndex(ridx, res.Dims)
	}
	return res, res.Dims
}

// ISO p.151
func encode(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	al, lok := L.(apl.Array)
	ar, rok := R.(apl.Array)

	// If L or R is empty, return (⍴L,⍴R)⍴0.
	_, ae := L.(apl.EmptyArray)
	_, re := R.(apl.EmptyArray)
	if ae || re {
		var shape []int
		if lok {
			if ae {
				shape = []int{0}
			} else {
				shape = append(shape, al.Shape()...)
			}
		}
		if rok {
			if re {
				shape = append(shape, 0)
			} else {
				shape = append(shape, ar.Shape()...)
			}
		}
		sv := apl.IntArray{Dims: []int{len(shape)}, Ints: shape}
		zv := apl.IntArray{Dims: []int{1}, Ints: []int{0}}
		return rho2(a, sv, zv)
	}

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

	// L is a vector and R is a scalar.
	ls := al.Shape()
	if len(ls) == 1 && rok == false {
		ravelL := make([]apl.Value, ls[0])
		for i := range ravelL {
			ravelL[i] = al.At(i).Copy()
		}
		return encodeVecScalar(a, ravelL, R)
	}

	return encodeArray(a, al, R)
}

// encodeVecScalar returns L⊤R for vector L and scalar R.
func encodeVecScalar(a *apl.Apl, L []apl.Value, R apl.Value) (apl.Value, error) {
	eq := arith2("=", compare("="))
	fsub := arith2("-", sub2)
	fdiv := arith2("÷", div2)
	mod := arith2("|", abs2)

	// Two vectors Z (len ⍴A) and C (len 1+⍴A)
	Z := make([]apl.Value, len(L))
	C := make([]apl.Value, len(L)+1)
	C[len(C)-1] = R
	for i := len(Z) - 1; i >= 0; i-- {
		// Z[i] ← L[i] ⊤ C[i+1]
		v, err := mod(a, L[i], C[i+1])
		if err != nil {
			return nil, err
		}
		Z[i] = v.Copy()

		// If L[i] is 0: C[i] ← 0
		a0, err := eq(a, L[i], apl.Int(0))
		if err != nil {
			return nil, err
		}
		if iszero := a0.(apl.Bool); iszero == true {
			C[i] = apl.Int(0)
		} else {
			// Otherwise: C[i] ← (C[i+1] - Z[i])÷A[i]
			d, err := fsub(a, C[i+1], Z[i])
			if err != nil {
				return nil, err
			}
			d, err = fdiv(a, d, L[i])
			if err != nil {
				return nil, err
			}
			C[i] = d.Copy()
		}
	}
	return a.UnifyArray(apl.MixedArray{Dims: []int{len(Z)}, Values: Z}), nil
}

func encodeArray(a *apl.Apl, al apl.Array, R apl.Value) (apl.Value, error) {
	ls := al.Shape()
	rs := []int{}
	ar, rok := R.(apl.Array)
	if rok {
		rs = ar.Shape()
	}

	fdiv := arith2("/", div2)
	flor := arith1("⌊", min)
	fmul := arith2("×", mul2)
	fsub := arith2("-", sub2)
	mod := arith2("|", abs2)
	eq := arith2("=", compare("="))

	// The shape of the result is the catenation of the shapes of L and R.
	shape := make([]int, len(ls)+len(rs))
	copy(shape[:len(ls)], ls)
	copy(shape[len(ls):], rs)
	res := apl.MixedArray{Dims: shape}
	res.Values = make([]apl.Value, apl.Prod(res.Dims))

	// enc represents r in the given radix power vector and sets the result to vec.
	enc := func(rad []apl.Value, r apl.Value, vec []apl.Value) error {
		var p apl.Value
		for i := range rad {
			p = apl.Int(1)
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
			vec[i] = v.Copy()
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
		zerold, err := eq(a, rad[0], apl.Int(0))
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
			res.Values[i*n+off] = vec[i].Copy()
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
			rad[i] = v.Copy()
		}
		return nil
	}

	// Number of iterations over L omitting the first axis
	NL := 1
	if len(ls) > 1 {
		NL = apl.Prod(ls[1:])
	}
	// Number of iterations over R
	NR := 0
	if rok {
		NR = ar.Size()
	}
	// Number of result elements divided by length of first axis
	NN := 1
	if len(shape) > 1 {
		NN = apl.Prod(shape[1:])
	}
	rad := make([]apl.Value, shape[0])
	off := 0
	for i := 0; i < NL; i++ {
		// Build radix vec from the first axis of L
		for k := 0; k < len(rad); k++ {
			rad[k] = al.At(k*NL + i).Copy()
		}
		powerradix(rad)

		if rok == false {
			if err := apply(rad, R, NN, off); err != nil {
				return nil, err
			}
			off++
		}
		for k := 0; k < NR; k++ {
			if err := apply(rad, ar.At(k), NN, off); err != nil {
				return nil, err
			}
			off++
		}
	}
	return a.UnifyArray(res), nil
}

/* There should be a simpler algorithm for encodeArray:
// ISO p.151
// The shape of the result is ⍴L,⍴R
shape := apl.CopyShape(al)
if rok {
	shape = append(shape, apl.CopyShape(ar)...)
}
res := apl.GeneralArray{Dims: shape}
res.Values = make([]apl.Value, apl.ArraySize(res))

// Ravel list of R
var r []apl.Value
if rok == false {
	r = []apl.Value{R}
} else {
	r = make([]apl.Value, apl.ArraySize(ar))
	for i := 0; i < len(r); i++ {
		r[i], _ = ar.At(i)
	}
}
N := len(r)

// Ravel-along-axis one (see ISO p.23)
// A1: Vector item i of ravel-along-axis one of L
A1 := make([]apl.Value, apl.ArraySize(al)/ls[0])
var err error
for i := 0; i < ls[0]; i++ {
	off := i * len(A1)
	for k := range A1 {
		A1[k], err = al.At(off + k)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("A1", A1)
	for j, B1 := range r {
		P := j + N*i
		v, err := encodeVecScalar(a, A1, B1)
		if err != nil {
			return nil, err
		}
		res.Values[P] = v
	}
}
return res, nil
*/
