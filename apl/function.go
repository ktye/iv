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
	if f.left != nil {

		// Special case for modified assignments.
		// Defer evaluation of the left argument.
		if d, ok := f.right.(*derived); ok && d.op == "←" {
			l = assignment{f.left}
		} else {
			l, err = f.left.Eval(a)
			if err != nil {
				return nil, err
			}
		}
	}
	r, err = f.right.Eval(a)
	if err != nil {
		return nil, err
	}

	// Special case: the last function in a selective assignment uses Select instead of Call.
	if _, ok := f.right.(numVar); ok && f.selection {
		if p, ok := f.Function.(Primitive); ok == false {
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
func (p Primitive) Select(a *Apl, L, R Value) (IndexArray, error) {
	if handles := a.primitives[p]; handles == nil {
		return IndexArray{}, fmt.Errorf("primitive function %s does not exist", p)
	} else {
		for _, h := range handles {
			if l, r, ok := h.To(a, L, R); ok {
				return h.Select(a, l, r)
			}
		}
	}
	if L == nil {
		return IndexArray{}, fmt.Errorf("primitive is not implemented: %s %T ", p, R)
	}
	return IndexArray{}, fmt.Errorf("primitive is not implemented: %T %s %T ", L, p, R)
}

// PrimitiveHandler is the interface that implementations of primitive functions satisfy.
type PrimitiveHandler interface {
	Domain
	Call(*Apl, Value, Value) (Value, error)
	Select(*Apl, Value, Value) (IndexArray, error)
	Doc() string
}
