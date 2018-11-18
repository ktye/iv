package numbers

import (
	"fmt"
	"math"
	"math/cmplx"
	"strings"

	"github.com/ktye/iv/apl"
)

type Complex complex128

// String formats a Float as a string.
// If the format string contains a single %, it is passed to fmt
// with the complex arguments.
// If it contains an "a", two format strings are assumed and magnitude
// and degree are passed to fmt.
// Otherwise real and imag parts are passed.
// By default - is replaced with ¯, expept if the format string
// starts with -.
// Examples:
//	"%.3f", "%ga%.0f", "-%v", "%.5fJ%.5f"
func (c Complex) String(a *apl.Apl) string {
	format, minus := getformat(a, c, "%va%v")
	var s string
	if strings.Count(format, "%") == 1 {
		s = fmt.Sprintf(format, complex128(c))
	} else {
		a, b := real(c), imag(c)
		if strings.Index(format, "a") != -1 {
			a, b = cmplx.Polar(complex128(c))
			b *= 180.0 / math.Pi
			if b < 0 {
				b += 360
			}
			if a == 0 {
				b = 0
			}
			if b == -0 || b == 360 {
				b = 0
			}
		}
		s = fmt.Sprintf(format, a, b)
	}
	if minus == false {
		s = strings.Replace(s, "-", "¯", 1)
	}
	return s
}

// ParseComplex parses a Complex from a string.
// The number may be given as MAGNITUDEaANGLE with the angle in degree,
// or as realJimag.
// Both parts are parsed as Floats.
// If neither "a" or "J" are within the string, it is parsed with 0 imag part.
func ParseComplex(s string) (apl.Number, bool) {
	var c Complex
	if idx := strings.Index(s, "a"); idx != -1 {
		mag, ok := ParseFloat(s[:idx])
		if ok == false {
			return c, false
		}
		deg, ok := ParseFloat(s[idx+1:])
		if ok == false {
			return c, false
		}
		f := float64(mag.(Float))
		switch deg.(Float) {
		case 0:
			return Complex(complex(f, 0)), true
		case 90:
			return Complex(complex(0, f)), true
		case 180:
			return Complex(complex(-f, 0)), true
		case 270:
			return Complex(complex(0, -f)), true
		}
		return Complex(cmplx.Rect(f, float64(deg.(Float))/180.0*math.Pi)), true
	} else if idx := strings.Index(s, "J"); idx != -1 {
		re, ok := ParseFloat(s[:idx])
		if ok == false {
			return c, false
		}
		im, ok := ParseFloat(s[idx+1:])
		if ok == false {
			return c, false
		}
		return Complex(complex(float64(re.(Float)), float64(im.(Float)))), true
	} else {
		if n, ok := ParseFloat(s); ok == false {
			return c, false
		} else {
			return Complex(complex(float64(n.(Float)), 0)), true
		}
	}
}

// ToIndex converts a Complex to an int, if an exact conversion is possible.
func (c Complex) ToIndex() (int, bool) {
	if imag(complex128(c)) != 0 {
		return 0, false
	}
	r := real(complex128(c))
	n := int(r)
	if float64(n) == r {
		return n, true
	}
	return 0, false
}

func (c Complex) Add() (apl.Value, bool) {
	return Complex(cmplx.Conj(complex128(c))), true
}
func (c Complex) Add2(R apl.Value) (apl.Value, bool) {
	return c + R.(Complex), true
}

func (c Complex) Sub() (apl.Value, bool) {
	return -c, true
}
func (c Complex) Sub2(R apl.Value) (apl.Value, bool) {
	return c - R.(Complex), true
}

func (c Complex) Mul() (apl.Value, bool) {
	if c == 0 {
		return c, true
	}
	r := cmplx.Abs(complex128(c))
	return Complex(complex(real(c)/r, imag(c)/r)), true
}
func (c Complex) Mul2(R apl.Value) (apl.Value, bool) {
	return c * R.(Complex), true
}

func (c Complex) Div() (apl.Value, bool) {
	r := Complex(complex(1, 0) / complex128(c))
	if e, ok := isException(r); ok {
		return e, true
	}
	return r, true
}
func (c Complex) Div2(b apl.Value) (apl.Value, bool) {
	r := Complex(complex128(c) / complex128(b.(Complex)))
	if e, ok := isException(r); ok {
		return e, true
	}
	return r, true
}

func (c Complex) Pow() (apl.Value, bool) {
	return Complex(cmplx.Exp(complex128(c))), true
}
func (c Complex) Pow2(R apl.Value) (apl.Value, bool) {
	return Complex(cmplx.Pow(complex128(c), complex128(R.(Complex)))), true
}

func (c Complex) Log() (apl.Value, bool) {
	l := cmplx.Log(complex128(c))
	if e, ok := isException(Complex(l)); ok {
		return e, true
	}
	return Complex(l), true
}
func (c Complex) Log2(R apl.Value) (apl.Value, bool) {
	l := cmplx.Log(complex128(c))
	r := cmplx.Log(complex128(R.(Complex)))
	if e, ok := isException(Complex(l)); ok {
		return e, true
	}
	if e, ok := isException(Complex(r)); ok {
		return e, true
	}
	return Complex(r) / Complex(l), true
}
