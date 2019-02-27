package apl

import (
	"fmt"
	"io"
	"strings"

	"github.com/ktye/iv/apl/scan"
)

type parser struct {
	a      *Apl
	tokens []scan.Token
	stack  []item
	pos    int
}

const (
	noun        class = 1 << iota // array expression    (A)
	verb                          // function expression (f)
	adverb                        // monadic operator    (/)
	conjunction                   // dyadic operator     (.)
)

// Item is an element of the parse stack.
// It contains a expr with an associated class.
type item struct {
	e     expr
	class class
}
type class int

func (c class) String() string {
	s := "Af/."
	for i := range s {
		if c&class(1<<uint(i)) != 0 {
			return string(s[i])
		}
	}
	return "?"
}

// Parse parses the tokens to a program, which is a slice of expressions.
func (p *parser) parse(tokens []scan.Token) (Program, error) {

	var prog Program
	var itm item
	var err error
	for {
		tokens, err = p.nextStatement(tokens)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		itm, err = p.parseStatement()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		prog = append(prog, itm.e)
	}
	if len(prog) == 0 {
		return nil, fmt.Errorf("empty program")
	}
	return prog, nil
}

// nextStatement extracts the next statements from tokens and sets it to the parser.
// It returns the remaining tokens.
// Statements are separated by diamond tokens.
// Diamonds within lambda expressions are skipped.
func (p *parser) nextStatement(tokens []scan.Token) ([]scan.Token, error) {
	if len(tokens) == 0 {
		return nil, io.EOF
	}
	b := 0
	for i, t := range tokens {
		if t.T == scan.LeftBrace {
			b++
		} else if t.T == scan.RightBrace {
			b--
			if b < 0 {
				return nil, fmt.Errorf("too many }")
			}
		}
		if b == 0 && t.T == scan.Diamond {
			p.tokens = tokens[:i]
			return tokens[i+1:], nil
		}
	}
	p.tokens = tokens
	return nil, nil
}

// Statement parses a statement, that is a complete sentence from right to left
// and returns a single expression in tree form.
//
// It pulls the next token from the right side and converts it to a stack item.
// The item is pushed to the stack and a stack reduction step is executed.
// At the end a single item must be left on the stack, an is returned.
//
// (), [], {} with their child tokens are extracted and parsed separately.
// Their result is pushed on the stack.
//
// Only items with class A, f, / or . are passed to the stack.
//
// This is similar to the parser described in Iverson, A dictionary of APL 1987, sec I.
// With one difference:
// Instead of evaluating during reduction, the reduction assembles an expression
// and returns it's root node.
func (p *parser) parseStatement() (item, error) {
	// The stack stores items in reverse order. Pushing appends to the end.
	p.stack = make([]item, 0, 20)

	var parseError error
	push := func(i item, last bool) {
		p.stack = append(p.stack, i)
		parseError = p.reduce(last)
	}

	for {
		if parseError != nil {
			return item{}, parseError
		}

		t := p.pull()
		switch t.T {
		case scan.Endl:
			if len(p.stack) == 0 {
				return item{}, nil // TODO: Should this return an error?
			}
			if err := p.reduce(true); err != nil {
				return item{}, err
			}
			p.linkFuncAssign()
			return p.stack[0], nil

		// A symbol may be a primitive function, a dyadic or a monadic operator.
		case scan.Symbol:

			if _, ok := p.a.primitives[Primitive(t.S)]; ok {
				push(item{e: Primitive(t.S), class: verb}, false)
			} else if ops, ok := p.a.operators[t.S]; ok {
				i := item{e: &derived{op: t.S}, class: adverb}
				if ops[0].DyadicOp() == true {
					i.class = conjunction
				}
				if t.S == "∘" {
					i = p.specialJot(i)
				}
				push(i, false)
			} else {
				return item{}, fmt.Errorf("unknown symbol: %s", t.S)
			}

		case scan.Number, scan.String, scan.Chars:
			e, err := p.collectArray(t)
			if err != nil {
				return item{}, err
			}
			push(item{e: e, class: noun}, false)

		case scan.Identifier:
			i := item{class: verb}
			if ok, fok := isVarname(t.S); ok == false {
				return item{}, fmt.Errorf("illegal variable name: %s", t.S)
			} else if fok == false {
				e, err := p.collectArray(t)
				if err != nil {
					return item{}, err
				}
				i.e = e
				i.class = noun
			} else {
				i.e = fnVar(t.S)
			}
			push(i, false)

		case scan.Self:
			push(item{e: self{}, class: verb}, false)

		case scan.LeftParen, scan.LeftBrack, scan.LeftBrace:
			return item{}, fmt.Errorf("unexpected opening %s", t.S)

		case scan.RightParen, scan.RightBrack, scan.RightBrace:
			i, err := p.subStatement(t.T)
			if err != nil {
				return item{}, fmt.Errorf("%s", err)
			}
			push(i, false)

		case scan.Colon:
			return item{}, fmt.Errorf("unexpected : outside {}")

		case scan.Semicolon:
			return item{}, fmt.Errorf("unexpected ; outside []")

		default:
			return item{}, fmt.Errorf("unknown token %s", t.S)
		}
	}
	return item{}, fmt.Errorf("illegal parser state") // Should not be reached.
}

