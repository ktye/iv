package apl

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

// Object is a compound type that has keys and values.
//
// Values are accessed by indexing with keys.
//	Object[Key]
// To set a key, use indexed assignment:
//	Object[Name]←X
// This also works for vectors
//	Object[`k1`k2`k3] ← 5 6 7
//
// Keys are returned by #Object.
// Number of keys can also be obtained by ⍴Object.
//
// Indexing by vector returns a Dict with the specified keys.
//	Object["key1" "key2"].
//
// Method calls (calling a function stored in a key) or a go method
// for an xgo object cannot be applied directly:
//	Object[`f] R  ⍝ cannot be parsed
// Instead, assign it to a function variable, or commute:
//	f←Object[`f] ⋄ f R
//      Object[`f]⍨R
type Object interface {
	Keys() []Value
	At(*Apl, Value) Value
	Set(*Apl, Value, Value) error
}

// Dict is a dictionary object.
type Dict struct {
	K []Value
	M map[Value]Value
}

func (d *Dict) Keys() []Value {
	return d.K
}

func (d *Dict) At(a *Apl, key Value) Value {
	if d.M == nil {
		return nil
	}
	return d.M[key]
}

// Set updates the value for the given key, or creates a new one,
// if the key does not exist.
// Keys must be valid variable names.
func (d *Dict) Set(a *Apl, key Value, v Value) error {
	if d.M == nil {
		d.M = make(map[Value]Value)
	}
	if _, ok := d.M[key]; ok == false {
		d.K = append(d.K, key)
	}
	d.M[key] = v
	return nil
}

func (d *Dict) String(a *Apl) string {
	var buf strings.Builder
	tw := tabwriter.NewWriter(&buf, 1, 0, 1, ' ', 0)
	for _, k := range d.K {
		fmt.Fprintf(tw, "%s:\t%s\n", k.String(a), d.M[k].String(a))
	}
	tw.Flush()
	s := buf.String()
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}
