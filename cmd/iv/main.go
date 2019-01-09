// iv stream processor
package main

import (
	"fmt"
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

	var rank int = 1
	var begin string
	var λ string
	var err error

	// If an argument starts with -r or -b, it's a flag.
	// There is not space between the flag and the option.
	args := os.Args[1:]
	for len(args) > 0 {
		s := args[0]
		if strings.HasPrefix(s, "-r") != -1 {
			rank, err = strconv.Atoi(s[2:])
			fatal(err)
		} else if strings.HasPrefix(s, "-b") {
			begin = strconv.Atoi(s[2:])
		} else {
			break
		}
		args = args[1:]
	}
	λ = strings.Join(args, " ")
	if len(λ) == 0 {
		fatal(fmt.Errorf("command line lambda argument is missing"))
	}

	// Start an interpreter instance with the iv extension package.
	a := apl.New(os.Stdout)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)
	iv.Register(a)

	if begin != "" {
		fatal(a.ParseAndEval(begin))
	}
	fatal(a.ParseAndEval("iv←{" + λ + "}"))
	fatal(a.ParseAndEval(fmt.Printf("IvC←iv→r %d", rank)))
	fatal(a.ParseAndEval("{(¯1↑⍵) iv 1↑⍵}¨IvC"))
}

func fatal(err error) {
	if err != nil {
		return fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`
Usage
	iv -rRANK -bBEGIN lambda function < data
Example
	iv -r2 '⍵,+/⍵'
	`)
	os.Exit(1)
}
