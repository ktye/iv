package apl

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// Int is the Integer type. It is used for numbers an indexes.
type Int int

func (i Int) ToIndex() (int, bool) {
	return int(i), true
}

// String formats an integer as a string.
// The format string is passed to fmt and - is replaced by ¯,
// except if the first rune is -.
func (i Int) String(f Format) string {
	format := f.Fmt[reflect.TypeOf(i)]
	minus := false
	if len(format) > 1 && format[0] == '-' {
		minus = true
		format = format[1:]
	}
	if format == "" {
		switch f.PP {
		case -8:
			format = "0%o"
		case -16:
			format = "0x%X"
		default:
			format = "%v"
		}
	}
	s := fmt.Sprintf(format, i)
	if minus == false {
		s = strings.Replace(s, "-", "¯", 1)
	}
	return s
}

func (i Int) Copy() Value {
	return i
}

func (i Int) Export() reflect.Value {
	return reflect.ValueOf(int(i))
}

// ParseInt parses an integer. It replaces ¯ with -, then uses ParseInt.
// Decimal, Octal (0x..) and Hexadecimal (0..) formats are supported.
func ParseInt(s string) (Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
	if n, err := strconv.ParseInt(s, 0, 64); err == nil {
		return Int(n), true
	}
	return Int(0), false
}

func (i Int) Less(R Value) (Bool, bool) {
	return Bool(i < R.(Int)), true
}

func (i Int) Add() (Value, bool) {
	return i, true
}
func (i Int) Add2(R Value) (Value, bool) {
	return i + R.(Int), true
}

func (i Int) Sub() (Value, bool) {
	return -i, true
}
func (i Int) Sub2(R Value) (Value, bool) {
	return i - R.(Int), true
}

func (i Int) Mul() (Value, bool) {
	if i > 0 {
		return Int(1), true
	} else if i < 0 {
		return Int(-1), true
	}
	return Int(0), true
}
func (i Int) Mul2(R Value) (Value, bool) {
	return i * R.(Int), true
}

func (i Int) Div() (Value, bool) {
	if i == 1 {
		return Int(1), true
	} else if i == -1 {
		return Int(-1), true
	}
	return nil, false
}
func (a Int) Div2(b Value) (Value, bool) {
	n := int(b.(Int))
	if n == 0 {
		return nil, false
	}
	r := int(a) / n
	if r*n == int(a) {
		return Int(r), true
	}
	return nil, false
}

func (i Int) Pow() (Value, bool) {
	if i == 0 {
		return Int(1), true
	}
	return nil, false
}
func (i Int) Pow2(R Value) (Value, bool) {
	return nil, false
}

func (i Int) Log() (Value, bool) {
	return nil, false
}
func (i Int) Log2() (Value, bool) {
	return nil, false
}

func (i Int) Abs() (Value, bool) {
	if i < 0 {
		return -i, true
	}
	return i, true
}

func (i Int) Ceil() (Value, bool) {
	return i, true
}
func (i Int) Floor() (Value, bool) {
	return i, true
}

func (i Int) Gamma() (Value, bool) {
	// 20 is the limit for int64.
	if i < 0 || i > 20 {
		return nil, false
	} else if i == 0 {
		return Int(1), true
	}
	n := 1
	for k := 1; k <= int(i); k++ {
		n *= k
	}
	return Int(n), true
}
func (L Int) Gamma2(r Value) (Value, bool) {
	m1exp := func(n Int) Int {
		if n%2 == 0 {
			return 1
		}
		return -1
	}
	R := r.(Int)
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
		return rg.(Int) / (lg.(Int) * rlg.(Int)), true
	} else if L >= 0 && R >= 0 && R-L < 0 {
		return Int(0), true
	} else if L >= 0 && R < 0 && R-L < 0 {
		v, ok := L.Gamma2(L - (1 + R))
		if ok == false {
			return nil, false
		}
		return m1exp(L) * v.(Int), true
	} else if L < 0 && R >= 0 && R-L >= 0 {
		return Int(0), true
	} else if L < 0 && R < 0 && R-L >= 0 {
		al1 := 1 + L
		if al1 < 0 {
			al1 = -al1
		}
		v, ok := (-(R + 1)).Gamma2(al1)
		if ok == false {
			return nil, false
		}
		return m1exp(R-L) * v.(Int), true
	} else if L < 0 && R < 0 && R-L < 0 {
		return Int(0), true
	}
	return nil, false
}

func (L Int) Gcd(R Value) (Value, bool) {
	l := big.NewInt(int64(L))
	r := big.NewInt(int64(R.(Int)))
	return Int(big.NewInt(0).GCD(nil, nil, l, r).Int64()), true
}
