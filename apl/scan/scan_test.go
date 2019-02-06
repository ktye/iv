package scan

import (
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	// Symbols are registered only when importing funcs/.
	// This does not happen when running the test, so we set them manually.
	symbols := make(map[rune]string)
	for _, r := range "+≥⍵≡/?," {
		symbols[r] = string(r)
	}

	testCases := []struct {
		input string
		exp   []Token
	}{
		{"", nil},
		{"1⋄2", []Token{Token{T: Number, S: "1"}, Token{T: Diamond, S: "⋄"}, Token{T: Number, S: "2"}}},
		{".5", []Token{Token{T: Number, S: ".5"}}},
		{"1", []Token{Token{T: Number, S: "1"}}},
		{"1.23", []Token{Token{T: Number, S: "1.23"}}},
		{"1J2", []Token{Token{T: Number, S: "1J2"}}},
		{"`alpha`beta", []Token{Token{T: String, S: "alpha"}, Token{T: String, S: "beta"}}},
		{"1.23 pkg→name+3", []Token{
			Token{T: Number, S: "1.23"},
			Token{T: Identifier, S: "pkg→name"},
			Token{T: Symbol, S: "+"},
			Token{T: Number, S: "3"},
		}},
		{"¯1.0E¯6a123.8", []Token{Token{T: Number, S: "¯1.0E¯6a123.8"}}},
		{"¯8", []Token{Token{T: Number, S: "¯8"}}},
		{`"a⍝b"+8.2⍝comment`, []Token{
			Token{T: String, S: `a⍝b`},
			Token{T: Symbol, S: "+"},
			Token{T: Number, S: "8.2"},
		}},
		{`+ alpha ≥3.23 "x\"yz"`, []Token{
			Token{T: Symbol, S: "+"},
			Token{T: Identifier, S: "alpha"},
			Token{T: Symbol, S: "≥"},
			Token{T: Number, S: "3.23"},
			Token{T: String, S: `x"yz`},
		}},
		{`⋄ ⋄1.23E¯5 4.234  0.234⍵`, []Token{
			Token{T: Diamond, S: "⋄"},
			Token{T: Diamond, S: "⋄"},
			Token{T: Number, S: "1.23E¯5"},
			Token{T: Number, S: "4.234"},
			Token{T: Number, S: "0.234"},
			Token{T: Symbol, S: "⍵"},
		}},
		{`{⍵≡0: A[2;3]}`, []Token{
			Token{T: LeftBrace, S: "{"},
			Token{T: Symbol, S: "⍵"},
			Token{T: Symbol, S: "≡"},
			Token{T: Number, S: "0"},
			Token{T: Colon, S: ":"},
			Token{T: Identifier, S: "A"},
			Token{T: LeftBrack, S: "["},
			Token{T: Number, S: "2"},
			Token{T: Semicolon, S: ";"},
			Token{T: Number, S: "3"},
			Token{T: RightBrack, S: "]"},
			Token{T: RightBrace, S: "}"},
		}},
		{`{⍵∇1}`, []Token{
			Token{T: LeftBrace, S: "{"},
			Token{T: Symbol, S: "⍵"},
			Token{T: Self, S: "∇"},
			Token{T: Number, S: "1"},
			Token{T: RightBrace, S: "}"},
		}},
	}

	var scn Scanner
	scn.SetSymbols(symbols)
	for _, tc := range testCases {
		if got, err := scn.Scan(tc.input); err != nil {
			t.Fatalf("%q: %s", tc.input, err)
		} else {
			if len(got) != len(tc.exp) {
				t.Fatalf("%q: got %d Tokens, expected %d", tc.input, len(got), len(tc.exp))
			}
			for i, e := range tc.exp {
				g := got[i]
				if g.T != e.T || g.S != e.S {
					t.Fatalf("%q: got %+v, expected %+v", tc.input, g, e)
				}
			}
		}
	}
}

func TestScanString(t *testing.T) {
	testCases := [][2]string{
		// Double quoted strings with backslash escapes.
		{`"alpha"`, `alpha`},
		{`"alpha beta"`, `alpha beta`},
		{`"alpha\nbeta"`, "alpha\nbeta"},
		{`"alpha\\beta"`, "alpha\\beta"},
		{`"alpha\nbeta\r"`, "alpha\nbeta\r"},
		{`"alpha\nbeta\r"trailing`, "alpha\nbeta\r"},
		{`"al\"ha"`, `al"ha`},
		{`"\u263a"`, "☺"},

		// Single quoted strings, with double escapes.
		{`'a'`, "a"},
		{`'a'trail`, "a"},
		{`'alpha'`, `alpha`},
		{`'al''pha'`, `al'pha`},
		{`'al''p\nha'`, `al'p\nha`},
		{`'al''p\nha'trailing`, `al'p\nha`},

		// Backtick strings.
		{"`alpha", "alpha"},
		{"`alpha trailing", "alpha"},
		{"`alpha`trailing", "alpha"},
		{"`alpha}trailing", "alpha"},
		{"`alpha]trailing", "alpha"},
		{"`alpha⋄trailing", "alpha"},
		{"`alpha#trailing", "alpha"},
		{"`alpha\nbeta", "alpha"},
		{"`alpha\tbeta", "alpha"},
		{"`alpha\rbeta", "alpha"},
		{"`a\\l\"'", "a\\l\"'"},
	}
	for _, tc := range testCases {
		in := tc[0]
		exp := tc[1]
		got, err := ReadString(strings.NewReader(in))
		if err != nil {
			t.Fatal(err)
		}
		if got != exp {
			t.Fatalf("expected %q, got %q", exp, got)
		}
	}
}
