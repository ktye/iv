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
			return apl.String(R.String(a.Format)), nil
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

	register(primitive{
		symbol: "⍎",
		doc:    "execute, evaluate expression",
		Domain: Monadic(IsString(nil)),
		fn:     execute,
	})
	register(primitive{
		symbol: "⍎",
		doc:    "parse data",
		Domain: Dyadic(Split(nil, IsString(nil))),
		fn:     parseData,
	})
}

// Format converts the argument to string.
// If L is a number it is used as the precision (sets PP).
// If L is a string L is used as a format string.
// Special formatting is used, if the string is "csv", "json", "mat" or "x".
func format(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	f := apl.Format{
		PP:  a.Format.PP,
		Fmt: make(map[reflect.Type]string),
	}
	for k, v := range a.Format.Fmt {
		f.Fmt[k] = v
	}

	if n, ok := L.(apl.Number); ok {
		if i, ok := n.ToIndex(); ok {
			f.PP = i
		}
	} else if s, ok := L.(apl.String); ok {
		switch s {
		case "csv":
			return formatCsv(f, nil, R)
		case "json":
			f.PP = -2
		case "mat":
			f.PP = -3
		case "x":
			f.PP = -16
		default:
			t := reflect.TypeOf(R)
			f.Fmt[t] = string(s)
		}
	}
	return apl.String(R.String(f)), nil
}

// L is an object and R a Table.
// Corresponding values of L are used as format arguments to values in R.
// If L contains the key CSV, formatCSV is used.
func formatTable(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	t := R.(apl.Table)
	d := L.(apl.Object)
	if d.At(apl.String("CSV")) != nil {
		return formatCsv(a.Format, d, R)
	}
	var b bytes.Buffer
	if err := t.WriteFormatted(a.Format, d, &b); err != nil {
		return nil, err
	}
	return apl.String(b.Bytes()), nil
}

// formatCSV formats R in csv format.
// R must be a rank 2 array or a table.
// If L with corresponding keys.
func formatCsv(f apl.Format, L apl.Object, R apl.Value) (apl.Value, error) {
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
				records[k] = ar.At(idx).String(f)
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
		if err := t.Csv(f, L, &b); err != nil {
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
		fmt.Fprintln(a.GetOutput(), v.String(a.Format))
	}
	return values[len(values)-1], nil
}

// ParseData parses data from strings that has been written with ¯1⍕V.
// L may be "A", "D" or "T" for array, dict or table.
// If L is a value of type array, dict or table it is used as a prototype with stricter requirements.
func parseData(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	var p apl.Value
	ls, ok := L.(apl.String)
	if ok == false {
		p = L
		if _, ok := L.(apl.Array); ok {
			ls = "A"
		} else if _, ok := L.(*apl.Dict); ok {
			ls = "D"
		} else if _, ok := L.(apl.Table); ok {
			ls = "T"
		} else {
			return nil, fmt.Errorf("parse data: left argument is an unknown prototype %T", L)
		}
	}
	rs := R.(apl.String)
	switch ls {
	case "A":
		return a.ParseArray(p, string(rs))
	case "D":
		return a.ParseDict(p, string(rs))
	case "T":
		return a.ParseTable(p, string(rs))
	}
	return nil, fmt.Errorf("parse data: left argument is an unknown type: %s", ls)
}
