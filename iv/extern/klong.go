package extern

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ktye/iv/iv"
)

// Klong is the klong backend to iv.
// 	t3x.org/klong
// The interface runs klong as an external process.
// The binary kg must be on the path.
//
// Syntax:
// Iv assigns the variables
//	"." current array (ndimensional array)
//	"N" record number (int)
//	"E" termination level (int)
//	We cannot use "_" as in ivy, as this is not allowed as an identifier.
// Klong may assign to
//	"NEXT" a boolean or integer which is not zero to skip the rest of the rules.
//
// Example
// 	echo 1 2 3 | iv -akg '.,(+/.)%#.'
//	[1 2 3 2.0]
type Klong struct {
	external
	rank  int
	ack   int
	out   io.Writer
	nr    int
	rules int
}

func newKlong(rank int) (iv.Apl, error) {
	kg := Klong{
		external: external{program: "kg", args: []string{"-q"}},
		rank:     rank,
	}
	if err := kg.start(); err != nil {
		return nil, err
	}
	return &kg, nil
}

func (kg *Klong) SetOut(stdout, stderr io.Writer) {
	kg.out = stdout
	go func() {
		io.Copy(stderr, kg.external.stderr)
	}()
}

// Parse writes the commands in s to klong and prints the results.
func (kg *Klong) Parse(s, name string) error {
	return kg.send(s, kg.out)
}

// AddRule writes the rule to klong and discards the output.
func (kg *Klong) AddRule(cond string, expr string) error {
	rl := fmt.Sprintf("rule%dcond::{%s}", kg.rules, cond)
	if err := kg.send(rl, ioutil.Discard); err != nil {
		return err
	}
	ex := fmt.Sprintf("rule%dex::{%s}", kg.rules, expr)
	if err := kg.send(ex, ioutil.Discard); err != nil {
		return err
	}
	kg.rules++
	return nil
}

func (kg *Klong) ParseScalar(s string) (iv.Scalar, error) {
	return parseFloat(s)
}

func (kg *Klong) Execute(shape []int, vector []iv.Scalar, term int, eof bool) error {
	kg.nr++
	var buf strings.Builder
	fmt.Fprintf(&buf, ". :: %v :^ %v\n", shape, vector) // Assign the reshape of vector to .
	fmt.Fprintf(&buf, "N :: %d\n", kg.nr)               // Assign the record number to N
	fmt.Fprintf(&buf, "E :: %d\n", term)                // Assing termination level to E
	eofnum := 0
	if eof {
		eofnum = 1
	}
	fmt.Fprintf(&buf, "EOF :: %d\n", eofnum) // Assign EOF.
	fmt.Fprintf(&buf, "NEXT :: 0\n")         // Reset next.
	kg.send(buf.String(), ioutil.Discard)    // Send new state to klong.

	for i := 0; i < kg.rules; i++ {
		rl := fmt.Sprintf("rule%dcond()", i)
		ex := fmt.Sprintf("rule%dex()", i)

		// Check condition.
		if res, err := kg.evalCondition(rl); err != nil {
			return err
		} else if res == false {
			continue
		}

		// Execute rule expression.
		if err := kg.send(ex, kg.out); err != nil {
			return err
		}

		// Check value of NEXT.
		if next, err := kg.evalCondition("NEXT"); err != nil {
			return err
		} else if next {
			continue
		}

		// We ignore the value of E.
		// Klong can always exit itself.
	}
	return nil
}

func init() {
	iv.Register("kg", newKlong)
	iv.Register("klong", newKlong)
}
