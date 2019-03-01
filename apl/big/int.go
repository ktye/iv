package big

import (
	"math/big"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

var int0 = big.NewInt(0)
var int1 = big.NewInt(1)

// Int is a big integer of arbitrary size.
type Int struct {
	*big.Int
}

func (i Int) String(f apl.Format) string {
	// TODO formats
	s := i.Int.String()
	if s[0] == '-' {
		return "¯" + s[1:]
	}
	return s
}
func (i Int) Copy() apl.Value {
	r := big.NewInt(0)
	r = r.Set(i.Int)
	return Int{r}
}

func ParseInt(s string) (apl.Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
	i := new(big.Int)
	i, ok := i.SetString(s, 10)
	if ok == false {
		return nil, false
	}
	return Int{i}, true
}

func (i Int) ToIndex() (int, bool) {
	if i.Int.IsInt64() == false {
		return 0, false
	}
	n64 := i.Int.Int64()
	if n := int(n64); int64(n) == n64 {
		return n, true
	}
	return 0, false
}

func intToRat(i apl.Number) (apl.Number, bool) {
	r := new(big.Rat)
	r = r.SetInt(i.(Int).Int)
	return Rat{r}, true
}

func (i Int) Equals(R apl.Value) (apl.Bool, bool) {
	return i.Int.Cmp(R.(Int).Int) == 0, true
}

func (i Int) Less(R apl.Value) (apl.Bool, bool) {
	return i.Int.Cmp(R.(Int).Int) < 0, true
}

func (i Int) Add() (apl.Value, bool) {
	return i, true
}
func (i Int) Add2(R apl.Value) (apl.Value, bool) {
	z := new(big.Int)
	z = z.Add(i.Int, R.(Int).Int)
	return Int{z}, true
}

func (i Int) Sub() (apl.Value, bool) {
	z := new(big.Int)
	z = z.Neg(i.Int)
	return Int{z}, true
}
func (i Int) Sub2(R apl.Value) (apl.Value, bool) {
	z := new(big.Int)
	z = z.Sub(i.Int, R.(Int).Int)
	return Int{z}, true
}

func (i Int) Mul() (apl.Value, bool) {
	return Int{big.NewInt(int64(i.Int.Sign()))}, true
}
func (i Int) Mul2(R apl.Value) (apl.Value, bool) {
	z := new(big.Int)
	z = z.Mul(i.Int, R.(Int).Int)
	return Int{z}, true
}

func (i Int) Div() (apl.Value, bool) {
	if i.Int.Cmp(int0) == 0 {
		return numbers.Inf, true
	}
	b := big.NewInt(1)
	b = b.Div(b, i.Int)
	m := new(big.Int)
	m = m.Mul(b, i.Int)
	if m.Cmp(i.Int) == 0 {
		return Int{b}, true
	}
	return nil, false
}
func (i Int) Div2(R apl.Value) (apl.Value, bool) {
	L0 := false
	R0 := false
	if i.Int.Cmp(int0) == 0 {
		L0 = true
	}
	if R.(Int).Int.Cmp(int0) == 0 {
		R0 = true
	}
	if L0 && R0 {
		return numbers.NaN, true
	} else if L0 {
		return Int{big.NewInt(0)}, true
	} else if R0 {
		return numbers.Inf, true
	}
	b := new(big.Int)
	b = b.Div(i.Int, R.(Int).Int)
	m := new(big.Int)
	m = m.Mul(b, R.(Int).Int)
	if m.Cmp(i.Int) == 0 {
		return Int{b}, true
	}
	return nil, false
}

func (i Int) Pow() (apl.Value, bool) {
	if i.Int.Cmp(int0) == 0 {
		return Int{big.NewInt(1)}, true
	}
	return nil, false
}
func (i Int) Pow2(R apl.Value) (apl.Value, bool) {
	if c := R.(Int).Int.Cmp(int0); c == 0 {
		return Int{big.NewInt(1)}, true
	} else if c < 0 {
		return nil, false
	}
	z := new(big.Int)
	z = z.Exp(i.Int, R.(Int).Int, nil)
	return Int{z}, true
}

func (i Int) Abs() (apl.Value, bool) {
	if i.Int.Sign() < 0 {
		return i.Sub()
	}
	return i, true
}

func (i Int) Ceil() (apl.Value, bool) {
	return i, true
}
func (i Int) Floor() (apl.Value, bool) {
	return i, true
}

func (i Int) Gamma() (apl.Value, bool) {
	m, ok := i.ToIndex()
	if ok == false {
		return nil, false
	} else if m == 0 {
		return apl.Int(1), true
	} else if m < 0 || m > 300 { // where should be the limit?
		return nil, false
	}
	n := new(big.Int).SetInt64(1)
	t := new(big.Int)
	for k := 1; k <= m; k++ {
		t.SetInt64(int64(k))
		n = n.Mul(n, t)
	}
	return Int{n}, true
}

func (L Int) Gcd(R apl.Value) (apl.Value, bool) {
	return Int{big.NewInt(0).GCD(nil, nil, L.Int, R.(Int).Int)}, true
}
