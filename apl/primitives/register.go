package primitives

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ktye/iv/apl"
)

func Register(a *apl.Apl) {
	for _, p := range primitives {
		a.RegisterPrimitive(apl.Primitive(p.symbol), p)
	}
}

var primitives []primitive

func register(p primitive) {
	// Add source path to documentation.
	_, fn, line, _ := runtime.Caller(1)
	if idx := strings.Index(fn, "apl/primitives"); idx != -1 {
		fn = fn[idx:]
	}
	p.doc += fmt.Sprintf("\t%s:%d", fn, line)

	primitives = append(primitives, p)
}

type primitive struct {
	apl.Domain
	symbol string
	doc    string
	fn     func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error)
	sel    func(*apl.Apl, apl.Value, apl.Value) (apl.IndexArray, error)
}

func (p primitive) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) { return p.fn(a, L, R) }
func (p primitive) Select(a *apl.Apl, L, R apl.Value) (apl.IndexArray, error) {
	if p.sel == nil {
		return apl.IndexArray{}, fmt.Errorf("primitive %s cannot be used in selective assignment", p.symbol)
	}
	return p.sel(a, L, R)
}
func (p primitive) Doc() string { return p.doc }
