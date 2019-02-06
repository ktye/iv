package primitives

import (
	"bytes"
	"encoding/csv"
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
	register(primitive{
		symbol: "⍕",
		doc:    "format, convert to string",
		Domain: Dyadic(Split(IsObject(nil), IsTable(nil))),
		fn:     formatTable,
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
// If L is the string "csv", csv encoding is used.
// If L is a string and R a Number or uniform numeric array, L is used as a format string.
// If L is ¯1, R is formatted with Marshal, if it implements an Marshaler.
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
		if s == "csv" {
			return formatCsv(a, nil, R)
		}
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

	if m, ok := R.(apl.Marshaler); ok && textenc {
		return apl.String(m.Marshal(a)), nil
	}

	return apl.String(R.String(a)), nil
}

// L is an object and R a Table.
// Corresponding values of L are used as format arguments to values in R.
// If L contains the key CSV, formatCSV is used.
func formatTable(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	t := R.(apl.Table)
	d := L.(apl.Object)
	if d.At(a, apl.String("CSV")) != nil {
		return formatCsv(a, L, R)
	}
	var b bytes.Buffer
	if err := t.WriteFormatted(a, d, &b); err != nil {
		return nil, err
	}
	return apl.String(b.Bytes()), nil
}

// formatCSV formats R in csv format.
// R must be a rank 2 array or a table.
// If L with corresponding keys.
func formatCsv(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	var b bytes.Buffer
	w := csv.NewWriter(&b)

	ar, ok := R.(apl.Array)
	if ok {
		shape := ar.Shape()
		if len(shape) != 2 {
			return nil, fmt.Errorf("format csv: R must be rank 2: shape is %v", shape)
		}
		if shape[0] == 0 || shape[1] == 0 {
			return apl.String(""), nil
		}
		records := make([]string, shape[1])
		idx := 0
		for i := 0; i < shape[0]; i++ {
			for k := 0; k < shape[1]; k++ {
				records[k] = ar.At(idx).String(a)
				idx++
			}
			if err := w.Write(records); err != nil {
				return nil, fmt.Errorf("format csv: %s", err)
			}
		}
		w.Flush()
		return apl.String(b.Bytes()), nil
	} else if t, ok := R.(apl.Table); ok {
		var b bytes.Buffer
		if err := t.Csv(a, L, &b); err != nil {
			return nil, err
		}
		return apl.String(b.Bytes()), nil
	}
	return nil, fmt.Errorf("format csv: unexpected type: %T", R)
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
