package apl

import (
	"fmt"
	"io"
	"sort"
	"strings"
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

// registerSymbol adds single rune symbols for the parser.
func (a *Apl) registerSymbol(s string) {
	if r, w := utf8.DecodeRuneInString(s); w == len(s) {
		a.symbols[r] = s
	}
}

// RegisterPackage adds an external package to apl.
func (a *Apl) RegisterPackage(name string, m map[string]Value) {
	a.pkg[name] = &env{parent: nil, vars: m}
}

// Doc writes the documentation of all registered primitives and operators to the writer.
func (a *Apl) Doc(w io.Writer) {
	fmt.Fprintln(w, "## Primitive functions")
	fmt.Fprintln(w, "```")
	{

		s := make([]struct {
			symbol Primitive
			doc    string
		}, len(a.primitives))
		i := 0
		for symbol, handlers := range a.primitives {
			h := handlers[0]
			s[i].symbol = symbol
			s[i].doc = h.Doc()
			i++
		}
		sort.Slice(s, func(i, j int) bool { return s[i].doc < s[j].doc })
		for _, k := range s {
			symbol := k.symbol
			handlers := a.primitives[k.symbol]
			fmt.Fprintf(w, "%s\t\t\n", symbol)
			for _, h := range handlers {
				dom := h.String(a)
				domain := fmt.Sprintf("L%sR  %s", symbol, dom)
				if strings.Index(dom, "L") == -1 {
					domain = fmt.Sprintf("%sR  %s", symbol, dom)
				}
				fmt.Fprintf(w, "\t%s\n\t%s\t\n", h.Doc(), domain)
			}
			fmt.Fprintf(w, "\t\t\n")
		}

	}
	fmt.Fprintln(w, "```")

	fmt.Fprintln(w, "## Operators")
	fmt.Fprintln(w, "```")
	{
		s := make([]struct {
			symbol string
			doc    string
		}, len(a.operators))
		i := 0
		for symbol, ops := range a.operators {
			h := ops[0]
			s[i].symbol = symbol
			s[i].doc = h.Doc()
			i++
		}
		sort.Slice(s, func(i, j int) bool { return s[i].doc < s[j].doc })
		for _, k := range s {
			symbol := k.symbol
			handlers := a.operators[k.symbol]
			fmt.Fprintf(w, "%s\t\t\n", symbol)
			for _, h := range handlers {
				dom := h.String(a)
				domain := fmt.Sprintf("LO%sRO  %s", symbol, dom)
				if strings.Index(dom, "LO") == -1 {
					domain = fmt.Sprintf("%sRO  %s", symbol, dom)
				}
				fmt.Fprintf(w, "\t%s\n\t%s\t\n", h.Doc(), domain)
			}
			fmt.Fprintf(w, "\t\t\n")
		}

		fmt.Fprintln(w, "```")
	}
}
