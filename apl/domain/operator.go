package domain

import "github.com/ktye/iv/apl"

// MonadicOp is used to define a monadic operator.
func MonadicOp(l SingleDomain) apl.Domain {
	return monop{l}
}

type monop struct {
	left SingleDomain
}

func (m monop) To(a *apl.Apl, L, R apl.Value) (apl.Value, apl.Value, bool) {
	if m.left == nil {
		return L, R, true
	}
	if v, ok := m.left.To(a, L); ok {
		return v, R, true
	}
	return L, R, false
}
func (m monop) DyadicOp() bool { return false }
func (m monop) String(f apl.Format) string {
	if m.left == nil {
		return "LO any"
	}
	return "LO " + m.left.String(f)
}

// DyadicOp is used to define a dyadic operator.
func DyadicOp(child apl.Domain) apl.Domain {
	return dyop{child}
}

type dyop struct {
	child apl.Domain
}

func (d dyop) To(a *apl.Apl, L, R apl.Value) (apl.Value, apl.Value, bool) {
	if L == nil {
		return L, R, false
	}
	if d.child == nil {
		return L, R, true
	}
	return d.child.To(a, L, R)
}
func (d dyop) DyadicOp() bool { return true }
func (d dyop) String(f apl.Format) string {
	if d.child == nil {
		return "any"
	}
	return d.child.String(f)
}
