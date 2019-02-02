package u

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

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
func sam(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	// Callback for button-2 exec:
	// Open files from the apl file system.
	exec := func(sam *ui.Sam, s string) {
		if strings.HasPrefix(s, "/") == false {
			return
		}
		if strings.ContainsRune(s, '\n') {
			return
		}
		if strings.HasSuffix(s, "/") == false {
			// sam.Cmd.Write([]byte("w" + s + "\n"))
		}
		out := a.GetOutput()
		defer func() {
			a.SetOutput(out)
		}()
		a.SetOutput(ioutil.Discard)

		e := `"` + s + `"` + " u→sam <`" + s // TODO: s should be quoted, when quoted strings are supported.
		if err := a.ParseAndEval(e); err != nil {
			sam.Cmd.AppendText(err.Error())
		}
	}

	e := apl.EmptyArray{}
	if window == nil || window.Top.W == nil {
		return nil, fmt.Errorf("edt can only be called in a graphical session")
	}
	cmd, edt := read(a, R)
	if len(cmd) > 0 && cmd[len(cmd)-1] != '\n' {
		cmd += "\n"
	}
	cmd += "⍎ \nq\n"

	if L != nil {
		if s, ok := L.(apl.String); ok {
			cmd += "w" + string(s)
		}
	}

	sam := ui.NewSam(window)
	sam.Cmd.SetText(rope.New(cmd))
	sam.Edt.SetText(edt)
	sam.SetExec(exec)
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
			sam.Cmd.AppendText("")

			a.Assign("Dot", apl.String(t))
			if err := a.ParseAndEval(args); err != nil {
				sam.Cmd.AppendText(err.Error())
			}

		},
	}
	window.Top.W = sam
	window.Resize()
	return e, nil
}

// read returns the text for sam's command and edit window based
// on the value v.
func read(a *apl.Apl, R apl.Value) (cmd string, edt rope.Rope) {
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
		cmd := fmt.Sprintf("from channel (%d %s)\n", n, ts)
		return cmd, edt
	case apl.String:
		if val := a.Lookup(string(v)); val != nil {
			s := fmt.Sprintf("var %s %s\n⍎ %s← ⍎Dot\n", v, reflect.TypeOf(val).String(), v)
			return s, rope.New(val.String(a))
		}
		return reflect.TypeOf(R).String(), rope.New(R.String(a))
	default:
		shape := ""
		if ar, ok := R.(apl.Array); ok {
			shape = fmt.Sprintf(" %v", ar.Shape())
		}
		return reflect.TypeOf(R).String() + shape, rope.New(R.String(a))
	}
}

func samCmd(t []scan.Token) []scan.Token {
	return append([]scan.Token{scan.Token{T: scan.Identifier, S: "u→sam"}}, t...)
}
