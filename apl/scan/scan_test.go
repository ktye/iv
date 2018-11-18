package scan

import (
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
		{"¯1.0E¯6a123.8", []Token{Token{T: Number, S: "¯1.0E¯6a123.8"}}},
		{"¯8", []Token{Token{T: Number, S: "¯8"}}},
		{`+ alpha ≥3.23 "x""yz"`, []Token{
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
	}

	var scn Scanner
	scn.SetSymbols(symbols)
	for _, tc := range testCases {
		if got, err := scn.Scan(tc.input); err != nil {
			t.Fatalf("%s: %s", tc.input, err)
		} else {
			if len(got) != len(tc.exp) {
				t.Fatalf("got %d Tokens, expected %d", len(got), len(tc.exp))
			}
			for i, e := range tc.exp {
				g := got[i]
				if g.T != e.T || g.S != e.S {
					t.Fatalf("got %+v, expected %+v", g, e)
				}
			}
		}
	}
}
