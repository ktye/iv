package iv

import (
	"fmt"
	"io"
	"os"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/primitives"
	// TODO "github.com/ktye/iv/apl/operators"
)

// AplIv implements the Apl interface using iv/apl as a backend.
type AplIv struct {
	a      *apl.Apl
	rules  []aplRule
	nr     int64
	stderr io.Writer
}

type aplRule struct {
	cond apl.Program
	expr apl.Program
}

func newAplIv(rank int) (Apl, error) {
	// α is not the APL alpha symbol (⍺), which is not accepted
	// by go as an identifier.
	var α AplIv
	α.a = apl.New(os.Stdout)
	primitives.Register(α.a)
	// TODO operators.Register(α.a)
	return &α, nil
}

func (α *AplIv) SetOut(stdout, stderr io.Writer) {
	α.a.SetOutput(stdout)
	// iv/apl does not write to stderr, it returns
	// errors instead.
	α.stderr = stderr
}

func (α *AplIv) Parse(s, name string) error {
	_, err := α.a.Parse(s)
	return err
}

func (α *AplIv) AddRule(cond, expr string) error {
	if expr == "" {
		return fmt.Errorf("cannot add an empty rule")
	}
	n := len(α.rules) + 1
	var r aplRule
	if cond == "" {
		cond = "1"
	}
	if p, err := α.a.Parse(cond); err != nil {
		return fmt.Errorf("condition %d: %s", n, err)
	} else if len(p) != 1 {
		return fmt.Errorf("condition %d: must be a single expression, not %d", n, len(p))
	} else {
		r.cond = p
	}

	if p, err := α.a.Parse(expr); err != nil {
		return fmt.Errorf("rule %d: %s", n, err)
	} else {
		r.expr = p
	}
	α.rules = append(α.rules, r)
	return nil
}

func (α *AplIv) ParseScalar(s string) (Scalar, error) {
	// Iv/apl accepts both ¯, and -.
	if n, err := apl.ParseNumber(s); err == nil {
		return n, nil
	}

	// If it cannot be parsed as a number, treat it as a string.
	return apl.String(s), nil
}

func (α *AplIv) Execute(dims []int, array []Scalar, term int, eof bool) error {
	// Assign the current array to _.
	ar := apl.GeneralArray{
		Values: make([]apl.Value, len(array)),
		Dims:   dims,
	}
	for i, a := range array {
		ar.Values[i] = a.(apl.Value)
	}
	α.a.Assign("_", ar)

	// Assign the current record number N,
	// termination level E,
	// End of input marker EOF
	α.nr++
	α.a.Assign("N", apl.Int(α.nr))
	α.a.Assign("E", apl.Int(term))
	α.a.Assign("EOF", apl.Bool(eof))

	boolean := func(v apl.Value) (bool, error) {
		if b, ok := v.(apl.Bool); ok {
			return bool(b), nil
		} else if i, ok := v.(apl.Int); ok {
			return i == 1, nil
		}
		return false, fmt.Errorf("expected bool or int, got %T", v)
	}

	// Execute rules sequentially.
	α.a.Assign("NEXT", apl.Bool(false))
	for i, r := range α.rules {
		// Test if rule condition applies.
		// r.cond contains a single apl.expr, verified by AddRule.
		if cond, err := r.cond[0].Eval(α.a); err != nil {
			return fmt.Errorf("rule#%d: condition: %s", i+1, err)
		} else if b, err := boolean(cond); err != nil {
			fmt.Errorf("rule#%d: %s", i+1, err)
		} else if b == false {
			continue
		}

		// Execute the rule.
		if err := α.a.Eval(r.expr); err != nil {
			return fmt.Errorf("rule#%d: %s", i+1, err)
		}

		// Check if NEXT has been assigned.
		if next := α.a.Lookup("NEXT"); next == nil {
			return fmt.Errorf("variable NEXT does not exist")
		} else if b, err := boolean(next); err != nil {
			return fmt.Errorf("NEXT: %s", err)
		} else if b == true {
			break
		}
	}
	return nil
}
