// APL stream processor.
//
// Usage
//	cat data | iv COMMANDS
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	"github.com/ktye/iv/cmd"
)

var stdin io.ReadCloser = os.Stdin

func main() {
	if len(os.Args) < 2 {
		fatal(fmt.Errorf("arguments expected"))
	}
	a := newApl(stdin)
	fatal(cmd.Iv(a, strings.Join(os.Args[1:], " "), os.Stdout))
}

func newApl(r io.ReadCloser) *apl.Apl {
	stdin = r
	a := apl.New(nil)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)

	// Add a minimal io package that's sole purpose is to allow
	// to read from stdin.
	pkg := map[string]apl.Value{
		"r": apl.ToFunction(readfd),
	}
	a.RegisterPackage("io", pkg)
	return a
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func readfd(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	fd, ok := R.(apl.Int)
	if ok == false {
		return nil, fmt.Errorf("iv/io read: argument must be 0 (stdin)")
	}
	if fd != 0 {
		return nil, fmt.Errorf("monadic <: right argument must be 0 (stdin)")
	}
	return apl.LineReader(stdin), nil
}
