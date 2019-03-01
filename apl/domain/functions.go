package domain

import (
	"github.com/ktye/iv/apl"
)

func Function(child SingleDomain) SingleDomain {
	return function{child}
}

type function struct{ child SingleDomain }

func (f function) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if _, ok := V.(apl.Function); ok == false {
		return V, false
	}
	if f.child == nil {
		return V, true
	}
	if v, ok := f.child.To(a, V); ok {
		return v, true
	}
	return V, false
}
func (f function) String(af apl.Format) string {
	if f.child == nil {
		return "function"
	}
	return "function " + f.child.String(af)
}

func IsPrimitive(p string) SingleDomain {
	return primitive(p)
}

type primitive string

func (p primitive) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if pf, ok := V.(apl.Primitive); ok && string(pf) == string(p) {
		return V, true
	}
	return V, false
}
func (p primitive) String(f apl.Format) string {
	return string(p)
}
