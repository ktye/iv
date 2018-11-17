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
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"

	"github.com/ktye/iv/aplextra/help"
	ivduit "github.com/ktye/iv/duit"
)

func main() {
	var fontsize = 18
	var quiet bool
	var extra bool
	flag.IntVar(&fontsize, "fontsize", fontsize, "size of built-in font")
	flag.BoolVar(&quiet, "quiet", false, "dont show welcome message")
	flag.BoolVar(&extra, "extra", true, "register all packages in aplextra")
	flag.Parse()

	// Start APL.
	a := apl.New(nil)
	/* TODO
	if extra {
		aplextra.RegisterAll(a)
	} else {
		funcs.Register(a)
		operators.Register(a)
	}
	*/
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)

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
