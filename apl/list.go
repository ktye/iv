package apl

import (
	"fmt"
	"strings"
)

// List is a collection of items, possibly nested.
// It also acts as a vector (a rank 1 array) but cannot be reshaped.
type List []Value

func (l List) String(a *Apl) string {
	var buf strings.Builder
	buf.WriteRune('(')
	for i := range l {
		buf.WriteString(l[i].String(a))
		buf.WriteRune(';')
	}
	buf.WriteRune(')')
	return buf.String()
}

func (l List) At(i int) (Value, error) {
	if i < 0 || i >= len(l) {
		return nil, fmt.Errorf("index out of range")
	}
	return l[i], nil
}

func (l List) Shape() []int {
	return []int{len(l)}
}

func (l List) Size() int {
	return len(l)
}

/*
// funcExpr is a dummy expression that is used to put the already evaluated
// function value argument into a derived operator expression as an operand.
type funExpr struct {
	Value
}

func (f funExpr) Eval(a *Apl) (Value, error) {
	return f.Value, nil
}
*/

// Evaluate evaluates an expression stored in a list in tree form.
// It is not called by the parser directly (that would be Eval), but
// from the ⍎ primitive.
func (l List) Evaluate(a *Apl) (Value, error) {
	return aeval(a, l)
}

// Aeval evaluates an argument element from a parse tree in list form.
// - a string is interpreted as a variable identifier and dereferenced
// - a list with a single element is disclosed
// - a list with multiple elements evaluates as a function followed 
//   one or two arguments. The function and their arguments are evaluated
//   and applied.
// - other types are returned
func aeval(a *Apl, v Value) (Value, error) {
	if s, ok := v.(String); ok {
		if ok, _ := isVarname(s); ok {
			r := a.Lookup(s)
			if r == nil {
				return nil, fmt.Errorf("variable undeclared: %s", s)
			}
			return r, nil
		}
		return nil, fmt.Errorf("not a variable: %s", s)
	}
	if l, ok := v.(List); ok == false {
		return v, nil
	} else if len(l) == 1 {
		return l[0], nil
	} else if len(l) > 3 {
		return nil, fmt.Errorf("eval: list is too long: %d", len(l))
	} else {
		// f L R
		// Evaluate the function first, than R, last L.
		f, err := feval(a, l[0])
		if err != nil {
			return nil, err
		}
		var R Value
		if len(l) == 3 {
			R, err = aeval(a, l[2])
			if err != nil {
				return nil, err
			}
		}
		L, err := aeval(a, l[1])
		if err != nil {
			return nil, err
		}
		return f.Call(a, L, R)
	}
}

// Feval evaluates a value to a function value, or returns an error,
// if it is not a verb.
// It accepts:
// - strings that refer to function variables
// - primitives
// - operator expressions in list form, including
//   λ (lambda) and ⍦ (train) operators
func feval(a *Apl, v Value) (Function, error) {
	if s, ok := v.(String); ok {
		if ok, isf := isVarname(s); ok && isf {
			f := a.Lookup(s)
			if f == nil {
				return nil, fmt.Errorf("function variable undeclared: %s", s)
			}
			return f, nil
		}
		return nil, fmt.Errorf("not a function variable: %s", s)
	}
	l, ok := v.(List)
	if ok == false {
		return nil, fmt.Errorf("not a function: %T", v)
	}
	if p, ok := v.(Primitive); ok {
		return p, nil
	}
	if d, ok := v.(*derived); ok && d.lo == nil {
		ops, ok := a.operators[d.op]
		if ok == false {
			return l, nil
		}
		if ops[0].DyadicOp() {
			if len(l) != 3 {
				return nil, fmt.Errorf("dyadic op (%s) needs 2 args, not %d", d.op, len(l)-1)
			}
// We do not know what follows an operator:
// It can be a function or an noun.
// How should that be evaluated? with aeval or feval?

		} else {
			if len(l) != 2 {
				return nil, fmt.Errorf("monadic op (%s) needs 1 arg, not %d", d.op, len(l)-1)
			}
		}
		//
}

	/*
// Evaluate evaluates an expression stored in a list in tree form.
// It is not called by the parser directly (that would be Eval), but
// from the ⍎ primitive.
func (l List) Evaluate(a *Apl) (Value, error) {

	if len(l) < 1 {
		return l, nil
	}
	eval := func(v Value) (Value, error) {
		if lv, ok := v.(List); ok {
			return lv.Evaluate(a)
		}
		return v, nil
	}
	evalf := func(v Value) (expr, error) {
		if _, ok := v.(Function); ok {
			return funExpr{v}, nil
		}
		if lv, ok := v.(List); ok {
			fv, err := lv.Evaluate(a)
			if err != nil {
				return nil, err
			}
			if _, ok := fv.(Function); ok {
				return funExpr{fv}, nil
			}
		}
		return nil, fmt.Errorf("evaluate list: expected function: %T", v)
	}
	var fn Function
	var args List
	if d, ok := l[0].(*derived); ok {
		if d.lo == nil {
			ops, ok := a.operators[d.op]
			if ok == false {
				return l, nil
			}
			if ops[0].DyadicOp() {
				if len(l) == 4 || len(l) == 5 {
					ro, err := evalf(l[2])
					if err != nil {
						return l, nil
					}
					lo, err := evalf(l[1])
					if err != nil {
						return l, nil
					}
					fn = &derived{
						op: d.op,
						lo: lo,
						ro: ro,
					}
					args = l[3:]
				}
			} else if len(l) == 3 || len(l) == 4 {
				lo, err := evalf(l[1])
				if err != nil {
					return l, nil
				}
				fn = &derived{
					op: d.op,
					lo: lo,
				}
				args = l[2:]
			}
		}
	} else if f, ok := l[0].(Function); ok {
		fn = f
		args = l[1:]
	}
	if fn != nil {
		R, err := eval(args[len(args)-1])
		if err != nil {
			return nil, err
		}
		var L Value
		if len(args) == 2 {
			lv, err := eval(args[0])
			if err != nil {
				return nil, err
			}
			L = lv
		}
		return fn.Call(a, L, R)
	}
	return l, nil
}
*/

type list []expr

func (l list) Eval(a *Apl) (Value, error) {
	lst := make(List, len(l))
	var err error
	for i := range lst {
		lst[i], err = l[i].Eval(a)
		if err != nil {
			return nil, err
		}
	}
	return lst, nil
}

func (l list) String(a *Apl) string {
	var buf strings.Builder
	buf.WriteRune('(')
	for i := range l {
		buf.WriteString(l[i].String(a))
		buf.WriteRune(';')
	}
	buf.WriteRune(')')
	return buf.String()
}
