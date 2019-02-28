package apl

import (
	"fmt"
	"strings"
)

// Assignment contains the unevaluated left argument of an assignment.
// It may be an identifier or an expression containing an identifier.
// The identifier may be followed by a modifying function.
type assignment struct {
	e expr
}

func (as assignment) String(a *Apl) string {
	return "specification " + as.e.String(a)
}
func (as assignment) Copy() Value { return as } // This does not do a deep copy.

func (as assignment) Eval(a *Apl) (Value, error) {
	return as, nil
}

// Assignment is the evaluated left part of an assignment.
// It contains the Identifier, it's indexes and a modification function.
type Assignment struct {
	Identifier  string
	Identifiers []string // Multiple identifiers for vector assignment
	Indexes     Value    // Should be convertible to an Index vector
	Modifier    Function
}

func (as *Assignment) Copy() Value {
	r := Assignment{
		Identifier: as.Identifier,
	}
	if as.Identifiers != nil {
		r.Identifiers = make([]string, len(as.Identifiers))
		copy(r.Identifiers, as.Identifiers)
	}
	if as.Indexes != nil {
		r.Indexes = as.Indexes.Copy()
	}
	if as.Modifier != nil {
		r.Modifier = as.Modifier
	}
	return &r
}

func (as *Assignment) String(a *Apl) string {
	s := ""
	if as.Indexes != nil {
		s = "indexed/selective "
	}
	if as.Modifier != nil {
		s += "modified "
	}
	id := as.Identifier
	if as.Identifiers != nil {
		id = strings.Join(as.Identifiers, " ")
	}
	return "assignment to " + id
}

// EvalAssign evalutes the left part of an assignment and
// returns it as an Assignment value.
// It handles indexed, selective and modified assignment.
func evalAssign(a *Apl, e expr, modifier Function) (Value, error) {
	as := Assignment{
		Modifier: modifier,
	}

	// A function assignment can have no selections or modifications.
	if f, ok := e.(fnVar); ok {
		as.Identifier = string(f)
		return &as, nil
	}

	// Modified assignment masks the left argument in an assignment expr.
	if as, ok := e.(assignment); ok {
		e = as.e
	}

	// Vector assignment can only contain a vector of numVars.
	if ae, ok := e.(array); ok {
		as.Identifiers = make([]string, len(ae))
		for i, v := range ae {
			if nv, ok := v.(numVar); ok {
				as.Identifiers[i] = nv.name
			} else {
				return nil, fmt.Errorf("vector assignment can contain only numVars: %T", v)
			}
		}
		return &as, nil
	}

	// The identifier is the right-most argument in the expression.
	// The selection function (if present) is the function left to the Identifier.
	selection := false
	r := e
search:
	for {
		switch v := r.(type) {
		case numVar:
			as.Identifier = v.name
			break search
		case *function:
			r = v.right
			v.selection = true
			selection = true
		default:
			return nil, fmt.Errorf("unknown type in assignment expression: %T", r)
		}
	}

	if selection {
		if idx, err := e.Eval(a); err != nil {
			return nil, err
		} else {
			as.Indexes = idx
		}
	}

	return &as, nil
}
