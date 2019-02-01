package u

import (
	"fmt"
	"reflect"

	"github.com/eaburns/T/edit"
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
	if l := cmd.Len(); l == 0 {
		cmd = rope.New("⍎ \nq\n")
	} else if r := rope.Slice(cmd, l-1, l); r.String() == "\n" {
		cmd = rope.Append(cmd, rope.New("⍎ \nq\n"))
	} else {
		cmd = rope.Append(cmd, rope.New("\n⍎ \nq\n"))
	}

	sam := ui.NewSam(window)
	sam.Cmd.SetText(cmd)
	sam.Edt.SetText(edt)
	save := window.Top.W

	dot := func(addr string) (string, bool) {
		if sam.Edt.MarkAddr(addr) != nil {
			return "", false
		}
		return sam.Edt.Selection(), true
	}

	sam.Commands = map[string]func(*ui.Sam, string){
		"q": func(s *ui.Sam, args string) {
			window.Top.W = save
		},
		"⍎": func(s *ui.Sam, args string) {
			t, _ := dot(".")
			if t == "" {
				t, _ = dot(",")
			}
			dsave := a.Lookup("Dot")
			out := a.GetOutput()
			defer func() {
				a.Assign("Dot", dsave)
				a.SetOutput(out)
			}()
			a.SetOutput(sam.Cmd)

			a.Assign("Dot", apl.String(t))
			_, err := sam.Cmd.Edit.Edit(`/\n$/`)
			if _, ok := err.(edit.NoCommandError); !ok {
				sam.Cmd.Write([]byte{'\n'})
			}
			if err := a.ParseAndEval(args); err != nil {
				sam.Cmd.Write([]byte(err.Error() + "\n"))
			}
		},
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
			s := fmt.Sprintf("var %s %s\n⍎ %s← ⍎Dot\n", v, reflect.TypeOf(val).String(), v)
			return rope.New(s), rope.New(val.String(a))
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
