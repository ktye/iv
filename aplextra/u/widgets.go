package u

import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/xgo"
	"github.com/ktye/ui"
)

// Button returns a button with the given string as a label.
func button(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	label, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("button: right argument must be a string")
	}
	icon := apl.String("")
	if L != nil {
		icon, ok = L.(apl.String)
		if ok == false {
			return nil, fmt.Errorf("button: left argument icon must be a string")
		}
	}
	b := &ui.Button{Text: string(label), Icon: string(icon)}
	return xgo.Value(reflect.ValueOf(b)), nil
}

// Split returns a split widget with L and R as childs.
func split(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if L == nil {
		return nil, fmt.Errorf("u split: function must be called dyadically")
	}
	k1, ok := toWidget(L)
	if ok == false {
		return nil, fmt.Errorf("u split: left argument must be a widget")
	}
	k2, ok := toWidget(R)
	if ok == false {
		return nil, fmt.Errorf("u split: right argument must be a widget")
	}

	s := &ui.Split{Kids: ui.NewKids(k1, k2)}
	return xgo.Value(reflect.ValueOf(s)), nil
}
