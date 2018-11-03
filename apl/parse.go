package apl

import (
	"fmt"

	"github.com/ktye/iv/apl/scan"
)

type parser struct {
	a      *Apl
	tokens []scan.Token
	pos    int
}

func (p *parser) parse(tokens []scan.Token) (Program, error) {
	p.tokens = tokens
	p.pos = 0
	exprs, err := p.expressionList()
	return Program(exprs), err
}

func (p *parser) next() scan.Token {
	if p.pos >= len(p.tokens) {
		return scan.Token{}
	}
	p.pos++
	return p.tokens[p.pos-1]
}

func (p *parser) peek() scan.Token {
	if p.pos >= len(p.tokens) {
		return scan.Token{}
	}
	return p.tokens[p.pos]
}

// expressionList:
//	expr
//	expr ⋄ expr
func (p *parser) expressionList() ([]expr, error) {
	enter("expressionList", p.peek())
	defer leave("expressionList")

	var exprs []expr
	e, ok, err := p.expr()
	if err != nil {
		return nil, err
	} else if ok == false {
		return nil, fmt.Errorf("expected next expression, got: %s", p.peek())
	}
	if e != nil {
		exprs = []expr{e}
	}

	if t := p.peek(); t.T == scan.Diamond {
		p.next()
		more, err := p.expressionList()
		if err != nil {
			return nil, err
		} else {
			exprs = append(exprs, more...)
			return exprs, nil
		}
	} else if t.T == scan.Endl {
		return exprs, nil
	} else {
		return nil, fmt.Errorf("unexpected token: %s", t)
	}
}

// expr
//	operand
//	function
//	operand [ idxSpec ]    (indexing)
//	operand function expr  (dyadic context)
func (p *parser) expr() (expr, bool, error) {
	enter("expr", p.peek())
	defer leave("expr")

	switch t := p.peek(); t.T {
	case scan.RightBrack, scan.Colon, scan.Semicolon:
		return nil, false, nil
	}

	save := p.pos

	left, ook, err := p.operand()
	if err != nil {
		return nil, false, err
	}

	// An index specification follows the operand.
	if t := p.peek(); t.T == scan.LeftBrack && ook {
		p.next()
		if spec, err := p.idxSpec(); err != nil {
			return nil, false, err
		} else if t := p.peek(); t.T != scan.RightBrack {
			return nil, false, fmt.Errorf("index spec is not followed by ], but %s", t)
		} else {
			p.next()
			f := Primitive("[")
			return bind(f, left, spec)
		}
	} else if ook && (t.T == scan.RightBrack || t.T == scan.Semicolon) {
		return left, true, nil
	}

	f, fok, err := p.function(false)
	if err != nil {
		return nil, false, err
	}

	if ook == false && fok == false {
		return nil, false, nil
	}
	if ook == false {
		return f, true, nil
	}

	t := p.peek()
	switch t.T {
	case scan.Diamond, scan.Colon, scan.RightParen, scan.RightBrace, scan.Endl:
		if ook == true && fok == true {
			return nil, false, fmt.Errorf("unexpected %T after operand.function", t)
		}
		if ook == false {
			return f, true, nil
		}
		return left, true, nil
	}

	if ook == false {
		return nil, false, fmt.Errorf("function alone, was not followed by end but: %T", t)
	}

	right, ok, err := p.expr()
	if err != nil {
		return nil, false, err
	} else if ok == false {
		p.pos = save
		return nil, false, nil
	}

	return bind(f, left, right)
}

// operand
//	function expr (monadic context)
//	variable
//	number
//	string
//	chars
//	vector
//	TODO operand [...]
func (p *parser) operand() (expr, bool, error) {
	enter("operand", p.peek())
	defer leave("operand")

	save := p.pos
	if f, ok, err := p.function(false); err != nil {
		return nil, false, err
	} else if ok {
		if right, ok, err := p.expr(); err != nil {
			return nil, false, err
		} else if ok == false {
			p.pos = save
			return nil, false, nil
		} else {
			if fn, _, err := bind(f, nil, right); err != nil {
				return nil, false, err
			} else {
				return fn, true, nil
			}
		}
	}

	t := p.peek()
	if ok, isfunc := isVarname(t.S); ok && isfunc && t.T == scan.Identifier {
		p.next()
		return fnVar(t.S), true, nil
	}

	switch t.T {
	case scan.Identifier, scan.Number, scan.String, scan.Chars, scan.LeftParen:
		if ar, ok, err := p.array(); err != nil {
			return nil, false, err
		} else if ok {
			return ar, true, nil
		} else {
			return nil, false, nil
		}
	}

	return nil, false, nil
}

