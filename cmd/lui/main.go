// Lui is a gui frontend to APL\iv
package main

import (
	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/apl"
	apkg "github.com/ktye/iv/apl/a"
	"github.com/ktye/iv/apl/big"
	"github.com/ktye/iv/apl/http"
	"github.com/ktye/iv/apl/io"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	"github.com/ktye/iv/apl/rpc"
	aplstrings "github.com/ktye/iv/apl/strings"
	"github.com/ktye/iv/apl/xgo"
	"github.com/ktye/iv/aplextra/q"
	"github.com/ktye/iv/aplextra/u"
	"github.com/ktye/ui"
)

func main() {

	// Start APL and register all packages.
	a := apl.New(nil)
	numbers.Register(a)
	big.Register(a, "")
	primitives.Register(a)
	operators.Register(a)
	apkg.Register(a, "")
	io.Register(a, "")
	aplstrings.Register(a, "")
	xgo.Register(a, "")
	rpc.Register(a, "")
	http.Register(a, "")
	q.Register(a, "")
	u.Register(a, "")

	// Implement the repl interface.
	var interp interp
	repl := &ui.Repl{Reply: true}
	repl.SetText(rope.New("\t"))
	interp.repl = repl
	repl.Interp = &interp

	// Connect APL with the repl.
	interp.apl = a
	a.SetOutput(repl)

	// Run the ui main loop.
	u.Loop(repl)
}

type interp struct {
	apl  *apl.Apl
	repl *ui.Repl
}

func (i *interp) Eval(s string) {
	i.repl.Write([]byte{'\n'})
	p, err := i.apl.ParseLines(s)
	if err == nil {
		err = i.apl.Eval(p)
	}
	if err != nil {
		i.repl.Write([]byte(err.Error() + "\n"))
	}
	i.repl.Write([]byte("\t"))
	i.repl.Edit.MarkAddr("$")
}

// Cancel is ignored.
// Currently there is no way to interrupt APL.
func (i *interp) Cancel() {}
