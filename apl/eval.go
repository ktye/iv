package apl

import (
	"fmt"
	"strings"
)

// Program contains a slice of parsed expressions.
type Program []expr

// Eval executes an apl program.
// It can be called in a loop for every line of input.
func (a *Apl) Eval(p Program) error {
	v := make([]string, len(p))
	if a.debug {
		for i, e := range p {
			v[i] = e.String(a)
		}
		fmt.Fprintln(a.stdout, strings.Join(v, " ⋄ "))
	}

	for i, expr := range p {
		v, err := expr.Eval(a)
		if err != nil {
			return err
		}
		if i == 0 {
			if a.debug {
				fmt.Fprintf(a.stdout, "%T\n", v)
			}
			if isAssignment(expr) == false {
				fmt.Fprintln(a.stdout, v.String(a))
			}
		}
	}
	return nil
}

// EvalProgram evaluates all expressions in the program and returns the values.
func (a *Apl) EvalProgram(p Program) ([]Value, error) {
	res := make([]Value, len(p))
	for i, e := range p {
		if v, err := e.Eval(a); err != nil {
			return nil, err
		} else {
			res[i] = v
		}
	}
	return res, nil
}

func (p Program) String(a *Apl) string {
	v := make([]string, len(p))
	for i := range p {
		v[i] = p[i].String(a)
	}
	return strings.Join(v, "⋄")
}

type expr interface {
	Eval(*Apl) (Value, error)
	String(*Apl) string
}

func isAssignment(e expr) bool {
	// This works as long as the symbol is not overloaded
	// with something else.
	if fn, ok := e.(*function); ok && fn != nil {
		if p, ok := fn.Function.(Primitive); ok && p == "←" {
			return true
		}
	}
	return false
}
