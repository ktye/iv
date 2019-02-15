package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
	ex "os/exec"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/scan"
)

// Stdin is exported to be overwritable by tests.
var Stdin io.ReadCloser = os.Stdin

// read reads from a file
func read(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	name, ok := R.(apl.String)
	if ok == false {
		// If R is 0, it reads from stdin.
		if num, ok := R.(apl.Number); ok {
			if n, ok := num.ToIndex(); ok && n == 0 {
				return apl.LineReader(Stdin), nil
			}
		}
		return nil, fmt.Errorf("io read: expect file name %T", R)
	}
	f, err := Open(string(name))
	if err != nil {
		return nil, err
	}
	return apl.LineReader(f), nil // LineReader closes the file.
}

// exec executes a program and sends the output through a channel.
// If called dyadically it uses R as an input, that can be a channel or a Value.
// If the program starts with a slash, it's location is looked up in the file system.
// TODO: should all arguments starting with a slash be replaced?
func exec(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	r := R
	var in io.Reader
	if L != nil {
		r = L
		c, ok := R.(apl.Channel)
		if ok {
			in = bufio.NewReader(apl.NewChannelReader(a, c))
		} else {
			in = strings.NewReader(R.String(a))
		}
	}

	v, ok := domain.ToStringArray(nil).To(a, r)
	if ok == false {
		return nil, fmt.Errorf("io exec: argv must be strings: %T", r)
	}
	argv := v.(apl.StringArray).Strings
	if len(argv) == 0 {
		return nil, fmt.Errorf("io exec: argv empty")
	}

	// If the command starts with a slash, we may relocate it.
	if strings.HasPrefix(argv[0], "/") {
		fsys, mpt, err := lookup(argv[0])
		if err != nil {
			return nil, err
		}
		if f, ok := fsys.(fs); ok == false {
			return nil, fmt.Errorf("exec: %s: file system is not an os fs: %s", argv[0], fsys.String())
		} else {
			relpath := strings.TrimPrefix(argv[0], mpt)
			argv[0] = f.path(relpath)
		}
	}

	cmd := ex.Command(argv[0], argv[1:]...)
	cmd.Stdin = in
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c := apl.LineReader(out)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return c, nil
}

// Load reads the file R and executes it.
// It returns an error, if R is not a file.
// If L is given, the file is executed in a new environment and the resulting variables
// are copied to a package with the name of L.
func load(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	s, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("io l: argument must be a file name: %T", R)
	}

	f, err := Open(string(s))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if L == nil {
		if err := a.EvalFile(f, string(s)); err != nil {
			return nil, err
		}
		return apl.EmptyArray{}, nil
	}

	l, ok := L.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("io l: left argument must be a package name: %T", L)
	}

	if err := a.LoadPkg(f, string(s), string(l)); err != nil {
		return nil, err
	}
	return apl.EmptyArray{}, nil
}

func lCmd(t []scan.Token) []scan.Token {
	// /l `/file.apl `pkg
	if len(t) == 2 && t[0].T == scan.String && t[1].T == scan.String {
		return []scan.Token{t[1], scan.Token{T: scan.Identifier, S: "io→l"}, t[0]}
	}
	l := scan.Token{T: scan.Identifier, S: "io→l"}
	return append([]scan.Token{l}, t...)
}
