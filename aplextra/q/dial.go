package q

import (
	"fmt"

	"github.com/ktye/iv/apl"
	kdb "github.com/sv/kdbgo"
)

func Dial(host string, port int) (apl.Value, error) {
	c, err := kdb.DialKDB(host, port, "")
	if err != nil {
		return nil, err
	}
	return Conn{c}, nil
}

type Conn struct {
	*kdb.KDBConn
}

func (c Conn) Copy() apl.Value { return c }
func (c Conn) String(a *apl.Apl) string {
	if c.KDBConn == nil {
		return "q not connected"
	}
	return fmt.Sprintf("q connection to %s:%s", c.Host, c.Port)
}

func (c Conn) Close() (apl.Value, error) {
	if c.KDBConn == nil {
		return nil, fmt.Errorf("not connected")
	}
	if err := c.KDBConn.Close(); err != nil {
		return nil, err
	}
	return apl.Int(1), nil
}

func (c Conn) Call(a *apl.Apl, cmd string, args apl.List) (apl.Value, error) {
	if c.KDBConn == nil {
		return nil, fmt.Errorf("not connected")
	}

	k := make([]*kdb.K, len(args))
	for i, v := range args {
		q, err := FromAPL(a, v)
		if err != nil {
			return nil, err
		}
		k[i] = q
	}

	if res, err := c.KDBConn.Call(cmd, k...); err != nil {
		return nil, err
	} else {
		return ToAPL(res)
	}
}
