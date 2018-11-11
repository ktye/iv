package apl

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"strings"
)

type Scalar expr
type Bool bool
type Int int64
type Float float64
type Complex complex128
type String string

func ParseNumber(s string) (Scalar, error) {
	s = strings.Replace(s, "¯", "-", -1)
	parseFloats := func(s string, idx int) (float64, float64, bool) {
		var err error
		var a, b float64
		if a, err = strconv.ParseFloat(s[:idx], 64); err != nil {
			return 0, 0, false
		}
		if b, err = strconv.ParseFloat(s[idx+1:], 64); err != nil {
			return 0, 0, false
		}
		return a, b, true
	}
	if idx := strings.Index(s, "@"); idx != -1 {
		amp, ang, ok := parseFloats(s, idx)
		if ok {
			return Complex(cmplx.Rect(amp, ang/180.0*math.Pi)), nil
		}
	} else if idx := strings.Index(s, "J"); idx != -1 {
		re, im, ok := parseFloats(s, idx)
		if ok {
			return Complex(complex(re, im)), nil
		}
	} else {
		if n, err := strconv.Atoi(s); err == nil {
			return Int(n), nil
		}
		if n, err := strconv.ParseFloat(s, 64); err == nil {
			return Float(n), nil
		}
	}
	return nil, fmt.Errorf("cannot parse number: %s", s)
}

func IsScalar(v Value) bool {
	switch v.(type) {
	case Bool, Int, Float, Complex, String:
		return true
		// TODO assignment?
	}
	return false
}

// ScalarArray wraps a Scalar into a size 1 array.
func ScalarArray(v Value) Value {
	return GeneralArray{
		Values: []Value{v},
		Dims:   []int{1},
	}
}

// ScalarValue returns a scalar value.
// The value could be a vector of size 1.
// It returns false, if the value is a larger array.
func ScalarValue(v Value) (Value, bool) {
	if a, ok := v.(Array); ok {
		if shape := a.Shape(); len(shape) == 1 && shape[0] == 1 {
			s, _ := a.At(0)
			return s, true
		} else {
			return v, false
		}
	}
	return v, true
}

/* TODO remove
// CompareScalars compares two scalar values and returns
// if a == b, and if a < b.
// Values are converted to the same type before comparison.
// It returns ErrNaN for floating point comparision.
// Complex number are compared for equality but return ErrCmpCmplx.
func CompareScalars(a, b Value) (bool, bool, error) {
	var err error
	a, b, err = SameNumericTypes(a, b)
	if err != nil {
		if as, ok := a.(String); ok {
			if bs, ok := b.(String); ok {
				return as == bs, as < bs, nil
			}
		}
		return false, false, fmt.Errorf("cannot compare %T and %T", a, b)
	}
	switch a := a.(type) {
	case Bool:
		return a == b.(Bool), a == true && b.(Bool) == false, nil
	case Int:
		return a == b.(Int), a < b.(Int), nil
	case Float:
		if math.IsNaN(float64(a)) || math.IsNaN(float64(b.(Float))) {
			return false, false, ErrNaN
		}
		return a == b.(Float), a < b.(Float), nil
	case Complex:
		if cmplx.IsNaN(complex128(a)) || cmplx.IsNaN(complex128(b.(Complex))) {
			return false, false, ErrNaN
		}
		return a == b.(Complex), false, ErrCmpCmplx
	default:
		return false, false, fmt.Errorf("cannot compare %T and %T", a, b)
	}
}

var ErrNaN error
var ErrCmpCmplx error
*/

/* TODO remove (ported to domain.Bool)
// ToBool converts numeric scalar or single element array to bool.
// It returns false if the value is not 0 or 1.
func ToBool(v Value) (bool, bool) {
	if n, ok := ToInt(v); ok {
		if n == 0 {
			return false, true
		} else if n == 1 {
			return true, true
		}
	}
	return false, false
}

// TODO remove (ported to domain.Int)
// ToInt converts a numeric scalar or a single element array to an int.
// It uptypes Bool and downtypes Float and Complex if they have no fractional
// or imaginary part.
// If v is not convertable, ToInt returns an error.
func ToInt(v Value) (int, bool) {
	v, ok := ScalarValue(v)
	if ok == false {
		return 0, false
	}
	switch v := v.(type) {
	case Bool:
		if v {
			return 1, true
		}
		return 0, true
	case Int:
		return int(v), true
	case Float:
		i := int(float64(v))
		if Float(i) == v {
			return i, true
		}
		return 0, false
	case Complex:
		c := complex128(v)
		if imag(c) != 0 {
			return 0, false
		}
		i := int(real(c))
		if float64(i) == real(c) {
			return i, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// TODO remove (ported to domain)
// ToFloat converts a numeric scalar or single element array to float64.
// It uptypes Bool and Int and downtypes Complex, if the imaginary part is zero.
// If v is not convertable, ToComplex returns false.
func ToFloat(v Value) (float64, bool) {
	v, ok := ScalarValue(v)
	if ok == false {
		return 0, false
	}
	switch v := v.(type) {
	case Bool:
		if v {
			return float64(1), true
		}
		return float64(0), true
	case Int:
		return float64(v), true
	case Float:
		return float64(v), true
	case Complex:
		if imag(complex128(v)) == 0 {
			return real(complex128(v)), true
		}
	}
	return 0, false
}

// TODO remove: ported to domain.
// ToComplex converts a numeric scalar or a single element array to complex128.
// It uptypes Bool, Int and Float.
// If v is not convertable, ToComplex returns false.
func ToComplex(v Value) (complex128, bool) {
	v, ok := ScalarValue(v)
	if ok == false {
		return 0, false
	}
	switch v := v.(type) {
	case Bool:
		if v {
			return complex(1, 0), true
		}
		return complex(0, 0), true
	case Int:
		return complex(float64(v), 0), true
	case Float:
		return complex(float64(v), 0), true

	case Complex:
		return complex128(v), true
	default:
		return 0, false
	}
}

var typeOrder map[reflect.Type]int

// TODO remove
// SameNumericTypes converts a or b to the higher numeric type.
// If both types are identical, they are returned, even if not numeric.
// Otherwise it returns an error, if any of both is not Bool, Int, Float or Complex.
func SameNumericTypes(a, b Value) (Value, Value, error) {
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		return a, b, nil
	}
	av := typeOrder[reflect.TypeOf(a)]
	bv := typeOrder[reflect.TypeOf(b)]
	if av == 0 || bv == 0 {
		return nil, nil, fmt.Errorf("cannot convert %T and %T to the same type", a, b)
	}
	if av == bv {
		return a, b, nil
	}
	swap := false
	if av > bv {
		av, bv = bv, av
		a, b = b, a
		swap = true
	}

	if av == 1 {
		i := Int(0)
		if a.(Bool) == true {
			i = 1
		}
		if bv == 2 {
			a = i
		} else if bv == 3 {
			a = Float(i)
		} else if bv == 4 {
			a = Complex(complex(float64(i), 0))
		}
	} else if av == 2 {
		if bv == 3 {
			a = Float(a.(Int))
		} else if bv == 4 {
			a = Complex(complex(float64(a.(Int)), 0))
		}
	} else if av == 3 {
		a = Complex(complex(float64(a.(Float)), 0))
	}

	if swap {
		a, b = b, a
	}
	return a, b, nil
}
*/

