// Package io provides input and output streams.
//
// Linking it into APL leads to an unsafe system.
// See README.md
package io

import (
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/scan"
)

// Register adds the io package to the interpreter.
// This will provide access to the file system and allows to start external processes.
func Register(a *apl.Apl, name string) {
	if name == "" {
		name = "io"
	}
	pkg := map[string]apl.Value{
		"cd":     apl.ToFunction(cd),
		"e":      apl.ToFunction(env),
		"l":      apl.ToFunction(load),
		"r":      apl.ToFunction(read),
		"x":      apl.ToFunction(exec),
		"mount":  apl.ToFunction(mount),
		"umount": apl.ToFunction(umount),
	}
	cmd := map[string]scan.Command{
		"cd": toCommand(cdCmd),
		"l":  toCommand(lCmd),
		"m":  toCommand(mCmd),
	}
	a.AddCommands(cmd)
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
	RegisterProtocol("var", varfs{Apl: a})
}
