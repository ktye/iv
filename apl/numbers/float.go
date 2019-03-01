package numbers

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/ktye/iv/apl"
)

type Float float64

// String formats a Float as a string.
// The format string is passed to fmt and - is replaced by ¯,
// except if the first rune is -.
func (n Float) String(f apl.Format) string {
	format, minus := getformat(f, n)
	if format == "" {
		switch prec := f.PP; {
		case prec == 0:
			format = "%.6G"
		case prec == -16:
			format = "%b"
		case prec < 0:
			format = "%v"
		default:
			format = fmt.Sprintf("%%.%dG", prec)
		}
	}
	s := fmt.Sprintf(format, float64(n))
	if minus == false {
		s = strings.Replace(s, "-", "¯", -1)
	}
	return s
}
func (f Float) Copy() apl.Value { return f }

// ParseFloat parses a Float. It replaces ¯ with -, then uses ParseFloat.
// A trailing . is stripped, so that "2." is parsed as a float.
func ParseFloat(s string) (apl.Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
	if n, err := strconv.ParseFloat(strings.TrimSuffix(s, "."), 64); err == nil {
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

func (f Float) Less(R apl.Value) (apl.Bool, bool) {
	return apl.Bool(f < R.(Float)), true
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
		return apl.Int(1), true
	} else if f < 0 {
		return apl.Int(-1), true
	}
	return apl.Int(0), true
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

func (f Float) Abs() (apl.Value, bool) {
	return Float(math.Abs(float64(f))), true
}
func (f Float) Abs2(R apl.Value) (apl.Value, bool) {
	return Float(math.Remainder(float64(f), float64(R.(Float)))), true
}

func (f Float) Floor() (apl.Value, bool) {
	return Float(math.Floor(float64(f))), true
}
func (f Float) Ceil() (apl.Value, bool) {
	return Float(math.Ceil(float64(f))), true
}

func (f Float) Gamma() (apl.Value, bool) {
	y := Float(math.Gamma(float64(f) + 1))
	if e, ok := isException(y); ok {
		return e, true
	}
	return y, true
}

func beta(a, b float64) float64 {
	ga, sa := math.Lgamma(a)
	gb, sb := math.Lgamma(b)
	gab, sab := math.Lgamma(a + b)
	sn := float64(sa * sb * sab)
	return sn * math.Exp(ga+gb-gab)
}
func (L Float) Gamma2(R apl.Value) (apl.Value, bool) {
	// Dyalog: Beta(X,Y) ←→ ÷Y×(X-1)!X+Y-1
	// Solving for R and L, with: R=X+Y-1 and L=X-1
	// L!R = (R over L) = 1/((R-L)*beta(R-L, L+1))
	r := float64(R.(Float))
	l := float64(L)
	f := Float(1.0 / ((r - l) * beta(r-l, l+1)))
	if e, ok := isException(f); ok {
		return e, true
	}
	return f, true
}

func (L Float) PiTimes() (apl.Value, bool) {
	return Float(math.Pi) * L, true
}
func (L Float) Trig(R apl.Value) (apl.Value, bool) {
	x := float64(R.(Float))
	var y float64
	switch L {
	case 0:
		y = math.Sqrt(1.0 - x*x)
	case -1:
		y = math.Asin(x)
	case 1:
		y = math.Sin(x)
	case -2:
		y = math.Acos(x)
	case 2:
		y = math.Cos(x)
	case -3:
		y = math.Atan(x)
	case 3:
		y = math.Tan(x)
	case -4:
		y = 0
		if x != -1 {
			y = (x + 1.0) * math.Sqrt((x-1.0)/(x+1.0))
		}
	case 4:
		y = math.Sqrt(1.0 + x*x)
	case -5:
		y = math.Asinh(x)
	case 5:
		y = math.Sinh(x)
	case -6:
		y = math.Acosh(x)
	case 6:
		y = math.Cosh(x)
	case -7:
		y = math.Atanh(x)
	case 7:
		y = math.Tanh(x)
	case -8:
		y = -math.Sqrt(x*x - 1.0)
	case 8:
		y = math.Sqrt(x*x - 1.0)
	case -9, 9:
		y = x // 9: real part
	case -10:
		y = x
	case 10:
		y = math.Abs(x)
	case -11:
		return nil, false
	case 11:
		y = x // imag part
	case -12:
		return nil, false
	case 12: // phase
		y = 1
		if x == 0 {
			y = 0
		} else if x < 0 {
			y = -1
		}
	default:
		return nil, false
	}
	return Float(y), true
}

func (L Float) Gcd(R apl.Value) (apl.Value, bool) {
	l := math.Abs(float64(L))
	r := math.Abs(float64(R.(Float)))

	ab, lok := big.NewRat(1, 1).SetString(fmt.Sprintf("%v", l))
	cd, rok := big.NewRat(1, 1).SetString(fmt.Sprintf("%v", r))
	if lok == false || rok == false {
		return nil, false
	}
	a := ab.Num()
	b := ab.Denom()
	c := cd.Num()
	d := cd.Denom()
	ad := big.NewInt(0).Mul(a, d)
	cb := big.NewInt(0).Mul(c, b)
	gcd := big.NewInt(0).GCD(nil, nil, ad, cb)
	bd := big.NewInt(0).Mul(b, d)
	rat := big.NewRat(1, 1).SetFrac(gcd, bd)
	f, _ := rat.Float64()
	return Float(f), true
}
