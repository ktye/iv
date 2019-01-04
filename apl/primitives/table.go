package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

// table1 tries to apply the elementary function returned to each column of a table
// or each value in a dict.
// If the argument is an object, a dict is returned.
func table1(symbol string, fn func(*apl.Apl, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	scalar := arith1(symbol, fn)
	array := array1(symbol, fn)
	return func(a *apl.Apl, _ apl.Value, R apl.Value) (apl.Value, error) {
		istable := false
		var src apl.Object
		if t, ok := R.(apl.Table); ok {
			istable = true
			src = t.Dict
		} else {
			src = R.(apl.Object)
		}

		keys := src.Keys()
		d := apl.Dict{K: make([]apl.Value, len(keys)), M: make(map[apl.Value]apl.Value)}
		var err error
		for i, k := range keys {
			d.K[i] = k // TODO: copy?
			v := src.At(a, k)
			if v == nil {
				return nil, fmt.Errorf("missing value for key %s", k.String(a))
			}
			ar, ok := v.(apl.Array)
			if ok {
				v, err = array(a, nil, ar)
				if err != nil {
					return nil, err
				}
			} else {
				v, err = scalar(a, nil, v)
				if err != nil {
					return nil, err
				}
			}
			d.M[k] = v
		}

		if istable {
			rows := 0
			if len(keys) > 0 {
				v := d.At(a, keys[0])
				if ar, ok := v.(apl.Array); ok {
					shape := ar.Shape()
					if len(shape) == 1 {
						rows = shape[0]
					}
				}
			}
			return apl.Table{Dict: &d, Rows: rows}, nil
		}
		return &d, nil
	}
}
