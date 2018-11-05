package apl

import "sort"

// TODO use these placeholders in the documentations:
// L Left argument
// R Right argument
// F Function
// LO Left operand
// RO Right operand
// MOP Monadic operator
// DOP Dyadic operator

// GetDoc returns the documentation for the function.
// It is an empty string, if it does not exist.
func (a *Apl) GetDoc(key string) string {
	return a.doc[key]
}

// GetDocKeys returns all registered doc keys.
func (a *Apl) GetDocKeys() []string {
	var keys []string
	for key := range a.doc {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// ListAllPrimitives returns all registered Primitive functions.
func (a *Apl) ListAllPrimitives() []string {
	var l []string
	for p := range a.primitives {
		l = append(l, string(p))
	}
	sort.Strings(l)
	return l
}

// ListAllOperators returns all registered operators.
func (a *Apl) ListAllOperators() []string {
	var l []string
	for o := range a.operators {
		l = append(l, string(o))
	}
	sort.Strings(l)
	return l
}
