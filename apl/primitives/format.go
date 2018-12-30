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
	// TODO: dyadic ⍕: format with specification.

	register(primitive{
		symbol: "⍎",
		doc:    "execute, evaluate expression",
		Domain: Monadic(IsString(nil)),
		fn:     execute,
	})
	// TODO: dyadic ⍎: execute with namespace.
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
