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
	"github.com/ktye/iv/apl/funcs"
	"github.com/ktye/iv/apl/image"
	"github.com/ktye/iv/apl/operators"
	ivduit "github.com/ktye/iv/duit"
)

func main() {
	var fontsize = 18
	var quiet bool
	flag.IntVar(&fontsize, "fontsize", fontsize, "size of built-in font")
	flag.BoolVar(&quiet, "quiet", false, "dont show welcome message")
	flag.Parse()

	// Start APL.
	a := apl.New(nil)
	funcs.Register(a)
	operators.Register(a)
	image.Register(a)
	// ... add more packages here ...

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
	content := welcome
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

// The nice keyboard is from gnu-apl.
// Is there any other font than APL385 under which it displays correctly?
// I had to remove the back-tick.
// Added symbols are: πσδø
const welcome = `APL\iv
╔════╦════╦════╦════╦════╦════╦════╦════╦════╦════╦════╦════╦════╦═════════╗
║ ~⍨ ║ !¡ ║ @€ ║ #£ ║ $⍧ ║ %  ║ ^  ║ &  ║ *⍂ ║ (⍱ ║ )⍲ ║ _≡ ║ +⌹ ║         ║
║  ◊ ║ 1¨ ║ 2¯ ║ 3< ║ 4≤ ║ 5= ║ 6≥ ║ 7> ║ 8≠ ║ 9∨ ║ 0∧ ║ -× ║ =÷ ║ BACKSP  ║
╠════╩══╦═╩══╦═╩══╦═╩══╦═╩══╦═╩══╦═╩══╦═╩══╦═╩══╦═╩══╦═╩══╦═╩══╦═╩══╦══════╣
║       ║ Q¿ ║ W⌽ ║ E⍷ ║ R  ║ T⍉ ║ Y¥ ║ U  ║ I⍸ ║ O⍥ ║ P⍟ ║ {π ║ }  ║  |⍀  ║
║  TAB  ║ q? ║ w⍵ ║ e∊ ║ r⍴ ║ t∼ ║ y↑ ║ u↓ ║ i⍳ ║ o○ ║ p⋆ ║ [← ║ ]→ ║  \⍝  ║
╠═══════╩═╦══╩═╦══╩═╦══╩═╦══╩═╦══╩═╦══╩═╦══╩═╦══╩═╦══╩═╦══╩═╦══╩═╦══╩══════╣
║ (CAPS   ║ A⊖ ║ S  ║ D  ║ F⍫ ║ G⍒ ║ H⍋ ║ J⍤ ║ K⌺ ║ L⍞ ║ :σ ║ "δ ║         ║
║  LOCK)  ║ a⍺ ║ s⌈ ║ d⌊ ║ f_ ║ g∇ ║ h∆ ║ j∘ ║ k' ║ l⎕ ║ ;⊢ ║ '⊣ ║ RETURN  ║
╠═════════╩═══╦╩═══╦╩═══╦╩═══╦╩═══╦╩═══╦╩═══╦╩═══╦╩═══╦╩═══╦╩═══╦╩═════════╣
║             ║ Z  ║ X  ║ C⍝ ║ Vø ║ B⍎ ║ N⍕ ║ M⌶ ║ <⍪ ║ >⍙ ║ ?⌿ ║          ║
║  SHIFT      ║ z⊂ ║ x⊃ ║ c∩ ║ v∪ ║ b⊥ ║ n⊤ ║ m| ║ ,⌷ ║ .⍎ ║ /⍕ ║  SHIFT   ║
╚═════════════╩════╩════╩════╩════╩════╩════╩════╩════╩════╩════╩══════════╝
        ` // This is the start position (tabs are not supported for now).
