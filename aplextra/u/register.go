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
	"github.com/ktye/ui"
)

// Register adds the ui package to the interpreter with the default name u.
func Register(a *apl.Apl, name string) {
	if name == "" {
		name = "u"
	}
	pkg := map[string]apl.Value{
		"kb":     apl.ToFunction(kb),
		"top":    apl.ToFunction(top),
		"f":      apl.ToFunction(setCallback),
		"button": apl.ToFunction(button),
		"split":  apl.ToFunction(split),
	}
	a.RegisterPackage(name, pkg)
}

// Kb returns the APL keyboard layout as a string.
func kb(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return apl.String(ui.AplKeyboard{}.String()), nil
}