// pull returns the last from the parsers tokens and removes it from the buffer.
// If there is no token, the empty token with type scan.Endl is returned.
func (p *parser) pull() scan.Token {
	if len(p.tokens) == 0 {
		return scan.Token{}
	}
	t := p.tokens[len(p.tokens)-1]
	p.tokens = p.tokens[:len(p.tokens)-1]
	return t
}

// subStatement parses a parenthesized substatement.
// Parens may be (), [] or {}.
func (p *parser) subStatement(right scan.Type) (item, error) {

	left := map[scan.Type]scan.Type{
		scan.RightBrace: scan.LeftBrace,
		scan.RightParen: scan.LeftParen,
		scan.RightBrack: scan.LeftBrack,
	}[right]

	// Pull until matching left paren. The right paren is not present anymore.
	var tokens []scan.Token
	l := 1
	for {
		t := p.pull()
		tokens = append(tokens, t)
		switch t.T {
		case scan.Endl:
			return item{}, fmt.Errorf("unmatched %s", left.String())
		case right:
			l++
		case left:
			l--
			if l == 0 {
				tokens = tokens[:len(tokens)-1]
				goto rev
			}
		}
	}
rev:
	if len(tokens) == 0 {
		if left == scan.LeftBrace {
			return item{e: &lambda{}, class: verb}, nil
		} else if left == scan.LeftBrack {
			// Empty brackets are allowed in: A[]←9
			return item{e: idxSpec(nil), class: adverb}, nil
		}
		return item{}, fmt.Errorf("empty: %s%s", left.String(), right.String())
	}

	// Reverse tokens.
	for i := 0; i < len(tokens)/2; i++ {
		k := len(tokens) - i - 1
		tokens[i], tokens[k] = tokens[k], tokens[i]
	}

	// Create a new parser for the substatement and return it's result.
	q := &parser{a: p.a, tokens: tokens}

	switch left {
	case scan.LeftParen:
		if tokens[len(tokens)-1].T == scan.Semicolon {
			return q.parseList()
		}
		return q.parseStatement()
	case scan.LeftBrack:
		return q.parseBrackets()
	case scan.LeftBrace:
		return q.parseLambda()
	default:
		return item{}, fmt.Errorf("unknown paranthesis type %T", left) // This should not happen.
	}
}

