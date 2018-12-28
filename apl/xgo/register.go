package xgo

import (
	"reflect"
	"strings"

	"github.com/ktye/iv/apl"
)

func Register(a *apl.Apl) {
	pkg := map[string]apl.Value{
		"t": New(reflect.TypeOf(T{})),
		"s": New(reflect.TypeOf(S{})),
		"i": New(reflect.TypeOf(I(0))),
	}
	a.RegisterPackage("xgo", pkg)
}

type I int

// T is an example struct with methods with pointer receivers.
type T struct {
	A string
	I int
	F float64
	C complex128
	V []string
}

func (t *T) Inc() {
	t.I++
}

func (t *T) Join(sep string) (int, string) {
	s := strings.Join(t.V, sep)
	return len(t.V), s
}

// S is an example struct with a method without pointer receiver.
type S struct {
	A int
	B int
}

func (s S) Sum() int {
	return s.A + s.B
}
