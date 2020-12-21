// APL interpreter.
//
// Usage
//	apl < INPUT
//	apl FILES...
package main

import (
	"fmt"
	"os"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/big"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	"github.com/ktye/iv/cmd"
)

func main() {
	a := newApl()
	a.SetOutput(os.Stdout)
	if err := cmd.Apl(a, os.Stdin, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newApl() *apl.Apl {
	a := apl.New(nil)
	numbers.Register(a)
	big.Register(a, "")
	primitives.Register(a)
	operators.Register(a)
	return a
}