// parseList parses a list expression, that always evaluates as a noun.
func (p *parser) parseList() (item, error) {
	l := p.splitTokens(scan.Semicolon, []scan.Type{scan.LeftParen}, []scan.Type{scan.RightParen})
	if len(l) > 0 {
		// Remove empty tail caused by trailing semicolon
		l = l[:len(l)-1]
	}
	lst := make(list, len(l))
	for i := range l {
		q := &parser{a: p.a, tokens: l[i]}
		it, err := q.parseStatement()
		if err != nil {
			return item{}, err
		}
		lst[i] = it.e
	}
	return item{e: lst, class: noun}, nil
}

// ParseBrackets parses the expression within brackets [...].
// This may be an index or axis specification.
func (p *parser) parseBrackets() (item, error) {
	//l := p.splitTokens(scan.Semicolon, scan.LeftBrack, scan.RightBrack)
	// Allow lists in index specifications. If nested brackets are needed, splitTokens mut be modified
	// to use multiple nesting tokens.
	l := p.splitTokens(scan.Semicolon, []scan.Type{scan.LeftParen, scan.LeftBrack}, []scan.Type{scan.RightParen, scan.RightBrack})
	spec := make(idxSpec, len(l))
	for i := range l {
		q := &parser{a: p.a, tokens: l[i]}
		it, err := q.parseStatement()
		if err != nil {
			return item{}, err
		}
		spec[i] = it.e
		if spec[i] == nil {
			spec[i] = EmptyArray{}
		}
	}
	// A bracket parses as an adverb, that will be corrected in p.reduce.
	return item{e: spec, class: adverb}, nil
}

// ParseLambda parses the lambda expression given in the parser's token list.
// A lambda expression contains a list of guarded expressions:
//	{ guardList }
// The outer braces are not present anymore in the parsers's tokens.
// Lambdas are calles dfns in dyalog: DyaProg p. 131
func (p *parser) parseLambda() (item, error) {
	// Entries of the guardList are separated by diamonds.
	l := p.splitTokens(scan.Diamond, []scan.Type{scan.LeftBrace}, []scan.Type{scan.RightBrace})
	body := make(guardList, len(l))
	for i := range l {
		q := &parser{a: p.a, tokens: l[i]}
		ge, ternary, err := q.guardExpr()
		if err != nil {
			return item{}, err
		}
		body[i] = ge
		if ternary != nil && i != len(l)-1 {
			return item{}, fmt.Errorf("lambda: ternary is only allowed as the last item")
		} else if ternary != nil {
			body = append(body, &guardExpr{e: ternary})
		}
	}
	// Remove an empty last entry, that happens if a diamond is inserted before }.
	if len(body) > 0 {
		last := body[len(body)-1]
		if last.cond == nil && last.e == nil {
			body = body[:len(body)-1]
		}
	}
	return item{e: &lambda{body}, class: verb}, nil
}

// GuardExpr parses a guarded expression, which is part of a lambda expression.
//	cond:expr
//	cond:expr:expr2 (short ternary form, only for the last in the list).
func (p *parser) guardExpr() (*guardExpr, expr, error) {
	l := p.splitTokens(scan.Colon, []scan.Type{scan.LeftBrace}, []scan.Type{scan.RightBrace})
	if len(l) > 3 {
		return nil, nil, fmt.Errorf("lambda has too many colons")
	}
	ge := &guardExpr{}
	for i := range l {
		q := &parser{a: p.a, tokens: l[i]}
		item, err := q.parseStatement()
		if err != nil {
			return nil, nil, err
		}

		if i == 0 && len(l) > 1 {
			ge.cond = item.e
		} else if i == 0 || i == 1 {
			ge.e = item.e
		} else if i == 2 {
			return ge, item.e, nil
		}
	}
	return ge, nil, nil
}

