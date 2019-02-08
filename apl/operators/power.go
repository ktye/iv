package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "⍣",
		Domain:  DyadicOp(Split(Function(nil), nil)),
		doc:     "power",
		derived: power,
	})
}

// TODO: should there be a limit? How to set/change it?
const powerlimit = 1000

func power(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		f := f.(apl.Function)
		gn, isgf := g.(apl.Function)
		to := ToIndex(nil)
		if isgf == false {
			// RO g is not a function but an integer.
			nv, ok := to.To(a, g)
			if ok == false {
				return nil, fmt.Errorf("power: non-function RO must be an integer: %T", g)
			}
			n := int(nv.(apl.Int))
			if n < 0 {
				return nil, fmt.Errorf("power: function inverse is not implemented")
			} else if n == 0 {
				return R, nil
			}
			var err error
			v := R
			for i := 0; i < n; i++ {
				v, err = f.Call(a, L, v)
				if err != nil {
					return nil, err
				}
			}
			return v, nil
		} else {
			// RO g is a function.
			var err error
			var fR, v apl.Value
			r := R
			m := 0
			for {
				if m > powerlimit {
					return nil, fmt.Errorf("power: recusion limit exceeded")
				}
				m++
				fR, err = f.Call(a, L, r)
				if err != nil {
					return nil, err
				}
				v, err = gn.Call(a, fR, r)
				if err != nil {
					return nil, err
				}
				nv, ok := to.To(a, v)
				if ok == false {
					return nil, fmt.Errorf("power: gY must be an integer: %T", v)
				}
				n := int(nv.(apl.Int))

				if n == 1 {
					return fR, nil
				}
				r = fR
			}
		}
	}
	return function(derived)
}
