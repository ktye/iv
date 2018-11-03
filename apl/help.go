package apl

var doc map[string]string

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
	return doc[key]
}
