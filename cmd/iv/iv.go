// Program iv is a command line interface to the iv stream processor.
//
// The program reads data from a text stream on stdin
// and executes the given iv rules.
//
// It supports multiple APL backends: iv/apl, ivy, k (kona), j and klong.
// Kona, j and klong need to be installed as external binaries.
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ktye/iv/complete"
	"github.com/ktye/iv/iv"
	_ "github.com/ktye/iv/iv/extern" // load all backends
)

const usage = `
Usage
	iv [OPTIONS] RULE1 RULE2 ... < input
	iv -f file < input
Options
	-f FILE   read program from iv file
	-r RANK   execution rank (default 1)
                  call APL with current array for each:
                  0: scalar, 1: vector, 2: matrix, ...
        -l LIB    execute library file before start
	-b BEGIN  execute BEGIN statements at start
	-u        check hard uniformity of input data (default: unset)
		  subsequent arrays of rank r must have the same shape
	-n        don't use convenience input wrapper,
                  which replaces whitespace with a single newline
	-a PROG   use PROG as the APL backend (default iv/apl)
`

func main() {
	var file string
	var rank int = 1
	var lib string
	var begin string
	var uniform bool
	var linemode bool = true
	var backend string
	var err error

	// See iv/complete/bash.go on how to install bash completion for iv.
	if len(os.Args) > 1 && os.Args[1] == "-complete-bash" {
		complete.Bash(os.Args[2:])
		os.Exit(0)
	}

	options, args := getOptions("hhfrlbuna", "x-xxxx--x")

	if _, ok := options["-h"]; ok {
		fmt.Println(usage)
		os.Exit(0)
	}

	// Set options given on the command line.
	file = options["-f"]
	if s := options["-r"]; s != "" {
		rank, err = strconv.Atoi(s)
		if err != nil {
			fatalf("-r=rank: %s\n", err)
		}
	}
	lib = options["-l"]
	begin = options["-b"]
	_, uniform = options["-u"]
	if _, ok := options["-n"]; ok {
		linemode = false
	}
	backend = options["-a"]

	// Parse the program from a file, f may be nil.
	var f *iv.File
	if file != "" {
		if r, err := os.Open(file); err != nil {
			fatalf("%s\n", err)
		} else {
			defer r.Close()
			f, err = iv.ParseFile(r, file, rank, uniform, begin)
		}
	}

	// Parse remaining options as rules from the command line.
	if f != nil && len(args) > 0 {
		// While we allow -rRANK, -u and -bBEGIN on the command line
		// mixed with a program file, we do not allow additional rules.
		fatalf("-f FILE is given with additional rules on the command line")
	} else if f == nil {
		f = &iv.File{
			Rank:    rank,
			Begin:   begin,
			Uniform: uniform,
		}
	}
	for _, s := range args {
		var r [2]string
		if idx := strings.Index(s, ":"); idx == -1 {
			r = [2]string{"1", s}
		} else {
			r[0] = s[:idx]
			r[1] = s[idx+1:]
		}
		f.Rules = append(f.Rules, r)
	}

	// We normalize text input by default, unless -n is given.
	var r *bufio.Reader
	if linemode {
		r = bufio.NewReader(iv.TabularText(os.Stdin))
	} else {
		r = bufio.NewReader(os.Stdin)
	}

	// Start the apl interpreter.
	v, err := iv.New(f.Rank, f.Uniform, backend)
	if err != nil {
		fatalf("%s\n", err)
	}

	// Connect outputs.
	v.SetOut(os.Stdout, os.Stderr)

	// Run lib file.
	if lib != "" {
		if b, err := ioutil.ReadFile(lib); err != nil {
			fatalf("%s\n", err)
		} else {
			if err := v.Parse(string(b), lib); err != nil {
				fatalf("%s\n", err)
			}
		}
	}

	// Run begin block.
	if f.Begin != "" {
		if err := v.Parse(f.Begin, "<BEGIN>"); err != nil {
			fatalf("%s\n", err)
		}
	}

	// Add rules.
	for _, r := range f.Rules {
		if err := v.AddRule(r[0], r[1]); err != nil {
			fatalf("%s\n", err)
		}
	}

	// Read data as text from stdin.
	next := iv.TextStream{
		Reader:    r,
		Parse:     v.TextParser(),
		Separator: '\n',
		Rank:      rank,
	}
	v.SetNext(&next)

	// Let's go.
	if err := v.Run(); err != nil {
		fatalf("%s\n", err)
	}
}

// We use a custom flag parser because we want a short command line, e.g. -r2 -bBEGIN.
func getOptions(pattern, opts string) (map[string]string, []string) {
	args := os.Args[1:]
	getarg := func() (string, string) {
		if len(args) == 0 {
			return "", ""
		}
		a := args[0]
		for i, p := range pattern {
			prefix := "-" + string(p)
			if strings.HasPrefix(a, prefix) {
				a = strings.TrimPrefix(a, prefix)
				if opts[i] == '-' && len(a) != 0 {
					fatalf("argument %s does not take an option", prefix)
				} else if opts[i] == '-' {
					args = args[1:]
					return prefix, ""
				} else if opts[i] == 'x' && len(a) != 0 {
					// -xY or -x=Y
					if a[0] == '=' {
						a = a[1:]
					}
					args = args[1:]
					return prefix, a
				} else if opts[i] == 'x' && len(a) == 0 {
					// -x Y
					args = args[1:]
					if len(args) == 0 {
						fatalf("argument %s requires an option", prefix)
					}
					o := args[0]
					args = args[1:]
					return prefix, o
				}
			}
		}
		return "", ""
	}

	options := make(map[string]string)
	for {
		if a, o := getarg(); a == "" {
			break
		} else {
			options[a] = o
		}
	}
	return options, args
}

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
