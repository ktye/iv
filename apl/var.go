package apl

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/ktye/iv/apl/scan"
)

// Identifier is a the Value of a variable identifier.
// It evaluates to itself, not the stored Value.
type Identifier string

func (id Identifier) String(a *Apl) string {
	return string(id)
}

func (id Identifier) Eval(a *Apl) (Value, error) {
	return id, nil
}

// Assign assigns a value to a variable with the given name.
func (a *Apl) Assign(name string, v Value) error {
	return a.AssignEnv(name, v, nil)
}

// AssignEnv assigns a variable in the given environment.
func (a *Apl) AssignEnv(name string, v Value, env *env) error {
	ok, isfunc := isVarname(name)
	if ok == false {
		return fmt.Errorf("variable name is not allowed: %s", name)
	}

	// Assignment to the special variable ⎕ prints the value.
	if name == "⎕" {
		fmt.Fprintf(a.stdout, "%s\n", v.String(a))
		return nil
	} else if name == "⎕IO" {
		if n, ok := v.(Number); ok {
			if b, ok := a.Tower.ToBool(n); ok {
				a.Origin = 0
				if b {
					a.Origin = 1
				}
				return nil
			}
		}
		return fmt.Errorf("cannot set index origin: %T", v)
	}

	if _, ok := v.(Function); ok && isfunc != true {
		return fmt.Errorf("cannot assign a function to an uppercase variable")
	} else if ok == false && isfunc == true {
		return fmt.Errorf("only functions can be assigned to lowercase variables")
	}

	if env == nil {
		env = a.env
	}

	// Special case: Default left argument in lambda expressions:
	// Do not overwrite the given argument.
	if name == "⍺" && env.vars["⍺"] != nil {
		return nil
	}

	env.vars[name] = v
	return nil
}

// Lookup returns the value stored under the given variable name.
// It returns nil, if the variable does not exist.
// Variables are lexically scoped.
func (a *Apl) Lookup(name string) Value {
	v, _ := a.LookupEnv(name)
	return v
}

// LookupEnv returns the value of a variable and a pointer to the environment,
// where it was found.
func (a *Apl) LookupEnv(name string) (Value, *env) {
	if name == "⎕IO" {
		return Index(a.Origin), nil
	}

	e := a.env
	for {
		v, ok := e.vars[name]
		if ok {
			return v, e
		}
		if e.parent == nil {
			break
		}
		e = e.parent
	}
	return nil, nil
}

// NumVar contains the identifier to a value.
// The name is upper case and does not evaluate to a function.
// NumVar evaluates to the stored value or to an Identifier if it is undeclared.
type numVar struct {
	name string
}

func (v numVar) String(a *Apl) string {
	return v.name
}

func (v numVar) Eval(a *Apl) (Value, error) {
	x := a.Lookup(v.name)
	if x == nil {
		return Identifier(v.name), nil
	}
	return x, nil
}

// FnValue contains the identifier to a function value.
// It's name is lowercase.
// FnVar evaluates to the stored value or to an Identifier if it is undeclared.
type fnVar string

func (f fnVar) String(a *Apl) string {
	return string(f)
}

func (f fnVar) Eval(a *Apl) (Value, error) {
	return f, nil
}

func (f fnVar) Call(a *Apl, l, r Value) (Value, error) {
	x := a.Lookup(string(f))
	if x == nil {
		return Identifier(f), nil
	}
	fn, ok := x.(Function)
	if ok == false {
		return nil, fmt.Errorf("value in function variable is not a function: %T", x)
	}
	if fn == nil {
		return nil, fmt.Errorf("value in function variable %s is nil", string(f))
	}
	return fn.Call(a, l, r)
}

// isVarname returns if the string is allowed as a variable name and
// referes to a number or function value.
func isVarname(s string) (ok, isfunc bool) {
	if s == "" {
		return false, false
	}
	upper := false
	for i, r := range s {
		if scan.AllowedInVarname(r, i == 0) == false {
			return false, false
		}
		if i == 0 && unicode.IsUpper(r) || strings.IndexRune("⎕⍺⍵", r) != -1 {
			upper = true
		}
	}
	return true, upper == false
}
