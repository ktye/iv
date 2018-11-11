package domain

import "github.com/ktye/iv/apl"

// IsString accepts strings
func IsString(child SingleDomain) SingleDomain {
	return stringtype{child}
}

type stringtype struct{ child SingleDomain }

func (s stringtype) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if v, ok := V.(apl.String); ok {
		if s.child == nil {
			return v, true
		}
		return s.child.To(a, v)
	}
	return V, false
}
func (s stringtype) String(a *apl.Apl) string {
	if s.child == nil {
		return "string"
	}
	return "string" + " " + s.child.String(a)
}