// function
//	primitive
//	lambda
//	variable (lowercase)
//	array|function operator
//	array|function operator function
//
// If asRightOperand is true, the function being parsed is taken as a right operand to an operator.
// Due to the right operand binding, this is the immediate expression and not a derived function.
func (p *parser) function(asRightOperand bool) (expr, bool, error) {
	enter("function", p.peek())
	defer leave("function")

	switch t := p.peek(); t.T {
	case scan.Colon, scan.RightBrack, scan.Semicolon:
		return nil, false, nil
	}

	// While an operator follows, it is a derived function.
	collectOperators := func(f expr) (expr, error) {
		for {
			if op, isop := p.operator(); isop {
				if d, err := p.derived(f, op); err != nil {
					return nil, err
				} else {
					f = d
				}
			} else {
				return f, nil
			}
		}
	}

	if asRightOperand == false {
		// If next is an array, it might be a derived function.
		save := p.pos
		if ar, ok, err := p.array(); err != nil {
			return nil, false, err
		} else if ok {
			if op, isop := p.operator(); isop {
				if d, err := p.derived(ar, op); err != nil {
					return nil, false, err
				} else {
					f, err := collectOperators(d)
					return f, true, err
				}
			} else {
				p.pos = save
				return nil, false, nil
			}
		}
	}

	// Next must be a lambda, primitive function or a function variable.
	var f expr
	if λ, ok, err := p.lambda(); err != nil {
		return nil, false, err
	} else if ok {
		f = λ
	}
	if f == nil {
		t := p.peek()
		switch t.T {
		case scan.Symbol:
			if _, ok := p.a.primitives[Primitive(t.S)]; ok {
				p.next()
				f = Primitive(t.S)
			}

		case scan.Identifier:
			if ok, isfunc := isVarname(t.S); ok == false {
				return nil, false, fmt.Errorf("identifier is not allowed as a variable name: %s", t.S)
			} else if isfunc {
				p.next()
				// A function identifier is only a function, if no assignment follows.
				if t := p.peek(); t.T == scan.Symbol && t.S == "←" {
					p.pos--
					return nil, false, nil
				}
				f = fnVar(t.S)
			}
		}
	}
	if f == nil {
		return nil, false, nil
	}

	if asRightOperand {
		return f, true, nil
	}

	f, err := collectOperators(f)
	return f, true, err
}

// Derived returns a derived function from the left operand, the operator and a possible right operand,
// if the operator is dyadic.
func (p *parser) derived(lo expr, opsym string) (expr, error) {
	enter("derived", p.peek())
	defer leave("derived")

	op, ok := p.a.operators[opsym]
	if ok == false {
		return nil, fmt.Errorf("not an operator: %s", opsym)
	}

	if op.IsDyadic() == false {
		return &derived{
			op: opsym,
			lo: lo,
		}, nil
	}

	// A dyadic operator needs a function (or an array).
	// In this implementation, we allow only functions as right operands.

	ro, isfunc, err := p.function(true)
	if err != nil {
		return nil, err
	}
	if isfunc == false {
		return nil, fmt.Errorf("dyadic operator expected function as right operand, got %T", p.peek())
	}

	return &derived{
		op: opsym,
		lo: lo,
		ro: ro,
	}, nil
}

// idxSpec
//	expr
//	idxSpec ; expr
func (p *parser) idxSpec() (idxSpec, error) {
	enter("idxSpec", p.peek())
	defer leave("idxSpec")

	var idx idxSpec
	for {
		if e, ok, err := p.expr(); err != nil {
			return nil, err
		} else if ok == false {
			if p.peek().T == scan.Semicolon {
				p.next()
				idx = append(idx, EmptyArray{})
			} else {
				if p.peek().T == scan.RightBrack {
					idx = append(idx, EmptyArray{})
				}
				break
			}
		} else {
			idx = append(idx, e)
			t := p.peek()
			if t.T == scan.RightBrack {
				break
			} else if t.T != scan.Semicolon {
				return nil, fmt.Errorf("index specification is not separated by ; but %v", t)
			}
			p.next()
		}
	}
	if idx == nil {
		return nil, fmt.Errorf("empty index specification") // Or is this valid?
	}
	return idx, nil
}

// Bind attaches the left and right argument to a function.
// Left may be nil in a monadic context.
// The bool return value is a dummy only.
func bind(fn expr, left, right expr) (expr, bool, error) {
	if right == nil {
		return nil, false, fmt.Errorf("bind: right argument is nil")
	}
	f, ok := fn.(Function)
	if ok == false {
		return nil, false, fmt.Errorf("bind: not a function: (%T) %v", fn, fn)
	}
	return &function{
		Function: f,
		left:     left,
		right:    right,
	}, true, nil
}

// lambda
//	{ guardList }
func (p *parser) lambda() (expr, bool, error) {
	enter("lambda", p.peek())
	defer leave("lambda")

	t := p.peek()
	if t.T != scan.LeftBrace {
		return nil, false, nil
	}
	p.next()

	body, ok, err := p.guardList()
	if err != nil {
		return nil, false, err
	} else {

		// A body must be terminated by }
		t = p.next()
		if t.T != scan.RightBrace {
			return nil, false, fmt.Errorf("missing } after lambda body, got %s", t)
		}

		if ok == false {
			return &lambda{}, true, nil // Empty lambda expression.
		} else {
			return &lambda{body: body}, true, nil
		}
	}
}

