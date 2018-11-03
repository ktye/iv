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
type Operator interface {
	IsDyadic() bool
	Apply(l, r Value) FunctionHandle
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
	o, ok := a.operators[d.op]
	if ok == false {
		return "<unknown operator>"
	}
	right := ""
	if o.IsDyadic() {
		right = fmt.Sprintf(" %s", d.ro.String(a))
	}
	left := "?"
	if d.lo != nil {
		left = fmt.Sprintf("%s", d.lo.String(a))
	}
	return fmt.Sprintf("(%s %s%s)", left, d.op, right)
}

func (d *derived) Call(a *Apl, l, r Value) (Value, error) {
	o, ok := a.operators[d.op]
	if ok == false {
		return nil, fmt.Errorf("operator %s does not exist", d.op)
	}

	h := o.Apply(d.lo, d.ro)
	if ok, v, err := h(a, l, r); ok == false {
		var ds string
		if o.IsDyadic() {
			ds = fmt.Sprintf("(%T %s %T)", d.lo, d.op, d.ro)
		} else {
			ds = fmt.Sprintf("(%T %s)", d.lo, d.op)
		}
		if l == nil {
			return nil, fmt.Errorf("derived monadic function %s cannot handle %T", ds, r)
		}
		return nil, fmt.Errorf("derived dyadic function %s cannot handle %T, %T", ds, l, r)

	} else {
		return v, err
	}
}
