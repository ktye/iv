package apl

import (
	"fmt"
	"strings"
)

// IdxSpec is the value of an index specification.
// For an axis specification this must be a length 1 slice
// and the value is an apl.Float.
// For an index specification this is an array and each value
// is either an EmptyArray or an IntArray.
// IdxSpec is origin dependend.
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
func (x IdxSpec) Copy() Value {
	r := make(IdxSpec, len(x))
	for i := range r {
		r[i] = x[i].Copy()
	}
	return r
}

func (x IdxSpec) Eval(a *Apl) (Value, error) {
	return x, nil
}

// Shape returns the result shape of the index specification applied to an array with shape src.
func (x IdxSpec) Shape(src []int) ([]int, error) {
	return nil, fmt.Errorf("TODO apl.IdxSpec.Shape")
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

// Axis combines the right argument with an axis specification.
type Axis struct {
	R Value
	A Value
}

func (ax Axis) String(a *Apl) string {
	return fmt.Sprintf("[%s]%s", ax.A.String(a), ax.R.String(a))
}
func (ax Axis) Copy() Value {
	r := Axis{}
	if ax.R != nil {
		r.R = ax.R.Copy()
	}
	if ax.A != nil {
		r.A = ax.A.Copy()
	}
	return r
}
