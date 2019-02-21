package p

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
)

func TestPlot(t *testing.T) {
	a := apl.New(ioutil.Discard)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)
	Register(a, "")

	b, err := ioutil.ReadFile("test.apl")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	for i, s := range lines {
		if s == "" || strings.HasPrefix(s, "⍝") {
			continue
		}
		if err := a.ParseAndEval(s); err != nil {
			t.Fatalf("test.apl:%d: %s", i+1, err)
		}
	}
}
