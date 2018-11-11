package apl

import (
	"fmt"
	"unicode/utf8"
)

// RegistersPrimitive attaches the primitive handler h to the symbol p.
// If the symbol exists already, it is overloaded.
// When the function is applied, the last registered handle is tested
// first, if the arguments match to the domain of the handler.
func (a *Apl) RegisterPrimitive(p Primitive, h PrimitiveHandler) {
	a.primitives[p] = append([]PrimitiveHandler{h}, a.primitives[p]...)
	a.registerSymbol(string(p))
}

// RegisterOperator registers s as the symbol for the operator.
func (a *Apl) RegisterOperator(s string, op Operator) error {
	if op == nil {
		return fmt.Errorf("cannot register a nil operator to %s", s)
	}
	if ops, ok := a.operators[s]; ok && ops[0].DyadicOp() != op.DyadicOp() {
		return fmt.Errorf("cannot register operator %s with differing arity", s)
	}
	a.operators[s] = append([]Operator{op}, a.operators[s]...)
	a.registerSymbol(s)
	return nil
}

// RegisterDoc adds help text to a key in the documentation.
// This should be called by packages, that add primitives or operators.
func (a *Apl) RegisterDoc(key, help string) {
	if s := a.doc[key]; s == "" {
		a.doc[key] = help
	} else {
		a.doc[key] += "\n" + help
	}
}

// registerSymbol adds single rune symbols for the parser.
func (a *Apl) registerSymbol(s string) {
	if r, w := utf8.DecodeRuneInString(s); w == len(s) {
		a.symbols[r] = s
	}
}
