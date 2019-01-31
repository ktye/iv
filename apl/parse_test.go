package apl

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {

	testCases := []struct {
		in, exp string
	}{
		{"+1", "(+ 1)"},
		{"1", "1"},
		{"1 2", "(1 2)"},
		{`1 "alpha" 2`, `(1 "alpha" 2)`},
		{"+'e'-'Pete'", `(+ (e - ("P" "e" "t" "e")))`},
		// {"1 (2+3) 4", ""}, // not supported
		{"-1", "(- 1)"},
		{"¯2+3", "(¯2 + 3)"},
		{"1 2 3+4 5 6", "((1 2 3) + (4 5 6))"},
		{"(1+(1))", "(1 + 1)"},
		{"((1+1)+1)+1", "(((1 + 1) + 1) + 1)"},
		{"+", "+"},
		{"++1+1", "(+ (+ (1 + 1)))"},
		{"3+1/4", "(3 + ((1 /) 4))"},
		{"3+X←4", "(3 + ((X ←) 4))"},
		{"++(-1)/2", "(+ (+ (((- 1) /) 2)))"},
		{"+/1 2 3", "((+ /) (1 2 3))"},
		{"+.*/1 2 3", "(((+ . *) /) (1 2 3))"},
		{"f ← +", "((f ←) +)"},
		{"f ← +.*", "((f ←) (+ . *))"},
		{"1 2/3 4 5", "(((1 2) /) (3 4 5))"},
		{"1 2/[2]3 4 5", "((((1 2) /) ⍂ 2) (3 4 5))"},
		{"X ← +/ 3 4 5 + 1 2 3", "((X ←) ((+ /) ((3 4 5) + (1 2 3))))"},
		{"+.*/1", "(((+ . *) /) 1)"},
		{"+.*.*/1", "((((+ . *) . *) /) 1)"},
		{"+/+/+/1", "((+ /) ((+ /) ((+ /) 1)))"},
		{"X←3⋄X←4", "((X ←) 3)⋄((X ←) 4)"},
		{"math→sin 314", "(math→sin 314)"},
		{"1 2 3 math→add 4 5 6", "((1 2 3) math→add (4 5 6))"},
		{"1 2 3 math→Pi 4 5 6", "(1 2 3 math→Pi 4 5 6)"},
		{"A[1;2;3]", "([1;2;3] ⌷ A)"},
		{"A[1]", "([1] ⌷ A)"},
		{"A[1+1]", "([(1 + 1)] ⌷ A)"},
		{"A[1;2]", "([1;2] ⌷ A)"},
		{"A[1;]", "([1;] ⌷ A)"},
		{"A[;2]", "([;2] ⌷ A)"},
		{"A[;2;]", "([;2;] ⌷ A)"},
		{"A[;2;;]", "([;2;;] ⌷ A)"},
		{"A[1;2+2]", "([1;(2 + 2)] ⌷ A)"},
		{`(1;2;)`, "(1;2;)"},
		{`(1;(2;3;);)`, "(1;(2;3;);)"},
		{"{}", "{}"},
		{"{⍺+⍵}", "{(⍺ + ⍵)}"},
		{"{1:2}", "{1:2}"},
		{"{1:2⋄3:4⋄5}", "{1:2⋄3:4⋄5}"},
		{"{1:2:3}", "{1:2⋄3}"},
		{"{∇⍵}1", "({(∇ ⍵)} 1)"},
		{"{⍺+⍵}/1 2 3", "(({(⍺ + ⍵)} /) (1 2 3))"},
		{"{⍺{⍺+⍵}⍵}", "{(⍺ {(⍺ + ⍵)} ⍵)}"},
	}

	for i, tc := range testCases {
		a := New(os.Stdout)
		reg(a)

		p, err := a.Parse(tc.in)
		if err != nil {
			t.Fatalf("[%d] %s: %s", i+1, tc.in, err)
		}
		got := p.String(a)
		if got != tc.exp {
			t.Fatalf("[%d] %s:\nexpected:\n%s\ngot:\n%s", i+1, tc.in, tc.exp, got)
		}
	}
}

