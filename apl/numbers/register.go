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
		Uniform: makeUniform,
		SetPP:   func(pp [2]int) { setpp(pp, m) },
	}
	return t
}

func makeUniform(v []apl.Value) (apl.Value, bool) {
	if len(v) == 0 {
		return nil, false
	}
	if t := reflect.TypeOf(v[0]); t == reflect.TypeOf(apl.Index(0)) {
		return makeIndexArray(v), true
	} else if t := reflect.TypeOf(v[0]); t == reflect.TypeOf(Integer(0)) {
		return makeIntegerArray(v), true
	} else if t == reflect.TypeOf(Float(0.0)) {
		return makeFloatArray(v), true
	} else if t == reflect.TypeOf(Complex(0)) {
		return makeComplexArray(v), true
	} else if t == reflect.TypeOf(y0) {
		return makeTimeArray(v), true
	}
	return nil, false
}

// Setpp sets default width and precision for numeric types.
func setpp(pp [2]int, m map[reflect.Type]apl.Numeric) {
	f := ""
	if pp[0] == 0 {
		f = fmt.Sprintf("%%.%dG", pp[1])
	} else {
		f = fmt.Sprintf("%%%d.%dF", pp[0], pp[1])
	}
	if pp == [2]int{0, 0} {
		f = ""
	}
	formats := map[reflect.Type]string{
		reflect.TypeOf(Float(0)):   f,
		reflect.TypeOf(Complex(0)): f + "J" + f,
	}
	for t, n := range m {
		n.Format = formats[t]
		m[t] = n
	}
}
