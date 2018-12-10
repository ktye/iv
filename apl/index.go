package apl

import (
	"fmt"
	"strings"
)

// IdxSpec is the value of an index specification.
// For an axis specification this must be a length 1 slice
// and the value is an apl.Float.
// For an index specification this is an array and each value
// is either an EmptyArray or an IndexArray.
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

// TODO: make idxSpec implement an indexer(?) interface for index assignments.
