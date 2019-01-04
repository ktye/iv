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
			return dict2table(a, &d)
		}
		return &d, nil
	}
}

// tableBoth applies elementary function to two tables or two object.
func tableBoth(symbol string, fn func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	scalar := arith2(symbol, fn)
	array := array2(symbol, fn)
	return func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		istable := false
		var l, r apl.Object
		if t, ok := L.(apl.Table); ok {
			istable = true
			l = t.Dict
		} else {
			l = L.(apl.Object)
		}

		if t, ok := R.(apl.Table); istable && ok == false {
			return nil, fmt.Errorf("both tables and objects cannot be mixed")
		} else if ok == true {
			r = t.Dict
		} else {
			r = R.(apl.Object)
		}

		zero := func(v apl.Value) apl.Value {
			if u, ok := v.(apl.Uniform); ok {
				return u.Zero()
			}
			return apl.Index(0)
		}

		keys := l.Keys()
		d := apl.Dict{K: make([]apl.Value, len(keys)), M: make(map[apl.Value]apl.Value)}
		toArrays := arrays{}

		// Loop over all keys in the left dict.
		var err error
		var v apl.Value
		for i, k := range keys {
			lv := l.At(a, k)
			rv := r.At(a, k)
			if rv == nil {
				rv = zero(lv)
			}
			if la, ra, ok := toArrays.To(a, lv, rv); ok {
				v, err = array(a, la, ra)
			} else {
				v, err = scalar(a, la, ra)
			}
			if err != nil {
				return nil, err
			}
			d.K[i] = k
			d.M[k] = v
		}

		// Loop over all keys in the right dict, only use keys that are not present in left.
		keys = r.Keys()
		for _, k := range keys {
			lv := l.At(a, k)
			if lv != nil {
				continue
			}
			rv := r.At(a, k)
			lv = zero(rv)
			if la, ra, ok := toArrays.To(a, lv, rv); ok {
				v, err = array(a, la, ra)
			} else {
				v, err = scalar(a, la, ra)
			}
			if err != nil {
				return nil, err
			}
			d.K = append(d.K, k)
			d.M[k] = v
		}

		if istable {
			return dict2table(a, &d)
		}
		return &d, nil
	}
}

// tableAny applies elementary function to a combination of a scalar or array
// with a table or object.
func tableAny(symbol string, fn func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	scalar := arith2(symbol, fn)
	array := array2(symbol, fn)
	return func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		istable := false
		leftarray := false
		var o apl.Object
		if t, ok := L.(apl.Table); ok {
			istable = true
			o = t.Dict
		} else if d, ok := L.(apl.Object); ok {
			o = d
		} else if t, ok := R.(apl.Table); ok {
			istable = true
			leftarray = true
			o = t.Dict
		} else {
			leftarray = true
			o = R.(apl.Object)
		}

		keys := o.Keys()
		d := apl.Dict{K: make([]apl.Value, len(keys)), M: make(map[apl.Value]apl.Value)}
		toArrays := arrays{}

		var err error
		var v, r, l apl.Value
		for i, k := range keys {
			l = o.At(a, k) // TODO: copy?
			r = R
			if leftarray {
				l, r = L, l
			}
			if la, ra, ok := toArrays.To(a, l, r); ok {
				v, err = array(a, la, ra)
			} else {
				v, err = scalar(a, l, r)
			}
			if err != nil {
				return nil, err
			}
			d.K[i] = k
			d.M[k] = v
		}

		if istable {
			return dict2table(a, &d)
		}
		return &d, nil
	}
}

// isTableCat tests if one of the arguments is a table or object.
// It may be hidden in an axis on the right.
func isTableCat(a *apl.Apl, L, R apl.Value) (apl.Value, apl.Value, bool, bool) {
	first := false
	r := R
	if ax, ok := R.(apl.Axis); ok {
		r = ax.R
		first = true // we assume it's non-default.
	}
	if l, ok := L.(apl.Table); ok {
		return l, r, first, true
	} else if l, ok := L.(apl.Object); ok {
		return l, r, first, true
	} else if rr, ok := r.(apl.Table); ok {
		return l, rr, first, true
	} else if rr, ok := r.(apl.Object); ok {
		return l, rr, first, true
	}
	return L, R, false, false
}

// catenateTables catenates dicts or tables.
// At least one argument is a dict or table, the other may be a scalar or an array.
func catenateTables(a *apl.Apl, L, R apl.Value, first bool) (apl.Value, error) {
	if l, ok := L.(apl.Table); ok {
		if r, ok := R.(apl.Table); ok {
			return catenateTwoTables(a, l, r, first)
		} else if _, ok = R.(apl.Object); ok {
			return nil, fmt.Errorf("catenate: cannot mix object and table")
		}
	} else if l, ok := L.(apl.Object); ok {
		if r, ok := R.(apl.Object); ok {
			return catenateTwoTables(a, l, r, first)
		} else if _, ok = R.(apl.Table); ok {
			return nil, fmt.Errorf("catenate: cannot mix object and table")
		}
	}
	return nil, fmt.Errorf("TODO: cat tables")
}

// catenateTwoTables catenates dicts or tables.
// Both arguments are the same type.
func catenateTwoTables(a *apl.Apl, L, R apl.Value, first bool) (apl.Value, error) {
	_, istable := L.(apl.Table)
	var l, r apl.Object
	if istable {
		l = L.(apl.Table).Dict
		r = R.(apl.Table).Dict
	} else {
		l = L.(apl.Object)
		r = R.(apl.Object)
	}

	if istable {
		first = !first
	}

	keys := l.Keys()
	d := apl.Dict{K: make([]apl.Value, len(keys)), M: make(map[apl.Value]apl.Value)}
	if first {
		for i, k := range keys {
			d.K[i] = k // TODO: Copy?
			d.M[k] = l.At(a, k)
		}
		for _, k := range r.Keys() {
			if _, ok := d.M[k]; !ok {
				d.K = append(d.K, k)
			}
			d.M[k] = r.At(a, k)
		}
	} else {
		rkeys := r.Keys()
		if istable && len(keys) != len(rkeys) {
			return nil, fmt.Errorf("catenate table on first axis: tables have different number of columns")
		}
		for i, k := range keys {
			d.K[i] = k // TODO: Copy?
			lv := l.At(a, k)
			rv := r.At(a, k)
			if rv == nil {
				if istable {
					return nil, fmt.Errorf("catenate table on first axis: tables have different columns")
				}
				d.M[k] = lv
			} else {
				v, err := catenate(a, lv, rv)
				if err != nil {
					return nil, err
				}
				d.M[k] = v
			}
		}
		for _, k := range rkeys {
			if d.M[k] == nil {
				if istable {
					return nil, fmt.Errorf("catenate table on first axis: tables have different columns")
				}
				d.K = append(d.K, k)
				d.M[k] = r.At(a, k)
			}
		}
	}

	if istable {
		return dict2table(a, &d)
	}
	return &d, nil
}

func dict2table(a *apl.Apl, d *apl.Dict) (apl.Table, error) {
	rows := 0
	if len(d.K) > 0 {
		v := d.At(a, d.K[0])
		if ar, ok := v.(apl.Array); ok {
			shape := ar.Shape()
			if len(shape) == 1 {
				rows = shape[0]
			}
		}
	}
	return apl.Table{Dict: d, Rows: rows}, nil
}
