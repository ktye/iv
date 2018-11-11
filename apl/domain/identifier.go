package domain

import "github.com/ktye/iv/apl"

func IsIdentifier(child SingleDomain) SingleDomain {
	return identifier{child}
}

type identifier struct{ child SingleDomain }

func (id identifier) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	_, ok := V.(apl.Identifier)
	if ok == false {
		return V, false
	}
	return propagate(a, V, id.child)
}

func (id identifier) String(a *apl.Apl) string {
	name := "identfier"
	if id.child == nil {
		return name
	}
	return name + " " + id.child.String(a)
}
