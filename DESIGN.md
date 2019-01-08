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
It does not need to know much, just which Unicode rune APL considers as a symbol.
It is told so by:
```go
func (s *Scanner) SetSymbols(symbols map[rune]string)
```
The actual list of symbols depends on which packages are built-in the executable.

The scanner is also agnostic to what is a number. 
It knows how a number starts and how it ends.
Parsing it, is left as an exercise for the parser.
The latter delegates this task to the current numeric tower.

## Parser
The parser is `apl/parse.go`.
```go
func (p *parser) parse(tokens []scan.Token) (Program, error)
```
called by the interpreter in `apl/apl.go: func (a *Apl) Parse(line string) (Program, error)`

It reads tokens from the right side and pushes them on a stack.
After pushing a token, it tries to *reduce* the stack by assembling it into a nested tree.

This is where it differes from a classic APL parser (to my understanding).
APL would do the reduction by partial evaluation.

At the end of parsing a line, only a single item may be left. Hence the error message
```
cannot reduce expression
```
if it cannot.

The *Reduction* step tries to reduce operators in the expression from the left and nouns from the right.
It also has to handle bracket indexing and assignments specially.
I'm not sure if a real APL parser does it in a similar way.
I hope someone will enlighten me some day.

The result of parsing is a *Program*. This is a bad name.
It refers to a single line of APL input. Maybe *Phrase* would be better.

I played around with the idea of parsing directly into an `apl.List` value, similarly to what I think K does.
This would allow to do Lisp style manipulation at runtime.
It's deferred for the future.

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
An APL Value is implemented as a go interface.
Anything that can be printed into a string can act as an `apl.Value` (`apl/value.go`):
```go
type Value interface {
	String(*Apl) string
}
```
It's not yet decided if this is enough. Maybe a Copy method has to be added later.

APL is an array language:
```go
type Array interface {
	Value
	At(int) Value
	Shape() []int
	Size() int
}
```
An array is also implemented as an interface. This makes it possible to add all kinds of special implementations later.
The most prominent implementation is the `apl.MixedArray`
```go
type MixedArray struct {
	Values []Value
	Dims   []int
}
```
Technically it's Values can be anyting, but the default implementation of primary functions and operators implies scalars / atoms.

Functions are implemented as an interface, as well:
```go
type Function interface {
	Call(a *Apl, L Value, R Value) (Value, error)
}
```
If L is nil, it's a monadic call.

Primitive functions and derived functions are examples of actual types that implement it.
Another one could be a method of an external go type.


# Numbers
Numbers are added to the type system by an `apl.Tower`, see `apl/tower.go`.

A tower is a linear escalation of numeric types.
A dyadic call to an elementary function (`+×⍟⌊...`) implies that both values are on the same level.

The default tower is implemented in `apl/numbers/` and contains Integer, Float, Complex and Time.

Time is a little bit special here. It could have been implemented outside the number system.
Amoung other benefits, this allows parsing Dates directly.

Two other tower implementations are provided in `apl/big`:
A big tower with Int and Rat and a precise tower with Float and Complex containing any number of bits.

The big tower is nice for playing with numbers.
It's Rat type can be used in combination with ⌹ to solve linear equations loss-less.
I have seen this from Bob Smith and wanted it too.

The precise tower can be initialized with any precision.
It could come in handy when testing convergence of numerical algorithms.

The towers are only available if they are registered at compile time, so it's possible to use only a single one.
If multiple tower are present, they can be switched but not mixed at runtime.

# Overloading primitive functions and operators
Each APL symbol can be registered together with a handler multiple times.

Once it's registered as a function, it cannot be changed to an operator.
Operators also have to agree, if they are monadic or dyadic.
The first registrant decides.

At runtime APL decides which implementation to use starting with the last registrant.
The decision is done by testing if the argument values fall into the *Domain* of the function or operator.

The `Domain` can be based on the argument type, but is more general.
It is defined in `apl.domain.go` as:
```go
type Domain interface {
	To(*Apl, Value, Value) (Value, Value, bool)
	String(*Apl) string
}
```
The domain function of a primitive, receives both arguments L and R and can decide what to do:
- accept the arguments as they are and handle the call
- convert the arguments and handle the call
- reject responsibility and leave the values unchanged.

Dispatching operators is done the same way, but receive LO and RO as arguments.
When they have to make the decision if they handle the call, they don't know about L and R of the derived function.

To unify the implementation, there is a package 'apl/domain' that helps building domain functions by composing them:

As an example, see the implementation of *deal*, the dyadic version of `?`.
It's domain part reads like this:
```go
Domain: Dyadic(Split(ToScalar(ToIndex(nil)), ToScalar(ToIndex(nil)))),
```
This is translated verbally into:
- It must be a dyadic call (L may not be nil)
- L must be a scalar, or convertible to one (any single element array)
- L must be an integer of compatible as e.g. `3J0`
- The same rules apply to R

The domain part does not have to cover all cases.
The decision what code part to take can still be done within the handler.
But once the call is accepted, it cannot be delegated back to the dispatcher.

Many of the helper functions in package domain come in two versions: `Is*` and `To*`.
`ToArray` would convert a scalar to a single element array and accept, while `IsArray` would reject it.

A positive side effect of using package domain, is that it keeps the reference `REF.md`
up-to-date with a call to `go generate`. The part of `?` looks like that:
```
?                                              
   deal                                        apl/primitives/query.go:19
   L?R  L toscalar index R toscalar index      
   roll                                        apl/primitives/query.go:13
   ?R  any  
```
It may not be the final version, but it's definetly better then having an out-dated hand written description.

Using the domain package to build a domain function is optional.
Not all primitives do this. 
Sometimes it's better to implement a special case directly then to put everything in the general package.

# Lambda functions
Lambda functions `apl/lambda.go` should be mostly compatible with Dyalog's dfns.
```apl
{⍺×⍵}/2 3 4
S←0{⍺>20:⍺⋄⍵∇⎕←⍺+⍵}1
{⍵>1000:⍵⋄∇⍵+1}1
```
Guards are supported, recursion and tail calls (in contrast to the host language).
But there is room for improvement:
- currently everything has to fit on one line 
- no error guards
	- instead a simpler general method of error handling is planned:
	- `EXPR :: TRAP` individually for any line without an exception stack

# Go interface
# Streams and concurrency
