package extern

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ktye/iv/iv"
)

// Kona is the kona backend to iv.
// 	github.com/kevinlawler/kona
// The interface runs kona as an external process.
// The binary k must be on the path.
//
// Syntax:
// Iv assigns the variables
//	"X" current array (ndimensional array)
//	"N" record number (int)
//	"E" termination level (int)
//	We cannot use "_" as in ivy, as this is not allowed as an identifier.
// K may assign to
//	"NEXT" a boolean or integer which is not zero to skip the rest of the rules.
//
// Example
// 	echo 1 2 3 | iv -ak X,+/X
//	  1 2 3 6
type Kona struct {
	external
	rank  int
	ack   int
	out   io.Writer
	nr    int
	rules int
}

func newKona(rank int) (iv.Apl, error) {
	k := Kona{
		external: external{program: "k"},
		rank:     rank,
	}
	if err := k.start(); err != nil {
		return nil, err
	}
	return &k, nil
}

func (k *Kona) SetOut(stdout, stderr io.Writer) {
	k.out = stdout
	go func() {
		io.Copy(stderr, k.external.stderr)
	}()
}

// Parse writes the commands in s to k and prints the results.
func (k *Kona) Parse(s, name string) error {
	return k.send(s, k.out)
}

// AddRule writes the rule to k and discards the output.
func (k *Kona) AddRule(cond string, expr string) error {
	rl := fmt.Sprintf("rule%dcond:{%s}", k.rules, cond)
	if err := k.send(rl, ioutil.Discard); err != nil {
		return err
	}
	ex := fmt.Sprintf("rule%dex:{%s}", k.rules, expr)
	if err := k.send(ex, ioutil.Discard); err != nil {
		return err
	}
	k.rules++
	return nil
}

func (k *Kona) ParseScalar(s string) (iv.Scalar, error) {
	return parseFloat(s)
}

func (k *Kona) Execute(shape []int, vector []iv.Scalar, term int, eof bool) error {
	k.nr++
	var buf strings.Builder
	// Assign and reshape vector to X.
	fmt.Fprintf(&buf, "X : ")
	for _, v := range shape {
		fmt.Fprintf(&buf, "%v ", v)
	}
	fmt.Fprintf(&buf, "#")
	for _, v := range vector {
		fmt.Fprintf(&buf, "%v ", v)
	}
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "N : %d\n", k.nr) // Assign the record number to N
	fmt.Fprintf(&buf, "E : %d\n", term) // Assing termination level to E
	eofnum := 0
	if eof {
		eofnum = 1
	}
	fmt.Fprintf(&buf, "EOF : %d\n", eofnum) // Assign EOF.
	fmt.Fprintf(&buf, "NEXT : 0\n")         // Reset next.
	k.send(buf.String(), ioutil.Discard)    // Send new state to klong.

	for i := 0; i < k.rules; i++ {
		rl := fmt.Sprintf("rule%dcond()", i)
		ex := fmt.Sprintf("rule%dex()", i)

		// Check condition.
		if res, err := k.evalCondition(rl); err != nil {
			return err
		} else if res == false {
			continue
		}

		// Execute rule expression.
		if err := k.send(ex, k.out); err != nil {
			return err
		}

		// Check value of NEXT.
		if next, err := k.evalCondition("NEXT"); err != nil {
			return err
		} else if next {
			continue
		}
	}
	return nil
}

func init() {
	iv.Register("k", newKona)
	iv.Register("kona", newKona)
}
