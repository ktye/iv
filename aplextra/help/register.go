package help

import "github.com/ktye/iv/apl"

func Register(a *apl.Apl) {
	if err := a.Assign("help", help{}); err != nil {
		panic(err)
	}
	a.RegisterDoc("help", `help primitive function: help, documentation
Z←help R  R numeric, Z string
	R is the key in the doc database
	Z overview with all keys, primitives and operators
	
Z←help R  R string, Z string
	Z is the help text for the key R
	
Z←L help R  L string, R any, Z string
	Z is the accumulated help text for each entry who's first line
	matches the query
`)
	a.RegisterDoc("keyboard", "keyboard\n"+Keyboard)
}
