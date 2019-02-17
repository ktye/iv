// Lui is a gui frontend to APL\iv
package main

import (
	"fmt"
	"os"
	"strings"

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
	"github.com/ktye/iv/cmd"
	"github.com/ktye/ui"
)

const usage = `usage
     lui          ui mode
     lui -a ARGS  cmd/apl mode
     lui -i ARGS  cmd/iv mode
     lui -z ZIP   attach zip file
`

func main() {
	if len(os.Args) < 2 {
		gui(newApl())
		return
	}

	var err error
	switch os.Args[1] {
	case "-a":
		a := newApl()
		a.SetOutput(os.Stdout)
		a.SetImage(&u.Sxl{Writer: os.Stdout})
		err = cmd.Apl(a, os.Stdin, os.Args[2:])
	case "-i":
		a := newApl()
		a.SetImage(&u.Sxl{Writer: os.Stdout})
		err = cmd.Iv(a, strings.Join(os.Args[2:], " "), os.Stdout)
	case "-z":
		fmt.Println("TODO: attach zip file")
	default:
		fmt.Println(usage)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newApl() *apl.Apl {
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
	return a
}

// gui runs lui in gui mode.
func gui(a *apl.Apl) {
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
