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
	"github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
)

func main() {
	if len(os.Args) < 2 {
		fatal(fmt.Errorf("arguments expected"))
	}
	fatal(iv(strings.Join(os.Args[1:], " "), os.Stdout))
}

func iv(p string, w io.Writer) error {
	a := apl.New(w)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)

	a.RegisterPrimitive("<", apl.ToHandler(
		readfd,
		domain.Monadic(domain.ToIndex(nil)),
		"read fd",
	))
	return a.ParseAndEval(p)
}

var stdin io.ReadCloser = os.Stdin

// readfd is copied from apl/io.
func readfd(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	fd := int(R.(apl.Int))
	if fd != 0 {
		return nil, fmt.Errorf("monadic <: right argument must be 0 (stdin)")
	}
	return apl.LineReader(stdin), nil
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
