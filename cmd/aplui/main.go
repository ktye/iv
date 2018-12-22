// Aplui is a gui frontend to APL\iv
//
// Aplui uses the iv/duit widget and embedds APL385 Unicode font.
// Keystrokes are translated automatically and no special keyboard
// driver is needed.
//
// When pressing the ENTER key, the current line is interpreted
// and the result is appended to the end of the editor.
//
// Aplui displays image values on the top left corner over the
// input text. The image disappears at the next input event.
//
// Aplui builds as a single binary.
// On windows, build with: go build -ldflags -H=windowsgui
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ktye/duit"
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/big"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	aplstrings "github.com/ktye/iv/apl/strings"

	"github.com/ktye/iv/aplextra/help"
	ivduit "github.com/ktye/iv/duit"
)

func main() {
	var fontsize = 18
	var quiet bool
	var extra bool
	var bignum bool
	var prec uint
	flag.IntVar(&fontsize, "fontsize", fontsize, "size of built-in font")
	flag.BoolVar(&quiet, "quiet", false, "dont show welcome message")
	flag.BoolVar(&extra, "extra", true, "register all packages in aplextra")
	flag.BoolVar(&bignum, "big", false, "use big numbers int and rational")
	flag.UintVar(&prec, "prec", 0, "use multi precision floats and complex")
	flag.Parse()

	if bignum && prec != 0 {
		fmt.Println("only one of -big and -prec>0 is allowed")
	}

	// Start APL.
	a := apl.New(nil)
	if bignum {
		big.SetBigTower(a)
	} else if prec > 0 {
		big.SetPreciseTower(a, prec)
	} else {
		numbers.Register(a)
	}
	/* TODO
	if extra {
		aplextra.RegisterAll(a)
	} else {
		primitives.Register(a)
		operators.Register(a)
	}
	*/
	primitives.Register(a)
	operators.Register(a)
	aplstrings.Register(a)

	// Build the gui.
	registerFont(fontsize)
	dui, err := duit.NewDUI("APL\\iv", &duit.DUIOpts{
		FontName: fmt.Sprintf("APL385@%d", fontsize),
	})
	if err != nil {
		log.Fatal(err)
	}
	dui.Display.KeyTranslator = ivduit.AplKeyboard{}

	// Use a single apl widget as the only ui element.
	content := `APL\iv` + help.Keyboard + "        "
	if quiet {
		content = "        "
	}
	ui, err := ivduit.NewAPL(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	end := int64(len(content))
	ui.SetCursor(duit.Cursor{end, end})
	ui.Apl = a
	ui.Apl.SetOutput(ui)

	//dui.Top.UI = &duit.Box{Kids: duit.NewKids(print, edit)}
	dui.Top.UI = ui
	dui.Render()

	for {
		select {
		case e := <-dui.Inputs:
			dui.Input(e)

		case err, ok := <-dui.Error:
			if !ok {
				return
			}
			log.Print(err)
		}
	}
}
