package apl

import (
	"fmt"
	"strings"
)

// Train is a function train expression.
// It does not contain it's arguments like a Primitive or a fnVar.
// See DyaProg p 25.
type train []expr

func (t train) Eval(a *Apl) (Value, error) {
	return t, nil
}

func (t train) String(a *Apl) string {
	v := make([]string, len(t))
	for i := range t {
		v[i] = t[i].String(a)
	}
	return "(" + strings.Join(v, ", ") + ")"
}
func (t train) Copy() Value { return t }

func (t train) Call(a *Apl, L, R Value) (Value, error) {
	if len(t) < 2 {
		return nil, fmt.Errorf("cannot call short train, length %d", len(t))
	} else if len(t)%2 == 0 {
		// even number: f g h i j k → f(g h(i j k)) ⍝ atop(fork(fork))
		f := atop{}
		end := 1
		if len(t) == 2 {
			end = 2
		}
		for i, e := range t[:end] {
			if v, err := e.Eval(a); err != nil {
				return nil, err
			} else {
				f[i] = v
			}
		}
		if len(t) > 3 {
			f[1] = train(t[1:])
		}
		return f.Call(a, L, R)
	} else {
		// odd number: e f g h i j k → e f(g h(i j k)) ⍝ fork(fork(fork))
		f := fork{}
		end := 2
		if len(t) == 3 {
			end = 3
		}
		for i, e := range t[:end] {
			if v, err := e.Eval(a); err != nil {
				return nil, err
			} else {
				f[i] = v
			}
		}
		if len(t) > 3 {
			f[2] = train(t[2:])
		}
		return f.Call(a, L, R)
	}
}

type atop [2]Value

func (t atop) String(a *Apl) string {
	return fmt.Sprintf("(%s %s)", t[0].String(a), t[1].String(a))
}

func (t atop) Call(a *Apl, L, R Value) (Value, error) {
	g, ok := t[0].(Function)
	if ok == false {
		return nil, fmt.Errorf("atop: expected function g: %T", t[0])
	}
	h, ok := t[1].(Function)
	if ok == false {
		return nil, fmt.Errorf("atop: expected function h: %T", t[1])
	}

	// L may be nil.
	v, err := h.Call(a, L, R)
	if err != nil {
		return nil, err
	}
	return g.Call(a, nil, v)
}

type fork [3]Value

func (fk fork) String(a *Apl) string {
	return fmt.Sprintf("(%s %s %s)", fk[0].String(a), fk[1].String(a), fk[2].String(a))
}

func (fk fork) Call(a *Apl, L, R Value) (Value, error) {
	f, fok := fk[0].(Function)

	g, ok := fk[1].(Function)
	if ok == false {
		return nil, fmt.Errorf("fork: expected function g: %T", fk[1])
	}
	h, ok := fk[2].(Function)
	if ok == false {
		return nil, fmt.Errorf("fork: expected function h: %T", fk[2])
	}

	// Agh fork if f is not a function.
	var l Value
	l = fk[0]
	if fok {
		if v, err := f.Call(a, L, R); err != nil { // TODO copy?
			return nil, err
		} else {
			l = v
		}
	}
	r, err := h.Call(a, L, R) // TODO copy?
	if err != nil {
		return nil, err
	}
	return g.Call(a, l, r)
}
