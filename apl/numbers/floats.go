package numbers

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

// FloatArray is a uniform array of float64.
type FloatArray struct {
	Dims   []int
	Floats []float64
}

func (f FloatArray) String(a *apl.Apl) string {
	return apl.ArrayString(a, f)
}

func (f FloatArray) At(i int) apl.Value {
	return Float(f.Floats[i])
}

func (f FloatArray) Shape() []int {
	return f.Dims
}

func (f FloatArray) Size() int {
	return len(f.Floats)
}

func (f FloatArray) Zero() apl.Value {
	return Float(0.0)
}

func (f FloatArray) Set(i int, v apl.Value) error {
	if i < 0 || i > len(f.Floats) {
		return fmt.Errorf("index out of range")
	}
	if c, ok := v.(Float); ok {
		f.Floats[i] = float64(c)
		return nil
	}
	return fmt.Errorf("cannot assign %T to FloatArray", v)
}

func (f FloatArray) Make(shape []int) apl.Array {
	return FloatArray{
		Dims:   shape,
		Floats: make([]float64, prod(shape)),
	}
}

func makeFloatArray(v []apl.Value) FloatArray {
	f := make([]float64, len(v))
	for i, e := range v {
		f[i] = float64(e.(Float))
	}
	return FloatArray{
		Dims:   []int{len(v)},
		Floats: f,
	}
}

func (f FloatArray) Reshape(shape []int) apl.Value {
	res := FloatArray{
		Dims:   shape,
		Floats: make([]float64, prod(shape)),
	}
	k := 0
	for i := range res.Floats {
		res.Floats[i] = f.Floats[k]
		k++
		if k == len(f.Floats) {
			k = 0
		}
	}
	return res
}

func prod(shape []int) int {
	if len(shape) == 0 {
		return 0
	}
	n := shape[0]
	for i := 1; i < len(shape); i++ {
		n *= shape[i]
	}
	return n
}
