package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/xgo"
)

func init() {
	register(primitive{
		symbol: "#",
		doc:    "keys, methods",
		Domain: Monadic(nil),
		fn:     keys,
	})
	register(primitive{
		symbol: "#",
		doc:    "dict",
		Domain: Dyadic(nil),
		fn:     dict,
	})
}

// keys: R: object
func keys(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	methods := false
	if _, ok := R.(apl.Axis); ok {
		methods = true
		if r, _, err := splitAxis(a, R); err != nil {
			return nil, err
		} else {
			R = r
		}
	}
	obj, ok := R.(apl.Object)
	if ok == false {
		return nil, fmt.Errorf("keys: expected object: %T", R)
	}
	if methods {
		o, ok := obj.(xgo.Value)
		if ok == false {
			return nil, fmt.Errorf("methods: expected xgo.Value: %T", obj)
		}
		s := o.Methods()
		if s == nil {
			return apl.EmptyArray{}, nil
		}
		return apl.StringArray{Dims: []int{len(s)}, Strings: s}, nil
	} else {
		v := obj.Keys()
		return apl.MixedArray{
			Dims:   []int{len(v)},
			Values: v,
		}, nil
	}
}

func dict(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	al, ok := L.(apl.Array)

	if ok == false {
		return &apl.Dict{
			K: []apl.Value{L},
			M: map[apl.Value]apl.Value{
				L: R,
			},
		}, nil
	}

	ls := al.Shape()
	if len(ls) != 1 {
		return nil, fmt.Errorf("dict: left argument must be a vector")
	}

	ar, ok := R.(apl.Array)
	if ok == false {
		mr := apl.MixedArray{Dims: []int{ls[0]}, Values: make([]apl.Value, ls[0])}
		for i := range mr.Values {
			mr.Values[i] = R // TODO: copy?
		}
		ar = mr
	}
	rs := ar.Shape()
	if len(rs) != 1 || rs[0] != ls[0] {
		return nil, fmt.Errorf("dict: left and right arguments do not conform")
	}

	k := make([]apl.Value, al.Size())
	m := make(map[apl.Value]apl.Value)
	for i := 0; i < al.Size(); i++ {
		l := al.At(i)
		m[l] = ar.At(i)
		k[i] = l
	}
	return &apl.Dict{
		K: k,
		M: m,
	}, nil
}
