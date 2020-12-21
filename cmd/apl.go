// Package cmd contains shared code between cmd/apl, cmd/iv and cmd/lui
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ktye/iv/apl"
)

// Apl runs the interpreter in file mode if arguments are given, otherwise as a repl.
func Apl(a *apl.Apl, stdin io.Reader, args []string) error {
	// Execute files.
	if len(args) > 0 {
		for _, name := range args {
			var r io.Reader
			if name == "-" {
				r = stdin
				name = "stdin"
			} else {
				f, err := os.Open(name)
				if err != nil {
					return err
				}
				defer f.Close()
				r = f
			}
			return a.EvalFile(r, name)
		}
		return nil
	}

	// Run interactively.
	scanner := bufio.NewScanner(stdin)
	fmt.Printf("        ")
	for scanner.Scan() {
		s := scanner.Text()
		if err := a.ParseAndEval(s); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("        ")
	}
	return nil
}
