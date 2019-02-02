package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ktye/iv/apl"
)

func run(a *apl.Apl, r io.Reader, file string) (err error) {
	line := 0
	defer func() {
		if err != nil {
			err = fileError{file: file, line: line, err: err}
		}
	}()

	ok := true
	var p apl.Program
	b := apl.NewLineBuffer(a)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line++
		ok, err = b.Add(scanner.Text())
		if err != nil {
			return
		}

		if ok {
			p, err = b.Parse()
			if err != nil {
				return
			}

			err = a.Eval(p)
			if err != nil {
				return
			}
		}
	}
	if ok == false && b.Len() > 0 {
		return fmt.Errorf("multiline statement is not terminated")
	}
	return nil
}

type fileError struct {
	file string
	line int
	err  error
}

func (f fileError) Error() string {
	return fmt.Sprintf("%s:%d: %s", f.file, f.line, f.err.Error())
}
