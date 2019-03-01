// Package number defines the basic numeric types for APL.
// These are Integer, Float and Complex.
package numbers

import (
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
	m := make(map[reflect.Type]*apl.Numeric)
	m[reflect.TypeOf(Float(0))] = &apl.Numeric{
		Class:  0,
		Parse:  ParseFloat,
		Uptype: floatToComplex,
	}
	m[reflect.TypeOf(Complex(0))] = &apl.Numeric{
		Class: 1,
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
	m[reflect.TypeOf(Time{})] = &apl.Numeric{
		Class:  2,
		Parse:  ParseTime,
		Uptype: func(n apl.Number) (apl.Number, bool) { return n, false },
	}
	t := apl.Tower{
		Numbers: m,
		Import: func(n apl.Number) apl.Number {
			if b, ok := n.(apl.Bool); ok {
				if b {
					return Float(1)
				}
				return Float(0)
			} else if n, ok := n.(apl.Int); ok {
				return Float(n)
			}
			return n
		},
		Uniform: makeUniform,
	}
	return t
}

func makeUniform(v []apl.Value) (apl.Value, bool) {
	if len(v) == 0 {
		return nil, false
	}
	if t := reflect.TypeOf(v[0]); t == reflect.TypeOf(apl.Bool(false)) {
		return makeBoolArray(v), true
	} else if t := reflect.TypeOf(v[0]); t == reflect.TypeOf(apl.Int(0)) {
		return makeIndexArray(v), true
	} else if t == reflect.TypeOf(Float(0.0)) {
		return makeFloatArray(v), true
	} else if t == reflect.TypeOf(Complex(0)) {
		return makeComplexArray(v), true
	} else if t == reflect.TypeOf(y0) {
		return makeTimeArray(v), true
	}
	return nil, false
}

func makeBoolArray(v []apl.Value) apl.BoolArray {
	f := make([]bool, len(v))
	for i, e := range v {
		f[i] = bool(e.(apl.Bool))
	}
	return apl.BoolArray{
		Dims:  []int{len(v)},
		Bools: f,
	}
}
func makeIndexArray(v []apl.Value) apl.IntArray {
	f := make([]int, len(v))
	for i, e := range v {
		f[i] = int(e.(apl.Int))
	}
	return apl.IntArray{
		Dims: []int{len(v)},
		Ints: f,
	}
}

func getformat(f apl.Format, num apl.Value) (string, bool) {
	if f.Fmt == nil {
		return "", false
	}
	s := f.Fmt[reflect.TypeOf(num)]
	if len(s) > 0 && s[0] == '-' {
		return s[1:], true
	}
	if f.PP < -1 {
		return "", true // intended for external interchange (full prec, with -).
	}
	return s, false
}
