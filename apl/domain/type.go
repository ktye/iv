package domain

import (
	"reflect"

	"github.com/ktye/iv/apl"
)

// IsType tests if the value has a concrete type.
func IsType(t reflect.Type, child SingleDomain) SingleDomain {
	return typ{t, child, false}
}

// ToType tries converts the value to the given type.
func ToType(t reflect.Type, child SingleDomain) SingleDomain {
	return typ{t, child, true}
}

type typ struct {
	t     reflect.Type
	child SingleDomain
	conv  bool
}

func (t typ) String(f apl.Format) string {
	name := "type"
	if t.conv {
		name = "totype"
	}
	name += " " + t.t.String()
	if t.child == nil {
		return name
	}
	return name + " " + t.child.String(f)
}

func (t typ) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if reflect.TypeOf(V) == t.t {
		return propagate(a, V, t.child)
	}
	if t.conv == false {
		return V, false
	}
	// TODO: The function should also be used by the convert primitive (⌶).
	return V, false // TODO type conversions.
}
