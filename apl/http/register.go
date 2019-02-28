package http

import (
	"fmt"
	"net/http"

	"github.com/ktye/iv/apl"
)

func Register(a *apl.Apl, name string) {
	pkg := map[string]apl.Value{
		"get": apl.ToFunction(get),
	}
	if name == "" {
		name = "http"
	}
	a.RegisterPackage(name, pkg)
}

// http→get returns a channel to read strings from a http connection.
func get(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
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
