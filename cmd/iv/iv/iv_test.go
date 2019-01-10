package iv

import (
	"strings"
	"testing"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
)

func TestIv(t *testing.T) {
	testCases := []struct {
		data, prog, exp string
	}{
		{"", "1+1", "2\n"}, // warming up
		{"7 8", "C←iv→r 0 ⋄  {⍵}¨C", "((7;0;);(8;0;);)\n"},
		{"7 8\n9", "C←iv→r 0 ⋄  {⍵}¨C", "((7;0;);(8;1;);(9;1;);)\n"},
		{"7 8", "C←iv→r 1 ⋄  {⍵}¨C", "((7 8;0;);)\n"},
		{"3 4\n5 6", "C←iv→r 1 ⋄  {⍵}¨C", "((3 4;0;);(5 6;0;);)\n"},
		{"3\n4\n5", "C←iv→r 1 ⋄  {⍵}¨C", "((3;0;);(4;0;);(5;0;);)\n"},
		{"3 4\n5 6\n\n7 8", "C←iv→r 1 ⋄  {⍵}¨C", "((3 4;0;);(5 6;1;);(7 8;1;);)\n"},
		{"7", "C←iv→r 1 ⋄  {⍵}¨C", "((7;0;);)\n"},
		{"7", "C←iv→r 1 ⋄  {⍴⍵[1]}¨C", "(1;)\n"}, // make sure it's rank 1.
		{"3 4\n5 6", "C←iv→r 2 ⋄  {⍵}¨C", "(( 3 4\n 5 6;0;);)\n"},
		{"3 4\n5 6\n\n1 2\n3 4", "C←iv→r 2 ⋄  {⍵}¨C", "(( 3 4\n 5 6;0;);( 1 2\n 3 4;0;);)\n"},
	}

	for i, tc := range testCases {
		var buf strings.Builder
		a := apl.New(&buf)
		numbers.Register(a)
		primitives.Register(a)
		operators.Register(a)
		Register(a)
		Stdin = strings.NewReader(tc.data)

		if err := a.ParseAndEval(tc.prog); err != nil {
			t.Fatal(err)
		}
		if got := buf.String(); got != tc.exp {
			t.Fatalf("tc%d: exp: %q\n got %q\n", i, tc.exp, got)
		}
	}
}
