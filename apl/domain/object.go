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
func (s objtype) String(f apl.Format) string {
	if s.child == nil {
		return "object"
	}
	return "object" + " " + s.child.String(f)
}

// IsTable accepts objects
func IsTable(child SingleDomain) SingleDomain {
	return table{child}
}

type table struct{ child SingleDomain }

func (s table) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if _, ok := V.(apl.Table); ok {
		return propagate(a, V, s.child)
	}
	return V, false
}
func (s table) String(f apl.Format) string {
	if s.child == nil {
		return "table"
	}
	return "table" + " " + s.child.String(f)
}
