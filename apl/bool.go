package apl

import (
	"fmt"
)

type Bool bool

func (b Bool) String(a *Apl) string {
	if b {
		return "1"
	}
	return "0"
}

func (i Bool) Less(v Value) (Bool, bool) {
	j, ok := v.(Bool)
	if ok == false {
		return false, false
	}
	return i == false && j == true, true
}

func (i Bool) ToIndex() (int, bool) {
	if i {
		return 1, true
	}
	return 0, true
}

// BoolArray is a uniform array of type bool.
type BoolArray struct {
	Dims  []int
	Bools []bool
}

func (b BoolArray) String(a *Apl) string {
	return ArrayString(a, b)
}

func (b BoolArray) At(i int) (Value, error) {
	if i < 0 || i >= len(b.Bools) {
		return nil, fmt.Errorf("index out of range")
	}
	return Bool(b.Bools[i]), nil
}

func (b BoolArray) Shape() []int {
	return b.Dims
}

func (b BoolArray) Size() int {
	return len(b.Bools)
}

func (b BoolArray) Zero() Value {
	return Bool(false)
}

func (b BoolArray) Set(i int, v Value) error {
	if i < 0 || i > len(b.Bools) {
		return fmt.Errorf("index out of range")
	}
	if c, ok := v.(Bool); ok {
		b.Bools[i] = bool(c)
		return nil
	}
	return fmt.Errorf("cannot assign %T to BoolArray", v)
}

func makeBoolArray(v []Value) BoolArray {
	b := make([]bool, len(v))
	for i, e := range v {
		b[i] = bool(e.(Bool))
	}
	return BoolArray{
		Dims:  []int{len(v)},
		Bools: b,
	}
}

func (b BoolArray) Reshape(shape []int) Value {
	res := BoolArray{
		Dims:  shape,
		Bools: make([]bool, prod(shape)),
	}
	k := 0
	for i := range res.Bools {
		res.Bools[i] = b.Bools[k]
		k++
		if k == len(b.Bools) {
			k = 0
		}
	}
	return res
}
