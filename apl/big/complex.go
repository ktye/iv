package big

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

type Complex struct {
	re, im *big.Float
}

// String formats a complex as a string. The polar form is not supported.
func (c Complex) String(a *apl.Apl) string {
	// TODO parse MAGaDEG
	def := "%vJ%v"
	format, minus := getformat(a, c, "%vJ%v")
	if strings.Count(format, "%") != 2 {
		format = def
	}
	s := fmt.Sprintf(format, c.re, c.im)
	if minus == false {
		s = strings.Replace(s, "-", "¯", -1)
	}
	return s
}

func ParseComplex(s string, prec uint) (apl.Number, bool) {
	// If the number is in polar form, parse with numbers.Complex,
	// neglecting additional precision.
	if strings.Index(s, "a") != 0 {
		z, ok := numbers.ParseComplex(s)
		if ok == false {
			return nil, false
		}
		re := big.NewFloat(real(complex128(z.(numbers.Complex)))).SetPrec(prec)
		im := big.NewFloat(imag(complex128(z.(numbers.Complex)))).SetPrec(prec)
		return Complex{re, im}, true
	}
	s = strings.Replace(s, "¯", "-", -1)
	idx := strings.Index(s, "J")
	if idx == -1 {
		idx = strings.Index(s, "j")
	}
	if idx == -1 {
		z, _, err := big.NewFloat(0).SetPrec(prec).Parse(s, 10)
		if err != nil {
			return nil, false
		}
		return Complex{z, big.NewFloat(0).SetPrec(prec)}, true
	}

	re, _, err := big.NewFloat(0).SetPrec(prec).Parse(s[:idx], 10)
	if err != nil {
		return nil, false
	}
	im, _, err := big.NewFloat(0).SetPrec(prec).Parse(s[idx+1:], 10)
	if err != nil {
		return nil, false
	}
	return Complex{re, im}, true
}

func (c Complex) ToIndex() (int, bool) {
	if c.im.Sign() != 0 {
		return 0, false
	}
	return Float{c.re}.ToIndex()
}

func floatToComplex(f apl.Number) (apl.Number, bool) {
	z := f.(Float).cpy()
	return Complex{z, z.SetInt64(0)}, false
}

func (c Complex) cpy() Complex {
	re, im := new(big.Float), new(big.Float)
	return Complex{re.Copy(c.re), im.Copy(c.im)}
}

func (c Complex) Add() (apl.Value, bool) {
	z := c.cpy()
	z.im = z.im.Neg(z.im)
	return z, true
}
func (c Complex) Add2(R apl.Value) (apl.Value, bool) {
	z := c.cpy()
	r := R.(Complex)
	z.re = z.re.Add(z.re, r.re)
	z.im = z.im.Add(z.im, r.im)
	return z, true
}

func (c Complex) Sub() (apl.Value, bool) {
	z := c.cpy()
	z.re = z.re.Neg(z.re)
	z.im = z.im.Neg(z.im)
	return z, true
}
func (c Complex) Sub2(R apl.Value) (apl.Value, bool) {
	z := c.cpy()
	r := R.(Complex)
	z.re = z.re.Sub(z.re, r.re)
	z.im = z.im.Sub(z.im, r.im)
	return z, true
}

func (c Complex) Mul() (apl.Value, bool) {
	z := c.cpy()
	r := c.abs()
	if r.Sign() == 0 {
		return Complex{z.re.SetInt64(0), z.im.SetInt64(0)}, true
	}
	return Complex{z.re.Quo(z.re, r), z.im.Quo(z.im, r)}, true
}
func (L Complex) Mul2(R apl.Value) (apl.Value, bool) {
	// ac-bd + J (ad+bc)
	l, r := L.cpy(), R.(Complex).cpy()
	a, b := l.re, l.im
	c, d := r.re, r.im
	ac := new(big.Float).Copy(a)
	bd := new(big.Float).Copy(b)
	ad := new(big.Float).Copy(a)
	bc := new(big.Float).Copy(b)
	ac = ac.Mul(a, c)
	bd = bd.Mul(b, d)
	ad = ad.Mul(a, d)
	bc = bc.Mul(b, c)
	a = a.Sub(ac, bd)
	b = b.Add(ad, bc)
	return Complex{a, b}, true
}

