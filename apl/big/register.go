package big

import (
	"math/big"
	"reflect"

	"github.com/ktye/iv/apl"
)

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
	}
	if err := a.SetTower(t); err != nil {
		panic(err)
	}
}

// SetPreciseTower sets the numerical tower to Float with the given precision.
func SetPreciseTower(a *apl.Apl, prec uint) {
	m := make(map[reflect.Type]apl.Numeric)
	m[reflect.TypeOf(Float{})] = apl.Numeric{
		Class:  0,
		Parse:  func(s string) (apl.Number, bool) { return ParseFloat(s, prec) },
		Uptype: func(n apl.Number) (apl.Number, bool) { return n, false },
	}
	t := apl.Tower{
		Numbers: m,
		FromIndex: func(n int) apl.Number {
			return Float{big.NewFloat(float64(n)).SetPrec(prec)}
		},
	}
	if err := a.SetTower(t); err != nil {
		panic(err)
	}
}
