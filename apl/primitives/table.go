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
			d.K[i] = k.Copy()
			v := src.At(k)
			if v == nil {
				return nil, fmt.Errorf("missing value for key %s", k.String(apl.Format{}))
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
			d.M[k.Copy()] = v.Copy()
		}

		if istable {
			return dict2table(&d)
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
			return apl.Int(0)
		}

		keys := l.Keys()
		d := apl.Dict{K: make([]apl.Value, len(keys)), M: make(map[apl.Value]apl.Value)}
		toArrays := arrays{}

		// Loop over all keys in the left dict.
		var err error
		var v apl.Value
		for i, k := range keys {
			lv := l.At(k)
			rv := r.At(k)
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
			d.K[i] = k.Copy()
			d.M[k.Copy()] = v.Copy()
		}

		// Loop over all keys in the right dict, only use keys that are not present in left.
		keys = r.Keys()
		for _, k := range keys {
			lv := l.At(k)
			if lv != nil {
				continue
			}
			rv := r.At(k)
			lv = zero(rv)
			if la, ra, ok := toArrays.To(a, lv, rv); ok {
				v, err = array(a, la, ra)
			} else {
				v, err = scalar(a, la, ra)
			}
			if err != nil {
				return nil, err
			}
			d.K = append(d.K, k.Copy())
			d.M[k.Copy()] = v.Copy()
		}

		if istable {
			return dict2table(&d)
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
			l = o.At(k)
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
			d.K[i] = k.Copy()
			d.M[k.Copy()] = v.Copy()
		}

		if istable {
			return dict2table(&d)
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
		return L, rr, first, true
	} else if rr, ok := r.(apl.Object); ok {
		return L, rr, first, true
	}
	return L, R, false, false
}

// catenateTables catenates dicts or tables.
// At least one argument is a dict or table, the other may be a scalar or an array.
func catenateTables(a *apl.Apl, L, R apl.Value, first bool) (apl.Value, error) {
	istable := false
	leftarray := false
	var o apl.Object
	if l, ok := L.(apl.Table); ok {
		istable = true
		o = l.Dict
		if r, ok := R.(apl.Table); ok {
			return catenateTwoTables(a, l, r, first)
		} else if _, ok = R.(apl.Object); ok {
			return nil, fmt.Errorf("catenate: cannot mix object and table")
		}
	} else if l, ok := L.(apl.Object); ok {
		o = l
		if r, ok := R.(apl.Object); ok {
			return catenateTwoTables(a, l, r, first)
		} else if _, ok = R.(apl.Table); ok {
			return nil, fmt.Errorf("catenate: cannot mix object and table")
		}
	}

	// Catenate a table or object with a normal array or scalar.
	// The axis value is ignored.

	if r, ok := R.(apl.Table); ok {
		istable = true
		leftarray = true
		o = r.Dict
	} else if r, ok := R.(apl.Object); ok {
		leftarray = true
		o = r
	}

	var lv, rv apl.Value
	if leftarray {
		lv = L
	} else {
		rv = R
	}

	keys := o.Keys()
	d := apl.Dict{K: make([]apl.Value, len(keys)), M: make(map[apl.Value]apl.Value)}
	var err error
	for i, k := range keys {
		d.K[i] = k.Copy()
		v := o.At(k)
		if leftarray {
			rv = v
		} else {
			lv = v
		}
		v, err = catenate(a, lv, rv)
		if err != nil {
			return nil, err
		}
		d.M[k.Copy()] = v.Copy()
	}

	if istable {
		return dict2table(&d)
	}
	return &d, nil
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
			d.K[i] = k.Copy()
			d.M[k.Copy()] = l.At(k).Copy()
		}
		for _, k := range r.Keys() {
			if _, ok := d.M[k]; !ok {
				d.K = append(d.K, k.Copy())
			}
			d.M[k.Copy()] = r.At(k).Copy()
		}
	} else {
		rkeys := r.Keys()
		if istable && len(keys) != len(rkeys) {
			return nil, fmt.Errorf("catenate table on first axis: tables have different number of columns")
		}
		for i, k := range keys {
			d.K[i] = k.Copy()
			lv := l.At(k)
			rv := r.At(k)
			if rv == nil {
				if istable {
					return nil, fmt.Errorf("catenate table on first axis: tables have different columns")
				}
				d.M[k.Copy()] = lv.Copy()
			} else {
				v, err := catenate(a, lv, rv)
				if err != nil {
					return nil, err
				}
				d.M[k.Copy()] = v.Copy()
			}
		}
		for _, k := range rkeys {
			if d.M[k] == nil {
				if istable {
					return nil, fmt.Errorf("catenate table on first axis: tables have different columns")
				}
				d.K = append(d.K, k.Copy())
				d.M[k.Copy()] = r.At(k).Copy()
			}
		}
	}

	if istable {
		return dict2table(&d)
	}
	return &d, nil
}

