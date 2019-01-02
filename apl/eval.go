package apl

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Program contains a slice of parsed expressions.
type Program []expr

// Eval executes an apl program.
// It can be called in a loop for every line of input.
func (a *Apl) Eval(p Program) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", string(debug.Stack()))
		}
	}()
	v := make([]string, len(p))
	if a.debug {
		for i, e := range p {
			v[i] = e.String(a)
		}
		fmt.Fprintln(a.stdout, strings.Join(v, " ⋄ "))
	}

	var val Value
	for _, expr := range p {
		val, err = expr.Eval(a)
		if err != nil {
			return err
		}

		if a.debug {
			fmt.Fprintf(a.stdout, "%T\n", val)
		}
		if isAssignment(expr) == false {
			fmt.Fprintln(a.stdout, val.String(a))
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
	// Assignment is implemented as an operator.
	if fn, ok := e.(*function); ok && fn != nil {
		if d, ok := fn.Function.(*derived); ok && d.op == "←" {
			return true
		}
	}
	return false
}
