// Package io provides input and output streams.
//
// Linking it into APL leads to an unsafe system.
//
// Io overloads several primitive operators:
//	< string          returns a Channel reading from a file
//      < 0               returns a Channel reading from stdin
//	!`ls              execute program return a channel
//	!(`ls`-l)         same with arguments
//	`cat!A            same reading input from A (String method) or channel (pipe)
//	`file<channel     write to file (TODO)
//	`dst<<`src        copy idiom (TODO)
//	`log<!`prog       redirection (TODO)
//	⍎¨channel         returns a channel with values
package io

import (
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

// Register adds the io package to the interpreter.
// This will provide access to the file system and allows to start external processes.
func Register(a *apl.Apl, name string) {
	if name == "" {
		name = "io"
	}
	pkg := map[string]apl.Value{
		"r": apl.ToFunction(read),
		"x": apl.ToFunction(exec),
	}
	a.RegisterPackage(name, pkg)

	a.RegisterPrimitive("<", apl.ToHandler(
		read,
		domain.Monadic(domain.IsString(nil)),
		"read file",
	))
	a.RegisterPrimitive("<", apl.ToHandler(
		readfd,
		domain.Monadic(domain.ToIndex(nil)),
		"read fd",
	))
	a.RegisterPrimitive("!", apl.ToHandler(
		exec,
		domain.Monadic(domain.ToStringArray(nil)),
		"exec",
	))
	a.RegisterPrimitive("!", apl.ToHandler(
		exec,
		domain.Dyadic(domain.Split(domain.ToStringArray(nil), nil)),
		"exec",
	))

}
