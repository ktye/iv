package a

import (
	"time"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/scan"
)

// toCommand attaches the Rewrite method to the function.
type toCommand func([]scan.Token) []scan.Token

func (f toCommand) Rewrite(t []scan.Token) []scan.Token {
	return f(t)
}

func symbol(s string) scan.Token {
	return scan.Token{T: scan.Symbol, S: s}
}

// rw0 is a scan.Command that rewrite the symbol with a 0 argument.
// Example:
//	/q	is rewritten to a→q 0
type rw0 string

func (r rw0) Rewrite(t []scan.Token) []scan.Token {
	sym := scan.Token{T: scan.Identifier, S: "a→" + string(r)}
	num := scan.Token{T: scan.Number, S: "0"}
	tokens := make([]scan.Token, len(t)+2)
	tokens[0] = sym
	tokens[1] = num
	copy(tokens[2:], t)
	return tokens
}

// printvar prints a string representation of the value.
// If the value is a string that is a valid variable name, it is dereferenced.
// This allows to print the definition of lambda functions.
func printvar(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	s, ok := R.(apl.String)
	if !ok {
		return apl.String(R.String(a)), nil
	}
	v := a.Lookup(string(s))
	if v == nil {
		return s, nil
	}
	return apl.String(v.String(a)), nil
}

func printCmd(t []scan.Token) []scan.Token {
	return append([]scan.Token{scan.Token{T: scan.Identifier, S: "a→p"}}, t...)
}

// Timer is used to time an expression. It is called by the rewrite command /t
// If the argument is a time, it returns the elapsed duration since that time.
// Otherwise it returns the current time.
func timer(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	t, ok := R.(numbers.Time)
	if !ok {
		return numbers.Time(time.Now()), nil
	}
	dt := time.Since(time.Time(t))
	y0, _ := time.Parse("15h04", "00h00") // see apl/numbers/time.go
	return numbers.Time(y0.Add(dt)), nil
}

// timeCmd rewrites the tokens to calculate the duration.
//	T__← a→t 0 ⋄ [TOKENS] ⋄ a→t T__
func timeCmd(t []scan.Token) []scan.Token {
	t__ := scan.Token{T: scan.Identifier, S: "T__"}
	asn := symbol("←")
	tim := scan.Token{T: scan.Identifier, S: "a→t"}
	num := scan.Token{T: scan.Number, S: "0"}
	dia := scan.Token{T: scan.Diamond, S: "⋄"}

	tokens := []scan.Token{t__, asn, tim, num, dia}
	tokens = append(tokens, t...)
	return append(tokens, dia, tim, t__)
}
