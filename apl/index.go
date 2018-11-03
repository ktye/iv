package apl

import "strings"

// IdxSpec is the value of an index specification.
type IdxSpec []Value

func (x IdxSpec) String(a *Apl) string {
	if x == nil {
		return "[nil]"
	}
	v := make([]string, len(x))
	for i, e := range x {
		v[i] = e.String(a)
	}
	return "[" + strings.Join(v, ";") + "]"
}

func (x IdxSpec) Eval(a *Apl) (Value, error) {
	return x, nil
}

// IdxSpec represents an expression of an index specification.
// That is what is in square brackets following an array.
type idxSpec []expr

func (x idxSpec) String(a *Apl) string {
	if x == nil {
		return "[nil]"
	}
	v := make([]string, len(x))
	for i, e := range x {
		v[i] = e.String(a)
	}
	return "[" + strings.Join(v, ";") + "]"
}

func (x idxSpec) Eval(a *Apl) (Value, error) {
	idx := make(IdxSpec, len(x))
	for i, e := range x {
		if v, err := e.Eval(a); err != nil {
			return nil, err
		} else {
			idx[i] = v
		}
	}
	return idx, nil
}
