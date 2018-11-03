package apl

import (
	"fmt"
	"strings"
)

// lambda is a function expression in braces {...}.
// It is also known under the term dynamic function or dfn.
type lambda struct {
	body guardList
}

func (λ *lambda) String(a *Apl) string {
	if λ.body == nil {
		return "{}"
	}
	return fmt.Sprintf("{%s}", λ.body.String(a))
}

func (λ *lambda) Eval(a *Apl) (Value, error) {
	return λ, nil
}

func (λ *lambda) Call(a *Apl, l, r Value) (Value, error) {
	if λ.body == nil {
		return EmptyArray{}, nil
	}
	a.Assign("⍺", l)
	a.Assign("⍵", r)
	return λ.body.Eval(a)
}

// guardList is the body of a lambda expression.
// It represents a list of guarded expressions.
type guardList []*guardExpr

func (l guardList) String(a *Apl) string {
	v := make([]string, len(l))
	for i, g := range l {
		v[i] = g.String(a)
	}
	return strings.Join(v, "⋄")
}

// Eval evaluates the guardList.
// It checks the condition of each guardExpr.
// For the first condition which returns true, or a nil condition,
// the expression is evaluated. Sequent guarded expressions are ignored.
func (l guardList) Eval(a *Apl) (Value, error) {
	if len(l) == 0 {
		return EmptyArray{}, nil
	}
	for _, g := range l {
		if v, err := g.Eval(a); err != nil {
			return nil, err
		} else if v != nil {
			return v, nil
		}
	}
	return EmptyArray{}, nil // TODO: should it be an error, if all conditions are false?
}

// guardExpr contains a guarded expression.
// It's expressions is evaluated if the condition returns true or is nil.
type guardExpr struct {
	cond expr
	e    expr
}

func (g *guardExpr) String(a *Apl) string {
	if g.cond == nil {
		return g.e.String(a)
	} else {
		return g.cond.String(a) + ":" + g.e.String(a)
	}
}

// Eval evaluates a guarded expression.
// If the condition exists, it is evaluated and must return a bool.
// If the condition is nil or returns true, the expression is evaluated,
// otherwise nil is returned and no error.
func (g *guardExpr) Eval(a *Apl) (Value, error) {
	if g.cond == nil {
		return g.e.Eval(a)
	}

	if v, err := g.cond.Eval(a); err != nil {
		return nil, err
	} else if b, ok := v.(Bool); ok == false {
		return nil, fmt.Errorf("λ condition does not return a bool: %s", b.String(a))
	} else if ok == false {
		return nil, nil
	} else {
		return g.e.Eval(a)
	}
}
