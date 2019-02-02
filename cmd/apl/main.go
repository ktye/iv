// APL interpreter.
//
// Usage
//	apl < INPUT
//	apl FILES...
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
)

func main() {
	a := apl.New(os.Stdout)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)

	// Execute files.
	if len(os.Args) > 1 {
		for _, name := range os.Args[1:] {
			var r io.Reader
			if name == "-" {
				r = os.Stdin
				name = "stdin"
			} else {
				f, err := os.Open(name)
				fatal(err)
				defer f.Close()
				r = f
			}
			fatal(a.EvalFile(r, name))
		}
		os.Exit(0)
	}

	// Run interactively.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		if err := a.ParseAndEval(s); err != nil {
			fmt.Println(err)
		}
	}
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
