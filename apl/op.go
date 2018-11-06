package apl

import "fmt"

/* todo loop op (dyadic op)
(l)5⟳{2+2}r // loop 5 times over the lambda
(l)({2+2}⟳5)r // inverse args, same
(1 2 3)⟳{...}r // index args: use idx as alpha to lambda
l(x>0)⟳{...}r // if cond exec once
*/

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
// Apply receives the left and possibly right operands to the operator.
// It returns true and a derived function, if it can handle the input types.
//
// If multiple operator handlers are registerd for a symbol (operator overloading), they
// all must have the same arity.
// The first operator registered determines the arity that all others have to follow.
type Operator interface {
	IsDyadic() bool
	Apply(lo, ro Value) (bool, Function)
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
	if ops[0].IsDyadic() {
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
	if ops[0].IsDyadic() { // The first registerd operator decides arity.
		ro, err = d.ro.Eval(a)
		if err != nil {
			return nil, err
		}
	}
	lo, err = d.lo.Eval(a)
	if err != nil {
		return nil, err
	}

	for i := len(ops) - 1; i >= 0; i-- {
		if ok, df := ops[i].Apply(lo, ro); ok {
			return df.Call(a, l, r)
		}
	}
	return nil, fmt.Errorf("cannot handle operator %T %s %T", lo, d.op, ro)
}
