package complete

import "testing"

func TestFormat(t *testing.T) {
	testCases := []struct {
		line, exp string
	}{
		{"", ""},
		{"abcd", "abcd"},
		{"ab cd", "ab cd"},
		{"abcd ", "abcd "},
		{" abcd ", " abcd "},
		{"3+not x", "3+~x"},
		{"3a drop+ take me here", "3a↓+↑me here"},
		{`"don't take a s""tring"`, `"don't take a s""tring"`},
	}

	for _, tc := range testCases {
		if got := Format(tc.line); got != tc.exp {
			t.Fatalf("expected: %v, got %v", tc.exp, got)
		}
	}
}
