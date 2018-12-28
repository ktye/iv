package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/operators"
)

func init() {
	register(primitive{
		symbol: "⌹",
		doc:    "matrix inverse, domino",
		Domain: Monadic(ToArray(nil)),
		fn:     domino,
	})
	register(primitive{
		symbol: "⌹",
		doc:    "matrix divide, solve linear system, domino",
		Domain: Dyadic(Split(ToArray(nil), ToArray(nil))),
		fn:     domino2,
	})
}

// Domino is the general matrix divide for any numeric types.
// This gives exact results for rational numbers.
// There are faster methods for float and complex types as external packages.
func domino(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.EmptyArray); ok {
		return apl.EmptyArray{}, nil
	}
	ar := R.(apl.Array)
	rs := ar.Shape()
	if len(rs) > 2 {
		return nil, fmt.Errorf("matrix inverse: rank cannot be > 2: %d", len(rs))
	} else if len(rs) == 0 {
		return apl.EmptyArray{}, nil
	} else if len(rs) == 1 {
		f := array1("÷", div)
		return f(a, nil, R)
	}
	if rs[0] != rs[1] {
		return nil, fmt.Errorf("matrix inverse only works for a quadratic matrix")
	}

	n := rs[0]
	I := apl.IndexArray{Dims: []int{n, n}}
	I.Ints = make([]int, n*n)
	for i := 0; i < n; i++ {
		I.Ints[i*n+i] = 1
	}
	return domino2(a, I, ar)
}

func domino2(a *apl.Apl, RHS, R apl.Value) (apl.Value, error) {
	al := RHS.(apl.Array)
	ar := R.(apl.Array)
	ls := al.Shape()
	rs := ar.Shape()

	if apl.ArraySize(ar) == 1 {
		f := array2("÷", div2)
		return f(a, al, ar)
	}

	if len(rs) != 2 {
		return nil, fmt.Errorf("matrix divide: right argument matrix must have rank 2")
	}
	if rs[0] < rs[1] {
		return nil, fmt.Errorf("matrix divide: right argument matrix has more columns than rows")
	}
	if len(ls) == 1 {
		if rsh, ok := al.(apl.Reshaper); ok {
			v := rsh.Reshape([]int{ls[0], 1})
			al = v.(apl.Array)
		} else {
			return nil, fmt.Errorf("matrix divide: cannot reshape left argument vector: %T", al)
		}
		ls = al.Shape()
	} else if len(ls) != 2 {
		return nil, fmt.Errorf("matrix divide: left argument must have rank 2: %d", len(ls))
	}
	if ls[0] != rs[0] {
		return nil, fmt.Errorf("matrix divide: left and right matrices must have the same number of rows")
	}

	// For overdetermined systems, multiply with complex conjugate.
	// We have no special QR algorithm for the general case.
	if rs[0] > rs[1] {
		conj := array1("+", add)
		h, err := conj(a, nil, ar) // TODO copy?
		if err != nil {
			return nil, err
		}
		h, err = transpose(a, nil, h)
		if err != nil {
			return nil, err
		}
		AH := h.(apl.Array)

		mmul := operators.Scalarproduct(a, apl.Primitive("+"), apl.Primitive("×"))
		AHA, err := mmul.Call(a, AH, ar)
		if err != nil {
			return nil, err
		}
		AHB, err := mmul.Call(a, AH, al)
		if err != nil {
			return nil, err
		}

		ar = AHA.(apl.Array)
		al = AHB.(apl.Array)
		rs = ar.Shape()
		ls = al.Shape()
	}

	res := apl.MixedArray{Dims: []int{rs[0], ls[1]}}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	// A is a copy of ar as a 2d slice of Values.
	// It will be overwritten by LU.
	n := rs[0]
	A := make([][]apl.Value, n)
	for i := range A {
		A[i] = make([]apl.Value, n)
		for k := range A[i] {
			v, err := ar.At(i*n + k)
			if err != nil {
				return nil, err
			}
			A[i][k] = v
		}
	}

	// LU Decomposition overwrites M and returns the permutation matrix.
	P, err := lu(a, A)
	if err != nil {
		return nil, err
	}

	// TODO: solve for each column vector of al.
	b := make([]apl.Value, n)
	x := make([]apl.Value, n)
	for k := 0; k < ls[1]; k++ {
		// Copy column k of RHS to b.
		for i := 0; i < n; i++ {
			v, err := al.At(i*ls[1] + k)
			if err != nil {
				return nil, err
			}
			b[i] = v
		}

		// Solve for x.
		luSolve(a, A, b, P, x)

		// Copy x to result array.
		for i := 0; i < n; i++ {
			res.Values[i*ls[1]+k] = x[i] // TODO copy
		}
	}
	return res, nil
}

// LU decomposition.
func lu(a *apl.Apl, A [][]apl.Value) ([]int, error) {
	fabs := arith1("|", abs)
	fmul := arith2("×", mul2)
	fdiv := arith2("÷", div2)
	fsub := arith2("-", sub2)
	fless := arith2("<", compare("<"))

	n := len(A)
	P := make([]int, n)
	for i := range P {
		P[i] = i
	}

	var max apl.Value
	for i := 0; i < n; i++ {

		// Find row max.
		max = apl.Index(0)
		imax := i
		for k := i; k < n; k++ {
			absA, err := fabs(a, nil, A[k][i])
			if err != nil {
				return nil, err
			}
			lt, err := fless(a, max, absA)
			if err != nil {
				return nil, err
			}
			if lt.(apl.Bool) == true {
				max = absA
				imax = k
			}
		}

		// We do not compare against a tolerance, but against 0.
		if isEqual(a, apl.Index(0), max) {
			return nil, fmt.Errorf("matrix is singular")
		}

		if imax != i {
			P[i], P[imax] = P[imax], P[i]
			A[i], A[imax] = A[imax], A[i]
			// P[n]++ only needed to compute the determinant.
		}

		for j := i + 1; j < n; j++ {
			v, err := fdiv(a, A[j][i], A[i][i])
			if err != nil {
				return nil, err
			}
			A[j][i] = v
			for k := i + 1; k < n; k++ {
				v, err := fmul(a, A[j][i], A[i][k])
				if err != nil {
					return nil, err
				}
				v, err = fsub(a, A[j][k], v)
				if err != nil {
					return nil, err
				}
				A[j][k] = v
			}
		}
	}
	return P, nil
}

func luSolve(a *apl.Apl, A [][]apl.Value, b []apl.Value, P []int, x []apl.Value) error {
	fmul := arith2("×", mul2)
	fdiv := arith2("÷", div2)
	fsub := arith2("-", sub2)
	n := len(A)
	for i := 0; i < n; i++ {
		x[i] = b[P[i]]
		for k := 0; k < i; k++ {
			v, err := fmul(a, A[i][k], x[k])
			if err != nil {
				return err
			}
			v, err = fsub(a, x[i], v)
			if err != nil {
				return err
			}
			x[i] = v
		}
	}
	for i := n - 1; i >= 0; i-- {
		for k := i + 1; k < n; k++ {
			v, err := fmul(a, A[i][k], x[k])
			if err != nil {
				return err
			}
			v, err = fsub(a, x[i], v)
			if err != nil {
				return err
			}
			x[i] = v
		}
		v, err := fdiv(a, x[i], A[i][i])
		if err != nil {
			return err
		}
		x[i] = v
	}
	return nil
}