func (c Complex) Div() (apl.Value, bool) {
	z := c.cpy()
	z.re.SetInt64(1)
	z.im.SetInt64(0)
	return z.Div2(c)
}
func (l Complex) Div2(R apl.Value) (apl.Value, bool) {
	// This is not the classical Smith algorithm, but:
	// Michael Baudin, Robert Smith: A robust complex division in Scilab (2012)
	// https://archive.org/details/arxiv-1210.4539
	// It implements compdiv_improved from section 3,
	// not the scaled version compdiv_robost form sec 3.4.

	r := R.(Complex)
	if r.re.Sign() == 0 && r.im.Sign() == 0 {
		return numbers.Inf, true
	} else if l.re.Sign() == 0 && l.im.Sign() == 0 {
		return l.cpy(), true // zero
	}
	a := new(big.Float).Copy(l.re)
	b := new(big.Float).Copy(l.im)
	c := new(big.Float).Copy(r.re)
	d := new(big.Float).Copy(r.im)
	dd := new(big.Float).Copy(d) // dd = abs(d)
	if dd.Sign() < 0 {
		dd = dd.Neg(dd)
	}
	cc := new(big.Float).Copy(c) // cc = abs(c)
	if cc.Sign() < 0 {
		cc = cc.Neg(cc)
	}
	var e, f *big.Float
	if dd.Cmp(cc) <= 0 { // abs(d) <= abs(c)
		e, f = div_(a, b, c, d)
	} else {
		e, f = div_(b, a, d, c)
		f = f.Neg(f)
	}
	return Complex{e, f}, true
}

func div_(a, b, c, d *big.Float) (*big.Float, *big.Float) {
	r := new(big.Float).Quo(d, c) // r = d/c
	t := new(big.Float).Mul(d, r) // d*r
	t = t.Add(t, c)               // c + d*r
	one := new(big.Float).Copy(t).SetInt64(1)
	t = t.Quo(one, t) // t = 1/(c + d*r)
	e := new(big.Float)
	f := new(big.Float)
	if r.Sign() != 0 {
		e = e.Mul(b, r) // b * r
		e = e.Add(e, a) // a + b*r
		e = e.Mul(e, t) // e = (a + b*r)*t
		f = f.Mul(a, r) // a*r
		f = f.Sub(b, f) // b - a*r
		f = f.Mul(f, t) // f = (b - a*r)*t
	} else {
		e = e.Quo(b, c) // b/c
		e = e.Mul(e, d) // d * (b/c)
		e = e.Add(e, a) // a + d*(b/c)
		e = e.Mul(e, t) // e = (a + d*(b/c))*t
		f = f.Quo(a, c) // a/c
		f = f.Mul(f, d) // d*(a/c)
		f = f.Sub(b, f) // b - d*(a/c)
		f = f.Mul(f, t) // (b - d*(a/c))*t
	}
	return e, f
}

// TODO func (c Complex) Pow() (apl.Value, bool)
// r := math.Exp(real(x))
// s, c := math.Sincos(imag(x))
// return complex(r*c, r*s)

// TODO (c Complex) Pow2(R apl.Value) (apl.Value, bool)

// TODO (c Complex) Log, Log2

func (c Complex) Abs() (apl.Value, bool) {
	// This is a downtype. The tower needs to include Float.
	return Float{c.cpy().abs()}, true
}

func (c Complex) abs() *big.Float {
	// en.wikipedia.org/wiki/Hypot
	x, y := c.re, c.im
	if x.Sign() < 0 {
		x = x.Neg(x) // y = abs(y)
	}
	if y.Sign() < 0 {
		y = y.Neg(y) // y = abs(y)
	}
	if y.Cmp(x) > 0 {
		x, y = y, x
	}
	if x.Sign() == 0 {
		return y
	}
	y = y.Quo(y, x) // t = y / x
	y = y.Mul(y, y) // t*t
	one := new(big.Float).Copy(y)
	one = one.SetInt64(1)
	y = y.Add(y, one) // 1+t*t
	y = y.Sqrt(y)     // sqrt(1+t*t)
	y = x.Mul(x, y)   // x * sqrt(1+t*t)
	return y
}

// TODO Floor, Ceil.

// TODO port sin.go asin.go from ivy.

// TODO Gcd
