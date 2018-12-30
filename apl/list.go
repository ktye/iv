package apl

import (
	"fmt"
	"strings"
)

// List is a collection of items, possibly nested.
// It also acts as a vector (a rank 1 array) but cannot be reshaped.
type List []Value

func (l List) String(a *Apl) string {
	var buf strings.Builder
	buf.WriteRune('(')
	for i := range l {
		buf.WriteString(l[i].String(a))
		buf.WriteRune(';')
	}
	buf.WriteRune(')')
	return buf.String()
}

func (l List) At(i int) (Value, error) {
	if i < 0 || i >= len(l) {
		return nil, fmt.Errorf("index out of range")
	}
	return l[i], nil
}

func (l List) Shape() []int {
	return []int{len(l)}
}

func (l List) Size() int {
	return len(l)
}

type list []expr

func (l list) Eval(a *Apl) (Value, error) {
	lst := make(List, len(l))
	var err error
	for i := range lst {
		lst[i], err = l[i].Eval(a)
		if err != nil {
			return nil, err
		}
	}
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
