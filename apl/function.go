package apl

import "fmt"

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
}

// Eval calls the function with it's surrounding arugments.
func (f *function) Eval(a *Apl) (Value, error) {
	// Assignment is special, it does not evaluate the left argument.
	assignment := false
	if p, ok := f.Function.(Primitive); ok && p == "←" {
		assignment = true
	}

	var err error
	var l, r Value
	if f.left != nil {
		if assignment {
			if v, ok := f.left.(numVar); ok {
				l = Identifier(v.name)
			} else if v, ok := f.left.(fnVar); ok {
				l = Identifier(v)
			} else {
				return nil, fmt.Errorf("assignment to a non-variable: %T %s", f.left, f.left.String(a))
			}
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
// Implementing primitives does not involve using this structure.
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
func (p Primitive) Call(a *Apl, l, r Value) (Value, error) {
	if handles := a.primitives[p]; handles == nil {
		return nil, fmt.Errorf("primitive function %s does not exist", p)
	} else {
		for i := len(handles) - 1; i >= 0; i-- {
			if ok, v, err := handles[i].HandlePrimitive(a, l, r); ok {
				return v, err
			}
		}
		if l == nil {
			return nil, fmt.Errorf("primitive is not implemented for the type: %s %T ", p, r)
		}
		return nil, fmt.Errorf("primitive is not implemented for the type: %T %s %T ", l, p, r)
	}
}

// PrimitiveHandler is the interface that implementations of primitive functions satisfy.
// The handler decides based upon the input types if it is able to handle the function call
// and returns true if that is the case.
// If the it returns false the Value and error are ignored.
// Returning true and a non-nil error triggeres an evaluation error.
//
// The handler is used for both, monadic and dyadic calls.
// Monadic calls receive a nil value for the first (left) Value argument.
// It is in the responsitiliby of the function to accept, reject or raise an error.
type PrimitiveHandler interface {
	HandlePrimitive(*Apl, Value, Value) (bool, Value, error)
}
