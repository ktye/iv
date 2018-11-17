package primitives

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
)

//go:generate go run gen.go

var testCases = []struct {
	in, exp string
	compare func(string, string) bool
}{
	// Basics numbers and arithmetics.
	{"1", "1", nil},
	{"1+1", "2", nil},
	{"1-2", "¯1", nil}, // negative number
	{"¯1", "¯1", nil},
	{"1-¯2", "3", nil},
	{"1@90", "1@90", nil}, // a complex number
	{"1@60+1@300", "1@0", nil},
	{"1J1", "1.4142135623730951@45", nil},

	// Vectors.
	{"1 2 3", "1 2 3", nil},
	{"1+1 2 3", "2 3 4", nil},
	{"1 2 3+¯1", "0 1 2", nil},
	{"1 2 3+4 5 6", "5 7 9", nil},

	// Braces.
	{"1 2+3 4", "4 6", nil},
	// TODO: {"1 (2+3) 4", "1 5 6", nil},
	{"(1 2)+3 4", "4 6", nil},
	{"1×2+3×4", "14", nil},
	{"1×(2+3)×4", "20", nil},
	{"(3×2)+3×4", "18", nil},
	{"3×2+3×4", "42", nil},

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
	{"1000×(1+.06÷1 4 12 365)*10×1 4 12 365", "1790.8476965428547 1814.0184086689414 1819.3967340322804 1822.0289545386752", cmpFloats},
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

func testCompare(got, exp string, eql func(a, b string) bool) bool {
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

func cmpFloats(a, b string) bool {
	tol := 1.0E-9
	eq := func(x, y string) bool {
		x = strings.Replace(x, "¯", "-", -1)
		y = strings.Replace(y, "¯", "-", -1)
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return false
		}
		g, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return false
		}
		if e := math.Abs(f - g); e > tol {
			return false
		}
		return true
	}
	return testCompare(a, b, eq)
}

func TestApl(t *testing.T) {
	// Compare result with expectation but ignores differences in whitespace.

	for i, tc := range testCases {
		var buf strings.Builder
		a := apl.New(&buf)
		numbers.Register(a)
		Register(a)
		operators.Register(a)
		lines := strings.Split(tc.in, "\n")
		for k, s := range lines {
			t.Logf("\t%s", s)
			if err := a.ParseAndEval(s); err != nil {
				t.Fatalf("tc%d:%d: %s: %s\n", i+1, k+1, tc.in, err)
			}
		}
		got := buf.String()
		t.Log(got)
		cmp := tc.compare
		if cmp == nil {
			cmp = func(a, b string) bool {
				eq := func(x, y string) bool {
					return x == y
				}
				return testCompare(a, b, eq)
			}
		}
		if cmp(got, tc.exp) == false {
			t.Fatalf("tc%d:\nin>\n%s\ngot>\n%s\nexpected>\n%s", i+1, tc.in, got, tc.exp)
		}
	}
}
