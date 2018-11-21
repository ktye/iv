package primitives

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
)

//go:generate go run gen.go

type formatmap map[reflect.Type]string

var format5g formatmap = map[reflect.Type]string{
	reflect.TypeOf(numbers.Float(0)): "%.5g",
}
var formatJ formatmap = map[reflect.Type]string{
	reflect.TypeOf(numbers.Complex(0)): "%vJ%v",
}
var formatJR5 formatmap = map[reflect.Type]string{
	reflect.TypeOf(numbers.Float(0)):   "%.5g",
	reflect.TypeOf(numbers.Complex(0)): "%.5gJ%.5g",
}

var testCases = []struct {
	in, exp string
	formats map[reflect.Type]string
}{

	// Basics numbers and arithmetics.
	{"1", "1", nil},
	{"1+1", "2", nil},
	{"1-2", "¯1", nil}, // negative number
	{"¯1", "¯1", nil},
	{"1-¯2", "3", nil},
	{"1a90", "1a90", nil}, // a complex number
	{"1a60+1a300", "1a0", nil},
	{"1J1", "1.4142135623730951a45", nil},

	// Vectors.
	{"1 2 3", "1 2 3", nil},
	{"1+1 2 3", "2 3 4", nil},
	{"1 2 3+¯1", "0 1 2", nil},
	{"1 2 3+4 5 6", "5 7 9", nil},

	// Braces.
	{"1 2+3 4", "4 6", nil},
	{"1 (2+3) 4", "1 5 4", nil},
	{"(1 2)+3 4", "4 6", nil},
	{"1×2+3×4", "14", nil},
	{"1×(2+3)×4", "20", nil},
	{"(3×2)+3×4", "18", nil},
	{"3×2+3×4", "42", nil},

	// Comparison
	// Comparison tolerance is not implemented.
	{"1 2 3 4 5 > 2", "0 0 1 1 1", nil},     // greater than
	{"1 2 3 4 5 ≥ 3", "0 0 1 1 1", nil},     // greater or equal
	{"2 4 6 8 10<6", "1 1 0 0 0", nil},      // less than
	{"2 4 6 8 10≤6", "1 1 1 0 0", nil},      // less or equal
	{"1 2 3 ≠ 1.1 2 3", "1 0 0", nil},       // not equal
	{"3=3.1 3 ¯2 ¯3 3J0", "0 1 0 0 1", nil}, // equal
	{"2+2=2", "3", nil},                     // calculating with boolean values
	{"2×1 2 3=4 2 1", "0 2 0", nil},         // dyadic array
	{"-3<4", "¯1", nil},                     // monadic scalar
	{"-1 2 3=0 2 3", "0 ¯1 ¯1", nil},        // monadic array

	// Boolean logical
	{"0 1 0 1 ^ 0 0 1 1", "0 0 0 1", nil}, // and
	{"0 1 0 1 ∧ 0 0 1 1", "0 0 0 1", nil}, // accept both ^ and ∧
	{"0^0 0 1 1", "0 0 0 0", nil},         // or
	{"0 0 1 1∨0 1 0 1", "0 1 1 1", nil},   // or
	{"1∨0 1 0 1", "1 1 1 1", nil},         // or
	{"0 0 1 1⍱0 1 0 1", "1 0 0 0", nil},   // nor
	{"0 0 1 1⍲0 1 0 1", "1 1 1 0", nil},   // nand
	// {"15 1 2 7 ^ 35 1 4 0", "105 1 4 0", nil}, // least common multiple
	// {"2 3 4 ∧ 0j1 1j2 2j3", "0J2 3J6 8J12", nil},// least common multiple
	// {"2j2 2j4 ∧ 5j5 4j4", "10J10 ¯4J12", nil},// least common multiple
	// {"15 1 2 7 ∨ 35 1 4 0", "5 1 2 7", nil}, // greatest common divisor

	// Multiple expressions.
	{"1⋄2⋄3", "1", nil},

	// Iota and reshape.
	{"⍳5", "1 2 3 4 5", nil},       // index generation
	{"⍳0", "", nil},                // empty array
	{"⍴⍳5", "5", nil},              // shape
	{"⍴5", "", nil},                // shape of scalar is empty
	{"⍴⍴5", "0", nil},              // shape of empty is 0
	{"⍴⍳0", "0", nil},              // empty array has zero dimensions
	{"⍴⍴⍳0", "1", nil},             // rank of empty array is 1
	{"2 3⍴1", "1 1 1\n1 1 1", nil}, // shape

	// Magnitude, Residue, Ceil, Floor, Min, Max
	{"|1 ¯2 ¯3.2 2.2a20", "1 2 3.2 2.2", nil},                  // magnitude
	{"3 3 ¯3 ¯3|¯5 5 ¯4 4", "1 2 ¯1 ¯2", nil},                  // residue
	{"0.5|3.12 ¯1 ¯0.6", "0.12 0 0.4", format5g},               // residue
	{"¯1 0 1|¯5.25 0 2.41", "¯0.25 0 0.41", format5g},          // residue
	{"1j2|2j3 3j4 5j6", "1J1 ¯1J1 0J1", formatJ},               // complex residue
	{"4J6|7J10", "3J4", formatJ},                               // complex residue
	{"¯10 7J10 .3|17 5 10", "¯3 ¯5J7 0.1", formatJR5},          // residue
	{"⌊¯2.3 0.1 100 3.3", "¯3 0 100 3", nil},                   // floor
	{"⌊0.5 + 0.4 0.5 0.6", "0 1 1", nil},                       // floor
	{"⌊1j3.2 3.3j2.5 ¯3.3j¯2.5", "1J3 3J2 ¯3J¯3", formatJ},     // complex floor
	{"⌊1.5J2.5", "2J2", formatJ},                               // complex floor
	{"⌊1J2 1.2J2.5 ¯1.2J¯2.5", "1J2 1J2 ¯1J¯3", formatJ},       // complex floor
	{"⌈¯2.7 3 .5", "¯2 3 1", nil},                              // ceil
	{"⌈1.5J2.5", "1J3", formatJ},                               // complex ceil
	{"⌈1J2 1.2J2.5 ¯1.2J¯2.5", "1J2 1J3 ¯1J¯2", formatJ},       // complex ceil
	{"⌈¯2.3 0.1 100 3.3", "¯2 1 100 4", nil},                   // ceil
	{"⌈1.2j2.5 1.2j¯2.5", "1J3 1J¯2", formatJ},                 // ceil
	{"5⌊4 5 7", "4 5 5", nil},                                  // min
	{"¯2⌊¯3", "¯3", nil},                                       // min
	{"3.3 0 ¯6.7⌊3.1 ¯4 ¯5", "3.1 ¯4 ¯6.7", nil},               // min
	{"¯2.1 0.1 15.3 ⌊ ¯3.2 1 22", "¯3.2 0.1 15.3", nil},        // min
	{"5⌈4 5 7", "5 5 7", nil},                                  // max
	{"¯2⌈¯3", "¯2", nil},                                       // max
	{"3.3 0 ¯6.7⌈3.1 ¯4 ¯5", "3.3 0 ¯5", nil},                  // max
	{"¯2.01 0.1 15.3 ⌈ ¯3.2 ¯1.1 22.7", "¯2.01 0.1 22.7", nil}, // max

	// Match, Not match, tally, depth
	{"≡5", "0", nil},          // depth
	{"≡⍳0", "1", nil},         // depth for empty array
	{`≡"alpha"`, "0", nil},    // a string is a scalarin APLv.
	{"≢2 3 4⍴⍳10", "2", nil},  // tally
	{"≢2", "1", nil},          // tally
	{"≢⍳0", "0", nil},         // tally
	{"1 2 3≡1 2 3", "1", nil}, // match
	{"3≡1⍴3", "0", nil},       // match shape
	{`""≡⍳0`, "0", nil},       // match empty string
	{`''≡⍳0`, "1", nil},       // this is false in other APLs (here '' is an empty array).
	{"2.0-1.0≡1>0", "1", nil}, // compare numbers of different type
	{"1≢2", "1", nil},         // not match
	{"1≢1", "0", nil},         // not match
	{"3≢1⍴3", "1", nil},       // not match
	{`""≢⍳0`, "1", nil},       // not match

	// Array expressions.
	{"-⍳3", "¯1 ¯2 ¯3", nil},

	// Ravel, enlist, catenate, join
	// TODO ravel with axis
	// TODO laminate
	{",2 3⍴⍳6", "1 2 3 4 5 6", nil},     // ravel
	{"∊2 3⍴⍳6", "1 2 3 4 5 6", nil},     // enlist (identical for simple arrays)
	{"⍴,3", "1", nil},                   // scalar ravel
	{"⍴,⍳0", "0", nil},                  // ravel empty array
	{"1 2 3,4 5 6", "1 2 3 4 5 6", nil}, // catenate
	{`"abc",1 2`, `abc 1 2`, nil},
	{"(2 3⍴⍳6),2 2⍴7 8 9 10", "1 2 3 7 8\n4 5 6 9 10", nil},
	{"2 3≡2,3", "1", nil},                // catenate vector result
	{"(1 2 3,4 5 6)≡⍳6", "1", nil},       // catenate vector result
	{"0,2 3⍴1", "0 1 1 1\n0 1 1 1", nil}, // catenate scalar and array

	// Decode
	{"3⊥1 2 1", "16", nil},
	{"3⊥4 3 2 1", "142", nil},
	{"2⊥1 1 1 1", "15", nil},
	// {"24 60 60⊥2 23 12", "8592", nil}, // mixed radix

	// Basic operators.
	{"+/1 2 3", "6", nil},                            // plus reduce
	{"1 2 3 +.× 4 3 2", "16", nil},                   // scalar product
	{"(2 3⍴⍳6) +.× 3 2⍴5+⍳6", "52 58\n124 139", nil}, // matrix multiplication

	// Variable assignments.
	{"X←3", "", nil},          // assign a number
	{"-X←3", "¯3", nil},       // assign a value and use it
	{"X←3⋄X←4", "", nil},      // assign and overwrite
	{"X←3⋄⎕←X", "3", nil},     // assign and check
	{"f←+", "", nil},          // assign a function
	{"f←+⋄⎕←3 f 3", "6", nil}, // assign a function and apply
	{"X←4⋄⎕←÷X", "0.25", nil}, // assign and use it in another expr

	// Bracket indexing.
	//{"A←⍳6 ⋄ ⎕←A[1]", "x", nil}, // simple indexing

	// IBM APL Language, 3rd edition, June 1976.
	{"1000×(1+.06÷1 4 12 365)*10×1 4 12 365", "1790.8476965428547 1814.0184086689414 1819.3967340322804 1822.0289545386752", nil},
	// the original prints as: "1790.85 1413.02 1819.4 1822.03"
	{"Area ← 3×4\nX←2+⎕←3×Y←4\nX\nY", "12\n14\n4", nil},

	// Lambda expressions.
	{"{2×⍵}3", "6", nil},           // lambda in monadic context
	{"2{⍺+3{⍺×⍵}⍵+2}2", "14", nil}, // nested lambas
	{"2{(⍺+3){⍺×⍵}⍵+⍺{⍺+1+⍵}1+2}2", "40", nil},
	{"1{1+⍺{1+⍺{1+⍺+⍵}1+⍵}1+⍵}1", "7", nil},
	{"2{}4", "", nil},          // empty lambda expression ignores arguments
	{"{⍺×⍵}/2 3 4", "24", nil}, // TODO

	// Tool of thought.
}