// collectArray pulls tokens from the parser that form an array starting with the given right end token.
//
// TODO: this is not correct: Vector binding is not stronger than right operand binding
// APL2 p 36: LO DOP A B ←→ (LO DOP A) B not LO DOP (A B)
// Check if the token left to the first number is a DOP.
func (p *parser) collectArray(right scan.Token) (expr, error) {
	// Push back the right token.
	p.tokens = append(p.tokens, right)

	// The array is collected in reverse order.
	var ar array
loop:
	for {
		if len(p.tokens) == 0 {
			break
		}
		t := p.tokens[len(p.tokens)-1]
		switch t.T {
		case scan.Number:
			if n, err := p.a.Tower.Parse(t.S); err != nil {
				return nil, err
			} else {
				ar = append(ar, n)
			}

		case scan.String:
			ar = append(ar, String(t.S))

		case scan.Chars:
			runes := []rune(t.S)
			if len(runes) > 0 {
				chars := make(array, len(runes))
				for i := range chars {
					k := len(chars) - 1 - i
					chars[i] = String(string(runes[k]))
				}
				if ar != nil && len(chars) > 1 {
					return nil, fmt.Errorf("only scalars can be added to an array")
				}
				ar = append(ar, chars...)
			}

		case scan.Identifier:
			if ok, fok := isVarname(t.S); ok == false || fok == true {
				break loop
			}
			ar = append(ar, numVar{t.S})

		default:
			break loop
		}
		p.pull() // Remove the token that has just been processed.
	}

	if ar == nil {
		return EmptyArray{}, nil
	}

	// Reverse the array to the normal left to right order.
	for i := 0; i < len(ar)/2; i++ {
		n := len(ar) - i - 1
		ar[i], ar[n] = ar[n], ar[i]
	}

	// If there is only 1 item, return it unboxed.
	if len(ar) == 1 {
		return ar[0], nil
	}
	return ar, nil
}

// Reduce tries to reduce the partial right tail of the stack.
func (p *parser) reduce(last bool) error {
	// in := p.shortStack()
	// defer func() { fmt.Printf("reduce: %s → %s\n", in, p.shortStack()) }()

	if _, ok := p.leftItem(0).e.(idxSpec); ok {
		return nil
	}
	if err := p.resolveBrackets(); err != nil {
		return err
	}
	p.resolveOperators(last)
	p.resolveArrays(last)
	p.resolveFunctions(last)

	if last && len(p.stack) > 1 {
		return fmt.Errorf("cannot reduce expression")
	}
	return nil
}

// ResolveBrackets changes the order of bracket expressions.
// Brackets are parsed as idxSpec following an Identifier or as
// an axis specification following a primitive or an operator.
// IdxSpec is converted to a dyadic indexing function.
// An axis specification is converted to a dyadic operator.
func (p *parser) resolveBrackets() error {
	if len(p.stack) < 2 {
		return nil
	}
	if spec, ok := p.leftItem(1).e.(idxSpec); ok {
		l := p.leftItem(0)
		if id, ok := l.e.(numVar); ok && l.class == noun {
			fn := &function{
				Function: Primitive("⌷"),
				left:     spec,
				right:    id,
			}
			p.setLeft(1, item{e: fn, class: noun})
			p.removeLeft(0)
			return nil
		} else if _, ok := l.e.(Primitive); ok {
			if len(spec) != 1 {
				return fmt.Errorf("axis must hold a single expression, not %d", len(spec))
			}
			d := derived{
				lo: l.e,
				ro: spec[0],
				op: "⍂",
			}
			p.setLeft(1, item{e: &d, class: verb})
			p.removeLeft(0)
		} else if _, ok := l.e.(*derived); ok {
			// The axis specification following an operator is rewritten as a dyadic operator.
			// The operator is called "⍂" and as the left operand the axis spec is inserted.
			if len(spec) != 1 {
				return fmt.Errorf("axis must hold a single expression, not %d", len(spec))
			}
			d := derived{
				op: "⍂",
			}
			p.setLeft(1, item{e: &d, class: conjunction})
			p.insertLeft(1, item{e: spec[0], class: noun})
		} else if l.class == noun {
			// Axis specification follows an array or a noun expression.
			fn := &function{
				Function: Primitive("⌷"),
				left:     spec,
				right:    l.e,
			}
			p.setLeft(1, item{e: fn, class: noun})
			p.removeLeft(0)
			return nil
		} else {
			return fmt.Errorf("bracket expr following an %T %v\n", l.e, l.class)
		}
	}
	return nil
}

