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

// read reads from a file
func read(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	name, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("io read: expect file name %T", R)
	}
	f, err := Open(string(name))
	if err != nil {
		return nil, err
	}
	return apl.LineReader(f), nil // LineReader closes the file.
}

// readfd reads from a file descriptor (Index). Only 0 is allowed.
func readfd(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	fd := int(R.(apl.Index))
	if fd != 0 {
		return nil, fmt.Errorf("io readfd: argument must be 0 (stdin)")
	}
	return apl.LineReader(os.Stdin), nil
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
		fsys, err := lookup(argv[0])
		if err != nil {
			return nil, err
		}
		if f, ok := fsys.(fs); ok == false {
			return nil, fmt.Errorf("exec: %s: file system is not an os fs: %s", argv[0], fsys.String())
		} else {
			argv[0] = f.path(argv[0])
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

// load reads the file R and executes it.
// It returns an error, if R is not a file.
func load(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	s, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("io l: argument must be a file name: %T", R)
	}

	f, err := Open(string(s))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := a.EvalFile(f, string(s)); err != nil {
		return nil, err
	}
	return apl.EmptyArray{}, nil
}

func lCmd(t []scan.Token) []scan.Token {
	l := scan.Token{T: scan.Identifier, S: "io→l"}
	return append([]scan.Token{l}, t...)
}
