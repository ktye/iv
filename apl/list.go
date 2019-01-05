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

func (l List) At(i int) Value {
	return l[i]
}

func (l List) Shape() []int {
	return []int{len(l)}
}

func (l List) Size() int {
	return len(l)
}

func (l List) GetDeep(idx []int) (Value, error) {
	return l.getset(idx, nil)
}

func (l List) SetDeep(idx []int, v Value) error {
	_, err := l.getset(idx, v)
	return err
}

func (l List) getset(idx []int, v Value) (Value, error) {
	if len(idx) == 0 {
		return nil, fmt.Errorf("empty index")
	}
	for i, k := range idx {
		if k < 0 || k >= len(l) {
			return nil, fmt.Errorf("index out of range")
		}
		if i == len(idx)-1 {
			if v != nil {
				l[k] = v
				return nil, nil
			} else {
				return l[k], nil
			}
		}
		if lst, ok := l[k].(List); ok == false {
			return nil, fmt.Errorf("index is too deep")
		} else {
			l = lst
		}
	}
	return nil, fmt.Errorf("not reached")
}

type list []expr

func (l list) Eval(a *Apl) (Value, error) {
	lst := make(List, len(l))
	var err error
	for i := range lst {
		if l[i] == nil {
			lst[i] = EmptyArray{}
			continue
		}
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
