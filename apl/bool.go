package apl

import "fmt"

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

// Bools is a uniform array of type bool.
type Bools struct {
	Dims  []int
	Bools []bool
}

func (b Bools) String(a *Apl) string {
	return ArrayString(a, b)
}

func (b Bools) At(i int) (Value, error) {
	if i < 0 || i >= len(b.Bools) {
		return nil, fmt.Errorf("index out of range")
	}
	return Bool(b.Bools[i]), nil
}

func (b Bools) Shape() []int {
	return b.Dims
}

func (b Bools) Size() int {
	return len(b.Bools)
}

func (b Bools) Zero() interface{} {
	return false
}

func (b Bools) Reshape(shape []int) Value {
	res := Bools{
		Dims:  shape,
		Bools: make([]bool, prod(shape)),
	}
	k := 0
	for i := range b.Bools {
		res.Bools[i] = b.Bools[k]
		k++
		if k == len(b.Bools) {
			k = 0
		}
	}
	return res
}
