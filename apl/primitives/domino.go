package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
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

	/*
		L, U, P, err := lu(ar)
		if err != nil {
			return nil, err
		}

		luInverse()
	*/

	return nil, fmt.Errorf("TODO")
}

func domino2(a *apl.Apl, RHS, A apl.Value) (apl.Value, error) {
	al := RHS.(apl.Array)
	ar := A.(apl.Array)
	ls := al.Shape()
	rs := ar.Shape()

	if apl.ArraySize(ar) == 1 {
		f := array2("÷", div2)
		return f(a, al, ar)
	}

	if len(rs) != 2 || rs[0] != rs[1] {
		return nil, fmt.Errorf("matrix divide: only a quadratic right argument matrix is supported")
	}
	if len(ls) == 1 {
		if rsh, ok := al.(apl.Reshaper); ok {
			v := rsh.Reshape([]int{1, ls[0]})
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

	res := apl.GeneralArray{Dims: []int{rs[0], ls[1]}}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	/* TODO
	M := apl.GeneralArray{} // Copy ar
	P, err := lu(M)
	if err != nil {
		return nil, err
	}
	*/

	// TODO: solve for each column vector of al.

	return nil, fmt.Errorf("TODO")
}

func lu(A apl.GeneralArray) (P apl.IndexArray, err error) {
	return apl.IndexArray{}, fmt.Errorf("TODO: LU decomposition")
}

// TODO
// luSolve
// luInverse
