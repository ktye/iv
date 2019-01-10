package big

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

func Register(a *apl.Apl, name string) {
	pkg := map[string]apl.Value{
		"set": settower{},
	}
	if name == "" {
		name = "big"
	}
	a.RegisterPackage(name, pkg)
}

// SetBigTower sets the numerical tower to Int->Rat.
func SetBigTower(a *apl.Apl) {
	m := make(map[reflect.Type]apl.Numeric)
	m[reflect.TypeOf(Int{})] = apl.Numeric{
		Class:  0,
		Parse:  ParseInt,
		Uptype: intToRat,
	}
	m[reflect.TypeOf(Rat{})] = apl.Numeric{
		Class:  1,
		Parse:  ParseRat,
		Uptype: func(n apl.Number) (apl.Number, bool) { return n, false },
	}
	t := apl.Tower{
		Numbers: m,
		FromIndex: func(n int) apl.Number {
			return Int{big.NewInt(int64(n))}
		},
		Uniform: func(v []apl.Value) (apl.Value, bool) { return nil, false },
	}
	if err := a.SetTower(t); err != nil {
		panic(err)
	}
}

// SetPreciseTower sets the numerical tower to Float->Complex with the given precision.
func SetPreciseTower(a *apl.Apl, prec uint) {
	m := make(map[reflect.Type]apl.Numeric)
	m[reflect.TypeOf(Float{})] = apl.Numeric{
		Class:  0,
		Parse:  func(s string) (apl.Number, bool) { return ParseFloat(s, prec) },
		Uptype: floatToComplex,
	}
	m[reflect.TypeOf(Complex{})] = apl.Numeric{
		Class:  1,
		Parse:  func(s string) (apl.Number, bool) { return ParseComplex(s, prec) },
		Uptype: func(n apl.Number) (apl.Number, bool) { return n, false },
	}
	t := apl.Tower{
		Numbers: m,
		FromIndex: func(n int) apl.Number {
			return Float{big.NewFloat(float64(n)).SetPrec(prec)}
		},
		Uniform: func(v []apl.Value) (apl.Value, bool) { return nil, false },
	}
	if err := a.SetTower(t); err != nil {
		panic(err)
	}
}

type settower struct{}

func (s settower) String(a *apl.Apl) string { return "set" }

func (_ settower) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	n, ok := R.(apl.Number)
	if ok == false {
		return nil, fmt.Errorf("set needs a number (0, 1, 256...)")
	}
	idx, ok := n.ToIndex()
	if ok == false {
		return nil, fmt.Errorf("set needs an integer")
	}
	if idx < 0 {
		return nil, fmt.Errorf("set: precision must be positive")
	} else if idx == 0 {
		numbers.Register(a)
	} else if idx == 1 {
		SetBigTower(a)
	} else {
		SetPreciseTower(a, uint(idx))
	}
	return R, nil
}

func getformat(a *apl.Apl, num apl.Number, def string) (string, bool) {
	if a == nil {
		return def, false
	}
	if n, ok := a.Tower.Numbers[reflect.TypeOf(num)]; ok == false {
		return def, false
	} else {
		f := n.Format
		if f == "" {
			return def, false
		}
		if f[0] == '-' {
			return f[1:], true
		}
		return f, false
	}
}
