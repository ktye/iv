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

func (t TimeArray) String(f apl.Format) string {
	return apl.ArrayString(f, t)
}

func (t TimeArray) Copy() apl.Value {
	r := TimeArray{Dims: apl.CopyShape(t), Times: make([]time.Time, len(t.Times))}
	copy(r.Times, t.Times)
	return r
}

func (t TimeArray) At(i int) apl.Value {
	return Time(t.Times[i])
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

func (t TimeArray) Make(shape []int) apl.ArraySetter {
	return TimeArray{
		Dims:  shape,
		Times: make([]time.Time, prod(shape)),
	}
}

func makeTimeArray(v []apl.Value) TimeArray {
	t := make([]time.Time, len(v))
	for i, e := range v {
		t[i] = time.Time(e.(Time))
	}
	return TimeArray{
		Dims:  []int{len(v)},
		Times: t,
	}
}