func testCompare(got, exp string) bool {
	got = strings.TrimSpace(got)
	gotlines := strings.Split(got, "\n")
	explines := strings.Split(exp, "\n")
	if len(gotlines) != len(explines) {
		return false
	}
	for i, g := range gotlines {
		e := explines[i]
		gf := strings.Fields(g)
		ef := strings.Fields(e)
		if len(gf) != len(ef) {
			return false
		}
		for k := range gf {
			if gf[k] != ef[k] {
				return false
			}
		}
	}
	return true
}

func TestApl(t *testing.T) {
	// Compare result with expectation but ignores differences in whitespace.
	for i, tc := range testCases {
		var buf strings.Builder
		a := apl.New(&buf)
		numbers.Register(a)
		Register(a)
		operators.Register(a)

		// Set numeric formats.
		if tc.formats != nil {
			for t, f := range tc.formats {
				if num, ok := a.Tower.Numbers[t]; ok {
					num.Format = f
					a.Tower.Numbers[t] = num
				}
			}
		}

		lines := strings.Split(tc.in, "\n")
		for k, s := range lines {
			t.Logf("\t%s", s)
			if err := a.ParseAndEval(s); err != nil {
				t.Fatalf("tc%d:%d: %s: %s\n", i+1, k+1, tc.in, err)
			}
		}
		got := buf.String()
		t.Log(got)

		g := got
		g = spaces.ReplaceAllString(g, " ")
		g = newline.ReplaceAllString(g, "\n")
		g = strings.TrimSpace(g)
		if g != tc.exp {
			fmt.Printf("%q != %q\n", g, tc.exp)
			t.Fatalf("tc%d:\nin>\n%s\ngot>\n%s\nexpected>\n%s", i+1, tc.in, got, tc.exp)
		}
	}
}

var spaces = regexp.MustCompile(`  *`)
var newline = regexp.MustCompile(`\n *`)
