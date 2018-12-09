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
	// TODO: should a train evaluate as a train of trains,
	// if it has more than 3 arguments?
	// Or should that be done at Call?
	return t, fmt.Errorf("TODO: evaluate train")
}

func (t train) String(a *Apl) string {
	v := make([]string, len(t))
	for i := range t {
		v[i] = t[i].String(a)
	}
	return "(" + strings.Join(v, ", ") + ")"
}

func (t train) Call(a *Apl, L, R Value) (Value, error) {
	return nil, fmt.Errorf("TODO: call function train")
}
