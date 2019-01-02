package http

import (
	"fmt"
	"net/http"

	"github.com/ktye/iv/apl"
)

func Register(a *apl.Apl) {
	pkg := map[string]apl.Value{
		"get": get{},
	}
	a.RegisterPackage("http", pkg)
}

type get struct{}

func (_ get) String(a *apl.Apl) string {
	return "http get"
}

// http→get returns a channel to read strings from a http connection.
func (_ get) Call(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	addr, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("http get: right argument must be a string")
	}
	res, err := http.Get(string(addr))
	if err != nil {
		return nil, err
	}
	return apl.LineReader(res.Body), nil
}
