package iv

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

// TestNext tests the scalar parsing logic.
func TestNext(t *testing.T) {
	testCases := []struct {
		in   string
		exp  []int // slice of termination levels for each call to next.
		fail bool
	}{
		{"", []int{}, false},
		{"1", []int{0}, false},
		{"1 ", []int{0}, false},
		{"1 2", []int{1, 1}, false},
		{"1   2", []int{1, 1}, false},
		{"1\n2", []int{2, 2}, false},
		{"1\n\n2", []int{3, 3}, false},
		{"1\n2\n\n\n", []int{2, 2}, false},
		{i1, []int{1, 1, 2, 1, 1, 3, 1, 1, 2, 1, 1, 3}, false},
		{i1 + "\n", []int{1, 1, 2, 1, 1, 3, 1, 1, 2, 1, 1, 3}, false},
	}

	run := func(in string, exp []int) error {
		a := apl.New(ioutil.Discard)
		numbers.Register(a)
		p := &InputParser{Separator: '\n', Apl: a}
		p.Reader = bufio.NewReader(tabularText(strings.NewReader(in)))
		for i := 0; ; i++ {
			_, S, err := p.Next()
			if err == io.EOF && i == len(exp) {
				return nil
			} else if err != nil {
				return fmt.Errorf("%d: %s", i, err)
			} else if i >= len(exp) {
				return fmt.Errorf("%d: more input than expected", i)
			}
			if S != exp[i] {
				return fmt.Errorf("%d: got S=%d, exp: %d", i, S, exp[i])
			}
		}
		return nil
	}

	for i, tc := range testCases {
		err := run(tc.in, tc.exp)
		if err != nil && tc.fail == false {
			t.Fatalf("tc:%d %s", i, err)
		} else if err == nil && tc.fail == true {
			t.Fatalf("tc:%d should have failed", i)
		}
	}
}

var i1 = `1 2 3
4 5 6

7 8 9
1 2 3`
