package big

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/big/bigfloat"
	"github.com/ktye/iv/apl/numbers"
)

type Float struct {
	*big.Float
}

func (f Float) Copy() apl.Value {
	re := new(big.Float)
	re = re.Copy(f.Float)
	return Float{re}
}

func (f Float) String(af apl.Format) string {
	format, minus := getformat(af, f)
	if format == "" {
		if af.PP < 0 {
			format = "%v"
		} else if af.PP == 0 {
			format = "%.6G"
		} else {
			format = fmt.Sprintf("%%.%dG", af.PP)
		}
	}
	s := fmt.Sprintf(format, f.Float)
	if s == "-0" {
		s = "0"
	}
	if minus == false {
		s = strings.Replace(s, "-", "¯", -1)
	}
	return s
}

func ParseFloat(s string, prec uint) (apl.Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
	z, _, err := big.NewFloat(0).SetPrec(prec).Parse(s, 10)
	if err != nil {
		return nil, false
	}
	return Float{z}, true
}

func (f Float) ToIndex() (int, bool) {
	if f.IsInt() == false {
		return 0, false
	}
	i, _ := f.Float.Int64()
	n := int(i)
	if big.NewFloat(float64(n)).Cmp(f.Float) == 0 {
		return n, true
	}
	return 0, false
}

func (f Float) cpy() *big.Float {
	c := new(big.Float)
	return c.Copy(f.Float)
}

func (f Float) Equals(R apl.Value) (apl.Bool, bool) {
	return f.Float.Cmp(R.(Float).Float) == 0, true
}

func (f Float) Less(R apl.Value) (apl.Bool, bool) {
	return f.Float.Cmp(R.(Float).Float) < 0, true
}

func (f Float) Add() (apl.Value, bool) {
	return f, true
}
func (f Float) Add2(R apl.Value) (apl.Value, bool) {
	z := f.cpy()
	return Float{z.Add(z, R.(Float).Float)}, true
}

func (f Float) Sub() (apl.Value, bool) {
	z := f.cpy()
	return Float{z.Neg(f.Float)}, true
}
func (f Float) Sub2(R apl.Value) (apl.Value, bool) {
	z := f.cpy()
	return Float{z.Sub(z, R.(Float).Float)}, true
}

func (f Float) Mul() (apl.Value, bool) {
	return apl.Int(f.Float.Sign()), true
}
func (f Float) Mul2(R apl.Value) (apl.Value, bool) {
	z := f.cpy()
	return Float{z.Mul(z, R.(Float).Float)}, true
}

func (f Float) Div() (apl.Value, bool) {
	return Float{big.NewFloat(1)}.Div2(f)
}
func (f Float) Div2(R apl.Value) (apl.Value, bool) {
	if f.Float.IsInf() {
		return numbers.Inf, true
	}
	if R.(Float).Float.IsInf() {
		return numbers.NaN, true
	}
	lz := f.Float.Sign() == 0
	rz := R.(Float).Float.Sign() == 0
	if lz && rz {
		return numbers.NaN, true
	} else if lz {
		z := f.cpy().SetInt64(0)
		return Float{z}, true
	} else if rz {
		return numbers.Inf, true
	}
	return Float{f.cpy().Quo(f.Float, R.(Float).Float)}, true
}

func (f Float) Pow() (apl.Value, bool) {
	z := bigfloat.Exp(f.Float)
	if z.IsInf() {
		return numbers.Inf, true
	}
	return Float{z}, true
}
func (f Float) Pow2(R apl.Value) (apl.Value, bool) {
	if f.Float.Cmp(f.Float) < 0 {
		return nil, false
	}
	z := bigfloat.Pow(f.Float, R.(Float).Float)
	if z.IsInf() {
		return numbers.Inf, true
	}
	return Float{z}, true
}

func (f Float) Log() (apl.Value, bool) {
	if f.Float.Sign() < 0 {
		return nil, false
	}
	return Float{bigfloat.Log(f.Float)}, true
}
func (f Float) Log2(R apl.Value) (apl.Value, bool) {
	if f.Float.Sign() < 0 {
		return nil, false
	}
	r := R.(Float).Float
	if r.Sign() < 0 {
		return nil, false
	}
	logl := bigfloat.Log(f.Float)
	logr := bigfloat.Log(r)
	return Float{logr.Quo(logr, logl)}, true
}

func (f Float) Abs() (apl.Value, bool) {
	if f.Float.Sign() < 0 {
		return f.Sub()
	}
	return f, true
}

func (f Float) Ceil() (apl.Value, bool) {
	z, _ := f.Float.Float64()
	return Float{f.cpy().SetFloat64(math.Ceil(z))}, true
}
func (f Float) Floor() (apl.Value, bool) {
	z, _ := f.Float.Float64()
	return Float{f.cpy().SetFloat64(math.Floor(z))}, true
}

// TODO Trig

// TODO Gcd
