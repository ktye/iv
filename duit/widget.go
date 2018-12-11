// Duit provides an APL widget for duit.
package duit

import (
	"fmt"
	"image"
	"strings"
	"unicode/utf8"

	"github.com/ktye/duit"
	draw "github.com/ktye/duitdraw"
	"github.com/ktye/iv/apl"
	// TODO aplimg "github.com/ktye/iv/aplextra/image"
	"golang.org/x/mobile/event/key"
)

// Apl is a widget for duit.
// It embedds a duit.Edit widget with the following modifications:
//	- Pressing newline sends the current line to APL
//	  The result is appended to the bottom
//	- If the result value is an image, it is displayed
//	  on the top left corner.
//	  The image is not persistent. It will be overdrawn
//	  by any updates of the edit widget.
// An example application is cmd/aplui which also embedds
// APL385 unicode fonts and translates key events without
// the need for special keyboard drivers.
type Apl struct {
	*duit.Edit
	Apl *apl.Apl
}

// NewAPL initializes the underlying edit widget and returns it.
// APL is still unset and needs to be connected.
func NewAPL(f duit.SeekReaderAt) (*Apl, error) {
	e, err := duit.NewEdit(f)
	if err != nil {
		return nil, err
	}
	return &Apl{
		Edit: e,
	}, nil
}

// Key filters the keyboard events that are sent to duit.Edit.
// Modifications are:
//	Ctrl-c, Ctrl-x, Ctrl-v: use Windows conventions
//	Enter: execute APL command
func (ui *Apl) Key(dui *duit.DUI, self *duit.Kid, k rune, m draw.Mouse, orig image.Point) (r duit.Result) {

	// Remap Ctrl-c/x/v for copy paste.
	// duit.Edit uses a different convention.
	const Ctrl = 0x1f
	switch k {
	case 3: // Ctrl-c
		k = draw.KeyCmd + 'c'
	case 24: // Ctrl-x
		k = draw.KeyCmd + 'x'
	case 22: // Ctrl-v
		k = draw.KeyCmd + 'v'
	case 10: // Enter
		// Read the current line.
		// It seems we have to make two readers.
		// One that reads from the current position to the end of the line,
		// and one that reads backwards to the start.
		// Is there a simpler way to get the current line?
		off := ui.Edit.Cursor().Cur
		br := ui.Edit.ReverseEditReader(off)
		fr := ui.Edit.EditReader(off)
		_, rev, _ := br.Line(false)
		_, fwd, _ := fr.Line(false)
		nr := utf8.RuneCountInString(rev)
		nf := utf8.RuneCountInString(fwd)
		runes := make([]rune, nr+nf)
		i := nr
		for _, r := range rev {
			i--
			runes[i] = r
		}
		i = 0
		for _, r := range fwd {
			runes[nr+i] = r
			i++
		}
		line := strings.TrimSpace(string(runes))

		if len(line) > 0 {
			if ui.Apl == nil {
				fmt.Fprintf(ui, "\nAPL is not initialized\n")
				ui.ScrollCursor(dui)
				self.Draw = duit.Dirty
			} else {
				im := ui.execute(line)
				fmt.Fprintf(ui, "%s", prompt)
				ui.ScrollCursor(dui)

				// If the result is an image, draw it over everything else.
				// The update will be executed, after updating the editor.
				// It will disappear, as soon as the editor receives the
				// next event.
				if im != nil {
					if rgb, ok := im.(*image.RGBA); ok {
						dui.Call <- func() {
							m := dui.Display.MakeImage(rgb)
							dui.Display.ScreenImage.Draw(rgb.Bounds(), m, nil, image.ZP)
							dui.Display.Flush()
						}
					} else {
						fmt.Printf("image is not an *image.RGBA but %T\n", im)
					}
				}
				self.Draw = duit.Dirty
			}
		}
		return
	}

	return ui.Edit.Key(dui, self, k, m, orig)
}

// Execute the line by the APL interpreter.
// The output and a prompt is appended to the editor.
func (ui *Apl) execute(line string) image.Image {
	var im image.Image
	// If it's not the last line, print the input string.
	last := false
	cur := ui.Edit.Cursor().Cur
	fr := ui.Edit.EditReader(cur)
	n, _, eof := fr.Line(false)
	if eof {
		// The cursor was at EOF
		last = true
	} else {
		// We could be within the last line.
		// Try again.
		cur += int64(n)
		fr = ui.Edit.EditReader(cur)
		if _, _, eof := fr.Line(false); eof {
			last = true
		}
	}
	if last == false {
		fmt.Fprintf(ui, "%s", line)
	}
	ui.Write([]byte{'\n'})

	prog, err := ui.Apl.Parse(line)
	if err != nil {
		fmt.Fprintf(ui, "%s\n", err)
		return nil
	}
	vals, err := ui.Apl.EvalProgram(prog)
	if err != nil {
		fmt.Fprintf(ui, "%s\n", err)
		return nil
	}
	if len(vals) > 0 {
		/* TODO
		v := vals[0]
		if img, ok := v.(aplimg.Value); ok {
			im = img.Image
		} else {
		*/
		// We cannot handle tab's correctly.
		for _, v := range vals {
			s := strings.Replace(v.String(ui.Apl), "\t", "        ", -1)
			fmt.Fprintln(ui, s)
		}
		//}
	}
	return im
}

