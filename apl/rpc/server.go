package rpc

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/ktye/iv/apl"
)

// ListenAndServe puts APL into server mode.
// It accepts a single connection at a time from anyone.
func ListenAndServe(a *apl.Apl, addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Print("listen on ", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		} else {
			handle(a, conn)
		}
	}
}

type Request struct {
	Fn   string
	L, R apl.Value
}

type Response struct {
	Err string
	V   apl.Value
}

func handle(a *apl.Apl, cn net.Conn) {
	log.Print("conn ", cn.RemoteAddr())
	var req Request
	for {
		var res Response
		if err := gob.NewDecoder(cn).Decode(&req); err != nil {
			res.Err = err.Error()
		} else {
			if v, err := exec(a, req); err != nil {
				res.Err = err.Error()
			} else {
				res.V = v
			}
		}
		if err := gob.NewEncoder(cn).Encode(res); err != nil {
			log.Print(err)
			cn.Close()
			return
		}
	}
}

func exec(a *apl.Apl, req Request) (apl.Value, error) {
	if req.R == nil {
		return nil, fmt.Errorf("right argument is nil")
	}
	if p, err := a.Parse(req.Fn); err != nil {
		return nil, err
	} else if len(p) != 1 {
		return nil, fmt.Errorf("expected a single function expression: got %d", len(p))
	} else if v, err := p[0].Eval(a); err != nil {
		return nil, err
	} else if f, ok := v.(apl.Function); ok == false {
		return nil, fmt.Errorf("expr is not a function")
	} else {
		return f.Call(a, req.L, req.R)
	}
}
