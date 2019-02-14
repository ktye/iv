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
	"github.com/ktye/iv/apl/scan"
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
	a.AddCommands(map[string]scan.Command{"l": load{}})

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
	return a
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// The following functions are a subset of apl/io to support loading a library
// and streaming from stdin.

func readfd(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	fd := int(R.(apl.Int))
	if fd != 0 {
		return nil, fmt.Errorf("monadic <: right argument must be 0 (stdin)")
	}
	return apl.LineReader(stdin), nil
}
func read(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	name, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("read: expect file name %T", R)
	}
	f, err := os.Open(string(name))
	if err != nil {
		return nil, err
	}
	return apl.LineReader(f), nil // LineReader closes the file.
}

type load struct{}

func (c load) Rewrite(t []scan.Token) []scan.Token {
	if len(t) < 2 {
		return t
	}
	return append([]scan.Token{
		scan.Token{T: scan.Symbol, S: "⍎"},
		scan.Token{T: scan.Symbol, S: "¨"},
		scan.Token{T: scan.Symbol, S: "<"},
		t[0],
		scan.Token{T: scan.Diamond, S: "⋄"},
	}, t[1:]...)
}
