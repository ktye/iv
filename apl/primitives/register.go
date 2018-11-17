package primitives

import (
	"github.com/ktye/iv/apl"
)

func Register(a *apl.Apl) {
	for _, p := range primitives {
		a.RegisterPrimitive(apl.Primitive(p.symbol), p)
	}
}

var primitives []primitive

func register(p primitive) {
	primitives = append(primitives, p)
}

type primitive struct {
	apl.Domain
	symbol string
	doc    string
	fn     func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error)
}

func (p primitive) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) { return p.fn(a, L, R) }
func (p primitive) Doc() string                                        { return p.doc }
