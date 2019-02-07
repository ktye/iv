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
//	T→Col1
// Sorting by column
//	T[⍋T→Time]
// Selecting rows
//	T[⍸T→Qty>5]
type Table struct {
	*Dict
	Rows int
}

// String formats a table using a tabwriter.
// Each value is printed using by it's String method, same as ⍕V.
func (t Table) String(a *Apl) string {
	var b bytes.Buffer
	if err := t.WriteFormatted(a, nil, &b); err != nil {
		return ""
	}
	return string(b.Bytes())
}

// Marshal formats a table in a parsable but human readable form using a tabwriter.
// It is called by ¯1⍕T and formats each value by ¯1⍕V.
func (t Table) Marshal(a *Apl) string {
	var b bytes.Buffer
	if err := t.WriteFormatted(a, Index(-1), &b); err != nil {
		return ""
	}
	return string(b.Bytes())
}

// Csv writes a table in csv format.
// If L is nil, it uses ⍕V on each value.
// If L is a dict with conforming keys, it uses the values as left arguments to format (L[Key])⍕V
// for columns of the corresponding keys.
// If L is not a dict, it is used as the left argument to format for each key.
func (t Table) Csv(a *Apl, L Value, w io.Writer) error {
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
func (t Table) WriteFormatted(a *Apl, L Value, w io.Writer) error {
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

	marshal := func(v Value) string {
		if f.marshal {
			if m, ok := v.(Marshaler); ok {
				return m.Marshal(a)
			}
		}
		return v.String(a)
	}
	setnumformat := func(v, k Value) {
		if f.fmt == nil {
			return
		}
		// Reset default values for unspecified fields.
		a.Tower.SetPP(a.PP)

		s, ok := f.fmt[k]
		if ok == false {
			return
		}
		if _, ok := v.(Number); ok {
			t := reflect.TypeOf(v)
			if num, ok := a.Tower.Numbers[reflect.TypeOf(v)]; ok {
				num.Format = s
				a.Tower.Numbers[t] = num
			}
		}
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
					vec[j] = marshal(e)
				}
				r[i] = strings.Join(vec, " ")
				if f.marshal {
					r[i] = "[" + r[i] + "]"
				}
			} else {
				r[i] = marshal(v)
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

func (t Table) newFormatter(a *Apl, L Value) *tableFormatter {
	var f tableFormatter
	if L == nil {
		return &f
	}
	f.a = a
	f.pp = a.PP[:]                        // save pp
	f.num = make(map[reflect.Type]string) // save numeric formats
	for t, num := range a.Tower.Numbers {
		f.num[t] = num.Format
	}

	pp, err := a.toPP(L)
	if err == nil && pp[0] == 0 && pp[1] == -1 {
		f.marshal = true
	} else if err == nil {
		a.SetPP(IndexArray{Dims: []int{2}, Ints: pp[:]})
	}
	d, ok := L.(Object)
	if err == nil || ok == false {
		return &f
	}

	f.fmt = make(map[Value]string)
	keys := d.Keys()
	for _, k := range keys {
		v := d.At(a, k)
		if s, ok := v.(String); ok {
			f.fmt[k] = string(s)
		}
	}
	return &f
}

type tableFormatter struct {
	a       *Apl
	pp      []int
	marshal bool
	fmt     map[Value]string
	num     map[reflect.Type]string
}

func (f *tableFormatter) Close() {
	if f.pp != nil {
		f.a.SetPP(IndexArray{Dims: []int{2}, Ints: f.pp})
	}
	if f.num != nil {
		for t, s := range f.num {
			num := f.a.Tower.Numbers[t]
			num.Format = s
			f.a.Tower.Numbers[t] = num
		}
	}
}
