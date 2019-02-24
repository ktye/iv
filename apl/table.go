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
func (t Table) String(a *Apl) string {
	if a.PP == -2 || a.PP == -3 {
		return t.Dict.String(a)
	}
	var b bytes.Buffer
	if err := t.WriteFormatted(a, nil, &b); err != nil {
		return ""
	}
	return string(b.Bytes())
}

// Csv writes a table in csv format.
// If L is nil, it uses ⍕V on each value.
// If L is a dict with conforming keys, it uses the values as left arguments to format (L[Key])⍕V
// for columns of the corresponding keys.
func (t Table) Csv(a *Apl, L Object, w io.Writer) error {
	f := t.newFormatter(a, L)
	defer f.Close()

	cw := csv.NewWriter(w)
	if err := t.write(a, f, csvTable{cw}); err != nil {
		return err
	}
	cw.Flush()
	return nil
}

// WriteFormatted writes the table with a tablwriter.
// The format of the values is given by L in the same way as for Csv.
func (t Table) WriteFormatted(a *Apl, L Object, w io.Writer) error {
	f := t.newFormatter(a, L)
	defer f.Close()

	tw := tabwriter.NewWriter(w, 1, 0, 1, ' ', 0)
	if err := t.write(a, f, wsTable{tw}); err != nil {
		return err
	}
	return tw.Flush()
}

func (t Table) write(a *Apl, f *tableFormatter, rw rowWriter) error {
	keys := t.Keys()
	if len(keys) == 0 {
		return nil
	}

	r := make([]string, len(keys))
	for i := range r {
		r[i] = keys[i].String(a) // The header is always formatted with String.
	}
	if err := rw.writeRow(r); err != nil {
		return err
	}

	setnumformat := func(v, k Value) {
		if f.fmt == nil {
			return
		}
		t := reflect.TypeOf(v)
		s, ok := f.fmt[k]
		if ok == false {
			s = f.rst[t]
		}
		a.Fmt[t] = s
	}
	for n := 0; n < t.Rows; n++ {
		for i, k := range keys {
			v := t.At(a, k).(Array).At(n)
			setnumformat(v, k)
			if ar, ok := v.(Array); ok {
				size := ar.Size()
				vec := make([]string, size)
				for j := 0; j < size; j++ {
					e := ar.At(j)
					setnumformat(e, k)
					vec[j] = e.String(a)
				}
				r[i] = strings.Join(vec, " ")
				if a.PP < 0 {
					r[i] = "[" + r[i] + "]"
				}
			} else {
				r[i] = v.String(a)
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

func (t Table) newFormatter(a *Apl, L Object) *tableFormatter {
	var f tableFormatter
	if L == nil {
		return &f
	}
	f.a = a
	f.rst = make(map[reflect.Type]string)
	for t, s := range a.Fmt {
		f.rst[t] = s
	}

	f.fmt = make(map[Value]string)
	keys := L.Keys()
	for _, k := range keys {
		v := L.At(a, k)
		if s, ok := v.(String); ok {
			f.fmt[k] = string(s)
		}
	}
	return &f
}

type tableFormatter struct {
	a   *Apl
	rst map[reflect.Type]string
	fmt map[Value]string
}

func (f *tableFormatter) Close() {
	if f.a != nil {
		f.a.Fmt = f.rst
	}
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
