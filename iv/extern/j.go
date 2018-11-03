package extern

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ktye/iv/iv"
)

// J is the j backend to iv.
// 	http://www.jsoftware.com
// The interface runs j as an external process.
// The binary jconsole must be on the path.
//
// Syntax:
//	A J rule is a verb which is applied to the reshaped input vector.
// Iv assigns the variables
//	"X" current array (ndimensional array)
//	"N" record number (int)
//	"E" termination level (int)
// J may assign to
//	"NEXT" a boolean or integer which is not zero to skip the rest of the rules.
//
// Example
//	echo 1 -2 3 | iv -aj '(+/%#)'
//	   0.666667

type J struct {
	external
	rank  int
	ack   int
	out   io.Writer
	nr    int
	rules []jRule
}

func newJ(rank int) (iv.Apl, error) {
	j := J{
		external: external{program: "jconsole", args: []string{"-jprofile"}},
		rank:     rank,
	}
	if err := j.start(); err != nil {
		return nil, err
	}
	return &j, nil
}

func (j *J) SetOut(stdout, stderr io.Writer) {
	j.out = stdout
	go func() {
		io.Copy(stderr, j.external.stderr)
	}()
}

// Parse writes the commands in s to k and prints the results.
func (j *J) Parse(s, name string) error {
	return j.send(s, j.out)
}

// AddRule stores the rules but does not send them to j.
func (j *J) AddRule(cond string, expr string) error {
	j.rules = append(j.rules, jRule{cond, expr})
	return nil
}

// ParseScalar parses float values, but stores them as a string replacing - by _.
func (j *J) ParseScalar(s string) (iv.Scalar, error) {
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return nil, err
	} else {
		if s[0] == '-' {
			return "_" + s[1:], nil
		}
	}
	return s, nil
}

func (j *J) Execute(shape []int, vector []iv.Scalar, term int, eof bool) error {
	j.nr++
	var buf strings.Builder
	// Assign and reshape vector to X.
	fmt.Fprintf(&buf, "X =: ")
	for _, v := range shape {
		fmt.Fprintf(&buf, "%v ", v)
	}
	fmt.Fprintf(&buf, "$")
	for _, v := range vector {
		fmt.Fprintf(&buf, "%v ", v)
	}
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "N =: %d\n", j.nr) // Assign the record number to N
	fmt.Fprintf(&buf, "E =: %d\n", term) // Assing termination level to E
	eofnum := 0
	if eof {
		eofnum = 1
	}
	fmt.Fprintf(&buf, "EOF =: %d\n", eofnum) // Assign EOF.
	fmt.Fprintf(&buf, "NEXT =: 0\n")         // Reset next.
	j.send(buf.String(), ioutil.Discard)     // Send new state to klong.

	for _, r := range j.rules {
		// Check condition.
		if res, err := j.evalCondition(r.cond); err != nil {
			return err
		} else if res == false {
			continue
		}

		// Execute rule expression.
		if err := j.send(r.verb+" X", j.out); err != nil {
			return err
		}

		// Check value of NEXT.
		if next, err := j.evalCondition("NEXT"); err != nil {
			return err
		} else if next {
			continue
		}
	}
	return nil
}

func (j *J) send(s string, w io.Writer) error {
	j.ack++
	ack := fmt.Sprintf("'[ack %d]'", j.ack)
	fmt.Fprintln(j.stdin, s)
	fmt.Fprintln(j.stdin, ack)

	ack = ack[1 : len(ack)-1] // strip '' for the response.
	for j.scanner.Scan() {
		if t := j.scanner.Text(); strings.TrimSpace(t) == ack {
			return nil
		} else {
			fmt.Fprintln(w, t)
		}
	}
	if err := j.scanner.Err(); err != nil {
		return err
	}
	return fmt.Errorf("%s did not send ack %d", j.program, j.ack)
}

// evalCondition writes s to the program's input and expects an integer.
// It returns if it is != 0 or an error.
func (j *J) evalCondition(s string) (bool, error) {
	var buf strings.Builder
	if err := j.send(s, &buf); err != nil {
		return false, err
	}
	s = strings.TrimSpace(buf.String())
	if i, err := strconv.Atoi(s); err != nil {
		return false, fmt.Errorf("expected bool, got %s", err)
	} else {
		return i != 0, nil
	}
}

type jRule struct {
	cond string
	verb string
}

func init() {
	iv.Register("j", newJ)
	iv.Register("J", newJ)
}
