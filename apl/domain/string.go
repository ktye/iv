package domain

import (
	"github.com/ktye/iv/apl"
)

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

// IsStrings accepts uniform.Strings
func IsStrings(child SingleDomain) SingleDomain {
	return stringstype{child, false}
}

type stringstype struct {
	child   SingleDomain
	convert bool
}

func (s stringstype) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if _, ok := V.(apl.Strings); ok {
		return propagate(a, V, s.child)
	} else {
		if s.convert == false {
			return V, false
		}
		if str, ok := V.(apl.String); ok {
			return propagate(a, apl.Strings{
				Dims:    []int{1},
				Strings: []string{string(str)},
			}, s.child)

		} else if ar, ok := V.(apl.Array); ok {
			str := make([]string, ar.Size())
			for i := range str {
				if v, err := ar.At(i); err != nil {
					return V, false
				} else if sv, ok := v.(apl.String); ok {
					str[i] = string(sv)
				} else {
					return V, false
				}
			}
			return propagate(a, apl.Strings{
				Dims:    apl.CopyShape(ar),
				Strings: str,
			}, s.child)
		} else {
			return V, false
		}
	}
}
func (s stringstype) String(a *apl.Apl) string {
	name := "string array"
	if s.convert {
		name = "to string array"
	}
	if s.child == nil {
		return name
	}
	return name + " " + s.child.String(a)
}
