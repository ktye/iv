package apl

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

type String string

// String formats s with %s by default.
// The format can be changed in Format.String.
func (s String) String(a *Apl) string {
	return string(s)
}

func (s String) Eval(a *Apl) (Value, error) {
	return s, nil
}

func (s String) ReadFrom(a *Apl, r io.Reader) (Value, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return String(b), nil
}

// Less implements primitives.lesser to be used for comparison and sorting.
func (s String) Less(r Value) (Bool, bool) {
	b, ok := r.(String)
	if ok == false {
		return false, false
	}
	return s < b, true
}

func (s String) Export() reflect.Value {
	return reflect.ValueOf(string(s))
}

// StringArray is a uniform array of strings.
type StringArray struct {
	Dims    []int
	Strings []string
}

func (s StringArray) String(a *Apl) string {
	return ArrayString(a, s)
}

func (s StringArray) At(i int) Value {
	return String(s.Strings[i])
}

func (s StringArray) Shape() []int {
	return s.Dims
}

func (s StringArray) Size() int {
	return len(s.Strings)
}

func (s StringArray) Zero() Value {
	return String("")
}

func (s StringArray) Set(i int, v Value) error {
	if i < 0 || i > len(s.Strings) {
		return fmt.Errorf("index out of range")
	}
	if c, ok := v.(String); ok {
		s.Strings[i] = string(c)
		return nil
	}
	return fmt.Errorf("cannot assign %T to StringArray", v)
}

func (s StringArray) Make(shape []int) Array {
	return StringArray{
		Dims:    shape,
		Strings: make([]string, prod(shape)),
	}
}

func makeStringArray(v []Value) StringArray {
	str := make([]string, len(v))
	for i, e := range v {
		str[i] = string(e.(String))
	}
	return StringArray{
		Dims:    []int{len(v)},
		Strings: str,
	}
}

func (s StringArray) Reshape(shape []int) Value {
	res := StringArray{
		Dims:    shape,
		Strings: make([]string, prod(shape)),
	}
	k := 0
	for i := range res.Strings {
		res.Strings[i] = s.Strings[k]
		k++
		if k == len(s.Strings) {
			k = 0
		}
	}
	return res
}
