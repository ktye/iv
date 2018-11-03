package extern

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ktye/iv/iv"
)

// external defines the common part for apl backends, that run as external prococesses.
type external struct {
	program string   // name of the external program
	args    []string // arguments for the external program
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	scanner *bufio.Scanner
	ack     int
}

func (e *external) start() error {
	var err error
	e.cmd = exec.Command(e.program, e.args...)
	e.stdout, err = e.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	e.stderr, err = e.cmd.StderrPipe()
	if err != nil {
		return err
	}
	e.stdin, err = e.cmd.StdinPipe()
	if err != nil {
		return err
	}
	e.scanner = bufio.NewScanner(e.stdout)
	return e.cmd.Start()
}

// send writes s to the program's stdin followed by a special acknowledgement string.
// It returns after it received an acknowledgement.
// The program output is written to w.
func (e *external) send(s string, w io.Writer) error {
	e.ack++
	ack := fmt.Sprintf("\"[ack %d]\"", e.ack)
	fmt.Fprintln(e.stdin, s)
	fmt.Fprintln(e.stdin, ack)

	for e.scanner.Scan() {
		if t := e.scanner.Text(); strings.TrimSpace(t) == ack {
			return nil
		} else {
			//fmt.Println("<recv ", t)
			fmt.Fprintln(w, t)
		}
	}
	if err := e.scanner.Err(); err != nil {
		return err
	}
	return fmt.Errorf("%s did not send ack %d", e.program, e.ack)
}

// evalCondition writes s to the program's input and expects an integer.
// It returns if it is != 0 or an error.
func (e *external) evalCondition(s string) (bool, error) {
	var buf strings.Builder
	if err := e.send(s, &buf); err != nil {
		return false, err
	}
	s = strings.TrimSpace(buf.String())
	if i, err := strconv.Atoi(s); err != nil {
		return false, fmt.Errorf("expected bool, got %s", err)
	} else {
		return i != 0, nil
	}
}

// Float can be used as a Scalar which parses from a string.
type float float64

func parseFloat(s string) (iv.Scalar, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return float(f), nil
}
