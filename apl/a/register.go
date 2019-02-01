// Package a provides an interface to the internals of the interpreter and the go runtime.
//
// Some of the functions defined return information about the go runtime.
// This includes the parent application if the interpreter is built-in.
//
//	c 0 return number of CPUs
//	g 0    return number of go routines
//	m 0    return runtime.MemStats as a dictionary
//	v 0    return go version
package a

import (
	"fmt"
	"os"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/scan"
)

// Register adds the a package to the interpreter.
func Register(p *apl.Apl, name string) {
	if name == "" {
		name = "a"
	}
	pkg := map[string]apl.Value{
		"c": apl.ToFunction(cpus),
		"g": apl.ToFunction(goroutines),
		"h": apl.ToFunction(help),
		"m": apl.ToFunction(Memstats),
		"p": apl.ToFunction(printvar),
		"q": apl.ToFunction(quit),
		"t": apl.ToFunction(timer),
		"v": apl.ToFunction(goversion),
	}
	cmd := map[string]scan.Command{
		"h": rw0("h"),
		"p": toCommand(printCmd),
		"q": rw0("q"),
		"t": toCommand(timeCmd),
	}
	p.AddCommands(cmd)
	p.RegisterPackage(name, pkg)
}

// Quit accepts a string or a number.
// Nonempty strings return an exit code != 0 and print the error to stderr.
// Numbers are used as exit code.
func quit(p *apl.Apl, _, R apl.Value) (apl.Value, error) {
	if n, ok := R.(apl.Number); ok {
		if idx, ok := n.ToIndex(); ok == false {
			return nil, fmt.Errorf("a q: exit code must be convertible to int")
		} else {
			os.Exit(idx)
		}
	}
	if s, ok := R.(apl.String); ok {
		if s == "" {
			os.Exit(0)
		} else {
			fmt.Fprintln(os.Stderr, s)
			os.Exit(1)
		}
	}
	return nil, fmt.Errorf("a q: argument must be a string or an int: %T", R)
}
