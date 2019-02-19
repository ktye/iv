// Package q provides an rpc interface from APL to kdb
//
// Example:
//	In q listen on port 1993:
//		q)\p 1993
//	In APL connect:
//		C←q→dial ":1993"
//	Make a function call:
//		q→call (C; "sum"; 3 3⍴⍳9;)  ⍝ pass an array
//	(1 2 3;4 5 6;7 8 9;)                 ⍝ the result is a list
//
//		q→call (C; "{n where 2 = sum 0 = n mod\:/: n:1 + til x}"; 50;)
//	2 3 5 7 11 13 17 19 23 29 31 37 41 43 47
//
//		D←`a`b`c#1 2 3		⍝ pass a dictionary
//		q→call (C;"sum";D;)
//	6
//
//		q→call (C; "!"; `a`b`c ; (1;2;3;);)
//	a: 1				⍝ result is a dictionary
//	b: 2
//	c: 3
//
//		T←⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)
//		q→call (C; "sum"; T;)   ⍝ pass a table
//	a: 6
//	b: 15
//	c: 24
//
//		q→call (C; "{select a,c from x where b>4}"; T ;)
//	a c				⍝ result is a table
//	2 8
//	3 8
package q

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ktye/iv/apl"
	kdb "github.com/sv/kdbgo"
)

func Register(a *apl.Apl, name string) {
	if name == "" {
		name = "q"
	}
	pkg := map[string]apl.Value{
		"dial":  apl.ToFunction(dial),
		"call":  apl.ToFunction(call),
		"test":  apl.ToFunction(test),
		"close": apl.ToFunction(closeconn),
	}
	a.RegisterPackage(name, pkg)
}

func dial(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	s, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("q dial: argument must be a string")
	}
	hostport := string(s)
	idx := strings.Index(hostport, ":")
	if idx == -1 {
		return nil, fmt.Errorf("q dial: right argument must contain a colon")
	}
	host := hostport[:idx]
	port, err := strconv.Atoi(hostport[idx+1:])
	if err != nil {
		return nil, fmt.Errorf("q dial: cannot parse port number")
	}
	return Dial(host, port)
}

func call(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	lst, ok := R.(apl.List)
	if ok == false {
		return nil, fmt.Errorf("q call: argument must be a list: %T", R)
	}
	if len(lst) < 2 {
		return nil, fmt.Errorf("q call: argument is too short")
	}
	c, ok := lst[0].(Conn)
	if ok == false {
		return nil, fmt.Errorf("q call: first list argument must be a connection")
	}
	cmd, ok := lst[1].(apl.String)
	if ok == false {
		return nil, fmt.Errorf("q call: second argument must be a string")
	}
	if len(lst) == 2 {
		return c.Call(a, string(cmd), nil)
	}
	return c.Call(a, string(cmd), lst[2:])
}

func closeconn(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	c, ok := R.(Conn)
	if ok == false {
		return nil, fmt.Errorf("right argument must be a connection")
	}
	return c.Close()
}

func test(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	c, ok := R.(Conn)
	if ok == false {
		return nil, fmt.Errorf("q test: right argument must be a connection")
	}
	if c.KDBConn == nil {
		return nil, fmt.Errorf("q: not connected")
	}
	res, err := c.KDBConn.Call("til", kdb.Int(10))
	if err != nil {
		return nil, err
	}
	return apl.String(fmt.Sprintf("%v", res)), nil
}
