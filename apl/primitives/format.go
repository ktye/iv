package primitives

import (
	"encoding"
	"fmt"
	"reflect"

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
// If L is a number it is used as the precision (sets PP).
// If L is two numbers, it is used as width and precision (sets PP).
// If L is a string and R a Number or uniform numeric array, L is used as a format string.
// If L is ¯1, R is formatted with MarshalText, if it implements an encoding.TextMarshaller.
func format(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	textenc := false

	// With 1 or 2 integers, set temporarily set PP.
	toIdx := ToIndexArray(nil)
	if ia, ok := toIdx.To(a, L); ok {
		if idx, ok := ia.(apl.IndexArray); ok && len(idx.Ints) == 1 && idx.Ints[0] == -1 {
			textenc = true
		} else {
			save := a.PP
			defer func() {
				a.PP = save
				a.Tower.SetPP(save)
			}()
			if err := a.SetPP(L); err != nil {
				return nil, err
			}
		}
	}

	// L string, R numeric: set format numeric format.
	if s, ok := L.(apl.String); ok {
		var n apl.Number
		if num, ok := R.(apl.Number); ok {
			n = num
		} else if u, ok := R.(apl.Uniform); ok {
			z := u.Zero()
			if num, ok := z.(apl.Number); ok {
				n = num
			}
		}
		if n != nil {
			t := reflect.TypeOf(n)
			if numeric, ok := a.Tower.Numbers[t]; ok {
				save := numeric.Format
				numeric.Format = string(s)
				a.Tower.Numbers[t] = numeric
				defer func() {
					numeric.Format = save
					a.Tower.Numbers[t] = numeric
				}()
			}
		}
	}

	if m, ok := R.(encoding.TextMarshaler); ok && textenc {
		if b, err := m.MarshalText(); err != nil {
			return nil, err
		} else {
			return apl.String(b), nil
		}
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
