package xgo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ktye/iv/apl"
)

func Register(a *apl.Apl) {
	pkg := map[string]apl.Value{
		"t":      New(reflect.TypeOf(T{})),
		"s":      New(reflect.TypeOf(S{})),
		"i":      New(reflect.TypeOf(I(0))),
		"source": source{},
		"echo":   echo{},
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

// source returns a Channel to pull numbers from.
// It stops if the max value is reached or the channel is closed.
// It is used for demonstrating apl.Channel.
type source struct{}

func (_ source) String(a *apl.Apl) string {
	return "source"
}

func (s source) Call(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	num, ok := R.(apl.Number)
	if ok == false {
		return nil, fmt.Errorf("source: r must be a number")
	}
	n, ok := num.ToIndex()
	if ok == false || n <= 0 {
		return nil, fmt.Errorf("source: R must be a positive integer")
	}
	c := apl.NewChannel()
	go func(c apl.Channel) {
		for i := 0; i < n; i++ {
			select {
			case <-c[1]:
				close(c[0])
				return
			default:
				c[0] <- apl.Index(i)
			}
		}
		close(c[0])
	}(c)
	return c, nil
}

// echo returns a Channel which echos what was send to it.
type echo struct{}

func (_ echo) String(a *apl.Apl) string {
	return "echo"
}

func (_ echo) Call(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	p := ""
	if s, ok := R.(apl.String); ok == false || len(s) == 0 {
		return nil, fmt.Errorf("echo: R must be a non-empty string")
	} else {
		p = string(s)
	}
	c := apl.NewChannel()
	go func(c apl.Channel) {
		var buf []apl.Value
		s := "empty stack"
		for {
			select {
			case r, ok := <-c[1]:
				if ok == false {
					close(c[0])
					return
				}
				buf = append(buf, r)
				if len(buf) == 1 {
					s = p + r.String(a)
				}
			case c[0] <- apl.String(s):
				if len(buf) > 1 {
					copy(buf, buf[1:])
					buf = buf[:len(buf)-1]
					if len(buf) > 0 {
						s = p + buf[0].String(a)
					}
				} else if len(buf) == 1 {
					buf = buf[:0]
					s = p + "empty stack"
				}
			}
		}
		close(c[0])
	}(c)
	return c, nil
}
