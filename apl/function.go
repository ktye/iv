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

// Call looks up the function handle within registered primitives.
// It calls the handle with the left and right argument.
// Left is nil in a monadic context.
func (p Primitive) Call(a *Apl, l, r Value) (Value, error) {
	if handles := a.primitives[p]; handles == nil {
		return nil, fmt.Errorf("primitive function %s does not exist", p)
	} else {
		for i := len(handles) - 1; i >= 0; i-- {
			if ok, v, err := handles[i](a, l, r); ok {
				return v, err
			}
		}
		if l == nil {
			return nil, fmt.Errorf("primitive is not implemented for the type: %s %T ", p, r)
		}
		return nil, fmt.Errorf("primitive is not implemented for the type: %T %s %T ", l, p, r)
	}
}

// FunctionHandle evaluates a function call.
// It gets two operands, where the left may be nil in monadic context.
// If it cannot handle the given types, it should return false and the other
// return values are ignored.
// It should return true, if it feels responsible for the types,
// with a possible error, if e.g. the domain of the values are illegal.
type FunctionHandle func(*Apl, Value, Value) (bool, Value, error)
