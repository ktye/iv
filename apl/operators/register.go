package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

// Register adds the operators in this package to a.
func Register(a *apl.Apl) {
	for _, op := range operators {
		a.RegisterOperator(op.symbol, op)
	}
}

type operator struct {
	apl.Domain
	symbol    string
	doc       string
	derived   func(*apl.Apl, apl.Value, apl.Value) apl.Function
	selection func(*apl.Apl, apl.Value, apl.Value, apl.Value, apl.Value) (apl.IndexArray, error)
}

func (op operator) Doc() string { return op.doc }
func (op operator) Derived(a *apl.Apl, LO, RO apl.Value) apl.Function {
	return op.derived(a, LO, RO)
}
func (op operator) Select(a *apl.Apl, L, LO, RO, R apl.Value) (apl.IndexArray, error) {
	if op.selection == nil {
		return apl.IndexArray{}, fmt.Errorf("operator %s cannot be used in selective assignment", op.symbol)
	} else {
		return op.selection(a, L, LO, RO, R)
	}
}
func (op operator) DyadicOp() bool {
	if ar, ok := op.Domain.(arity); ok {
		return ar.DyadicOp()
	}
	// Domain must start with MonadicOp or DyadicOp
	panic(op.symbol + ": operator Domain must start with MonadicOp or DyadicOp")
}

type arity interface {
	DyadicOp() bool
}

var operators []operator

func register(op operator) {
	operators = append(operators, op)
}

// primitive implements a SingleDomain, which tests against the prmitive symbol.
type primitive string

func (p primitive) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if pf, ok := V.(apl.Primitive); ok && pf == apl.Primitive(p) {
		return V, true
	}
	return V, false
}
func (p primitive) String(a *apl.Apl) string { return string(p) }

// function is both a func and implements the apl.Function interface,
// by calling itself.
// It is used to wrap derived functions to satisfy apl.Function.
type function func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error)

func (f function) Call(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	return f(a, l, r)
}
