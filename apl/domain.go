package apl

// Domain represents the application domain of a function or operator.
// Calling To converts the left and right input arguments and returns true,
// if the types or values are compatible with the function or operator.
// Otherwise they return false.
// If To returns false, it must return the original input values.
// This is important for a Not combination to work properly.
// String is used in documentation for a concise type/value description.
//
// Standard Domain implementations and universal combination functions
// are in the domain package.
type Domain interface {
	To(*Apl, Value, Value) (Value, Value, bool)
	String(*Apl) string
}
