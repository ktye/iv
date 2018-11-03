package iv

import "io"

// Apl abstracts the apl backend used by iv.
//
// The APL script relies on the following variables,
// which must be implemented:
//	_:   current array with ρρ_ ←→ iv.rank1
//	E:   (int) termination level (default 0)
//	N:   (int) number of records
//	EOF: (bool) last call
// The APL script may respond to iv to control program flow
// by assigning to variables, which must be implemented:
//	NEXT: if != 0: skip next rules
//	E:    if it is a string, print error and exit.
//	      if it is the empty string or an int < 0, exit with 0.
type Apl interface {
	// Parse and execute apl expressions.
	// This may be called on startup for the begin block.
	Parse(string, string) error

	// AddRule adds a rule which consists of a conditional and a statement.
	// It should not be executed when called.
	AddRule(string, string) error

	// ParseScalar returns a Scalar from a string or an error if it does not match.
	// This should not change the state of the interpreter.
	ParseScalar(string) (Scalar, error)

	// Execute sets the current n-dimensional object with the given shape.
	// It also sets the dimension of the object and signals eof.
	Execute([]int, []Scalar, int, bool) error

	// SetOut sets stdout and stderr mainly used for testing.
	SetOut(stdout, stderr io.Writer)
}

// Scalar is a scalar object from the apl implementation, usually a numeric value.
type Scalar interface{}
