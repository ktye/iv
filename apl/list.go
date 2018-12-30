package apl

import (
	"fmt"
	"strings"
)

// List is a collection of items, possibly nested.
// It also acts as a vector (a rank 1 array) but cannot be reshaped.
// The vector view sees the list as a flat list with all elements.
type List struct {
	N int     // total count
	L []Value // children
}

func (l List) String(a *Apl) string {
	var buf strings.Builder
	buf.WriteRune('(')
	for i := range l.L {
		buf.WriteString(l.L[i].String(a))
		buf.WriteRune(';')
	}
	buf.WriteRune(')')
	return buf.String()
}

func (l List) At(i int) (Value, error) {
	if i < 0 || i >= len(l.L) {
		return nil, fmt.Errorf("index out of range")
	}
	return l.L[i], nil
}

func (l List) Shape() []int {
	return []int{l.N}
}

func (l List) Size() int {
	return l.N
}

type list []expr

func (l list) Eval(a *Apl) (Value, error) {
	var lst List
	L := make([]Value, len(l))
	var err error
	for i := range L {
		L[i], err = l[i].Eval(a)
		if err != nil {
			return nil, err
		}
		if v, ok := L[i].(List); ok {
			lst.N += v.N
		} else {
			lst.N += 1
		}
	}
	lst.L = L
	return lst, nil
}

func (l list) String(a *Apl) string {
	var buf strings.Builder
	buf.WriteRune('(')
	for i := range l {
		buf.WriteString(l[i].String(a))
		buf.WriteRune(';')
	}
	buf.WriteRune(')')
	return buf.String()
}
