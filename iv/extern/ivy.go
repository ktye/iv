// Extern contains different backends for iv.
// To activate them, blank import this package in the main application.
package extern

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/ktye/iv/iv"

	"robpike.io/ivy/config"
	"robpike.io/ivy/exec"
	"robpike.io/ivy/parse"
	"robpike.io/ivy/run"
	"robpike.io/ivy/scan"
	"robpike.io/ivy/value"
)

// Ivy is an implementation of the Apl interface using the robpike.io/ivy as a backend.
type Ivy struct {
	conf    config.Config
	context value.Context
	rules   []ivyRule
	nr      int64
}

func newIvy(rank int) (iv.Apl, error) {
	var rob Ivy
	rob.conf = config.Config{}
	rob.context = exec.NewContext(&rob.conf)
	return &rob, nil
}

func (rob *Ivy) Parse(s, name string) error {
	scanner := scan.New(rob.context, name, translateIvy(strings.NewReader(s)))
	parser := parse.NewParser(name, scanner, rob.context)

	for !run.Run(parser, rob.context, false) {
	}
	return nil // run does not return any errors
}

func (rob *Ivy) AddRule(conditional, statement string) error {
	if statement == "" {
		return fmt.Errorf("cannot add an empty rule")
	}

	n := len(rob.rules) + 1
	if conditional == "" {
		conditional = "1"
	}

	name := fmt.Sprintf("<rule#%d>", n)
	var r ivyRule
	if exprs := rob.parse(conditional, name); len(exprs) < 1 {
		return fmt.Errorf("%s conditional has no expressions", name)
	} else {
		r.conditional = exprs
	}

	r.statements = rob.parse(statement, name)
	rob.rules = append(rob.rules, r)
	return nil
}

func (rob *Ivy) parse(s, name string) []value.Expr {
	scanner := scan.New(rob.context, name, translateIvy(strings.NewReader(s)))
	parser := parse.NewParser(name, scanner, rob.context)
	var exprs []value.Expr
	for {
		if e, ok := parser.Line(); ok == false {
			return exprs
		} else if len(e) > 0 {
			exprs = append(exprs, e...)
		}
	}
}

// ParseScalar currently parses integers as ivy.Int and floats as ivy.BigFloat.
// We could also support strings or other number types.
func (rob *Ivy) ParseScalar(s string) (iv.Scalar, error) {
	if i, err := strconv.Atoi(s); err == nil {
		return value.Int(i), nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return value.BigFloat{big.NewFloat(f)}, nil
	}
	return nil, fmt.Errorf("cannot parse scalar: %s", s)
}

// Execute executes APL rules for the current input array.
// The array is assigned to _ and reshaped according to dims.
// Before calling the rules, it sets E to term and EOF.
func (rob *Ivy) Execute(dims []int, array []iv.Scalar, term int, eof bool) error {
	current := make(value.Vector, len(array))
	for i, v := range array {
		current[i] = v.(value.Value)
	}

	shape := make(value.Vector, len(dims))
	for i, v := range dims {
		shape[i] = value.Int(int64(v))
	}

	// Assign the current record number to N.
	rob.nr++
	rob.context.Assign("N", value.Int(rob.nr))

	// Assign the termination level to E.
	rob.context.Assign("E", value.Int(term))

	// Reshape the current data array and assign it to _.
	ndarray := rob.context.EvalBinary(shape, "rho", current)
	rob.context.Assign("_", ndarray)

	// Assign the variable "EOF", that it can be used in a rule.
	if eof {
		rob.context.Assign("EOF", value.Int(1))
	}

	// Execute all rules sequentially.
	rob.context.Assign("NEXT", value.Int(0))
	for i, r := range rob.rules {
		if e := rob.context.Eval(r.conditional); len(e) != 1 {
			return fmt.Errorf("rule#%d: conditional returns %d values instead of 1", i+1, len(e))
		} else if b, ok := e[0].(value.Int); ok == false {
			return fmt.Errorf("rule#%d: conditional is not a boolean", i+1)
		} else if b != 0 {
			if e := rob.context.Eval(r.statements); len(e) > 0 {
				rob.print(e)
			}
			if next := rob.context.Lookup("NEXT"); next == nil {
				return fmt.Errorf("variable 'NEXT' does not exist")
			} else if b, ok := next.(value.Int); ok && b != 0 {
				rob.context.Assign("NEXT", value.Int(0))
				break
			}
		} else {
			fmt.Printf("rule#%d: conditional is false", i+1)
		}
	}
	return nil
}