// TestMultiline test parsing of multi line lambda functions.
func TestMultiline(t *testing.T) {
	testCases := []struct {
		in, exp string
	}{
		{"1", "1"},
		{"+1", "(+ 1)"},
		{"1+1", "(1 + 1)"},
		{"1{}2", "(1 {} 2)"},
		{"{X←⍵\n2+⍵}", "{((X ←) ⍵)⋄(2 + ⍵)}"},
		{"{\n\tX←⍵\n\t2+⍵\n}", "{((X ←) ⍵)⋄(2 + ⍵)}"},
	}

	for i, tc := range testCases {
		a := New(os.Stdout)
		reg(a)

		p, err := a.ParseLines(tc.in)
		if err != nil {
			t.Fatalf("[%d] %s: %s", i+1, tc.in, err)
		}
		got := p.String(a)
		if got != tc.exp {
			t.Fatalf("[%d] %s:\nexpected:\n%s\ngot:\n%s", i+1, tc.in, tc.exp, got)
		}
	}
}

// For testing the parser we register just a couple of dummy primitives and two operators.
func reg(a *Apl) {
	for _, r := range "+-*!>" {
		a.RegisterPrimitive(Primitive(r), dummy)
	}
	a.RegisterOperator("/", mop{})
	a.RegisterOperator("←", mop{})
	a.RegisterOperator(".", dot{})
	a.RegisterOperator("⍂", brack{})
	if err := a.SetTower(newTower()); err != nil {
		panic(err)
	}
}

// Dummy primitive.
var dummy dummyPrimitive

type dummyPrimitive struct{}

func (d dummyPrimitive) Call(a *Apl, l, r Value) (Value, error) {
	return EmptyArray{}, nil
}
func (d dummyPrimitive) Select(a *Apl, l, r Value) (IndexArray, error) {
	return IndexArray{}, fmt.Errorf("not supported")
}
func (d dummyPrimitive) To(a *Apl, l, r Value) (Value, Value, bool) {
	return l, r, true
}
func (d dummyPrimitive) String(a *Apl) string { return "any" }
func (d dummyPrimitive) Doc() string          { return "dummy" }

var dummyfunc dummyFunction

type dummyFunction struct{}

func (d dummyFunction) Call(a *Apl, l, r Value) (Value, error) { return Index(1), nil }

// Monadic operators.
type mop struct{}

func (r mop) To(a *Apl, LO, RO Value) (Value, Value, bool) { return LO, RO, true }
func (r mop) String(a *Apl) string                         { return "any" }
func (r mop) DyadicOp() bool                               { return false }
func (r mop) Derived(a *Apl, lo, ro Value) Function        { return dummyfunc }
func (r mop) Select(a *Apl, l Value, lo Value, ro Value, R Value) (IndexArray, error) {
	return IndexArray{}, nil
}
func (r mop) Doc() string { return "reduce" }

// Dyadic operators.
type dot struct {
}

func (d dot) To(a *Apl, l, r Value) (Value, Value, bool) { return l, r, true }
func (d dot) String(a *Apl) string                       { return "any" }
func (d dot) DyadicOp() bool                             { return true }
func (d dot) Derived(a *Apl, lo, ro Value) Function      { return dummyfunc }
func (d dot) Select(a *Apl, l Value, lo Value, ro Value, r Value) (IndexArray, error) {
	return IndexArray{}, nil
}
func (d dot) Doc() string { return "dot" }

type brack struct{}

func (d brack) To(a *Apl, l, r Value) (Value, Value, bool) { return l, r, true }
func (d brack) String(a *Apl) string                       { return "any" }
func (d brack) DyadicOp() bool                             { return true }
func (d brack) Derived(a *Apl, lo, ro Value) Function      { return dummyfunc }
func (d brack) Select(a *Apl, l Value, lo Value, ro Value, r Value) (IndexArray, error) {
	return IndexArray{}, nil
}
func (d brack) Doc() string { return "bracket operator" }

func newTower() Tower {
	m := make(map[reflect.Type]Numeric)
	m[reflect.TypeOf(Index(0))] = Numeric{
		Class: 0,
		Parse: func(s string) (Number, bool) {
			s = strings.Replace(s, "¯", "-", -1)
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, false
			}
			return Index(n), true
		},
		Uptype: func(n Number) (Number, bool) { return nil, false },
	}
	t := Tower{
		Numbers: m,
	}
	return t
}
