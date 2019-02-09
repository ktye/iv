package apl

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestScanRankArray(t *testing.T) {

	testCases := []struct {
		in   string
		rank int
		out  []string
	}{
		{"", 0, []string{}},
		{"", 1, []string{}},
		{"", 2, []string{}},
		{"1 2 3\n4 5 6\n\n7 8 9\n1 2 3", 0, []string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "1", "2", "3",
		}},
		{"1 2 3\n4 5 6\n\n7 8 9\n1 2 3", 1, []string{
			"1 2 3 [3]",
			"4 5 6 [3]",
			"7 8 9 [3]",
			"1 2 3 [3]",
		}},
		{"1 2 3\n4 5 6\n\n7 8 9\n1 2 3", 2, []string{
			" 1 2 3\n 4 5 6 [2 3]",
			" 7 8 9\n 1 2 3 [2 3]",
		}},
		{"1 2 3\n4 5 7\n", -1, []string{
			" 1 2 3\n 4 5 7 [2 3]",
		}},
		{"1 2 3\n4 5 6", -1, []string{
			" 1 2 3\n 4 5 6 [2 3]",
		}},
		{"1 2 3\n4 5 6\n\n7 8 9\n1 2 3\n\n", -1, []string{
			" 1 2 3\n 4 5 6\n\n 7 8 9\n 1 2 3 [2 2 3]",
		}},
		{"1 2 3\n4 5 6\n\n7 8 9\n1 2 3", -1, []string{
			" 1 2 3\n 4 5 6\n\n 7 8 9\n 1 2 3 [2 2 3]",
		}},
		{"1 2 3\n4 5 6\n\n7 8 9\n1 2 3", 3, []string{
			" 1 2 3\n 4 5 6\n\n 7 8 9\n 1 2 3 [2 2 3]",
		}},
		{"1 2 3\n\n4 5 6\n\n7 8 9\n\n1 2 3", 1, []string{
			"1 2 3 [3]",
			"4 5 6 [3]",
			"7 8 9 [3]",
			"1 2 3 [3]",
		}},
		{"1 2 3\n\n4 5 6\n\n7 8 9\n\n1 2 3", 2, []string{
			" 1 2 3 [1 3]",
			" 4 5 6 [1 3]",
			" 7 8 9 [1 3]",
			" 1 2 3 [1 3]",
		}},
		{"1 2 3\n\n4 5 6\n\n7 8 9\n\n1 2 3", 3, []string{
			" 1 2 3\n\n 4 5 6\n\n 7 8 9\n\n 1 2 3 [4 1 3]",
		}},
		{"[[[1 2 3],[4 5 6]],[[7 8 9],[1 2 3]]]", 2, []string{
			" 1 2 3\n 4 5 6 [2 3]",
			" 7 8 9\n 1 2 3 [2 3]",
		}},
		{"[[1,2,3],[4,5,6]]", -1, []string{
			" 1 2 3\n 4 5 6 [2 3]",
		}},
		{"[1,2,3]", -1, []string{
			"1 2 3 [3]",
		}},
		{"[1,2,3;4,5,6]", -1, []string{
			" 1 2 3\n 4 5 6 [2 3]",
		}},
	}

	for k, tc := range testCases {
		a := New(nil)
		in := strings.NewReader(tc.in)
		for i := 0; ; i++ {
			v, err := a.ScanRankArray(in, tc.rank)
			if err == io.EOF && i == len(tc.out) {
				goto next
			} else if err == io.EOF {
				t.Fatalf("#%d: expected %d results, got %d", k, len(tc.out), i-1)
			} else if err != nil {
				t.Fatalf("#%d: %s", k, err)
			} else if i == len(tc.out) {
				t.Fatalf("#%d: too many output values: %v", k, v)
			} else {
				s := ""
				if _, ok := v.(Array); !ok {
					s = v.String(a)
				} else {
					s = fmt.Sprintf("%v %v", v.String(a), v.(Array).Shape())
				}
				if s != tc.out[i] {
					t.Fatalf("#%d:%d: expected:\n%s\ngot:\n%s", k, i, tc.out[i], s)
				}
			}
		}
	next:
	}
}