// ResolveOperators tries to convert operators into derived functions.
//
// Operator reduction is done from the left side of the stack.
// If last==true, test if the second token is an operator:
//	:/+	mop(0+) reduction
//	:.?+	dop(0) reduction
// In any case, test if the third token is an operator.
//	::/+	mop(1+) reduction
//	::.?+	dop(1) reduction
// ?  an item of any class
// +  zero or more items of any class
// :  an item that is f or A (not an operator)
//
// The result is always of class f: a derived function.
// Repeat until no reduction can be done.
//
// Operators have long scope on the left and short scope on the right:
// 	+.+.+.+.*  ←→ ((((+.+).+).+).*)  ←≠→  (+.(+.(+.(+.*))))
// See DyaProg p 21.
func (p *parser) resolveOperators(last bool) {
	for {
		ok1, ok2, ok3, ok4 := false, false, false, false
		if last {
			ok1 = p.mopReduce(0)
		}
		ok2 = p.mopReduce(1)
		if last {
			ok3 = p.dopReduce(0)
		}
		ok4 = p.dopReduce(1)
		if ok1 || ok2 || ok3 || ok4 {
			continue
		}
		return
	}
}

// mopReduction reduces a monadic operator on position i+1 from the left, i must be 0 or 1.
// If there was a reduction it continues with i+1.
// It returns true, if there was at least one reduction.
func (p *parser) mopReduce(i int) bool {
	//  :/+		mop(0) reduction
	//  ::/+	mop(1) reduction
	reduced := false
	for {
		op := adverb | conjunction
		if len(p.stack) < 2+i || p.leftItem(i+1).class != adverb {
			return reduced
		}
		c0, c1 := p.leftItem(0).class, p.leftItem(1).class
		if (c0&op != 0) || ((i == 1) && (c1&op != 0)) {
			return reduced
		}
		d := p.leftItem(i + 1).e.(*derived)
		if d.lo != nil {
			fmt.Println("dopReduce: overwriting LO")
		}
		d.lo = p.leftItem(i).e
		p.setLeft(i, item{e: d, class: verb})
		p.removeLeft(i + 1)
		reduced = true
		i++
	}
}

// dopReduction reduces a dyadic operator on position i+1 from the left,
// i must be 0 or 1.
// It returns true, if there was a reduction.
func (p *parser) dopReduce(i int) bool {
	//  :.?+	dop(0) reduction
	//  ::.?+	dop(1) reduction
	op := adverb | conjunction
	if len(p.stack) < 3+i || p.leftItem(i+1).class != conjunction {
		return false
	}
	c0, c1 := p.leftItem(0).class, p.leftItem(1).class
	if (c0&op != 0) || ((i == 1) && (c1&op != 0)) {
		return false
	}
	d := p.leftItem(i + 1).e.(*derived)
	if d.lo != nil {
		fmt.Println("dopReduce: overwriting LO")
	}
	if d.ro != nil {
		fmt.Println("dopReduce overwriting RO")
	}
	d.lo = p.leftItem(i).e
	d.ro = p.leftItem(i + 2).e
	p.setLeft(i, item{e: d, class: verb})
	p.removeLeft(i + 1)
	p.removeLeft(i + 1)

	return true
}

