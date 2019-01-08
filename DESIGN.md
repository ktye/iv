Rob Pike told us: *Try implementing an APL-like language yourself. 
The basics take just a day or two. You'll learn a lot.*

There is some truth in it with the right definition of *basic*. However he left out the part about addiction...

# PARSING APL
*Ivy* parses from left to right. 
When I finally understood, that a function and an operator are two different beasts (that took a while) 
and ivy is cheating here, the parser became more complex.
The bad thing: *it still worked*.

I could not find any description, on how to write an APL parser except some lines of APL I could not decipher.
But these suggested that a simpler solution should be possible.

The current parser parses right to left, but still separates parsing and evaluating.

This can only be done with restrictions on variable names:
- verbs are lowercase
- nouns are are uppercase
- operators are registered unicode runes
  - this could one day be extended within the current frame by special names for defined operators (greek letters? Prefixes?)

## Tokenization
The scanner is in `apl/scan/scan.go`. 
```go
func (s *Scanner) Scan(line string) ([]Token, error)
```
It does not need to know much, just what makes an APL symbol.
It is told so by:
```go
func (s *Scanner) SetSymbols(symbols map[rune]string)
```
The actual list of symbols depends on which packages are built-in the executable.

The scanner is also agnostic to what is a number. 
It knows how a number starts and how it ends.
Parsing it, is an exercise left for the parser.
The latter delegates this task to the current numeric tower.

## Parser
The parser is `apl/parse.go`.
```go
func (p *parser) parse(tokens []scan.Token) (Program, error)
```
called by the interpreter in `apl/apl.go: func (a *Apl) Parse(line string) (Program, error)`

It reads tokens from the right side and pushes them on the stack.
After each token it tries to *reduce* the expression by assembling a nested tree.

This is where it differes from a classic APL parser (to my understanding).
An APL would do the reduction by partial evaluation.

At the end of parsing a line, only a single item may be left. Hence the error message
```
cannot reduce expression
```
if it cannot.

The *Reduction* step tries to reduce operator expressions from the left 
and noun expressions from the right side of the partial list, plus some special handling for bracket indexing and assignments.
I'm not sure if a real APL parser does it in a similar way.
I hope someone will enlighten me some day.

The result of parsing is a *Program*. This is a bad name.
It refers to a single line of APL input. Maybe *Phrase* would be better.

I played around with the idea of parsing directly into an `apl.List` value, similarly to what I think K does.
This would allow to do Lisp style manipulation at runtime.
It's deferred for future work.

## Evaluation
Evaluation executes a *Program* by converting it into one or more *Values*.

The interpreter `cmd/apl` does this by calling:
```go
func (a *Apl) ParseAndEval(line string) error {
	if p, err := a.Parse(line); err != nil {
		return err
	} else {
		return a.Eval(p)
	}
}
```
which also prints the result.
Other interfaces exist in `apl/eval.go`

# Types and Values
An APL Value is implemented as an interface.
Anything that can be printed into a string can act as an `apl.Value` (`apl/value.go`):
```go
type Value interface {
	String(*Apl) string
}
```
It's not yet decided if this is enough. Maybe a Copy method could be added later.

APL is an array language:
```go
type Array interface {
	String(*Apl) string
	At(int) Value
	Shape() []int
	Size() int
}
```
Also an array is implemented as an interface. This makes it possible add all kinds of special implementations later.
The most prominent implementation of it is the `apl.MixedArray`
```go
type MixedArray struct {
	Values []Value
	Dims   []int
}
```
Technically it's Values can be anyting, but the default implementation of primary functions and operators implies scalars / atoms.

Functions are also implemented as an interface.
```go
type Function interface {
	Call(a *Apl, L Value, R Value) (Value, error)
}
```
If L is nil, it's a monadic call.

Primitive functions and derived functions are examples of actual types that implement it.
Another one could be a method of an external go type.


# Numbers
Numbers are added to the type system by implementing a `apl.Tower`, see `apl/tower.go`.
A tower is a linear escalation of numeric types.
A dyadic call to an elementary function (`+×⍟⌊...`) implies that both values are on the same level.

The default tower is implemented in `apl/numbers/` and contains Integer, Float, Complex and Time.

Time is a little bit special here. It could have been implemented outside the number system.
But this allows direct parsing amoung benefits.

Two other tower implementations are provided in `apl/big`:
A big tower with Int and Rat and a precise tower with Float and Complex containing any number of bits.

The big tower is nice for playing with numbers.
It's Rat type can be used in combination with ⌹ to solve linear equations loss-less.
I have seen this from Bob Smith and wanted it too.

The precise tower could come in handy when testing numerical algorithms.

The towers are only available if they are registered at compile time, so it's possible to use only a single one.
If multiple tower are present, they can be switched but not mixed at runtime.

# Overloading primitive functions and operators
# Lambda functions
# Go interface
# Streams and concurrency
