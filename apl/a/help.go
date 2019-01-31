package a

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/ktye/iv/apl"
)

// help returns the help text in a channel.
func help(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Commands:")
	for _, c := range a.Scanner.Commands() {
		fmt.Fprintf(&buf, " %s", c)
	}
	fmt.Fprintf(&buf, "\n\n")

	a.Doc(&buf)
	return apl.LineReader(ioutil.NopCloser(&buf)), nil
}
