package p

import (
	"reflect"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/xgo"
	"github.com/ktye/plot"
)

func newLine(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	l := plot.Line{}
	if i, ok := R.(apl.Int); ok {
		l.Id = int(i)
	}
	return Line{xgo.Value(reflect.ValueOf(&l))}, nil
}

// plotAddLine adds the line L to plot P and returns P.
// Calls can be chained: L1+L2+P.
func plotAddLine(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	P := R.(Plot)
	line := L.(Line)
	p := P.p()
	p.Lines = append(p.Lines, *line.l())
	return P, nil
}