// Format is used by the stringers of default types.
type Format struct {
	Bool    string
	Int     string
	Float   string
	Complex string
	String  string
	Minus   bool
}

// String formats a Bool as 1 or 0 by default.
// This can be changed to true/false by setting a.format.Bool to "%t".
func (b Bool) String(a *Apl) string {
	if a != nil && a.format.Bool != "" {
		return fmt.Sprintf(a.format.Bool, bool(b))
	}
	if b {
		return "1"
	}
	return "0"
}

func (b Bool) Eval(a *Apl) (Value, error) {
	return b, nil
}

func (b Bool) toFloat() Float {
	if b {
		return Float(1)
	}
	return Float(0)
}

// String formats i as a decimal string and and uses ¯ for negative numbers.
// The format can be changed in a.format.Int.
// A normal minus sign is used, by setting a.format.Minus = true.
func (i Int) String(a *Apl) string {
	format := "%v"
	if a != nil && a.format.Int != "" {
		format = a.format.Int
	}
	s := fmt.Sprintf(format, int64(i))
	if a == nil || a.format.Minus == false {
		return strings.Replace(s, "-", "¯", 1)
	}
	return s
}

func (i Int) Eval(a *Apl) (Value, error) {
	return i, nil
}

// String formats f using %v by default an replacing ¯ for -.
// The format can be changed in a.format.Float.
// A normal minus sign is used, by setting a.format.Minus = true.
func (f Float) String(a *Apl) string {
	format := "%v"
	if a != nil && a.format.Float != "" {
		format = a.format.Float
	}
	s := fmt.Sprintf(format, float64(f))
	if a == nil || a.format.Minus == false {
		return strings.Replace(s, "-", "¯", 1)
	}
	return s
}
func (f Float) Eval(a *Apl) (Value, error) {
	return f, nil
}

// String formats c as %v@%v by default which uses magnitude and degree
// and replaces ¯ for -.
// This can be changed by setting a.format.Complex.
// If the format string contains a single %, the complex number is passed.
// Otherwise if it does not contain an @, real and imaginary parts are used.
// Examples are: "%v", "%.3f", "%.3f@%.0f°", "%vJ%v".
// A normal minus sign is used, by setting a.format.Minus = true.
func (c Complex) String(a *Apl) string {
	format := "%v@%v"
	if a != nil && a.format.Complex != "" {
		format = a.format.Complex
	}
	var s string
	if strings.Count(format, "%") == 1 {
		s = fmt.Sprintf(format, complex128(c))
	} else {
		a, b := real(c), imag(c)
		if strings.Index(format, "@") != -1 {
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
	if a == nil || a.format.Minus == false {
		return strings.Replace(s, "-", "¯", 1)
	}
	return s
}

func (c Complex) Eval(a *Apl) (Value, error) {
	return c, nil
}

// String formats s with %s by default.
// The format can be changed in Format.String.
func (s String) String(a *Apl) string {
	if a == nil {
		// For program output, replace single " by doubles.
		return strings.Replace(string(s), `"`, `""`, -1)
	}
	format := "%s"
	if a.format.String != "" {
		format = a.format.String
	}
	return fmt.Sprintf(format, string(s))
}

func (s String) Eval(a *Apl) (Value, error) {
	return s, nil
}

/* TODO remove
func init() {
	typeOrder = make(map[reflect.Type]int)
	typeOrder[reflect.TypeOf(Bool(false))] = 1
	typeOrder[reflect.TypeOf(Int(0))] = 2
	typeOrder[reflect.TypeOf(Float(0))] = 3
	typeOrder[reflect.TypeOf(Complex(0))] = 4
	ErrNaN = errors.New("value is NaN")
	ErrCmpCmplx = errors.New("complex values cannot be ordered")
}
*/
