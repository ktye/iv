package rpc

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

// Register adds the rpc package to the interpreter.
// Example:
//	Start a server on local port 1966 (cmd/apl):
//		apl ":1966"
//	On a different process, run a normal APL session:
//		C←rpc→dial ":1966"
//		rpc→call (C; "+/"; 5; ⍳10;)
//	15 20 25 30 35 40
// The rpc call evaluates the function string in the remove environment
// and calls it with the local values on the remote process.
//
// APL value that should be transfered over the wire need to be
// registerd to the gob package.
// See init.go for values that are already registerd.
func Register(a *apl.Apl, name string) {
	pkg := map[string]apl.Value{
		"dial":  dial{},
		"call":  call{},
		"close": closeconn{},
	}
	if name == "" {
		name = "rpc"
	}
	a.RegisterPackage(name, pkg)
}

type dial struct{}

func (_ dial) String(a *apl.Apl) string {
	return "rpc dial"
}

func (_ dial) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if L != nil {
		return nil, fmt.Errorf("rpc dial must be called monadically")
	}
	if s, ok := R.(apl.String); ok == false {
		return nil, fmt.Errorf("rpc dial: argument must be a string")
	} else {
		return Dial(string(s))
	}
}

type call struct{}

func (_ call) String(a *apl.Apl) string {
	return "rpc call"
}

func (_ call) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if L != nil {
		return nil, fmt.Errorf("rpc call must be called monadically")
	}
	lst, ok := R.(apl.List)
	if ok == false {
		return nil, fmt.Errorf("rpc call: argument must be a list: %T", R)
	}
	if len(lst) < 3 {
		return nil, fmt.Errorf("rpc call: argument list is too short")
	}
	if len(lst) > 4 {
		return nil, fmt.Errorf("rpc call: argument list is too long")
	}
	c, ok := lst[0].(Conn)
	if ok == false {
		return nil, fmt.Errorf("rpc call: first list argument must be a connection")
	}
	f, ok := lst[1].(apl.String)
	if ok == false {
		return nil, fmt.Errorf("rpc call: second list argument must be a string")
	}

	var Larg, Rarg apl.Value
	if len(lst) == 3 {
		Rarg = lst[2]
	} else {
		Larg = lst[2]
		Rarg = lst[3]
	}
	return c.Call(string(f), Larg, Rarg)
}

type closeconn struct{}

func (_ closeconn) String(a *apl.Apl) string {
	return "rpc close"
}

func (_ closeconn) Call(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	c, ok := R.(Conn)
	if ok == false {
		return nil, fmt.Errorf("right argument must be a connection")
	}
	return c.Close()
}
