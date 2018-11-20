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
func (c Complex) Div2(R apl.Value) (apl.Value, bool) {
	// Algorithm for robust complex division as described in
	// Robert L. Smith: Algorithm 116: Complex division. Commun. ACM 5(8): 435 (1962).
	// adapted from: github.com/Quasilyte/go-complex-nums-emulation/blob/master/complex64.go
	r := R.(Complex)
	if r.re.Sign() == 0 && r.im.Sign() == 0 {
		return numbers.NaN, true
	}

	r1 := new(big.Float).Copy(c.re)
	i1 := new(big.Float).Copy(c.im)
	r2 := new(big.Float).Copy(r.re)
	i2 := new(big.Float).Copy(r.im)

	ar2 := new(big.Float).Copy(r2)
	ar2 = ar2.Abs(ar2) // ar1 = abs(r1)
	ai2 := new(big.Float).Copy(r2)
	ai2 = ai2.Abs(ai2) // ai2 = abs(i2)

	ratio := new(big.Float).Copy(r1)
	denom := new(big.Float).Copy(r1)
	if ar2.Cmp(ai2) > 0 { // abs(r2) >= abs(i2)
		ratio = ratio.Quo(i1, r2) // i1 / r2
		i2 = i2.Mul(ratio, i2)    // ratio * i2
		denom = denom.Add(r2, i2) // r2 + ratio*i2
		if denom.Sign() == 0 {
			return numbers.NaN, true
		}
		e := new(big.Float).Copy(i1)
		e = e.Mul(e, ratio) // i1*ratio
		e = e.Add(e, r1)    // r1 + i1*ratio
		e = e.Quo(e, denom) // (r1 + i1*ratio) / denom
		f := new(big.Float).Copy(r1)
		f = f.Mul(f, ratio) // r1*ratio
		f = f.Sub(i1, f)    // i1 - r1*ratio
		f = f.Quo(f, denom)
		return Complex{e, f}, true
	} else {
		ratio = ratio.Quo(r2, i2)    // r2 / i2
		denom = denom.Mul(ratio, r2) //  ratio*r2
		denom = denom.Add(denom, i1) // i2 + ratio*r2
		if denom.Sign() == 0 {
			return numbers.NaN, true
		}
		e := new(big.Float).Copy(r1)
		e = e.Mul(e, ratio) // r1 * ratio
		e = e.Add(e, i1)    // r1*ratio + i1
		e = e.Quo(e, denom) // (r1*ratio + i1) / denom
		f := new(big.Float).Copy(i1)
		f = f.Mul(f, ratio) // i1 * ratio
		f = f.Sub(f, r1)    // i1*ratio - r1
		f = f.Quo(f, denom) // (i1*ratio - r1) / denom
		return Complex{e, f}, true
	}
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
