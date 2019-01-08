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
It reads tokens from the right side and pushes them on the stack.
After each token it tries to *reduce* the expression by assembling a nested tree.

This is where it differes from a classic APL parser (to my understanding).
An APL would do the reduction by partial evaluation.

At the end of parsing a line, only a single item may be left. Hence the error message
```
cannot reduce expression
```
if it fails to do so.

The *Reduction* step tries to reduce operator expressions from the left 
and noun expressions from the right side of the partial list, plus some special handling for bracket indexing and assignments.
I don't know if a real APL parser does it in a similar way.
I hope someone will enlighten me some day.
