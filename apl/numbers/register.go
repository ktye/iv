// Package number defines the basic numeric types for APL.
// These are Integer, Float and Complex.
package numbers

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ktye/iv/apl"
)

// Register sets the default numeric tower Integer->Float->Complex.
func Register(a *apl.Apl) {
	if err := a.SetTower(newTower()); err != nil {
		panic(err)
	}

	pkg := map[string]apl.Value{
		"fmt": setformat{},
	}
	a.RegisterPackage("numbers", pkg)
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
		Class: 2,
		Parse: ParseComplex,
		Uptype: func(n apl.Number) (apl.Number, bool) {
			// Uptype converts a number to seconds, if the imag part is 0
			if imag(complex128(n.(Complex))) != 0 {
				return nil, false
			}
			d := time.Duration(int64(1e9 * real(complex128(n.(Complex)))))
			return Time(y0.Add(d)), true
		},
	}
	m[reflect.TypeOf(Time{})] = apl.Numeric{
		Class:  3,
		Parse:  ParseTime,
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

type setformat struct{}

func (s setformat) String(a *apl.Apl) string { return "fmt" }

func (_ setformat) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if L == nil {
		return nil, fmt.Errorf("fmt must be called dyadically")
	}
	n, ok := L.(apl.Number)
	if ok == false {
		return nil, fmt.Errorf("fmt: left argument must be a number")
	}
	s, ok := R.(apl.String)
	if ok == false {
		return nil, fmt.Errorf("fmt: right argument must be a string")
	}
	t := reflect.TypeOf(n)
	num, ok := a.Tower.Numbers[t]
	if ok == false {
		return nil, fmt.Errorf("fmt: %T is not a member of the current tower", n)
	}
	num.Format = string(s)
	a.Tower.Numbers[t] = num
	return apl.String(fmt.Sprintf("%T: %q\n", n, s)), nil
}
