package apl

import "fmt"

// Assignment is the evaluates left operand of the derived function from an assignment.
// It contains the Identifier, possible Indexes and a modification function.
type Assignment struct {
	Var Identifier
	Idx IndexArray
	Mod Function
}

func (as Assignment) String(a *Apl) string {
	modified := ""
	indexed := ""
	if as.Mod != nil {
		modified = "modified "
	}
	if len(as.Idx.Ints) > 0 {
		indexed = "indexed "
	}
	return modified + indexed + "assignment to " + string(as.Var)
}

// EvalAssign evaluates a derived expression as a special case, if the operator is an assigment.
// It does not evaluate the identifier.
// It returns the value of the left operand.
// There are several cases:
//	A ← 3		simple assignment
//	A[1;1] ← 3	indexed assigment
//	X f← 3		modified assignment
//	X[3] f← 3	indexed modified assignment
//	A←B←C←D←1	multiple assignment
// Not implemented is:
//	(A B) ← 1 2	multiple assignment / vector specification
//	A[⊂1 1]←101	choose/reach indexed assignment
//	(f A) ← 3	selective assignment
//	(EXP X)[I]←Y	combined indexed and selective assignment
//	(EXP X)f←Y	selective modified assignment
//
// Expressions such as (f A) ← B are evaluated to an indexAssign
func (d *derived) evalAssign() (Value, error) {
	if v, ok := d.lo.(numVar); ok {
		// Simple assignmant: A ← 3
		return Assignment{Var: Identifier(v.name)}, nil
	} else if v, ok := d.lo.(fnVar); ok {
		// Simple function assignmant: f ← 3
		return Assignment{Var: Identifier(v)}, nil
	}
	return nil, fmt.Errorf("cannot assign to: %T", d.lo)

	/*
		// Indexed assignment is converted to a function application:
		//	A[idx] ← Y ←→ (idx ⌷ A) ← Y
		// It's converted to ??
		if f, ok := d.lo.(*function); ok {
			if p, ok := f.Function.(Primitive); ok && p == "⌷" {
			return nil, fmt.Errorf("TODO: assign to function %T", f.Function)
		}

		// TODO: partial evaluation for selective specification,
		// indexed assigments, multiple assignments, ...
		return nil, fmt.Errorf("cannot assign to: %T", d.lo)
	*/
}
