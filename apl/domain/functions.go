package domain

import "github.com/ktye/iv/apl"

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
func (f function) String(a *apl.Apl) string {
	if f.child == nil {
		return "function"
	}
	return "function " + f.child.String(a)
}
