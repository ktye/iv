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
	write := func(val Value) {
		switch v := val.(type) {
		case Image:
			if a.stdimg != nil {
				a.stdimg.WriteImage(v)
			} else {
				fmt.Fprintln(a.stdout, val.String(a.Format))
			}
		default:
			fmt.Fprintln(a.stdout, val.String(a.Format))
		}
	}

	var val Value
	for _, expr := range p {
		val, err = expr.Eval(a)
		if err != nil {
			return err
		}
		if isAssignment(expr) == false {
			switch v := val.(type) {
			case Channel:
				i := 0
				for e := range v[0] {
					if i == 0 {
						i++
						if _, ok := e.(Image); ok && a.stdimg != nil {
							a.stdimg.StartLoop()
							defer a.stdimg.StopLoop()
						}
					}
					write(e)
				}
			default:
				write(val)
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

func (p Program) String(f Format) string {
	v := make([]string, len(p))
	for i := range p {
		v[i] = p[i].String(f)
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
	String(f Format) string
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
