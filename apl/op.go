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
	Select(*Apl, Value, Value, Value, Value) (IntArray, error)
	Doc() string
}

// derived is a function which is derived from an operator and one or two arguments,
// which may be functions or arrays
type derived struct {
	op string
	// operands of the derived expression
	lo  expr                                       // left operand
	ro  expr                                       // right operand
	sel func(*Apl, Value, Value) (IntArray, error) // selection function for reduce and scan
}

func (d *derived) Eval(a *Apl) (Value, error) {
	return d, nil
}

func (d *derived) String(f Format) string {
	left := ""
	right := ""
	if d.lo == nil && d.ro == nil {
		return d.op
	}
	if d.lo != nil {
		left = d.lo.String(f) + " "
	}
	if d.ro != nil {
		right = " " + d.ro.String(f)
	}
	return "(" + left + d.op + right + ")"
}
func (d *derived) Copy() Value { return d }

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
		// Modified assignment contains the expr with the identifier in the left argument,
		// otherwise it is the LO.
		if l == nil {
			lo, err = evalAssign(a, d.lo, nil)
			if err != nil {
				return nil, err
			}
		} else {
			as, ok := l.(assignment)
			if ok == false {
				return nil, fmt.Errorf("modified assignment: expected assignment target expr on the left: %T", l)
			}
			var f Function
			if d.lo != nil {
				if pf, ok := d.lo.(Function); ok {
					f = pf
				} else {
					return nil, fmt.Errorf("modifier is not a function: %T", d.lo)
				}
			}
			lo, err = evalAssign(a, as, f)
			if err != nil {
				return nil, err
			}
			l = nil
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

func (d *derived) Select(a *Apl, L, R Value) (Value, error) {
	ops, ok := a.operators[d.op]
	if ok == false || len(ops) == 0 || ops[0] == nil {
		return nil, fmt.Errorf("operator %s does not exist", d.op)
	}

	if ops[0].DyadicOp() && d.op != "⍂" {
		// Scan and reduce are monadic, indexing can be used.
		return nil, fmt.Errorf("dyadic operators cannot be used in selective assignments")
	}

	var RO, LO Value
	var err error
	LO, err = d.lo.Eval(a)
	if err != nil {
		return nil, err
	}
	if ops[0].DyadicOp() { // All registerd operators have the same arity.
		RO, err = d.ro.Eval(a)
		if err != nil {
			return nil, err
		}
	}

	for _, op := range ops {
		if LO, RO, ok := op.To(a, LO, RO); ok {
			return op.Select(a, L, LO, RO, R)
		}
	}
	return nil, fmt.Errorf("cannot select with operator %T %T %s %T %T", L, LO, d.op, RO, R)
}
