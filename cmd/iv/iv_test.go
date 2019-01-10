package main

import (
	"strings"
	"testing"

	"github.com/ktye/iv/cmd/iv/iv"
)

func TestIv(t *testing.T) {
	testCases := []struct {
		args    []string
		in, exp string
	}{
		{[]string{"⍵"}, "1\n1", "1\n1\n"},
		{[]string{"⍺ ⍵"}, "1 2 3\n4 5 6\n", "1\n1\n"},
	}

	for i, tc := range testCases {
		iv.Stdin = strings.NewReader(tc.in)
		var out strings.Builder
		err := run(&out, tc.args)
		if err != nil {
			t.Fatalf("tc%d: %s", i, err)
		}
		got := out.String()
		if got != tc.exp {
			t.Fatalf("tc%d: got: %q\nexp: %q", i, got, tc.exp)
		}
	}
}
