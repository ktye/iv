// Package u provides ui elements to APL
//
// cmd/lui is a graphical user interface that shows a single Repl widget.
// u can be used to change the user interface during runtime from APL.
//
// u is a wrapper for github.com/ktye/ui which draws directly into go images
// and uses shiny as a backend to map windows and propaget key and mouse events.
//
// Functions
//	kb 0            APL keyboard layout as a string
//	top ⍳0          return the top widget
//	top T           set the top widget to T
//	button "label"  create a button with a label,  L (string): icon
//	A split B       create a split widget with childs A B
//	W f (fn;)       attach a callback to a widget
//
// Example
//	B← u→button "B"           ⍝ create a button
//	B u→f ({⎕←55};)          ⍝ attach a callback function
//	u→top B u→split u→top⍳0   ⍝ split the top widget and add the button next to the console
package u

import (
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/scan"
	"github.com/ktye/ui"
)

// Register adds the ui package to the interpreter with the default name u.
func Register(a *apl.Apl, name string) {
	if name == "" {
		name = "u"
	}
	pkg := map[string]apl.Value{
		"button": apl.ToFunction(button),
		"cls":    apl.ToFunction(cls),
		"sam":    apl.ToFunction(sam),
		"f":      apl.ToFunction(setCallback),
		"kb":     apl.ToFunction(kb),
		"split":  apl.ToFunction(split),
		"top":    apl.ToFunction(top),
	}
	cmd := map[string]scan.Command{
		"c": rw0("cls"),
		"e": toCommand(samCmd),
		"k": rw0("kb"),
	}
	a.AddCommands(cmd)
	a.RegisterPackage(name, pkg)
}

// Kb returns the APL keyboard layout as a string.
func kb(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return apl.String(ui.AplKeyboard{}.String()), nil
}

// Cls clears the repl window.
func cls(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	// We assume /c has been called from a toplevel repl.
	e := apl.EmptyArray{}
	if window == nil {
		return e, nil
	}
	repl, ok := window.Top.W.(*ui.Repl)
	if !ok {
		return e, nil
	}
	_, err := repl.Edit.Edit(",d")
	return e, err
}

// copied from apl/a/commands.go
type rw0 string

func (r rw0) Rewrite(t []scan.Token) []scan.Token {
	sym := scan.Token{T: scan.Identifier, S: "u→" + string(r)}
	num := scan.Token{T: scan.Number, S: "0"}
	tokens := make([]scan.Token, len(t)+2)
	tokens[0] = sym
	tokens[1] = num
	copy(tokens[2:], t)
	return tokens
}

type toCommand func([]scan.Token) []scan.Token

func (f toCommand) Rewrite(t []scan.Token) []scan.Token {
	return f(t)
}
