package apl

import "fmt"

// Operators take functions or arrays as operands and produce derived functions.
// An operator can be monadic or dyadic but is never ambivalent.
// Their derived functions are.
//
// Operators have long scope on the left and short scope on the right.
//
// Example
//	/ is a monadic operator
//	+/ is a derived function (summation), which is monadic: +/1 2 3 4
//	2+/ is a derived function (n-wise summation), which is dyadic: 2 +/ 1 2 3 4
//
// Derived receives the left and possibly right operands to the operator
// and returns a derived function.
// It is called only, if Domain.To returns true.
//
// If multiple operator handlers are registerd for a symbol (operator overloading), they
// all must have the same arity.
// The first operator registered determines the arity that all others have to follow.
type Operator interface {
	Domain
	DyadicOp() bool
	Derived(*Apl, Value, Value) Function
	Doc() string
}

// derived is a function which is derived from an operator and one or two arguments,
// which may be functions or arrays
type derived struct {
	op string
	// operands of the derived expression
	lo expr // left operand
	ro expr // right operand
}

func (d *derived) Eval(a *Apl) (Value, error) {
	return d, nil
}

func (d *derived) String(a *Apl) string {
	ops, ok := a.operators[d.op]
	if ok == false {
		return "<unknown operator>"
	}
	right := ""
	if ops[0].DyadicOp() {
		right = fmt.Sprintf(" %s", d.ro.String(a))
	}
	left := "?"
	if d.lo != nil {
		left = fmt.Sprintf("%s", d.lo.String(a))
	}
	return fmt.Sprintf("(%s %s%s)", left, d.op, right)
}

// Call tries to call a derived function.
// l and r are the left and right values to the derived function.
// The left and right operands are stored at d.lo and d.ro.
// The call tries all registered handlers for the given operator in reverse
// registration order until a handler accepts to build a derived function, which
// is then called with l and r.
func (d *derived) Call(a *Apl, l, r Value) (Value, error) {
	ops, ok := a.operators[d.op]
	if ok == false || len(ops) == 0 || ops[0] == nil {
		return nil, fmt.Errorf("operator %s does not exist", d.op)
	}

	// Evaluate the operands.
	var ro, lo Value
	var err error
	if ops[0].DyadicOp() { // All registerd operators have the same arity.
		ro, err = d.ro.Eval(a)
		if err != nil {
			return nil, err
		}
	}

	// Assignment is special: It does not evaluate the Identifier.
	if d.op == "←" {
		lo, err = d.evalAssign()
		if err != nil {
			return nil, err
		}
	} else {
		lo, err = d.lo.Eval(a)
		if err != nil {
			return nil, err
		}
	}

	for _, op := range ops {
		if LO, RO, ok := op.To(a, lo, ro); ok {
			return op.Derived(a, LO, RO).Call(a, l, r)
		}
	}
	return nil, fmt.Errorf("cannot handle operator %T %s %T", lo, d.op, ro)
}
