// iv stream processor
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	"github.com/ktye/iv/cmd/iv/iv"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	iv.Stdin = os.Stdin
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(w io.Writer, args []string) error {
	var rank int = 1
	var begin string
	var end string
	var λ string
	var err error

	// If an argument starts with -r, -b, -e it's a flag.
	// There is not space between the flag and the option.
	for len(args) > 0 {
		s := args[0]
		if strings.HasPrefix(s, "-r") {
			rank, err = strconv.Atoi(s[2:])
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(s, "-b") {
			begin = s[2:]
		} else if strings.HasPrefix(s, "-e") {
			end = s[2:]
		} else {
			break
		}
		args = args[1:]
	}
	λ = strings.Join(args, " ")
	if len(λ) == 0 {
		return fmt.Errorf("command line lambda argument is missing")
	}

	// Start an interpreter instance with the iv extension package.
	a := apl.New(w)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)
	iv.Register(a)

	program := []string{
		begin,
		"iv←{" + λ + "}",
		fmt.Sprintf("IvC←iv→r %d", rank),
		"IvN←0",
		"IvS←{IvN=0:⎕←(¯1↑⍺) iv 1↑⍺⋄⎕←(¯1↑⍵) iv 1↑⍵}/IvC",
		end,
	}

	for _, line := range program {
		if line != "" {
			if err := a.ParseAndEval(line); err != nil {
				return err
			}
		}
	}
	return nil
}

func usage() {
	fmt.Println(`
Usage
	iv -rRANK -bBEGIN -eEND lambda function < data
Example
	iv '+/⍵' < data
	iv -r2 '⍵,+/⍵' < data
	`)
	os.Exit(1)
}