// ResolveArray resolves arrays from the right end.
//
// Arrays are reduced from the right end of the stack.
// The pattern at the right side is always reduced to class A.
// If last==true:	reduction on the tail (result is always A)
//	fA		fA
//	AfA    		AfA
// In any case:
//	+ffA		fA
//	+/fA		fA
//	+fAfA		AfA
//	+!AAfA		AfA
// +: zero or more items of any type
// !: item of type A, f or / (everything but a DOP).
//
// The reduction is repeated until no pattern is found.
func (p *parser) resolveArrays(last bool) {
	for {
		fok := p.reducefA(last)
		aok := p.reduceAfA(last)
		if fok == false && aok == false {
			return
		}
	}
}

func (p *parser) reducefA(last bool) bool {
	// fA
	// +ffA
	// +/fA
	if len(p.stack) < 2 {
		return false
	}
	r := p.rightItem(0)
	f := p.rightItem(1)
	x := conjunction
	if len(p.stack) > 2 {
		x = p.rightItem(2).class
	}
	if last && len(p.stack) == 2 {
		x = verb
	}
	op := verb | adverb
	if r.class == noun && f.class == verb && (x&op != 0) {
		fn := &function{
			Function: f.e.(Function),
			right:    r.e,
		}
		p.setRight(1, item{e: fn, class: noun})
		p.removeRight(0)
		return true
	}
	return false
}

func (p *parser) reduceAfA(last bool) bool {
	// AfA
	// +fAfA
	// +!AAfA
	if len(p.stack) < 3 {
		return false
	}
	l := p.rightItem(2)
	f := p.rightItem(1)
	r := p.rightItem(0)

	x := conjunction
	y := noun
	if len(p.stack) > 4 {
		x = p.rightItem(4).class
		y = p.rightItem(3).class
	} else if last == false {
		return false
	}
	reduce := func() {
		fn := &function{
			Function: f.e.(Function),
			left:     l.e,
			right:    r.e,
		}
		p.setRight(2, item{e: fn, class: noun})
		p.removeRight(0)
		p.removeRight(0)
	}

	if !(l.class == noun && f.class == verb && r.class == noun) {
		return false
	}
	if last == true && len(p.stack) > 2 {
		reduce()
		return true
	}
	if y == verb || (y == noun && x != conjunction) {
		reduce()
		return true
	}
	return false
}

// ResolveFunction resolves functions from the right end.
//
// Functions are reduced from the right end of the stack.
// Two functions on the right are combined to function trains,
// if they match this patterns:
// 	+!ff
// The reduction is called train reduction.
// It is repeated until the pattern is not found again.
// On the last call, also
//	Af
// is reduced to a train, if f is already a train with an even length.
// This builds an Agh fork. It needs parenthesis.
func (p *parser) resolveFunctions(last bool) {
	for {
		if p.reduceff(last) == false {
			break
		}
	}
	// Reduce Agh fork.
	if last && len(p.stack) == 2 {
		r0 := p.rightItem(0)
		r1 := p.rightItem(1)
		if t, ok := r0.e.(train); ok && r1.class == noun && len(t)%2 == 0 {
			t = append(train{r1.e}, t...)
			p.setRight(1, item{e: t, class: verb})
			p.removeRight(0)
		}
	}
}

func (p *parser) reduceff(last bool) bool {
	if len(p.stack) < 2 {
		return false
	}
	// +!ff
	r0 := p.rightItem(0)
	r1 := p.rightItem(1)
	c := conjunction
	if len(p.stack) > 2 {
		c = p.rightItem(2).class
	}
	if (r0.class == verb && r1.class == verb) && ((len(p.stack) == 2 && last) || c != conjunction) {
		if t, ok := r0.e.(train); ok {
			t = append(train{r1.e}, t...)
			p.setRight(1, item{e: t, class: verb})
		} else {
			t = train{r1.e, r0.e}
			p.setRight(1, item{e: t, class: verb})
		}
		p.removeRight(0)
		return true
	}
	return false
}

