package numbers

import (
	"fmt"
	"time"

	"github.com/ktye/iv/apl"
)

// TimeArray is a uniform array of time.Time.
type TimeArray struct {
	Dims  []int
	Times []time.Time
}

func (t TimeArray) String(a *apl.Apl) string {
	return apl.ArrayString(a, t)
}

func (t TimeArray) At(i int) (apl.Value, error) {
	if i < 0 || i >= len(t.Times) {
		return nil, fmt.Errorf("index out of range")
	}
	return Time(t.Times[i]), nil
}

func (t TimeArray) Shape() []int {
	return t.Dims
}

func (t TimeArray) Size() int {
	return len(t.Times)
}

func (t TimeArray) Zero() apl.Value {
	return Time(y0)
}

func (t TimeArray) Set(i int, v apl.Value) error {
	if i < 0 || i > len(t.Times) {
		return fmt.Errorf("index out of range")
	}
	if c, ok := v.(Time); ok {
		t.Times[i] = time.Time(c)
		return nil
	}
	return fmt.Errorf("cannot assign %T to TimeArray", v)
}

func (t TimeArray) Reshape(shape []int) apl.Value {
	res := TimeArray{
		Dims:  shape,
		Times: make([]time.Time, prod(shape)),
	}
	k := 0
	for i := range res.Times {
		res.Times[i] = t.Times[k]
		k++
		if k == len(t.Times) {
			k = 0
		}
	}
	return res
}
