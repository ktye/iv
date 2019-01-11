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
)

// read reads from a file
func read(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	name, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("io read: expect file name %T", R)
	}
	f, err := os.Open(string(name))
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
