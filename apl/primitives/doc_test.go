package primitives

import (
	"io"
	"os"
	"testing"
	"text/tabwriter"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	aplstrings "github.com/ktye/iv/apl/strings"
)

func TestDoc(t *testing.T) {
	if testing.Short() {
		a := apl.New(os.Stdout)
		numbers.Register(a)
		Register(a)
		operators.Register(a)
		aplstrings.Register(a)

		var w io.Writer = os.Stdout
		tw := tabwriter.NewWriter(w, 2, 0, 2, ' ', 0)
		a.Doc(tw)
		tw.Flush()
	}
}
