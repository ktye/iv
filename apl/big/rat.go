package big

import (
	"math"
	"math/big"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

var rat0 *big.Rat = big.NewRat(0, 1)
var rat1 *big.Rat = big.NewRat(1, 1)

type Rat struct {
	*big.Rat
}

func (r Rat) String(a *apl.Apl) string {
	s := strings.Replace(r.Rat.String(), "¯", "-", -1)
	if strings.HasSuffix(s, "/1") {
		return s[:len(s)-2]
	}
	return strings.Replace(s, "/", "r", 1)
}

func (r Rat) ToIndex() (int, bool) {
	if r.Rat.IsInt() == false {
		return 0, false
	}
	i := r.Rat.Num()
	return Int{i}.ToIndex()
}

func ParseRat(s string) (apl.Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
	s = strings.Replace(s, "r", "/", 1)
	r := new(big.Rat)
	r, ok := r.SetString(s)
	if ok == false {
		return nil, false
	}
	return Rat{r}, true
}

func (l Rat) Equals(R apl.Value) (apl.Bool, bool) {
	return l.Rat.Cmp(R.(Rat).Rat) == 0, true
}

func (l Rat) Less(R apl.Value) (apl.Bool, bool) {
	return l.Rat.Cmp(R.(Rat).Rat) < 0, true
}

func (l Rat) Add() (apl.Value, bool) {
	return l, true
}
func (l Rat) Add2(R apl.Value) (apl.Value, bool) {
	z := new(big.Rat)
	z = z.Add(l.Rat, R.(Rat).Rat)
	return Rat{z}, true
}

func (l Rat) Sub() (apl.Value, bool) {
	z := new(big.Rat)
	z = z.Neg(l.Rat)
	return Rat{z}, true
}
func (l Rat) Sub2(R apl.Value) (apl.Value, bool) {
	z := new(big.Rat)
	z = z.Sub(l.Rat, R.(Rat).Rat)
	return Rat{z}, true
}

func (l Rat) Mul() (apl.Value, bool) {
	return Int{big.NewInt(int64(l.Rat.Sign()))}, true
}
func (l Rat) Mul2(R apl.Value) (apl.Value, bool) {
	z := new(big.Rat)
	z = z.Mul(l.Rat, R.(Rat).Rat)
	return Rat{z}, true
}

func (r Rat) Div() (apl.Value, bool) {
	if r.Rat.Cmp(rat0) == 0 {
		return numbers.Inf, true
	}
	z := new(big.Rat)
	z = z.Inv(r.Rat)
	return Rat{z}, true
}
func (l Rat) Div2(R apl.Value) (apl.Value, bool) {
	L0 := false
	R0 := false
	if l.Rat.Cmp(rat0) == 0 {
		L0 = true
	}
	if R.(Rat).Rat.Cmp(rat0) == 0 {
		R0 = true
	}
	if L0 && R0 {
		return numbers.NaN, true
	} else if L0 {
		return Int{big.NewInt(0)}, true
	} else if R0 {
		return numbers.Inf, true
	}
	z := new(big.Rat)
	z = z.Quo(l.Rat, R.(Rat).Rat)
	return Rat{z}, true
}

func (r Rat) Pow() (apl.Value, bool) {
	if r.Rat.Cmp(rat0) == 0 {
		return Int{big.NewInt(1)}, true
	}
	return nil, false
}
func (l Rat) Pow2(R apl.Value) (apl.Value, bool) {
	neg := false
	r := R.(Rat).Rat
	if c := r.Cmp(rat0); c == 0 {
		return Int{big.NewInt(1)}, true
	} else if c < 0 {
		neg = true
		r = r.Neg(r)
	} else if l.Rat.Cmp(rat0) == 0 {
		return Int{big.NewInt(0)}, true
	}

	if r.IsInt() == false {
		return nil, false
	}

	e := r.Num()
	a := l.Rat.Num()
	b := l.Rat.Denom()
	ae := new(big.Int)
	be := new(big.Int)
	ae = ae.Exp(a, e, nil)
	be = be.Exp(b, e, nil)
	z := new(big.Rat)
	z.SetFrac(ae, be)
	if neg {
		z = z.Inv(z)
	}
	return Rat{z}, true
}

// Log, Log2 is not defined

func (r Rat) Abs() (apl.Value, bool) {
	if r.Rat.Sign() < 0 {
		return r.Sub()
	}
	return r, true
}

func (r Rat) Ceil() (apl.Value, bool) {
	f, _ := r.Rat.Float64()
	return Rat{new(big.Rat).SetFloat64(math.Ceil(f))}, true
}
func (r Rat) Floor() (apl.Value, bool) {
	f, _ := r.Rat.Float64()
	return Rat{new(big.Rat).SetFloat64(math.Floor(f))}, true
}
