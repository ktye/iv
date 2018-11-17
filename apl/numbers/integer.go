package numbers

import (
	"fmt"
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

const intmax = Integer(int(^uint(0) >> 1))
const intmin = Integer(-intmax - 1)

func (i Integer) ToIndex() (int, bool) {
	if i > intmax || i < intmin {
		return 0, false
	}
	return int(i), true
}

func intToFloat(i apl.Number) (apl.Number, bool) {
	return Float(i.(Integer)), true
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
