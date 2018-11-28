package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "¨",
		Domain:  MonadicOp(Function(nil)),
		doc:     "each, map",
		derived: each,
	})
}

func each(a *apl.Apl, LO, RO apl.Value) apl.Function {
	f := LO.(apl.Function)
	derived := func(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
		if l == nil {
			return each1(a, r, f)
		}
		return each2(a, l, r, f)

	}
	return function(derived)
}

func each1(a *apl.Apl, R apl.Value, f apl.Function) (apl.Value, error) {
	ar, ok := R.(apl.Array)
	if ok {
		if apl.ArraySize(ar) == 0 {
			// TODO: Fill function of LO should be applied
			// with the prototype of R.
			// The result has the same shape as R.
			return apl.EmptyArray{}, nil
		}
	} else {
		// Apply f to scalar R.
		return f.Call(a, nil, R)
	}

	res := apl.GeneralArray{Dims: apl.CopyShape(ar)}
	res.Values = make([]apl.Value, apl.ArraySize(res))

	for i := range res.Values {
		r, err := ar.At(i)
		if err != nil {
			return nil, err
		}
		v, err := f.Call(a, nil, r)
		if err != nil {
			return nil, err
		}
		if _, ok := v.(apl.Array); ok {
			return nil, fmt.Errorf("each: result must be a scalar")
		}
		res.Values[i] = v
	}
	return res, nil
}

func each2(a *apl.Apl, L, R apl.Value, f apl.Function) (apl.Value, error) {
	ar, rok := R.(apl.Array)
	al, lok := L.(apl.Array)
	var rs, ls []int

	if rok == false && lok == false {
		return f.Call(a, L, R)
	}
	if rok == true && apl.ArraySize(ar) == 0 {
		return apl.EmptyArray{}, nil // TODO fill function
	}
	if lok == true && apl.ArraySize(al) == 0 {
		return apl.EmptyArray{}, nil // TODO fill function
	}

	if rok == true {
		rs = ar.Shape()
	}
	if lok == true {
		ls = al.Shape()
	}

	if rok == true && lok == true {
		if len(ls) != len(rs) {
			return nil, fmt.Errorf("each: ranks L and R are different")
		}
		for i := range ls {
			if ls[i] != rs[i] {
				return nil, fmt.Errorf("each: shapes of L and R must conform")
			}
		}
	}

	var shape []int
	var lv, rv apl.Value
	if rok == true {
		shape = apl.CopyShape(ar)
	} else {
		shape = apl.CopyShape(al)
		rv = R
	}
	if lok == false {
		lv = L
	}

	res := apl.GeneralArray{Dims: shape}
	res.Values = make([]apl.Value, apl.ArraySize(res))
	var err error
	for i := range res.Values {
		if rok == true {
			rv, err = ar.At(i)
			if err != nil {
				return nil, err
			}
		}
		if lok == true {
			lv, err = al.At(i)
			if err != nil {
				return nil, err
			}
		}
		v, err := f.Call(a, lv, rv)
		if err != nil {
			return nil, err
		}
		if _, ok := v.(apl.Array); ok {
			return nil, fmt.Errorf("each: result must be a scalar")
		}
		res.Values[i] = v
	}
	return res, nil
}
