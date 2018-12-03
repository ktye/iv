package apl

import (
	"fmt"
	"io"

	"github.com/ktye/iv/apl/scan"
)

type parser struct {
	a      *Apl
	tokens []scan.Token
	stack  []item
	pos    int
}

const (
	noun        class = 1 + iota // array expression    (A)
	verb                         // function expression (f)
	adverb                       // monadic operator    (/)
	conjunction                  // dyadic operator     (.)
)

// Item is an element of the parse stack.
// It contains a expr with an associated class.
type item struct {
	e     expr
	class class
}
type class int

func (c class) String() string {
	if c > 0 && c <= 4 {
		s := "_Af/."
		return string(s[int(c)])
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

		// A symbol is a primitive function, a dyadic or a monadic operator.
		case scan.Symbol:

			if _, ok := p.a.primitives[Primitive(t.S)]; ok {
				push(item{e: Primitive(t.S), class: verb}, false)
			} else if ops, ok := p.a.operators[t.S]; ok {
				i := item{e: &derived{op: t.S}, class: adverb}
				if ops[0].DyadicOp() == true {
					i.class = conjunction
				}
				push(i, false)
			} else {
				return item{}, fmt.Errorf(":%d: unknown symbol: %s", t.Pos, t.S)
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
				return item{}, fmt.Errorf(":%d: illegal variable name: %s", t.Pos, t.S)
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

		case scan.LeftParen, scan.LeftBrack, scan.LeftBrace:
			return item{}, fmt.Errorf(":%d: unexpected opening %s", t.Pos, t.S)

		case scan.RightParen:
			i, err := p.subStatement(t, scan.RightParen)
			if err != nil {
				return item{}, err
			}
			push(i, false)

		case scan.RightBrack:
			return item{}, fmt.Errorf("TODO: parse []")

		case scan.RightBrace:
			return item{}, fmt.Errorf("TODO: parse {}")

		case scan.Colon:
			return item{}, fmt.Errorf(":%d: unexpected : outside {}", t.Pos)

		case scan.Semicolon:
			return item{}, fmt.Errorf(":%d: unexpected ; outside []", t.Pos)

		default:
			return item{}, fmt.Errorf(":%d: unknown token %s", t.Pos, t.S)
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
func (p *parser) subStatement(left scan.Token, rightT scan.Type) (item, error) {
	// Pull until matching left paren. The right paren is not present anymore.
	var tokens []scan.Token
	l := 1
	for {
		t := p.pull()
		tokens = append(tokens, t)
		switch t.T {
		case scan.Endl:
			return item{}, fmt.Errorf(":%d: unmatched %s", left.Pos, left.S)
		case rightT:
			l++
		case left.T:
			l--
			if l == 0 {
				tokens = tokens[:len(tokens)-1]
			}
			break
		}
	}

	// Reverse tokens.
	for i := 0; i < len(tokens)/2; i++ {
		n := len(tokens) - i - 1
		tokens[i], tokens[n] = tokens[n], tokens[i]
	}

	// Create a new parser for the substatement and return it's result.
	q := parser{a: p.a, tokens: tokens}
	return q.parseStatement()
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
	for {
		if len(p.tokens) == 0 {
			break
		}
		t := p.tokens[len(p.tokens)-1]
		switch t.T {
		case scan.Number:
			if n, err := p.a.Tower.Parse(t.S); err != nil {
				return nil, fmt.Errorf(":%d: %s", t.Pos, err)
			} else {
				ar = append(ar, n)
			}

		case scan.String:
			ar = append(ar, String(t.S))

		case scan.Chars:
			runes := []rune(t.S)
			chars := make(array, len(runes))
			for i := len(chars) - 1; i >= 0; i-- {
				chars[i] = String(string(runes[i]))
			}
			if ar != nil && len(chars) > 1 {
				return nil, fmt.Errorf(":%d: only scalars can be added to an array", t.Pos)
			}
			ar = append(ar, chars...)

		case scan.Identifier:
			if ok, fok := isVarname(t.S); ok == false || fok == true {
				break
			}
			ar = append(ar, numVar{t.S})

		default:
			goto leave
		}
		p.pull() // Remove the token that has just been processed.
	}
leave:
	if ar == nil {
		return nil, fmt.Errorf(":%d: cannot collect array", right.Pos) // This should not happen.
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
	//p.printStack()
	in := p.shortStack()
	defer func() { fmt.Printf("reduce: %s → %s\n", in, p.shortStack()) }()

	p.resolveOps(last)
	p.resolveArray(last)
	p.resolveFunctions(last)

	if last && len(p.stack) > 1 {
		return fmt.Errorf("cannot reduce expression")
	}
	return nil
}

// ResolveOps tries to convert operators into derived functions.
//
// Operator reduction is done from the left side of the stack.
// If last==true, test if the second token is an operator:
//	:/+	mop(0) reduction
//	:.?+	dop(0) reduction
// Otherwise test if the third token is an operator.
//	::/+	mop(1) reduction
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
func (p *parser) resolveOps(last bool) {
	pos := 0
	if last == false {
		pos = 1
	}
	for {
		mok := p.mopReduce(pos)
		dok := p.dopReduce(pos)
		if mok == false && dok == false {
			return
		}
	}
}

// mopReduction reduces a monadic operator on position i+1 from the left,
// i must be 0 or 1.
// It returns if there was a reduction.
func (p *parser) mopReduce(i int) bool {
	//  :/+		mop(0) reduction
	//  ::/+	mop(1) reduction
	if len(p.stack) < 2+i || p.leftItem(i+1).class != verb {
		return false
	}
	d := p.leftItem(i + 1).e.(*derived)
	d.lo = p.leftItem(i).e
	p.setLeft(i, item{e: d, class: verb})
	p.removeLeft(i + 1)
	return true
}

// dopReduction reduces a dyadic operator on position i+1 from the left,
// i must be 0 or 1.
// It returns if there was a reduction.
func (p *parser) dopReduce(i int) bool {
	//  :.?+	dop(0) reduction
	//  ::.?+	dop(1) reduction
	if len(p.stack) < 3+i || p.leftItem(i+1).class != conjunction {
		return false
	}
	d := p.leftItem(i + 1).e.(*derived)
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
// Otherwise:
//	+ffA		fA
//	+/fA		fA
//	+!AAfA		AfA
// +: zero or more items of any type
// !: item of type A, f or / (everything but a DOP).
//
// The reduction is repeated until no pattern is found.
func (p *parser) resolveArray(last bool) {
	for {
		fok := p.reducefA(last)
		aok := p.reduceAfA(last)
		if fok == false && aok == false {
			return
		}
	}
}

/*
	for {
		if last && len(p.stack) == 1 {
			return nil
		}

		r := p.rightItem(0)
		if r.class != noun {
			return nil
		}
		if last && len(p.stack) == 2 {
			l := p.leftItem(0)
			switch l.class {
			case noun: // AA
				return fmt.Errorf("two consecutive arrays are not allowed")
			case verb:
				// fA
				f := &function{
					Function: l.e.(Function),
					right:    r.e,
				}
				p.stack = []item{item{e: f, class: noun}}
				return nil
			default:
				return fmt.Errorf("sole operator on left side") // This should be triggered already.
			}
		} else if last == true && len(p.stack) == 3 {
			// AfA
			i0 := p.leftItem(0)
			i1 := p.leftItem(1)
			if i0.class == noun && i1.class == verb {
				f := &function{
					Function: i1.e.(Function),
					left:     i0.e,
					right:    p.rightItem(0).e,
				}
				p.stack = []item{item{e: f, class: noun}}
				return nil
			}
			return fmt.Errorf("unexpected tail: %T %T %T", p.stack[2].e, p.stack[1].e, p.stack[0].e)
		} else if len(p.stack) > 2 && p.rightItem(2).class != conjunction {
			// error for: +fAA | +AAA | +/AA
			if r1 := p.rightItem(1); r1.class == noun {
				return fmt.Errorf("two consecutive arrays are not allowed")
			} else if r1.class == verb {
				// apply monad for: +ffA | +/fA not for: +AfA
				if p.rightItem(2).class != noun {
					f := &function{
						Function: r1.e.(Function),
						right:    p.rightItem(0).e,
					}
					p.setRight(1, item{e: f, class: noun})
					p.removeRight(0) // Remove last A
					continue
				}
			} else { // f/A
				return nil
			}
		} else if len(p.stack) > 3 && p.rightItem(3).class != conjunction {
			// dyadic application for: +fAfA, +AAfA, +/AfA
			r0 := p.rightItem(0)
			r1 := p.rightItem(1)
			r2 := p.rightItem(2)
			if r0.class == noun && r1.class == verb && r2.class == noun {
				f := &function{
					Function: r1.e.(Function),
					left:     r2.e,
					right:    r0.e,
				}
				p.setRight(2, item{e: f, class: noun})
				p.removeRight(0) // Remove fA
				p.removeRight(0)
				continue
			}
		}
		return nil
	}
	return nil
}
*/

func (p *parser) reducefA(last bool) bool {
	// fA
	// +ffA
	// +/fA
	if len(p.stack) < 2 {
		return false
	}
	r := p.rightItem(0)
	f := p.rightItem(1)
	x := item{class: conjunction}
	if len(p.stack) > 2 {
		x = p.rightItem(2).class
	}
	if last && len(p.stack) == 2 {
		x.class = verb
	}
	if r.class == noun && f.class == verb && (x.class == verb || x.class == adverb) {
		fn := &function{
			Function: f.e.(Function),
			right:    r.e,
		}
		p.setRight(1, item{e: f, class: noun})
		p.removeRight(0)
		return true
	}
	return false
}

func (p *parser) reduceAfA(last bool) bool {
	// AfA
	// +!AAfA
	if len(p.stack) < 3 {
		return false
	}
	l := p.rightItem(2)
	f := p.rightItem(1)
	r := p.rightItem(0)

	a := noun
	x := conjunction
	if len(p.stack) > 4 {
		x = p.righItem(4).class
		a = p.righItem(3).class
	} else if last == false {
		return false
	}
	if l.class == noun && f.class == verb && r.class == noun && a == noun && x != conjunction {
		fn := &function{
			Function: f.e.(Function),
			left:     l.e,
			right:    r.e,
		}
		p.setRight(2, item{e: f, class: noun})
		p.removeRight(0)
		p.removeRight(0)
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
func (p *parser) resolveFunction(last bool) error {
	if len(p.stack) == 0 {
		return fmt.Errorf("parse stack is empty")
	}
	for {
		if len(p.stack) == 1 {
			return nil
		}

		r0 := p.rightItem(0)
		if r0.class != verb {
			return nil
		}

		if last || (len(p.stack) > 2 && p.rightItem(2).class != conjunction) {
			r1 := p.rightItem(1)
			if r1.class == noun {
				// No preceeding array allowed: +AAF, +fAf, +/Af
				return fmt.Errorf("function expression is preceeded by array")
			}
			if r1.class == verb {
				// Combine functions: +Aff, +fff, +/ff.
				if t, ok := r0.e.(train); ok {
					t = append(train{r1.e}, t...)
					p.setRight(1, item{e: t, class: verb})
				} else {
					t = train{r1.e, r0.e}
					p.setRight(1, item{e: t, class: verb})
				}
				p.removeRight(0)
				continue
			}
			return nil
		}
		return nil
	}
	return nil
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
