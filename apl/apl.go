// Package apl implements an APL interpreter.
package apl

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/ktye/iv/apl/scan"
	// _ "github.com/ktye/iv/apl/funcs" // Register default funcs
)

// New starts a new interpreter.
func New(w io.Writer) *Apl {
	a := Apl{
		stdout:     w,
		env:        newEnv(),
		Origin:     1,
		PP:         0,
		Fmt:        make(map[reflect.Type]string),
		primitives: make(map[Primitive][]PrimitiveHandler),
		operators:  make(map[string][]Operator),
		symbols:    make(map[rune]string),
		pkg:        make(map[string]*env),
	}
	a.parser.a = &a
	return &a
}

// Apl stores the interpreter state.
type Apl struct {
	scan.Scanner
	parser
	stdout     io.Writer
	stdimg     ImageWriter
	Tower      Tower
	Origin     int
	PP         int
	Fmt        map[reflect.Type]string
	env        *env
	primitives map[Primitive][]PrimitiveHandler
	operators  map[string][]Operator
	symbols    map[rune]string
	pkg        map[string]*env
	scaninit   bool
	debug      bool
}

// LoadPkg loads a package from a file.
// It temporarily removes the current environment, executes the package file with EvalFile
// and stores the resulting environment in a package with the name of pkg.
func (a *Apl) LoadPkg(r io.Reader, file string, pkg string) (err error) {
	save := a.env
	a.env = newEnv()
	defer func() {
		a.env = save
	}()

	err = a.EvalFile(r, file)
	if err != nil {
		return err
	}
	a.pkg[pkg] = a.env
	return nil
}

// Parse parses a line of input into the current context.
// It returns a Program which can be evaluated.
func (a *Apl) Parse(line string) (Program, error) {
	tokens, err := a.Scan(line)
	if a.debug {
		fmt.Fprintf(a.stdout, "%s\n", scan.PrintTokens(tokens))
	}

	if err != nil {
		return nil, err
	}

	p, err := a.parse(tokens)
	if a.debug {
		fmt.Fprintf(a.stdout, "%s\n", p.String(a))
	}

	if err != nil {
		return nil, err
	} else {
		return p, nil
	}
}

func (a *Apl) ParseAndEval(line string) error {
	if p, err := a.Parse(line); err != nil {
		return err
	} else {
		return a.Eval(p)
	}
}

func (a *Apl) Scan(line string) ([]scan.Token, error) {
	// On the first call, the scanner needs to be told all symbols that
	// have been registered.
	if a.scaninit == false {
		m := make(map[rune]string)
		for r, s := range a.symbols {
			m[r] = s
		}
		a.SetSymbols(m)
		a.scaninit = true
	}
	return a.Scanner.Scan(line)
}

func (a *Apl) SetDebug(d bool) {
	a.debug = d
}

func (a *Apl) SetOutput(w io.Writer) {
	a.stdout = w
}

func (a *Apl) GetOutput() io.Writer {
	if a.stdout == nil {
		return ioutil.Discard
	}
	return a.stdout
}

func (a *Apl) SetImage(w ImageWriter) {
	a.stdimg = w
}

func newEnv() *env {
	return &env{vars: map[string]Value{}}
}