func dict2table(d *apl.Dict) (apl.Table, error) {
	rows := 0
	if len(d.K) > 0 {
		v := d.At(d.K[0])
		if ar, ok := v.(apl.Array); ok {
			shape := ar.Shape()
			if len(shape) == 1 {
				rows = shape[0]
			}
		}
	}
	return apl.Table{Dict: d, Rows: rows}, nil
}

func table2array(a *apl.Apl, t apl.Table) (apl.Array, error) {
	keys := t.Keys()
	rows := t.Rows
	res := apl.MixedArray{
		Dims:   []int{rows, len(keys)},
		Values: make([]apl.Value, len(keys)*rows),
	}
	n := len(keys)
	for k, key := range keys {
		col := t.At(key).(apl.Array)
		for i := 0; i < rows; i++ {
			v := col.At(i).Copy()
			res.Values[i*n+k] = v
		}
	}
	u, _ := a.Unify(res, true)
	return u, nil
}

// tableQuery applies the aggregation function to the columns of the table.
// The result is always a table.
// The function may return a scalar (row aggregation) or a vector.
// If a grouping is given, the aggregation is applied to each group.
//
// The function can be given as a single function or as a dictionary of functions.
// If it is given as a dictionary, the keys are used for the aggregated columns.
// Otherwise the resulting columns keep their names.
//
// The number of columns and functions must conform.
// A single function is applied to multiple columns,
// multiple functions may be applied to a column resulting in multiple aggregation columns,
// or the number of functions must match the number of columns.
//
// The group must be a single column key that selects the group column,
// or a function that is applied to the table.
// The function gets variables initialized with the column names, if they are strings:
//	{`w ⌊Date} rounds the Date column to weeks
// The group column should not be part of the aggregation.
// An anonymous group function always replaces the first column with the group result,
// before applying the aggregation.
func tableQuery(a *apl.Apl, t apl.Table, agg, grp apl.Value) (apl.Value, error) {

	keys := t.Keys()
	var gf apl.Function              // function to create group column
	var groupcol apl.Uniform         // group data column in input table
	var groupmap map[apl.Value][]int // map from group value to row indexes
	var groups []apl.Value           // distinct group values
	var groupres []apl.Value         // group column of result table
	var groupname apl.Value          // name of group column in result table
	if grp != nil {
		if o, ok := grp.(apl.Object); ok {
			if ks := o.Keys(); len(ks) != 1 {
				return nil, fmt.Errorf("groups object must have a single value")
			} else if f, ok := o.At(ks[0]).(apl.Function); ok {
				gf = f
				groupname = ks[0]
			} else {
				return nil, fmt.Errorf("groups object must contain a function: %T", o.At(ks[0]))
			}
			vars := make(map[string]apl.Value)
			for _, key := range keys {
				if s, ok := key.(apl.String); ok {
					vars[string(s)] = t.At(key).Copy()
				}
			}
			if len(keys) > 0 {
				keys = keys[1:]
			}

			v, err := a.EnvCall(gf, nil, apl.Int(t.Rows), vars)
			if err != nil {
				return nil, fmt.Errorf("group function: %s", err)
			}
			if ar, ok := v.(apl.Array); ok == false {
				return nil, fmt.Errorf("group function must return an array")
			} else if as := ar.Shape(); len(as) != 1 || as[0] != t.Rows {
				return nil, fmt.Errorf("group function result has wrong shape: %v != %d", as, t.Rows)
			} else if u, ok := a.Unify(ar, true); ok == false {
				return nil, fmt.Errorf("cannot unify group column")
			} else {
				groupcol = u.(apl.Uniform)
			}
		} else {
			vec := make([]apl.Value, 0, len(keys))
			for _, k := range keys {
				if k != grp {
					vec = append(vec, k)
				} else {
					groupname = k
				}
			}
			keys = vec
			if groupname == nil {
				return nil, fmt.Errorf("group does not exist")
			}
			groupcol = t.At(groupname).(apl.Uniform)
		}

		groupmap = make(map[apl.Value][]int)
		for i := 0; i < groupcol.Size(); i++ {
			v := groupcol.At(i)
			vec, ok := groupmap[v]
			vec = append(vec, i)
			groupmap[v] = vec
			if ok == false {
				groups = append(groups, v)
			}
		}
	}

	var names []apl.Value
	var fns []apl.Function
	if f, ok := agg.(apl.Function); ok {
		fns = make([]apl.Function, len(keys))
		for i := range fns {
			fns[i] = f
		}
		names = keys
	} else if o, ok := agg.(apl.Object); ok {
		names = o.Keys()
		if len(names) == 1 && len(keys) > 1 {
			ext := make([]apl.Value, len(keys))
			for i := range ext {
				ext[i] = names[0]
			}
			names = ext
		} else if len(keys) == 1 && len(names) > 1 {
			k := keys[0]
			keys = make([]apl.Value, len(names))
			for i := range keys {
				keys[i] = k
			}
		} else if len(names) != len(keys) {
			return nil, fmt.Errorf("number of aggregation functions and keys must conform")
		}
		fns = make([]apl.Function, len(names))
		for i, key := range names {
			f, ok := o.At(key).(apl.Function)
			if ok == false {
				return nil, fmt.Errorf("aggregation is not a function: %T", o.At(key))
			}
			fns[i] = f
		}
	} else {
		return nil, fmt.Errorf("aggregation functions must be passed in a dict: %T", agg)
	}

	// If no group is given, make a single one.
	if groupmap == nil {
		groupmap = make(map[apl.Value][]int)
		groups = []apl.Value{apl.Int(0)}
		idx := make([]int, t.Rows)
		for i := range idx {
			idx[i] = i
		}
		groupmap[apl.Int(0)] = idx
	}

	numrows := 0
	d := apl.Dict{M: make(map[apl.Value]apl.Value)}
	for k, key := range keys {
		column := t.At(key).(apl.Uniform)
		rescol := apl.MixedArray{Dims: []int{0}}
		f := fns[k]
		for _, gv := range groups {
			g := groupmap[gv]
			y := column.Make([]int{len(g)})
			for n, m := range g {
				y.Set(n, column.At(m).Copy())
			}
			r, err := f.Call(a, nil, y)
			if err != nil {
				return nil, err
			}
			ar, ok := r.(apl.Array)
			if ok == false {
				ar = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{r}}
			}
			for n := 0; n < ar.Size(); n++ {
				rescol.Values = append(rescol.Values, ar.At(n))
				rescol.Dims[0]++
			}
			if k == 0 {
				for n := 0; n < ar.Size(); n++ {
					groupres = append(groupres, gv)
				}
			}
		}
		u, ok := a.Unify(rescol, true)
		if ok == false {
			return nil, fmt.Errorf("column cannot be unified")
		}
		us := u.Shape()
		if len(us) != 1 {
			return nil, fmt.Errorf("aggregation result has rank %d", len(us))
		}

		d.K = append(d.K, names[k])
		d.M[names[k]] = u

		if k == 0 {
			numrows = us[0]
		} else if us[0] != numrows {
			return nil, fmt.Errorf("aggregation results in different number of rows")
		}
	}

	if groupname != nil {
		ug, ok := a.Unify(apl.MixedArray{Dims: []int{len(groupres)}, Values: groupres}, true)
		if ok == false {
			return nil, fmt.Errorf("cannot unify group column")
		}
		d.K = append([]apl.Value{groupname}, d.K...)
		d.M[groupname] = ug
	}

	return apl.Table{Rows: numrows, Dict: &d}, nil
}
