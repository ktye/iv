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
	register(operator{
		symbol:  "¨",
		Domain:  MonadicOp(IsPrimitive("<")),
		doc:     "channel each",
		derived: channelEach,
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
	if lst, ok := R.(apl.List); ok {
		return eachList(a, lst, f)
	}
	if c, ok := R.(apl.Channel); ok {
		return eachChannel(a, nil, c, f)
	}

	ar, ok := R.(apl.Array)
	if ok {
		if ar.Size() == 0 {
			// TODO: Fill function of LO should be applied
			// with the prototype of R.
			// The result has the same shape as R.
			return apl.EmptyArray{}, nil
		}
	} else {
		// Apply f to scalar R.
		return f.Call(a, nil, R)
	}

	res := apl.NewMixed(apl.CopyShape(ar))
	for i := range res.Values {
		v, err := f.Call(a, nil, ar.At(i))
		if err != nil {
			return nil, err
		}
		if _, ok := v.(apl.Array); ok {
			return nil, fmt.Errorf("each: result must be a scalar")
		}
		res.Values[i] = v.Copy()
	}
	return a.UnifyArray(res), nil
}

func eachList(a *apl.Apl, l apl.List, f apl.Function) (apl.Value, error) {
	res := make(apl.List, len(l))
	for i := range res {
		v, err := f.Call(a, nil, l[i])
		if err != nil {
			return nil, err
		}
		res[i] = v.Copy()
	}
	return res, nil
}

// EachChannel returns a channel and applies the function f to each value in the input channel.
// The result is written to the output channel.
// If f returns an EmptyArray, no output value is written.
// This can be used as a filter. Empty strings however are written.
func eachChannel(a *apl.Apl, L apl.Value, r apl.Channel, f apl.Function) (apl.Value, error) {
	return r.Apply(a, f, L, false), nil
}

// ChannelEach sends each value in R over a channel.
func channelEach(a *apl.Apl, _, _ apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		if L != nil {
			return nil, fmt.Errorf("channel each cannot be called dyadically")
		}
		var all []apl.Value
		c := apl.NewChannel()
		ar, ok := R.(apl.Array)
		if ok == false {
			all = []apl.Value{R.Copy()}
		} else {
			all = make([]apl.Value, ar.Size())
			for i := range all {
				all[i] = ar.At(i).Copy()
			}
		}
		go c.SendAll(all)
		return c, nil
	}
	return function(derived)
}

func each2(a *apl.Apl, L, R apl.Value, f apl.Function) (apl.Value, error) {
	if c, ok := R.(apl.Channel); ok {
		return eachChannel(a, L, c, f)
	}

	_, okl := L.(apl.List)
	_, okr := R.(apl.List)
	if okl || okr {
		return eachList2(a, L, R, f)
	}

	ar, rok := R.(apl.Array)
	al, lok := L.(apl.Array)
	var rs, ls []int

	if rok == false && lok == false {
		return f.Call(a, L, R)
	}
	if rok == true && ar.Size() == 0 {
		return apl.EmptyArray{}, nil // TODO fill function
	}
	if lok == true && al.Size() == 0 {
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

	res := apl.NewMixed(shape)
	for i := range res.Values {
		if rok == true {
			rv = ar.At(i)
		}
		if lok == true {
			lv = al.At(i)
		}
		v, err := f.Call(a, lv, rv)
		if err != nil {
			return nil, err
		}
		if _, ok := v.(apl.Array); ok {
			return nil, fmt.Errorf("each: result must be a scalar")
		}
		res.Values[i] = v.Copy()
	}
	return a.UnifyArray(res), nil
}

func eachList2(a *apl.Apl, L, R apl.Value, f apl.Function) (apl.Value, error) {
	l, lok := L.(apl.List)
	r, rok := R.(apl.List)
	size := 0
	if lok {
		size = len(l)
	}
	if rok {
		if len(r) != size {
			return nil, fmt.Errorf("each list: different list sizes")
		}
	}

	res := make(apl.List, size)
	for i := range res {
		lv := L
		rv := R
		if lok {
			lv = l[i]
		}
		if rok {
			rv = r[i]
		}
		v, err := f.Call(a, lv, rv)
		if err != nil {
			return nil, err
		}
		res[i] = v.Copy()
	}
	return res, nil
}
