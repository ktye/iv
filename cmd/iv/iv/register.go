// Package iv is an APL extension that sends subarrays over a channel
//
// Iv is used by the command line program cmd/iv, but can also be registered
// in any APL instance by calling
//	Register(a *Apl).
package iv

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

// Register adds iv to the interpreter.
// It provides two functions:
//
//	C2 ←iv→p C1
// C1: input channel (lines of text)
// C2: output channel List(Number value; separators)
//
// 	C3 ←(f;R;)iv→f C2
// f: lambda function receiving
//	⍺ termination level
//	⍵ sub-array
// R: rank of ⍵
// C3:
//
//	({+/⍵};2)iv→f iv→p <`data
func Register(a *apl.Apl, name string) {
	if name == "" {
		name = "iv"
	}
	pkg := map[string]apl.Value{
		"r": &InputParser{},
	}
	a.RegisterPackage(name, pkg)
}

func (_ *InputParser) String(a *apl.Apl) string {
	return "iv r"
}

// TODO: this should be converted to a channel L

var Stdin io.Reader

func (p *InputParser) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if p.Reader != nil {
		return nil, fmt.Errorf("iv: r can only be called once")
		// TODO: This could be made possible, but cmd/iv needs it only once.
	}
	if Stdin == nil {
		return nil, fmt.Errorf("iv: package is not initialized Stdin=nil")
	}

	toidx := domain.ToIndex(nil)
	if n, ok := toidx.To(a, R); ok == false {
		return nil, fmt.Errorf("iv r expects an int argument: rank, got %T", R)
	} else {
		p.Rank = int(n.(apl.Index))
		if p.Rank < 0 {
			return nil, fmt.Errorf("iv: rank must be > 0")
		}
	}

	// TODO: accept L as a channel to read input from
	p.Reader = bufio.NewReader(tabularText(Stdin))
	p.Apl = a
	p.Separator = '\n'

	c := apl.NewChannel()
	if p.Rank == 0 {
		go p.sendScalars(c)
	} else {
		go p.sendArrays(c)
	}
	return c, nil
}
