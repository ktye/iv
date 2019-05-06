package u

import (
	"fmt"
	"os"
	"reflect"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/xgo"
	"github.com/ktye/iv/cmd/lui/apl385"
	"github.com/ktye/ui"
)

var window *ui.Window

// Loop creates the window and runs the main loop.
// It does not return.
func Loop(top ui.Widget) {
	w := ui.New(nil)
	w.SetKeyTranslator(ui.AplKeyboard{})
	w.SetFont(apl385.TTF(), 20)
	w.Top.W = top
	w.Render()
	window = w
	for {
		select {
		case e := <-w.Inputs:
			w.Input(e)

		case err, ok := <-w.Error:
			if !ok {
				return
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func toWidget(v apl.Value) (ui.Widget, bool) {
	if val, ok := v.(xgo.Value); ok {
		w := reflect.Value(val).Interface()
		if t, ok := w.(ui.Widget); ok {
			return t, true
		}
	}
	return nil, false
}

// Top sets the top widget. If the widget is the empty array, it returns the top widget.
func top(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	if window == nil {
		return nil, fmt.Errorf("u: there is no window")
	}
	if _, ok := R.(apl.EmptyArray); ok {
		return xgo.Value(reflect.ValueOf(window.Top.W)), nil
	}

	if t, ok := toWidget(R); ok {
		// TODO: set the window concurrently in a callback.
		window.Top.W = t
		window.Top.Layout = ui.Dirty
		window.Top.Draw = ui.Dirty
		window.Render()
		return R, nil
	}
	return nil, fmt.Errorf("u top: cannot set top widget to %T", R)
}