// SpecialJot converts the item from a dyadic operator to a primitive function,
// if it follows a dot.
// This special case is applied for ∘ only, which is registered as a DOP,
// but used in: ∘.×
func (p *parser) specialJot(i item) item {
	if len(p.stack) < 1 {
		return i
	}
	dot := p.stack[len(p.stack)-1].e
	if d, ok := dot.(*derived); ok && d.op == "." {
		i = item{e: Primitive("∘"), class: verb}
	}

	return i
}

// RemoveLeft removes item i from the left side of the stack.
func (p *parser) removeLeft(l int) {
	i := len(p.stack) - 1 - l
	copy(p.stack[i:], p.stack[i+1:])
	p.stack = p.stack[:len(p.stack)-1]
}

// SetLeft overwrites the item at position i from left with the new item.
func (p *parser) setLeft(l int, i item) {
	p.stack[len(p.stack)-1-l] = i
}

// InsertLeft inserts an item at position i from the left.
func (p *parser) insertLeft(l int, e item) {
	i := len(p.stack) - l - 1
	p.stack = append(p.stack, item{})
	copy(p.stack[i+1:], p.stack[i:])
	p.stack[i] = e
}

// LeftItem returns the item number i from the left end, starting at 0.
func (p *parser) leftItem(i int) item {
	// tokens: a b c
	// stack: c b a
	return p.stack[len(p.stack)-1-i]
}

// RightItem returns the item number i from the right end, starting at 0.
func (p *parser) rightItem(i int) item {
	return p.stack[i]
}

// RemoveRight removes item i from the right side of the stack.
func (p *parser) removeRight(i int) {
	copy(p.stack[i:], p.stack[i+1:])
	p.stack = p.stack[:len(p.stack)-1]
}

// SetRight overwrites the item at position r from right with the new item.
func (p *parser) setRight(r int, i item) {
	p.stack[r] = i
}

// Split Tokens splits the parses's tokens at the given separator.
// It ignores sub-expressions nested by inc and dec tokens.
func (p *parser) splitTokens(at scan.Type, inc, dec []scan.Type) [][]scan.Token {
	any := func(t scan.Type, l []scan.Type) bool {
		for i := range l {
			if t == l[i] {
				return true
			}
		}
		return false
	}
	var l [][]scan.Token
	lv := 0
	n := 0
	for i, t := range p.tokens {
		if any(t.T, inc) {
			lv++
		} else if any(t.T, dec) {
			lv--
		}
		if lv == 0 && t.T == at {
			l = append(l, p.tokens[n:i])
			n = i + 1
		}
	}
	l = append(l, p.tokens[n:len(p.tokens)])
	return l
}

// LinkFuncAssign corrects a function assignment, that was parsed as a train.
//	f ← +
func (p *parser) linkFuncAssign() {
	if len(p.stack) != 1 {
		return
	}
	// The expression was parsed as a train of two functions:
	// 1. a derived assignment
	// 2: the function expression
	// It is converted to a monadic function, of the derived assignment
	// with the second verb as the right argument.
	if t, ok := p.stack[0].e.(train); ok && len(t) == 2 {
		if d, ok := t[0].(*derived); ok && d.op == "←" {
			p.stack = []item{
				item{
					e: &function{
						Function: d,
						right:    t[1],
					},
					class: verb,
				},
			}
		}
	}
}

func (p *parser) printStack() {
	for k, i := range p.stack {
		fmt.Printf("#%d: %s %s\n", k, i.class.String(), i.e.String(p.a))
	}
}

func (p *parser) shortStack() string {
	v := make([]byte, len(p.stack))
	k := 0
	for i := len(p.stack) - 1; i >= 0; i-- {
		s := p.stack[i].class.String()
		v[k] = s[0]
		k++
	}
	return string(v)
}

func (p *parser) tokenString() string {
	s := make([]string, len(p.tokens))
	for i, t := range p.tokens {
		s[i] = t.S
	}
	return strings.Join(s, " ")
}
