package apl

import (
	"bufio"
	"fmt"
	"io"
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
			err = fmt.Errorf("panic: %s\n%s", r, string(debug.Stack()))
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
			switch v := val.(type) {
			case Channel:
				for e := range v[0] {
					fmt.Fprintln(a.stdout, e.String(a))
				}
			case Image:
				if a.stdimg != nil {
					err = a.stdimg.WriteImage(v)
					if err != nil {
						return err
					}
				} else {
					fmt.Println(a.stdout, v.String(a))
				}
			default:
				fmt.Fprintln(a.stdout, val.String(a))
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

// EvalFile parses and evalutes from a reader.
// It handles multiline statements.
// The file argument is used only in the error message.
func (a *Apl) EvalFile(r io.Reader, file string) (err error) {
	line := 0
	defer func() {
		if err != nil {
			err = fileError{file: file, line: line, err: err}
		}
	}()

	ok := true
	var p Program
	b := NewLineBuffer(a)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line++
		ok, err = b.Add(scanner.Text())
		if err != nil {
			return
		}

		if ok {
			p, err = b.Parse()
			if err != nil {
				return
			}

			err = a.Eval(p)
			if err != nil {
				return
			}
		}
	}
	if ok == false && b.Len() > 0 {
		return fmt.Errorf("multiline statement is not terminated")
	}
	return nil
}

type fileError struct {
	file string
	line int
	err  error
}

func (f fileError) Error() string {
	return fmt.Sprintf("%s:%d: %s", f.file, f.line, f.err.Error())
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
