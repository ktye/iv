package apl

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

// Object is a compound type that has fields.
//
// Fields are accessed by
//	X←Object→Name
// and set by assigning to them
//	Object→Name←X
//
// Dict is the default object implementation.
// It is created, if the object does not exist.
// Another implementation is xgo.Value.
type Object interface {
	Field(*Apl, string) Value
	Set(*Apl, string, Value) error
}

// TODO
//	- delete keys
//	- set key order
//		e.g. by indexing: D←D[2 1 3]
type Dict struct {
	keys []string
	m    map[string]Value
}

// Field returns the value for the key.
// It returns nil, if the key does not exist.
func (d *Dict) Field(a *Apl, key string) Value {
	if d.m == nil {
		return nil
	}
	return d.m[key]
}

// Set updates the value for the given key, or creates a new one,
// if the key does not exist.
// Keys must be valid variable names.
func (d *Dict) Set(a *Apl, key string, v Value) error {
	ok, isfunc := isVarname(key)
	if ok == false {
		return fmt.Errorf("not a valid key name: %s", key)
	}
	if _, ok := v.(Function); ok && isfunc == false {
		return fmt.Errorf("function values can only be stored keys starting with a lowercase letter")
	} else if ok == false && isfunc == true {
		return fmt.Errorf("arrays can only be stored in keys starting with a capital letter")
	}
	if d.m == nil {
		d.m = make(map[string]Value)
	}
	if _, ok := d.m[key]; ok == false {
		d.keys = append(d.keys, key)
	}
	d.m[key] = v
	return nil
}

func (d *Dict) String(a *Apl) string {
	var buf strings.Builder
	tw := tabwriter.NewWriter(&buf, 1, 0, 1, ' ', 0)
	for _, k := range d.keys {
		fmt.Fprintf(tw, "%s:\t%s\n", k, d.m[k].String(a))
	}
	tw.Flush()
	s := buf.String()
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}
