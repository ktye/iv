package u

import (
	"fmt"
	"reflect"

	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/scan"
	"github.com/ktye/ui"
)

// Sam replaces the window with a full screen sam editor widget.
// On exit it puts the top window back.
//
// If it is a channel, all values are read, converted to strings and joined with newline.
//
// If the argument is a string that references an existing variable, the variable is used.
// Example print a lambda function f: /e`f
//
// Other values are converted to strings.
func sam(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	e := apl.EmptyArray{}
	if window == nil || window.Top.W == nil {
		return nil, fmt.Errorf("edt can only be called in a graphical session")
	}
	cmd, edt := read(a, R)

	sam := ui.NewSam(window)
	sam.SetTexts(cmd, edt)
	save := window.Top.W
	sam.Quit = func() ui.Event {
		window.Top.W = save
		return ui.Event{}
	}
	window.Top.W = sam
	window.Resize()
	return e, nil
}

// read returns the text for sam's command and edit window based
// on the value v.
func read(a *apl.Apl, R apl.Value) (cmd, edt rope.Rope) {
	switch v := R.(type) {
	case apl.Channel:
		edt = rope.New("")
		n := 0
		var t reflect.Type
		for e := range v[0] {
			if n == 0 {
				t = reflect.TypeOf(e)
			} else if t != nil {
				tt := reflect.TypeOf(e)
				if tt != t {
					fmt.Printf("types %v %v\n", tt, t)
					t = nil
				}
			}
			n++
			edt = rope.Append(edt, rope.New(e.String(a)+"\n"))
		}
		ts := "mixed"
		if t != nil {
			ts = t.String()
		}
		cmd := rope.New(fmt.Sprintf("from channel (%d %s)\n", n, ts))
		return cmd, edt
	case apl.String:
		if val := a.Lookup(string(v)); val != nil {
			return rope.New("var " + reflect.TypeOf(val).String()), rope.New(val.String(a))
		}
		return rope.New(reflect.TypeOf(R).String()), rope.New(R.String(a))
	default:
		shape := ""
		if ar, ok := R.(apl.Array); ok {
			shape = fmt.Sprintf(" %v", ar.Shape())
		}
		return rope.New(reflect.TypeOf(R).String() + shape), rope.New(R.String(a))
	}
}

func samCmd(t []scan.Token) []scan.Token {
	return append([]scan.Token{scan.Token{T: scan.Identifier, S: "u→sam"}}, t...)
}