func (rob *Ivy) print(values []value.Value) {
	writer := rob.conf.Output()

	if len(values) == 0 {
		return
	}
	if rob.conf.Debug("types") {
		for i, v := range values {
			if i > 0 {
				fmt.Fprint(writer, ",")
			}
			fmt.Fprintf(writer, "%T", v)
		}
		fmt.Fprintln(writer)
	}
	printed := false
	for _, v := range values {
		if _, ok := v.(parse.Assignment); ok {
			continue
		}
		s := v.Sprint(&rob.conf)
		if printed && len(s) > 0 && s[len(s)-1] != '\n' {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, s)
		printed = true
	}
	if printed {
		fmt.Fprintln(writer)
	}
}

func (rob *Ivy) SetOut(stdout, stderr io.Writer) {
	rob.conf.SetOutput(stdout)
	rob.conf.SetErrOutput(stderr)
}

type ivyRule struct {
	conditional []value.Expr
	statements  []value.Expr
}

// translateIvy returns a ByteReader which replaces APL runes with ivy syntax,
// when it is safe to do so.
func translateIvy(r io.RuneReader) *byteReader {
	return &byteReader{
		r:   r,
		buf: bytes.NewBuffer(nil),
	}
}

// ByteReader translates apl runs from the underlying RuneReader reader.
type byteReader struct {
	r   io.RuneReader
	buf *bytes.Buffer
}

// ReadByte returns the next translated byte from the ByteReader.
func (b *byteReader) ReadByte() (byte, error) {
	if b.buf.Len() == 0 {
		if r, _, err := b.r.ReadRune(); err != nil {
			return 0, err
		} else {
			if s, ok := replacers[r]; ok {
				b.buf.WriteString(s)
			} else {
				b.buf.WriteRune(r)
			}
		}
	}
	return b.buf.ReadByte()
}

// Replacers is a list of unique replacements from apl to ivy syntax.
var replacers = map[rune]string{
	'¯': "-",
	'?': "?",
	'⌈': " ceil ",
	'⌊': " floor ",
	'⍴': " rho ",
	'∼': " not ",
	'⍳': " iota ",
	'⋆': "**",
	'÷': "/",
	'⌹': " solve ",
	'⍟': " log ",
	'⌽': " rot ",
	'⊖': " flip ",
	'⍋': " up ",
	'⍒': " down ",
	'⍎': " ivy ",
	'⍕': " text ",
	// '⍉': "", does not exist in ivy
	// '!': "", factorial, does not exist in ivy
	'√': " sqrt ",
	'∈': " in ",
	'↑': " take ",
	'↓': " drop ",
	'⊥': " decode ",
	'⊤': " encode ",
	// '\\': " fill ", // this replaces backslash in strings: "\n"
	// '/': "sel",
	'≤': "<=",
	// '=': "==", // apl comparison (=) is not compatible with ivy.
	'≥': " >= ",
	'≠': " != ",
	'≪': "<<", // iv addition
	'≫': ">>", // iv addition
	'∨': " or ",
	'∧': " and ",
	'⍱': " nor ",
	'⍲': " nand ",
	'←': "=",
	// '∣': " abs ", not unique
	// '∣':  " mod ", not unique
	// '−': "-", negation cannot be replaced, otherwise −(1 2 3) formats as -1 2 3
	// '×': "sgn", not unique
	// '×': "*", // not unique
}

func init() {
	iv.Register("ivy", newIvy)
	iv.Register("y", newIvy)
}
