package funcs

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

func init() {
	register("=", compare("="))
	register("<", compare("<"))
	register(">", compare(">"))
	register("≠", compare("≠"))
	register("≤", compare("≤"))
	register("≥", compare("≥"))
	addDoc("=", `= primitive function: compare equality
Z←L=R
`)
	addDoc("<", `< primitive function: compare less than
Z←L<R
`)
	addDoc(">", `> primitive function: compare greater than
Z←L>R
`)
	addDoc("≠", `≠ primitive function: compare not equal
Z←L≠R
`)
	addDoc("≤", `≤ primitive function: compare less or equal
Z←L≤R
`)
	addDoc("≥", `≥ primitive function: compare greater or equal
Z←L≥R
`)
}

// Compare returns a dyadic function handle for the given comparison symbol.
func compare(s string) apl.FunctionHandle {
	cmp := func(a, b apl.Value) (bool, error) {
		eq, lt, err := apl.CompareScalars(a, b)
		// We treat comparison of NaN as an error in any case.
		if err == apl.ErrCmpCmplx {
			if s == "=" {
				return eq, nil
			} else if s == "≠" {
				return !eq, nil
			} else {
				return false, err
			}
		} else if err != nil {
			return false, err
		}
		switch s {
		case "=":
			return eq, nil
		case "<":
			return lt, nil
		case ">":
			return !eq && !lt, nil
		case "≠":
			return !eq, nil
		case "≤":
			return eq || lt, nil
		case "≥":
			return !lt, nil
		default:
			return false, fmt.Errorf("illegal comparision operator: %s", s)
		}

	}
	return func(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
		if l == nil {
			return false, nil, nil // compare cannot be used in monadic context.
		}
		if apl.IsScalar(l) && apl.IsScalar(r) {
			b, err := cmp(l, r)
			return true, apl.Bool(b), err
		}
		ar, isa := l.(apl.Array)
		br, isb := r.(apl.Array)
		if isa && isb {
			as := ar.Shape()
			bs := br.Shape()
			if len(as) != len(bs) {
				return true, nil, fmt.Errorf("cannot compare arrays of different size")
			} else {
				for i := range as {
					if as[i] != bs[i] {
						return true, nil, fmt.Errorf("cannot compare arrays of different size")
					}
				}
			}
		} else if isa {
			shape := make([]int, len(ar.Shape()))
			copy(shape, ar.Shape())
			values := make([]apl.Value, apl.ArraySize(ar))
			for i := range values {
				values[i] = r
			}
			x := apl.GeneralArray{
				Dims:   shape,
				Values: values,
			}
			br = x
		} else if isb {
			shape := make([]int, len(br.Shape()))
			copy(shape, br.Shape())
			values := make([]apl.Value, apl.ArraySize(br))
			for i := range values {
				values[i] = l
			}
			x := apl.GeneralArray{
				Dims:   shape,
				Values: values,
			}
			ar = x
		} else {
			return false, nil, nil
		}

		// Both ar and br are arryas of the same shape.
		shape := make([]int, len(ar.Shape()))
		copy(shape, ar.Shape())
		res := apl.Bitarray{
			Bits: make([]apl.Bool, apl.ArraySize(ar)),
			Dims: shape,
		}
		for i := range res.Bits {
			av, err := ar.At(i)
			if err != nil {
				return true, nil, err
			}
			bv, err := br.At(i)
			if err != nil {
				return true, nil, err
			}
			if b, err := cmp(av, bv); err != nil {
				return true, nil, err
			} else {
				res.Bits[i] = apl.Bool(b)
			}
		}
		return true, res, nil
	}
}
