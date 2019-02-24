package apl

import (
	"fmt"
)

// Function is any type that can be called, given it's left and right arguments.
// The left argument may be nil, in which case the function is called monadically.
// Currently this is implemented by *derived, *Primitive, fnVar and *lambda.
type Function interface {
	Call(*Apl, Value, Value) (Value, error)
}

// EnvCall calls a function in a new environment.
func (a *Apl) EnvCall(f Function, L, R Value, vars map[string]Value) (Value, error) {
	e := env{
		vars:   vars,
		parent: a.env,
	}
	save := a.env
	a.env = &e
	defer func() { a.env = save }()
	return f.Call(a, L, R)
}

// function wraps a Function with it's arguments.
type function struct {
	Function
	left, right expr
	selection   bool
}

// Eval calls the function with it's surrounding arugments.
func (f *function) Eval(a *Apl) (Value, error) {
	var err error
	var l, r Value

	// The right argument must be evaluated first.
	// Otherwise this A←1⋄A+(A←2) evaluates to 3,
	// but it should evaluate to 4.
	r, err = f.right.Eval(a)
	if err != nil {
		return nil, err
	}
	if f.left != nil {

		// Special case for modified assignments.
		// Defer evaluation of the left argument.
		if d, ok := f.Function.(*derived); ok && d.op == "←" {
			l = assignment{f.left}
		} else {
			l, err = f.left.Eval(a)
			if err != nil {
				return nil, err
			}
		}
	}

	// Special case: the last function in a selective assignment uses Select instead of Call.
	if _, ok := f.right.(numVar); ok && f.selection {
		if d, ok := f.Function.(*derived); ok == true {
			return d.Select(a, l, r)
		} else if p, ok := f.Function.(Primitive); ok == false {
			return nil, fmt.Errorf("cannot use %T in selective assignment", f.Function)
		} else {
			return p.Select(a, l, r)
		}
	}
	return f.Function.Call(a, l, r)
}

func (f *function) String(a *Apl) string {
	s := "nil"
	if f.Function != nil {
		switch p := f.Function.(type) {
		case *derived:
			s = p.String(a)
		case Primitive:
			s = string(p)
		case fnVar:
			s = string(p)
		case self:
			s = "∇"
		case *lambda:
			s = p.String(a)
		}
	}

	r := f.right.String(a)
	if f.left == nil {
		return fmt.Sprintf("(%s %s)", s, r)
	}
	l := f.left.String(a)
	return fmt.Sprintf("(%s %s %s)", l, s, r)
}

// Primitive is a primitive function expression.
// It may be a monadic or dyadic function.
// Primitives are defined and registered at compile time.
// Examples: + - × ⍴ ←.
// Default primitives are defined in package funcs, but others may
// be registered too.
// Multiple versions for the same symbol may be registered, which
// are tested in reverse sequence, until one takes over the responsibility.
// Implementing primitives does not involve using this type.
// It is exported, because it's used by operators to build derived functions.
type Primitive string

// Eval returns the primitive symbol itself.
func (p Primitive) Eval(a *Apl) (Value, error) {
	return p, nil
}

func (p Primitive) String(a *Apl) string {
	return string(p)
}

// Call looks up the handlers within registered primitives.
// It calls the handle with the left and right argument.
// Left is nil in a monadic context.
// If there are multiple handlers registered (primitive function overloading),
// they are tested in reverse registration order, until the first one takes the
// responsibility.
func (p Primitive) Call(a *Apl, L, R Value) (Value, error) {
	if handles := a.primitives[p]; handles == nil {
		return nil, fmt.Errorf("primitive function %s does not exist", p)
	} else {
		for _, h := range handles {
			if l, r, ok := h.To(a, L, R); ok {
				return h.Call(a, l, r)
			}
		}
	}
	if L == nil {
		return nil, fmt.Errorf("primitive is not implemented: %s %T ", p, R)
	}
	return nil, fmt.Errorf("primitive is not implemented: %T %s %T ", L, p, R)
}

// Select is similar to Call.
// It is used as a selection function in selective assignment.
// While the Call on these primitives would return selected values, Select returns the indexes of the values.
func (p Primitive) Select(a *Apl, L, R Value) (IntArray, error) {
	if handles := a.primitives[p]; handles == nil {
		return IntArray{}, fmt.Errorf("primitive function %s does not exist", p)
	} else {
		for _, h := range handles {
			if l, r, ok := h.To(a, L, R); ok {
				return h.Select(a, l, r)
			}
		}
	}
	if L == nil {
		return IntArray{}, fmt.Errorf("primitive is not implemented: %s %T ", p, R)
	}
	return IntArray{}, fmt.Errorf("primitive is not implemented: %T %s %T ", L, p, R)
}

// PrimitiveHandler is the interface that implementations of primitive functions satisfy.
type PrimitiveHandler interface {
	Domain
	Function
	Select(*Apl, Value, Value) (IntArray, error)
	Doc() string
}

// ToHandler wraps the arguments into a simplified primitive handler.
func ToHandler(f func(*Apl, Value, Value) (Value, error), d Domain, doc string) PrimitiveHandler {
	return pHandler{f, d, doc}
}

type pHandler struct {
	f func(*Apl, Value, Value) (Value, error)
	d Domain
	s string
}

func (h pHandler) Call(a *Apl, L, R Value) (Value, error) {
	return h.f(a, L, R)
}
func (h pHandler) To(a *Apl, L, R Value) (Value, Value, bool) {
	return h.d.To(a, L, R)
}
func (h pHandler) String(a *Apl) string {
	return h.d.String(a)
}
func (h pHandler) Doc() string {
	return h.s
}
func (h pHandler) Select(*Apl, Value, Value) (IntArray, error) {
	return IntArray{}, fmt.Errorf("function cannot be used in selective assignment")
}

// ToFunction can be used to cast a function with the right signature to a type that implements the Function interface.
type ToFunction func(*Apl, Value, Value) (Value, error)

func (f ToFunction) Call(a *Apl, L, R Value) (Value, error) {
	return f(a, L, R)
}

func (f ToFunction) String(a *Apl) string {
	return "anonymous function"
}
