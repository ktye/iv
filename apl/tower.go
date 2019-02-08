package apl

import (
	"fmt"
	"reflect"
)

type Tower struct {
	Numbers map[reflect.Type]*Numeric
	Import  func(v Number) Number       // Import Bool or Int
	Uniform func([]Value) (Value, bool) // Values must already be uniform.
	idx     []*Numeric
}

// SetTower sets the numerical tower.
// The default tower can be set by calling numbers.Register(a).
func (a *Apl) SetTower(t Tower) error {
	t.idx = make([]*Numeric, len(t.Numbers))
	for i := 0; i < len(t.Numbers); i++ {
		for _, n := range t.Numbers {
			if n.Class == i {
				m := n
				t.idx[i] = m
			}
		}
	}
	for c, p := range t.idx {
		if p == nil {
			return fmt.Errorf("not a valid tower: class %d is missing", c)
		}
	}
	a.Tower = t
	return nil
}

// Parse tries to parse a string as a Number, starting with the lowest number type.
func (t Tower) Parse(s string) (NumExpr, error) {

	// Bool and Index can be parsed directly.
	switch s {
	case "1b":
		return NumExpr{Bool(true)}, nil
	case "0b":
		return NumExpr{Bool(false)}, nil
	default:
		if n, ok := ParseInt(s); ok {
			return NumExpr{n}, nil
		}
	}

	if t.Numbers == nil || len(t.idx) != len(t.Numbers) {
		return NumExpr{}, fmt.Errorf("numeric tower is not initialized")
	}
	for _, num := range t.idx {
		if num.Parse == nil {
			continue
		}
		if n, ok := num.Parse(s); ok {
			return NumExpr{n}, nil
		}
	}
	return NumExpr{}, fmt.Errorf("cannot parse number: %s", s)
}

// SameType returns the two numbers with the same type.
// It uptypes the lower number type.
func (t Tower) SameType(a, b Number) (Number, Number, error) {
	at := reflect.TypeOf(a)
	bt := reflect.TypeOf(b)
	if at == bt {
		return a, b, nil
	}

	// Handle Bool and Index.
	if ab, ok := a.(Bool); ok {
		if _, ok := b.(Index); ok {
			return bool2int(ab), b, nil
		}
		a = t.Import(a)
		at = reflect.TypeOf(a)
	}
	if bb, ok := b.(Bool); ok {
		if _, ok := a.(Index); ok {
			return a, bool2int(bb), nil
		}
		b = t.Import(b)
		bt = reflect.TypeOf(b)
	}
	if _, ok := a.(Index); ok {
		a = t.Import(a)
		at = reflect.TypeOf(a)
	}
	if _, ok := b.(Index); ok {
		b = t.Import(b)
		bt = reflect.TypeOf(b)
	}

	na, ok := t.Numbers[at]
	if ok == false {
		return nil, nil, fmt.Errorf("numeric tower: unknown number type %T", a)
	}
	nb, ok := t.Numbers[bt]
	if ok == false {
		return nil, nil, fmt.Errorf("numeric tower: unknown number type %T", b)
	}
	for i := na.Class; i < nb.Class; i++ {
		a, ok = na.Uptype(a)
		if ok == false {
			// Uptype should return the original number if it fails.
			return nil, nil, fmt.Errorf("cannot uptype %T", a)
		}
		na = t.idx[i+1]
	}
	for i := nb.Class; i < na.Class; i++ {
		b, ok = nb.Uptype(b)
		if ok == false {
			// Uptype should return the original number if it fails.
			return nil, nil, fmt.Errorf("cannot uptype %T", b)
		}
		nb = t.idx[i+1]
	}
	return a, b, nil
}

func bool2int(b Bool) Index {
	if b {
		return Index(1)
	}
	return Index(0)
}

func (a *Apl) IsZero(n Number) bool {
	b, ok := a.Tower.ToBool(n)
	if ok == false {
		return false
	}
	return b == false
}

func (t *Tower) ToBool(n Number) (Bool, bool) {
	if idx, ok := n.ToIndex(); ok == false {
		return false, false
	} else if idx < 0 || idx > 1 {
		return false, false
	} else if idx == 0 {
		return false, true
	} else {
		return true, true
	}
}

func (t *Tower) ToNumeric(v Number) *Numeric {
	if _, ok := v.(Bool); ok {
		return &Numeric{
			Class: -2,
			Uptype: func(n Number) (Number, bool) {
				if b := n.(Bool); b {
					return Index(1), true
				}
				return Index(0), true
			},
		}
	}
	if _, ok := v.(Index); ok {
		return &Numeric{
			Class: -1,
			Uptype: func(n Number) (Number, bool) {
				return t.Import(n), true
			},
		}
	}
	if num, ok := t.Numbers[reflect.TypeOf(v)]; ok {
		return num
	}
	return nil
}

// Numeric is a member of the tower.
// Uptype converts a Number to the next higher class.
type Numeric struct {
	Class  int
	Parse  func(string) (Number, bool)
	Uptype func(Number) (Number, bool)
}

// Number is the interface that a numeric type must implement.
// It's a scalar numeric value that can be part of the current tower.
type Number interface {
	String(*Apl) string
	ToIndex() (int, bool)
}

// NumExpr wraps a Number to be used as an expression by the parser.
type NumExpr struct {
	Number
}

func (num NumExpr) Eval(a *Apl) (Value, error) {
	return num.Number, nil
}
