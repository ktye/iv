package apl

import (
	"fmt"
	"reflect"
	"strconv"
)

type Tower struct {
	Numbers   map[reflect.Type]Numeric
	FromIndex func(int) Number
	Uniform   func([]Value) (Value, bool) // Values must already be uniform.
	idx       []*Numeric
}

// SetTower sets the numerical tower.
// The default tower can be set by calling numbers.Register(a).
func (a *Apl) SetTower(t Tower) error {
	t.idx = make([]*Numeric, len(t.Numbers))
	for i := 0; i < len(t.Numbers); i++ {
		for _, n := range t.Numbers {
			if n.Class == i {
				m := n
				t.idx[i] = &m
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
		if i, err := strconv.Atoi(s); err == nil {
			return NumExpr{Index(i)}, nil
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
	if _, ok := a.(Bool); ok {
		if _, ok := b.(Index); ok {
			return indexFromBool(a.(Bool)), b, nil
		}
		a = t.FromBool(a.(Bool))
		at = reflect.TypeOf(a)
	} else if _, ok := a.(Index); ok {
		a = t.FromIndex(int(a.(Index)))
		at = reflect.TypeOf(a)
	}
	if _, ok := b.(Bool); ok {
		if _, ok := a.(Index); ok {
			return a, indexFromBool(b.(Bool)), nil
		}
		b = t.FromBool(b.(Bool))
		bt = reflect.TypeOf(b)
	} else if _, ok := b.(Index); ok {
		b = t.FromIndex(int(b.(Index)))
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
	pa := &na
	pb := &nb
	for i := na.Class; i < nb.Class; i++ {
		a, ok = pa.Uptype(a)
		if ok == false {
			// Uptype should return the original number if it fails.
			return nil, nil, fmt.Errorf("cannot uptype %T", a)
		}
		pa = t.idx[i+1]
	}
	for i := nb.Class; i < na.Class; i++ {
		b, ok = pb.Uptype(b)
		if ok == false {
			// Uptype should return the original number if it fails.
			return nil, nil, fmt.Errorf("cannot uptype %T", b)
		}
		pb = t.idx[i+1]
	}
	return a, b, nil
}

func indexFromBool(b Bool) Index {
	if b {
		return Index(1)
	}
	return Index(0)
}

func (t *Tower) FromBool(b Bool) Number {
	if b {
		return t.FromIndex(1)
	}
	return t.FromIndex(0)
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
