// Package rpc adds remove procedure calls to APL.
package rpc

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/ktye/iv/apl"
)

func Dial(address string) (Conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return Conn{}, err
	}
	return Conn{Conn: c}, nil
}

type Conn struct {
	net.Conn
}

func (c Conn) String(a *apl.Apl) string {
	remote := c.RemoteAddr()
	if remote == nil {
		return fmt.Sprintf("rpc→conn not connected")
	}
	return fmt.Sprintf("rpc→conn to %s", remote.String())
}
func (c Conn) Copy() apl.Value { return c }

func (c Conn) Close() (apl.Value, error) {
	if c.Conn == nil {
		return nil, fmt.Errorf("not connected")
	}
	if err := c.Conn.Close(); err != nil {
		return nil, err
	}
	return apl.Int(1), nil
}

func (c Conn) Call(f string, L, R apl.Value) (apl.Value, error) {
	req := Request{Fn: f, L: L, R: R}
	if err := gob.NewEncoder(c.Conn).Encode(req); err != nil {
		c.Conn.Close()
		return nil, err
	}
	var res Response
	if err := gob.NewDecoder(c.Conn).Decode(&res); err != nil {
		c.Conn.Close()
		return nil, err
	}
	if res.Err != "" {
		return nil, fmt.Errorf("%s", res.Err)
	} else if res.V == nil {
		return nil, fmt.Errorf("empty result")
	}
	return res.V, nil
}
