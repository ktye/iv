package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⍕",
		doc:    "format, convert to string",
		Domain: Monadic(nil),
		fn: func(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
			return apl.String(R.String(a)), nil
		},
	})
	register(primitive{
		symbol: "⍕",
		doc:    "format, convert to string",
		Domain: Dyadic(nil),
		fn:     format,
	})
	// TODO: dyadic ⍕: format with specification.

	register(primitive{
		symbol: "⍎",
		doc:    "execute, evaluate expression",
		Domain: Monadic(IsString(nil)),
		fn:     execute,
	})
	// TODO: dyadic ⍎: execute with namespace.
}

// Format converts the argument to string.
// If L is a number it is used as the precision.
// If L is two numbers, it is used as width and precision.
func format(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	// With 1 or 2 integers, set temporarily set PP.
	toIdx := ToIndexArray(nil)
	ia, ok := toIdx.To(a, L)
	if _, isempty := ia.(apl.EmptyArray); isempty {
		return nil, fmt.Errorf("format: left argument is empty")
	}
	if ok {
		idx := ia.(apl.IndexArray).Ints
		var pp [2]int
		if len(idx) == 1 {
			pp[1] = int(idx[0])
		} else if len(idx) == 2 {
			pp[0] = int(idx[0])
			pp[1] = int(idx[1])
		} else {
			return nil, fmt.Errorf("format: left argument has wrong number of integers")
		}
		a.Tower.SetPP(pp)
		defer func() {
			a.Tower.SetPP(pp)
		}()
	}
	return apl.String(R.String(a)), nil
}

// Execute evaluates the string in R.
// If it evaluates to multiple values, return the last but display all.
func execute(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	s := R.(apl.String)
	p, err := a.Parse(string(s))
	if err != nil {
		return nil, err
	}
	values, err := a.EvalProgram(p)
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return apl.EmptyArray{}, nil // Does this ever happen?
	}
	for _, v := range values[:len(values)-1] {
		// TODO: do not display shy values.
		fmt.Fprintln(a.GetOutput(), v.String(a))
	}
	return values[len(values)-1], nil
}
