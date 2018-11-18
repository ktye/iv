package big

import (
	"math/big"
	"reflect"

	"github.com/ktye/iv/apl"
)

// SetBigTower sets the numerical tower to Int->Rat.
func SetBigTower(a *apl.Apl) {
	m := make(map[reflect.Type]apl.Numeric)
	m[reflect.TypeOf(Int{int0})] = apl.Numeric{
		Class:  0,
		Parse:  ParseInt,
		Uptype: intToRat,
	}
	m[reflect.TypeOf(Rat{rat0})] = apl.Numeric{
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
