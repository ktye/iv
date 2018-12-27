package domain

import "github.com/ktye/iv/apl"

// IsObject accepts objects
func IsObject(child SingleDomain) SingleDomain {
	return objtype{child}
}

type objtype struct{ child SingleDomain }

func (s objtype) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if _, ok := V.(apl.Object); ok {
		return propagate(a, V, s.child)
	}
	return V, false
}
func (s objtype) String(a *apl.Apl) string {
	if s.child == nil {
		return "object"
	}
	return "object" + " " + s.child.String(a)
}
