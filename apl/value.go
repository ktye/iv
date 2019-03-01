package apl

import "io"

// Value is the result of an evaluation.
// Any type that implements the interface is a valid type for apl.
//
// The String method is used to display the value.
// It does not need to be unique or represent the input syntax.
// Mosty types have no input respresentation at all.
// They are the result of a computation.
//
// A Value may implement further interfaces such as Array, Uniform or Function.
type Value interface {
	String(Format) string
	Copy() Value
}

// VarReader is implemented by Values that are able to parse from a Reader.
// The ReadFrom method must return a new value of the same type.
// The function should be able to parse the format of it's String method.
// It is used by varfs in package io.
type VarReader interface {
	ReadFrom(*Apl, io.Reader) (Value, error)
}
