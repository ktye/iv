package u

import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/xgo"
	"github.com/ktye/ui"
)

// setCallback sets the callback function of a widget to an APL function.
// L ist a widget and R a function boxed in a list.
func setCallback(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	w, ok := toWidget(L)
	if ok == false {
		return nil, fmt.Errorf("u f: left argument must be a widget")
	}

	lst, ok := R.(apl.List)
	if ok == false {
		return nil, fmt.Errorf("u f: right argument must be a list")
	}
	if len(lst) == 1 {
		if fn, ok := lst[0].(apl.Function); ok {
			return setcb(a, w, fn)
		}
	}
	return nil, fmt.Errorf("u f: argument must be a list containing a single function")
}

func setcb(a *apl.Apl, w ui.Widget, f apl.Function) (apl.Value, error) {
	switch v := w.(type) {
	case *ui.Button:
		fmt.Printf("u: set callback of button %p\n", v)
		b := xgo.Value(reflect.ValueOf(v))
		v.Click = func() ui.Event {
			fmt.Println("u: click button")
			f.Call(a, nil, b)
			window.Top.Layout = ui.Dirty
			return ui.Event{Consumed: true}
		}
		return b, nil
	default:
		return nil, fmt.Errorf("u f: cannot set callback to a widget of type %T", w)
	}
}
