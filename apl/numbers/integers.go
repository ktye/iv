package numbers

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

// IntegerArray is a uniform array of int64.
type IntegerArray struct {
	Dims []int
	Ints []int64
}

func (f IntegerArray) String(a *apl.Apl) string {
	return apl.ArrayString(a, f)
}

func (f IntegerArray) At(i int) apl.Value {
	return Integer(f.Ints[i])
}

func (f IntegerArray) Shape() []int {
	return f.Dims
}

func (f IntegerArray) Size() int {
	return len(f.Ints)
}

func (f IntegerArray) Zero() apl.Value {
	return Integer(0)
}

func (f IntegerArray) Set(i int, v apl.Value) error {
	if i < 0 || i > len(f.Ints) {
		return fmt.Errorf("index out of range")
	}
	if c, ok := v.(Integer); ok {
		f.Ints[i] = int64(c)
		return nil
	}
	return fmt.Errorf("cannot assign %T to IntegerArray", v)
}

func (f IntegerArray) Make(shape []int) apl.Array {
	return IntegerArray{
		Dims: shape,
		Ints: make([]int64, prod(shape)),
	}
}

func makeIndexArray(v []apl.Value) apl.IndexArray {
	f := make([]int, len(v))
	for i, e := range v {
		f[i] = int(e.(apl.Index))
	}
	return apl.IndexArray{
		Dims: []int{len(v)},
		Ints: f,
	}
}

func makeIntegerArray(v []apl.Value) IntegerArray {
	f := make([]int64, len(v))
	for i, e := range v {
		f[i] = int64(e.(Integer))
	}
	return IntegerArray{
		Dims: []int{len(v)},
		Ints: f,
	}
}

func (f IntegerArray) Reshape(shape []int) apl.Value {
	res := IntegerArray{
		Dims: shape,
		Ints: make([]int64, prod(shape)),
	}
	k := 0
	for i := range res.Ints {
		res.Ints[i] = f.Ints[k]
		k++
		if k == len(f.Ints) {
			k = 0
		}
	}
	return res
}
