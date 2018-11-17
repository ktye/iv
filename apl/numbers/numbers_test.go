package numbers

import (
	"reflect"
	"testing"

	"github.com/ktye/iv/apl"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		s string
		n apl.Number
	}{
		{"1", Integer(1)},
		{"¯2", Integer(-2)},
		{"¯2.0", Float(-2)},
		{"3J¯2.0", Complex(complex(3, -2))},
		{"5@90", Complex(complex(0, 5))},
		{"3.12E¯2", Float(0.0312)},
		{".5", Float(0.5)},
		{"¯.3", Float(-0.3)},
	}

	a := apl.New(nil)
	Register(a)

	for _, tc := range testCases {
		ne, err := a.Tower.Parse(tc.s)
		if err != nil {
			t.Fatal(err)
		}
		n, err := ne.Eval(a)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.TypeOf(n) != reflect.TypeOf(tc.n) {
			t.Fatalf("%s: expected %T got %T", tc.s, tc.n, n)
		}
		if n != tc.n {
			t.Fatalf("%s: numbers are not equal: %v, %v", tc.s, tc.n, n)
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
		{Integer(1), Integer(2), Integer(1), Integer(2)},
		{Integer(0), Float(3), Float(0), Float(3)},
		{Float(3), Integer(4), Float(3), Float(4)},
		{Integer(2), Complex(3 + 1i), Complex(2), Complex(3 + 1i)},
		{Complex(1 + 2i), Float(3), Complex(1 + 2i), Complex(3)},
	}

	for _, tc := range testCases {
		c, d, err := a.Tower.SameType(tc.a, tc.b)
		if err != nil {
			t.Fatal(err)
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

func TestFromIndex(t *testing.T) {
	a := apl.New(nil)
	Register(a)

	testCases := []struct {
		i int
		n apl.Number
	}{
		{0, Integer(0)},
		{-1, Integer(-1)},
	}
	for _, tc := range testCases {
		n := a.Tower.FromIndex(tc.i)
		if reflect.TypeOf(n) != reflect.TypeOf(tc.n) {
			t.Fatal("wrong type")
		}
		if n != tc.n {
			t.Fatal("wrong value")
		}
	}
}
