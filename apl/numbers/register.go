// Package number defines the basic numeric types for APL.
// These are Integer, Float and Complex.
package numbers

import (
	"reflect"

	"github.com/ktye/iv/apl"
)

// Register sets the default numeric tower Integer->Float->Complex.
func Register(a *apl.Apl) {
	if err := a.SetTower(newTower()); err != nil {
		panic(err)
	}
}

func newTower() apl.Tower {
	m := make(map[reflect.Type]apl.Numeric)
	m[reflect.TypeOf(Integer(0))] = apl.Numeric{
		Class:  0,
		Parse:  ParseInteger,
		Uptype: intToFloat,
	}
	m[reflect.TypeOf(Float(0))] = apl.Numeric{
		Class:  1,
		Parse:  ParseFloat,
		Uptype: floatToComplex,
	}
	m[reflect.TypeOf(Complex(0))] = apl.Numeric{
		Class:  2,
		Parse:  ParseComplex,
		Uptype: func(n apl.Number) (apl.Number, bool) { return n, false },
	}
	t := apl.Tower{
		Numbers: m,
		FromIndex: func(n int) apl.Number {
			return Integer(n)
		},
	}
	return t
}
