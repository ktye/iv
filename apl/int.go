package apl

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// Index is the Integer type. It is used for numbers an indexes.
type Index int

func (i Index) ToIndex() (int, bool) {
	return int(i), true
}

// String formats an integer as a string.
// The format string is passed to fmt and - is replaced by ¯,
// except if the first rune is -.
func (i Index) String(a *Apl) string {
	format := a.Fmt[reflect.TypeOf(i)]
	minus := false
	if len(format) > 1 && format[0] == '-' {
		minus = true
		format = format[1:]
	}
	if format == "" {
		format = "%v"
	}
	s := fmt.Sprintf(format, i)
	if minus == false {
		s = strings.Replace(s, "-", "¯", 1)
	}
	return s
}

// ParseInt parses an integer. It replaces ¯ with -, then uses Atoi.
func ParseInt(s string) (Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
	if n, err := strconv.Atoi(s); err == nil {
		return Index(n), true
	}
	return Index(0), false
}

func (i Index) Less(R Value) (Bool, bool) {
	return Bool(i < R.(Index)), true
}

func (i Index) Add() (Value, bool) {
	return i, true
}
func (i Index) Add2(R Value) (Value, bool) {
	return i + R.(Index), true
}

func (i Index) Sub() (Value, bool) {
	return -i, true
}
func (i Index) Sub2(R Value) (Value, bool) {
	return i - R.(Index), true
}

func (i Index) Mul() (Value, bool) {
	if i > 0 {
		return Index(1), true
	} else if i < 0 {
		return Index(-1), true
	}
	return Index(0), true
}
func (i Index) Mul2(R Value) (Value, bool) {
	return i * R.(Index), true
}

func (i Index) Div() (Value, bool) {
	if i == 1 {
		return Index(1), true
	} else if i == -1 {
		return Index(-1), true
	}
	return nil, false
}
func (a Index) Div2(b Value) (Value, bool) {
	n := int(b.(Index))
	if n == 0 {
		return nil, false
	}
	r := int(a) / n
	if r*n == int(a) {
		return Index(r), true
	}
	return nil, false
}

func (i Index) Pow() (Value, bool) {
	if i == 0 {
		return Index(1), true
	}
	return nil, false
}
func (i Index) Pow2(R Value) (Value, bool) {
	return nil, false
}

func (i Index) Log() (Value, bool) {
	return nil, false
}
func (i Index) Log2() (Value, bool) {
	return nil, false
}

func (i Index) Abs() (Value, bool) {
	if i < 0 {
		return -i, true
	}
	return i, true
}

func (i Index) Ceil() (Value, bool) {
	return i, true
}
func (i Index) Floor() (Value, bool) {
	return i, true
}

func (i Index) Gamma() (Value, bool) {
	// 20 is the limit for int64.
	if i < 0 || i > 20 {
		return nil, false
	} else if i == 0 {
		return Index(1), true
	}
	n := 1
	for k := 1; k <= int(i); k++ {
		n *= k
	}
	return Index(n), true
}
func (L Index) Gamma2(r Value) (Value, bool) {
	m1exp := func(n Index) Index {
		if n%2 == 0 {
			return 1
		}
		return -1
	}
	R := r.(Index)
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
		return rg.(Index) / (lg.(Index) * rlg.(Index)), true
	} else if L >= 0 && R >= 0 && R-L < 0 {
		return Index(0), true
	} else if L >= 0 && R < 0 && R-L < 0 {
		v, ok := L.Gamma2(L - (1 + R))
		if ok == false {
			return nil, false
		}
		return m1exp(L) * v.(Index), true
	} else if L < 0 && R >= 0 && R-L >= 0 {
		return Index(0), true
	} else if L < 0 && R < 0 && R-L >= 0 {
		al1 := 1 + L
		if al1 < 0 {
			al1 = -al1
		}
		v, ok := (-(R + 1)).Gamma2(al1)
		if ok == false {
			return nil, false
		}
		return m1exp(R-L) * v.(Index), true
	} else if L < 0 && R < 0 && R-L < 0 {
		return Index(0), true
	}
	return nil, false
}

func (L Index) Gcd(R Value) (Value, bool) {
	l := big.NewInt(int64(L))
	r := big.NewInt(int64(R.(Index)))
	return Index(big.NewInt(0).GCD(nil, nil, l, r).Int64()), true
}
