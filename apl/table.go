package apl

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

// A Table is a transposed dictionary where each value is a vector
// with the same number of elements and unique type.
// Tables are constructed by transposing dictionaries T←⍉D
//
// Indexing tables selects rows:
//	T[⍳5]
// returns a table with the first 5 rows.
// Right arrow indexing selects columns, just like a dict.
//	T→Col1 // TODO update
// Sorting by column
//	T[⍋T→Time] // TODO update
// Selecting rows
//	T[⍸T→Qty>5] // TODO update
type Table struct {
	*Dict
	Rows int
}

// String formats a table using a tabwriter.
// Each value is printed using by it's String method, same as ⍕V.
func (t Table) String(f Format) string {
	if f.PP == -2 || f.PP == -3 {
		return t.Dict.String(f)
	}
	var b bytes.Buffer
	if err := t.WriteFormatted(f, nil, &b); err != nil {
		return ""
	}
	return string(b.Bytes())
}
func (t Table) Copy() Value {
	r := Table{Rows: t.Rows}
	if t.Dict != nil {
		r.Dict = t.Dict.Copy().(*Dict)
	}
	return r
}

// Csv writes a table in csv format.
// If L is nil, it uses ⍕V on each value.
// If L is a dict with conforming keys, it uses the values as left arguments to format (L[Key])⍕V
// for columns of the corresponding keys.
func (t Table) Csv(f Format, L Object, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := t.write(f, L, csvTable{cw}); err != nil {
		return err
	}
	cw.Flush()
	return nil
}

// WriteFormatted writes the table with a tablwriter.
// The format of the values is given by L in the same way as for Csv.
func (t Table) WriteFormatted(f Format, L Object, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 1, 0, 1, ' ', 0)
	if err := t.write(f, L, wsTable{tw}); err != nil {
		return err
	}
	return tw.Flush()
}

func (t Table) write(af Format, L Object, rw rowWriter) error {
	keys := t.Keys()
	if len(keys) == 0 {
		return nil
	}

	var colfmt map[Value]string
	if L != nil {
		colfmt = make(map[Value]string)
		for _, k := range L.Keys() {
			v := L.At(k)
			if s, ok := v.(String); ok {
				colfmt[k.Copy()] = string(s)
			}
		}
	}

	f := Format{
		PP:  af.PP,
		Fmt: make(map[reflect.Type]string),
	}
	for k, v := range af.Fmt {
		f.Fmt[k] = v
	}

	r := make([]string, len(keys))
	for i := range r {
		r[i] = keys[i].String(f) // The header is always formatted with String.
	}
	if err := rw.writeRow(r); err != nil {
		return err
	}

	for n := 0; n < t.Rows; n++ {
		for i, k := range keys {
			custom := ""
			if colfmt != nil {
				if s, ok := colfmt[k]; ok {
					custom = s
				}
			}
			v := t.At(k).(Array).At(n)
			t := reflect.TypeOf(v)
			if custom == "" {
				delete(f.Fmt, t)
			} else {
				f.Fmt[t] = custom
			}
			if ar, ok := v.(Array); ok {
				size := ar.Size()
				vec := make([]string, size)
				for j := 0; j < size; j++ {
					e := ar.At(j)
					vec[j] = e.String(f)
				}
				r[i] = strings.Join(vec, " ")
				if f.PP < 0 {
					r[i] = "[" + r[i] + "]"
				}
			} else {
				r[i] = v.String(f)
			}
		}
		if err := rw.writeRow(r); err != nil {
			return err
		}
	}
	return nil
}

type rowWriter interface {
	writeRow([]string) error
}

type csvTable struct {
	*csv.Writer
}

func (c csvTable) writeRow(records []string) error { return c.Writer.Write(records) }

type wsTable struct {
	*tabwriter.Writer
}

func (w wsTable) writeRow(records []string) error {
	_, err := fmt.Fprintln(w, strings.Join(records, "\t"))
	return err
}

func (a *Apl) ParseTable(prototype Value, s string) (Table, error) {
	if prototype != nil {
		_, ok := prototype.(Table)
		if ok == false {
			return Table{}, fmt.Errorf("ParseTable: prototype is not a table: %T", prototype)
		}
	}
	return Table{}, fmt.Errorf("TODO ParseTable")
}
