package apl

// Value is the result of an evaluation.
// Any type that implements the interface is a valid type for apl.
//
// The String method is used to display the value.
// It does not need to be unique or represent the input syntax.
// Mosty types have no input respresentation at all.
// They are the result of a computation.
//
// The *Apl argument may be used for formatting.
// A type should return a useful string if it is nil.
//
// If the Value implementation has the method
//	Call(*Apl, Value, Value) (Value, error)
// then it's also a Function.
//
// If it implements Array, it is an array.
type Value interface {
	String(*Apl) string

	// TODO: should there be a Copy or Clone interface?
	// All primitives would have to use it if present.

	// TODO: should we require a serialization interface?
	// or serialize optionally if a Value implements an Encoder?
}