const prompt = "        "

func (ui *Apl) Write(p []byte) (n int, err error) {
	ui.Edit.Append(p)
	return len(p), nil
}

// AplKeyboard is injected into the shiny event loop.
// It implements a duitdraw.KeyTranslator.
type AplKeyboard struct{}

func (t AplKeyboard) TranslateKey(e key.Event) rune {
	var r rune = -1
	// fmt.Printf("r=%d code=%d Modifiers=%d dir=%d\n", e.Rune, e.Code, e.Modifiers, e.Direction)
	if rn := e.Rune; e.Direction != key.DirRelease {

		// If the Alt-Gr key is pressed (the one right to the space bar),
		// I see e.Modifiers == 6. If shift is pressed as well, it is 7.
		if e.Modifiers == 6 || e.Modifiers == 7 {
			if r, ok := Keyboard[e.Code]; ok {
				if e.Modifiers == 6 {
					return r[0]
				}
				return r[1]
			}
		}

		if rn != -1 {
			r = rn
		} else {
			if rn, ok := keymap[e.Code]; ok {
				r = rn
			}
		}

	}
	// Shiny sends \r on Enter, duit expects \n.
	if r == '\r' {
		r = '\n'
	}
	return r
}

// This is only needed for the default code path of keyboard handling.
// It is copied from duitdraw.
var keymap = map[key.Code]rune{
	key.CodeHome:       draw.KeyHome,
	key.CodeUpArrow:    draw.KeyUp,
	key.CodePageUp:     draw.KeyPageUp,
	key.CodeLeftArrow:  draw.KeyLeft,
	key.CodeRightArrow: draw.KeyRight,
	key.CodeDownArrow:  draw.KeyDown,
	key.CodePageDown:   draw.KeyPageDown,
	key.CodeEnd:        draw.KeyEnd,
	//key.CodeDelete:     KeyDelete,
	key.CodeEscape: draw.KeyEscape,
	//key.CodeCmd:        KeyCmd,
}

// Keyboard maps from a key code to two runes.
// The first one is used if the Alt-Gr key is used,
// The sencond if both, Alt-Gr and Shift are used.
// See cmd/aplui/main.go: welcome for the keyboard layout.
var Keyboard = map[key.Code][2]rune{
	// Top row.
	53: {'⋄', '⍨'},
	30: {'¨', '¡'},
	31: {'¯', '€'},
	32: {'<', '£'},
	33: {'≤', '⍧'},
	34: {'=', '≢'},
	35: {'≥', 'τ'},
	36: {'>', 'η'},
	37: {'≠', '⍂'},
	38: {'∨', '⍱'},
	39: {'∧', '⍲'},
	45: {'×', '≡'},
	46: {'÷', '⌹'},

	// Second row.
	20: {'?', '¿'},
	26: {'⍵', '⌽'},
	8:  {'∊', '⍷'},
	21: {'⍴', 'λ'},
	23: {'∼', '⍉'},
	28: {'↑', '¥'},
	24: {'↓', 'μ'},
	12: {'⍳', '⍸'},
	18: {'○', '⍥'},
	19: {'⋆', '⍟'},
	47: {'←', 'π'},
	48: {'→', 'Ω'},
	49: {'⍝', '⍀'},

	// Third row.
	4:  {'⍺', '⊖'},
	22: {'⌈', 'σ'},
	7:  {'⌊', 'δ'},
	9:  {'_', '⍫'},
	10: {'∇', '⍒'},
	11: {'∆', '⍋'},
	13: {'∘', '⍤'},
	14: {'\'', '⌺'},
	15: {'⎕', '⍞'},
	51: {'⊢', -1},
	52: {'⊣', -1},

	// Fourth row.
	29: {'⊂', 'ζ'},
	27: {'⊃', 'ξ'},
	6:  {'∩', '⍝'},
	25: {'∪', 'ø'},
	5:  {'⊥', '⍎'},
	17: {'⊤', '⍕'},
	16: {'|', '⌶'},
	54: {'⌷', '⍪'},
	55: {'⍎', '⍙'},
	56: {'⍕', '⌿'},
}