// guardList
//	TODO: allow (local?) assignments
//	guardExpr
//	guardExpr:expr (short ternary form, must be last is guardList)
//	guardExpr ⋄ guardExpr
func (p *parser) guardList() (guardList, bool, error) {
	enter("guardList", p.peek())
	defer leave("guardList")

	var l guardList
	for {
		if e, ok, err := p.guardExpr(); err != nil {
			return nil, false, err
		} else if ok == false {
			break
		} else {
			l = append(l, e)

			// Check if there is only one guardExpr without conditional.
			if e.cond == nil && len(l) > 1 && l[len(l)-2].cond == nil {
				return nil, false, fmt.Errorf("λ: only one expr without conditional is allowed: %s", l.String(nil))
			}

			t := p.peek()
			if t.T == scan.Diamond {
				p.next()
			} else {
				break
			}
		}
	}
	if l == nil {
		return nil, false, nil
	}

	// Ternary form, a guardExpr in the form cond:expr:expr may be the last guardExpr.
	// cond:expr1:expr2 is equivalent to cond:expr1⋄expr2
	t := p.peek()
	if t.T == scan.Colon {
		p.next()
		if e, ok, err := p.expr(); err != nil {
			return nil, false, err
		} else if ok == false {
			return nil, false, fmt.Errorf("λ: ternary form expects expression, got %s", p.peek())
		} else {
			l = append(l, &guardExpr{
				e: e,
			})
		}
	}

	return l, true, nil
}

// guardExpr
//	expr
//	expr:expr
func (p *parser) guardExpr() (*guardExpr, bool, error) {
	enter("guardExpr", p.peek())
	defer leave("guardExpr")

	cond, ok, err := p.expr()
	if err != nil {
		return nil, false, err
	} else if ok == false {
		return nil, false, nil
	}

	t := p.peek()
	if t.T == scan.Colon {
		p.next()
		e, ok, err := p.expr()
		if err != nil {
			return nil, false, err
		} else if ok == false {
			return &guardExpr{e: cond}, true, nil
		} else {
			return &guardExpr{cond: cond, e: e}, true, nil
		}
	} else {
		return &guardExpr{e: cond}, true, nil
	}
}

func (p *parser) operator() (string, bool) {
	enter("operator", p.peek())
	defer leave("operator")

	if t := p.peek(); t.T == scan.Symbol {
		if _, ok := p.a.operators[t.S]; ok {
			p.next()
			return t.S, true
		}
	}
	return "", false
}

// array
//	chars
//	scalar ...
func (p *parser) array() (expr, bool, error) {
	enter("array", p.peek())
	defer leave("array")

	if chars, ok, err := p.chars(); err != nil {
		return nil, false, err
	} else if ok {
		return chars, true, nil
	}

	var ar array
	for {
		if n, ok, err := p.scalar(); err != nil {
			return nil, false, err
		} else if ok == false {
			break
		} else {
			ar = append(ar, n)
		}
	}
	if ar == nil {
		return nil, false, nil
	}
	if len(ar) == 1 {
		return ar[0], true, nil
	}
	return ar, true, nil
}

// chars
//	'abc'
func (p *parser) chars() (expr, bool, error) {
	enter("chars", p.peek())
	defer leave("chars")

	if t := p.peek(); t.T == scan.Chars {
		p.next()
		runes := []rune(t.S)
		ar := make(array, len(runes))
		for i := range ar {
			ar[i] = String(string(runes[i]))
		}
		return ar, true, nil
	}
	return nil, false, nil
}

// scalar
//	bool
//	integer
//	string
//	( expr )
func (p *parser) scalar() (expr, bool, error) {
	// Debugging is ommitted here.

	t := p.peek()
	switch t.T {
	case scan.Identifier:
		// TODO: what to do if the variable is an array?
		// This evaluates as an nested array in other APLs.
		// Should we flatten the array, or return an error?
		if ok, isfunc := isVarname(t.S); ok == false {
			return nil, false, fmt.Errorf("not a variable name: %s", t.S)
		} else if isfunc {
			return nil, false, nil
		}
		p.next()
		return numVar{t.S}, true, nil
	case scan.String:
		p.next()
		return String(t.S), true, nil
	case scan.LeftParen:
		// TODO: Is this ok to declare everything within parens as a scalar?
		p.next()
		if e, ok, err := p.expr(); err != nil {
			return nil, false, err
		} else if ok == false {
			return nil, false, fmt.Errorf("expected ( expr ) got: %s", p.peek())
		} else {
			t = p.next()
			if t.T != scan.RightParen {
				return nil, false, fmt.Errorf("expected ), found %v", t)
			} else {
				return e, true, nil
			}
		}
	case scan.Number:
		if n, err := ParseNumber(t.S); err != nil {
			return nil, false, err
		} else {
			p.next()
			return n, true, nil
		}
	}
	return nil, false, nil
}
