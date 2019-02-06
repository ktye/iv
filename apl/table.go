package apl

import (
	"fmt"
	"io"
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

func (t Table) String(a *Apl) string {
	var buf strings.Builder
	tw := tabwriter.NewWriter(&buf, 1, 0, 1, ' ', 0)
	keys := t.Keys()
	if len(keys) == 0 {
		return ""
	}

	for i, k := range keys {
		sep := "\t"
		if i == len(keys)-1 {
			sep = "\n"
		}
		fmt.Fprintf(tw, "%s%s", k.String(a), sep)
	}
	for n := 0; n < t.Rows; n++ {
		for i, k := range keys {
			sep := "\t"
			if i == len(keys)-1 {
				sep = "\n"
			}
			col := t.At(a, k)
			if col == nil {
				return "???"
			}
			fmt.Fprintf(tw, "%s%s", col.(Array).At(n).String(a), sep)
		}
	}
	tw.Flush()
	s := buf.String()
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}

func (t Table) Marshal(a *Apl) string {
	panic("TODO: table marshal text")
}

// Csv writes the table to w in csv format.
// L may be a dict with conforming keys of formatting values.
func (t Table) Csv(a *Apl, L Value, w io.Writer) error {
	return fmt.Errorf("TODO: table csv")
}

func (t Table) WriteFormatted(a *Apl, L Object, w io.Writer) error {
	return fmt.Errorf("TODO: table write formatted")
}
