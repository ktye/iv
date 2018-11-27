package numbers

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ktye/iv/apl"
)

type Integer int64

// String formats an integer as a string.
// The format string is passed to fmt and - is replaced by ¯,
// except if the first rune is -.
func (i Integer) String(a *apl.Apl) string {
	format, minus := getformat(a, i, "%d")
	s := fmt.Sprintf(format, int64(i))
	if minus == false {
		s = strings.Replace(s, "-", "¯", 1)
	}
	return s
}

// ParseInteger parses an integer. It replaces ¯ with -, then uses Atoi.
func ParseInteger(s string) (apl.Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
	if n, err := strconv.Atoi(s); err == nil {
		return Integer(n), true
	}
	return Integer(0), false
}

func (i Integer) ToIndex() (int, bool) {
	n := int(i)
	if i == Integer(n) {
		return n, true
	}
	return n, false
}

func intToFloat(i apl.Number) (apl.Number, bool) {
	return Float(i.(Integer)), true
}

func (i Integer) Less(R apl.Value) (apl.Bool, bool) {
	return apl.Bool(i < R.(Integer)), true
}

func (i Integer) Add() (apl.Value, bool) {
	return i, true
}
func (i Integer) Add2(R apl.Value) (apl.Value, bool) {
	return i + R.(Integer), true
}

func (i Integer) Sub() (apl.Value, bool) {
	return -i, true
}
func (i Integer) Sub2(R apl.Value) (apl.Value, bool) {
	return i - R.(Integer), true
}

func (i Integer) Mul() (apl.Value, bool) {
	if i > 0 {
		return Integer(1), true
	} else if i < 0 {
		return Integer(-1), true
	}
	return Integer(0), true
}
func (i Integer) Mul2(R apl.Value) (apl.Value, bool) {
	return i * R.(Integer), true
}

func (i Integer) Div() (apl.Value, bool) {
	if i == 1 {
		return Integer(1), true
	} else if i == -1 {
		return Integer(-1), true
	}
	return nil, false
}
func (a Integer) Div2(b apl.Value) (apl.Value, bool) {
	n := int64(b.(Integer))
	r := int64(a) / n
	if r*n == int64(a) {
		return Integer(r), true
	}
	return nil, false
}

func (i Integer) Pow() (apl.Value, bool) {
	if i == 0 {
		return Integer(1), true
	}
	return nil, false
}
func (i Integer) Pow2(R apl.Value) (apl.Value, bool) {
	return nil, false
}

func (i Integer) Log() (apl.Value, bool) {
	return nil, false
}
func (i Integer) Log2() (apl.Value, bool) {
	return nil, false
}

func (i Integer) Abs() (apl.Value, bool) {
	if i < 0 {
		return -i, true
	}
	return i, true
}

func (i Integer) Ceil() (apl.Value, bool) {
	return i, true
}
func (i Integer) Floor() (apl.Value, bool) {
	return i, true
}

func (i Integer) Gamma() (apl.Value, bool) {
	// 20 is the limit for int64.
	if i < 0 || i > 20 {
		return nil, false
	} else if i == 0 {
		return Integer(1), true
	}
	n := int64(1)
	for k := 1; k <= int(i); k++ {
		n *= int64(k)
	}
	return Integer(n), true
}
func (L Integer) Gamma2(r apl.Value) (apl.Value, bool) {
	m1exp := func(n Integer) Integer {
		if n%2 == 0 {
			return 1
		}
		return -1
	}
	R := r.(Integer)
	// This is the table from APL2 p 66
	if L >= 0 && R >= 0 && R-L >= 0 {
		lg, ok := L.Gamma()
		if ok == false {
			return nil, false
		}
		rg, ok := R.Gamma()
		if ok == false {
			return nil, false
		}
		rlg, ok := (R - L).Gamma()
		if ok == false {
			return nil, false
		}
		return rg.(Integer) / (lg.(Integer) * rlg.(Integer)), true
	} else if L >= 0 && R >= 0 && R-L < 0 {
		return Integer(0), true
	} else if L >= 0 && R < 0 && R-L < 0 {
		v, ok := L.Gamma2(L - (1 + R))
		if ok == false {
			return nil, false
		}
		return m1exp(L) * v.(Integer), true
	} else if L < 0 && R >= 0 && R-L >= 0 {
		return Integer(0), true
	} else if L < 0 && R < 0 && R-L >= 0 {
		al1 := 1 + L
		if al1 < 0 {
			al1 = -al1
		}
		v, ok := (-(R + 1)).Gamma2(al1)
		if ok == false {
			return nil, false
		}
		return m1exp(R-L) * v.(Integer), true
	} else if L < 0 && R < 0 && R-L < 0 {
		return Integer(0), true
	}
	return nil, false
}

func (L Integer) Gcd(R apl.Value) (apl.Value, bool) {
	l := big.NewInt(int64(L))
	r := big.NewInt(int64(R.(Integer)))
	return Integer(big.NewInt(0).GCD(nil, nil, l, r).Int64()), true
}
