// Package iv is an APL extension that sends subarrays over a channel
//
// Iv is used by the command line program cmd/iv, but can also be registered
// in any APL instance by calling
//	Register(a *Apl).
package iv

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

// TODO: This reads from stdin hard coded. That's ok for cmd/iv.
//       To make this useful as a package, it should use a reader.
//	 This is planned to be provided by a more general io package.

func Register(a *apl.Apl) {
	pkg := map[string]apl.Value{
		"r": &InputParser{},
	}
	a.RegisterPackage("iv", pkg)
}

func (_ *InputParser) String(a *apl.Apl) string {
	return "iv r"
}

func (p *InputParser) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if p.Reader != nil {
		return nil, fmt.Errorf("iv: r can only be called once")
		// TODO: This could be made possible, but cmd/iv needs it only once.
	}

	toidx := domain.ToIndex(nil)
	if n, ok := toidx.To(a, R); ok == false {
		return nil, fmt.Errorf("iv r expects an int argument: rank, got %T", R)
	} else {
		p.Rank = int(n.(apl.Index))
		if p.Rank <= 0 {
			return nil, fmt.Errorf("iv: rank must be > 0")
		}
	}

	// For testing: we accept a string L (a variable name) and read from it's string value
	// instead of stdin.
	// TODO: This should change to a channel.
	if L == nil {
		p.Reader = bufio.NewReader(tabularText(os.Stdin))
	} else {
		if v, ok := L.(apl.String); ok == false {
			return nil, fmt.Errorf("iv r: L must be a string")
		} else {
			p.Reader = bufio.NewReader(strings.NewReader(string(v)))
		}
	}
	p.a = a
	p.Separator = '\n'

	c := apl.NewChannel()
	go p.sendArrays(c)
	/*
		go func(c apl.Channel) {
			var buf []apl.Value
			s := "empty stack"
			for {
			}
			close(c[0])
		}(c)
	*/
	return c, nil
}
