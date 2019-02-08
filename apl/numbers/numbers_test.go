package numbers

import (
	"reflect"
	"testing"
	"time"

	"github.com/ktye/iv/apl"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		s string
		n apl.Number
	}{
		{"1", apl.Int(1)},
		{"1b", apl.Bool(true)},
		{"¯2", apl.Int(-2)},
		{"¯2.0", Float(-2)},
		{"2.", Float(2)},
		{"3J¯2.0", Complex(complex(3, -2))},
		{"5a90", Complex(complex(0, 5))},
		{"3.12E¯2", Float(0.0312)},
		{".5", Float(0.5)},
		{"¯.3", Float(-0.3)},
		{"2014.04.02", Time(time.Date(2014, 4, 2, 0, 0, 0, 0, time.UTC))},
		{"2014.04.02T09.37.22", Time(time.Date(2014, 4, 2, 9, 37, 22, 0, time.UTC))},
		{"10s", Time(y0.Add(10 * time.Second))},
	}

	a := apl.New(nil)
	Register(a)

	for k, tc := range testCases {
		ne, err := a.Tower.Parse(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		n, err := ne.Eval(a)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.TypeOf(n) != reflect.TypeOf(tc.n) {
			t.Fatalf("#%d: %s: expected %T got %T", k, tc.s, tc.n, n)
		}
		if n != tc.n {
			t.Fatalf("#%d: %s: numbers are not equal: %v, %v", k, tc.s, tc.n, n)
		}
	}
}

func TestSameType(t *testing.T) {
	a := apl.New(nil)
	Register(a)

	testCases := []struct {
		a, b apl.Number
		c, d apl.Number
	}{
		{apl.Bool(false), apl.Bool(true), apl.Bool(false), apl.Bool(true)},
		{apl.Bool(false), apl.Int(1), apl.Int(0), apl.Int(1)},
		{apl.Int(1), apl.Bool(false), apl.Int(1), apl.Int(0)},
		{apl.Int(1), apl.Int(2), apl.Int(1), apl.Int(2)},
		{apl.Bool(true), Float(3), Float(1), Float(3)},
		{apl.Int(0), Float(3), Float(0), Float(3)},
		{Float(3), apl.Int(4), Float(3), Float(4)},
		{apl.Int(2), Complex(3 + 1i), Complex(2), Complex(3 + 1i)},
		{Complex(1 + 2i), Float(3), Complex(1 + 2i), Complex(3)},
	}

	for n, tc := range testCases {
		c, d, err := a.Tower.SameType(tc.a, tc.b)
		if err != nil {
			t.Fatalf("#%d: %s", n, err)
		}
		if reflect.TypeOf(c) != reflect.TypeOf(d) {
			t.Fatalf("not the same type: %T %T", c, d)
		}
		if c != tc.c {
			t.Fatalf("expected %v got %v", tc.c, c)
		}
		if d != tc.d {
			t.Fatalf("expected %v got %v", tc.d, d)
		}
	}
}

func TestImport(t *testing.T) {
	a := apl.New(nil)
	Register(a)

	testCases := []struct {
		i apl.Number
		n apl.Number
	}{
		{apl.Bool(false), Float(0)},
		{apl.Bool(true), Float(1)},
		{apl.Int(0), Float(0)},
		{apl.Int(1), Float(1)},
		{apl.Int(2), Float(2)},
		{apl.Int(-1), Float(-1)},
	}
	for _, tc := range testCases {
		n := a.Tower.Import(tc.i)
		if reflect.TypeOf(n) != reflect.TypeOf(tc.n) {
			t.Fatal("wrong type")
		}
		if n != tc.n {
			t.Fatal("wrong value")
		}
	}
}
