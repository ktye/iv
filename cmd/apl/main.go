// Program apl ins the APL interpreter as a command line program.
//
// Usage
//	apl [-l LIB] < INPUT
//		interprete stdin
//	apl -fmt < INPUT
//		format input, no interpretation
//	apl [-l LIB] -acme < INPUT
//		format input, interprete, ignore all lines not starting with TAB
//	apl [-l LIB] -interactive
//		REPL using liner for console handling, and APL character input
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	// TODO "github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/complete"
)

type options struct {
	format      bool
	interactive bool
	acme        bool
	debug       bool
	lib         string
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	state       *apl.Apl
}

func main() {
	var opt options
	flag.BoolVar(&opt.format, "fmt", false, "replace APL symbol names in input")
	flag.BoolVar(&opt.interactive, "interactive", false, "use liner for console handling")
	flag.BoolVar(&opt.acme, "acme", false, "ignore all input except data following a tab until newline")
	flag.BoolVar(&opt.debug, "debug", false, "print debugging information")
	flag.StringVar(&opt.lib, "lib", "", "library to load on startup")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "too many command line arguments")
		os.Exit(1)
	}

	opt.stdin = os.Stdin
	opt.stdout = os.Stdout
	opt.stderr = os.Stderr

	if err := opt.run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func (opt *options) run() error {

	// Format only reformats all input, but does not evaluate.
	if opt.format {
		scanner := bufio.NewScanner(opt.stdin)
		for scanner.Scan() {
			s := scanner.Text()
			if opt.format {
				s = complete.Format(s)
			}
			fmt.Fprintln(opt.stdout, s)
		}
		return nil
	}

	opt.state = apl.New(opt.stdout)
	primitives.Register(opt.state)
	operators.Register(opt.state)
	opt.state.SetDebug(opt.debug)

	// Load library.
	if opt.lib != "" {
		if f, err := os.Open(opt.lib); err != nil {
			fmt.Fprintf(opt.stderr, "%s\n", err)
		} else {
			defer f.Close()

			line := 0
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				s := scanner.Text()
				line++
				if err := opt.state.ParseAndEval(s); err != nil {
					fmt.Fprintf(opt.stderr, "%s:%d: %s\n", opt.lib, line, err)
					os.Exit(1)
				}
			}
		}
	}

	// In interactive mode, input handling is done by liner.
	// Interpretion is done line by line.
	if opt.interactive {
		return opt.repl()
	}

	// Acme mode iterates over all lines and ignores all lines that
	// dont start with a tab.
	// Lines are formated, printed and evaluated.
	if opt.acme {
		opt.format = true // Acme always reformats.
		opt.stdin = &tabfilter{r: opt.stdin}
	}

	line := 0
	scanner := bufio.NewScanner(opt.stdin)
	for scanner.Scan() {
		s := scanner.Text()
		line++
		if opt.format {
			s = complete.Format(s)
		}
		if err := opt.state.ParseAndEval(s); err != nil {
			// In acme mode, the line number only counts lines
			// with leading tab.
			fmt.Fprintf(opt.stderr, "<stdin>:%d: %s\n", line, err)
			os.Exit(1)
		}
	}
	return nil
}
