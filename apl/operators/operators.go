package operators

import "github.com/ktye/iv/apl"

// Register adds the operators in this package to a.
func Register(a *apl.Apl) {
	for _, i := range operators {
		a.RegisterOperator(i.s, i.op)
	}
	for _, d := range doc {
		a.RegisterDoc(d[0], d[1])
	}
}

type dyadic struct{}

func (d dyadic) IsDyadic() bool { return true }

type monadic struct{}

func (m monadic) IsDyadic() bool { return false }

// function is both a func and implements the apl.Function interface,
// by calling itself.
// It is used to wrap derived functions to satisfy apl.Function.
type function func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error)

func (f function) Call(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	return f(a, l, r)
}

type regop struct {
	s  string
	op apl.Operator
}

var operators []regop

func register(s string, op apl.Operator) {
	operators = append(operators, regop{s, op})
}

var doc [][2]string

func addDoc(key, text string) {
	doc = append(doc, [2]string{key, text})
}
