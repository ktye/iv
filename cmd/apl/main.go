// APL interpreter.
//
// Usage
//	apl < INPUT
// Server mode
//	apl ADDR < INPUT
// Example
//	apl ":1966"
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	"github.com/ktye/iv/apl/rpc"
)

func main() {
	a := apl.New(os.Stdout)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)
	rpc.Register(a)

	line := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		line++
		// TODO: assemble multiline lambda expressions.
		if err := a.ParseAndEval(s); err != nil {
			fmt.Fprintf(os.Stderr, "%d: %s\n", line, err)
			os.Exit(1)
		}
	}

	if len(os.Args) > 1 {
		rpc.ListenAndServe(a, os.Args[1])
	}
	return
}
