package numbers

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ktye/iv/apl"
)

type Float float64

// String formats a Float as a string.
// The format string is passed to fmt and - is replaced by ¯,
// except if the first rune is -.
func (n Float) String(a *apl.Apl) string {
	format, minus := getformat(a, n, "%v")
	s := fmt.Sprintf(format, float64(n))
	if minus == false {
		s = strings.Replace(s, "-", "¯", 1)
	}
	return s
}

// ParseFloat parses a Float. It replaces ¯ with -, then uses ParseFloat.
func ParseFloat(s string) (apl.Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return Float(n), true
	}
	return Float(0), false
}

// ToIndex converts a Float to an int, if an exact conversion is possible.
func (f Float) ToIndex() (int, bool) {
	n := int(f)
	if Float(n) == f {
		return n, true
	}
	return 0, false
}

func floatToComplex(f apl.Number) (apl.Number, bool) {
	return Complex(complex(float64(f.(Float)), 0)), true
}

func (f Float) Add() (apl.Value, bool) {
	return f, true
}
func (f Float) Add2(R apl.Value) (apl.Value, bool) {
	return f + R.(Float), true
}

func (f Float) Sub() (apl.Value, bool) {
	return -f, true
}
func (f Float) Sub2(R apl.Value) (apl.Value, bool) {
	return f - R.(Float), true
}

func (f Float) Mul() (apl.Value, bool) {
	if f > 0 {
		return Integer(1), true
	} else if f < 0 {
		return Integer(-1), true
	}
	return Integer(0), true
}
func (f Float) Mul2(R apl.Value) (apl.Value, bool) {
	return f * R.(Float), true
}

func (f Float) Div() (apl.Value, bool) {
	n := 1.0 / float64(f)
	if e, ok := isException(Float(n)); ok {
		return e, true
	}
	return Float(n), true
}
func (f Float) Div2(b apl.Value) (apl.Value, bool) {
	n := Float(float64(f) / float64(b.(Float)))
	if e, ok := isException(n); ok {
		return e, true
	}
	return Float(n), true
}

func (f Float) Pow() (apl.Value, bool) {
	return Float(math.Exp(float64(f))), true
}
func (f Float) Pow2(R apl.Value) (apl.Value, bool) {
	return Float(math.Pow(float64(f), float64(R.(Float)))), true
}

func (f Float) Log() (apl.Value, bool) {
	l := math.Log(float64(f))
	if e, ok := isException(Float(l)); ok {
		return e, true
	}
	return Float(l), true
}
func (f Float) Log2(R apl.Value) (apl.Value, bool) {
	l := math.Log(float64(f))
	r := math.Log(float64(R.(Float)))
	if e, ok := isException(Float(l)); ok {
		return e, true
	}
	if e, ok := isException(Float(r)); ok {
		return e, true
	}
	return Float(r) / Float(l), true
}
