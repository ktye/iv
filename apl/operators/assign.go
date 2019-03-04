package operators

import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "←",
		Domain:  MonadicOp(nil),
		doc:     "assign, variable specification",
		derived: assign,
	})
}

func assign(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		as, ok := f.(*apl.Assignment)
		if ok == false {
			return nil, fmt.Errorf("cannot assign to %T", f)
		}
		if L != nil {
			return nil, fmt.Errorf("assign cannot be called dyadically")
		}

		if as.Identifiers != nil {
			if as.Indexes != nil {
				return nil, fmt.Errorf("vector and indexed assignment cannot exist simulaneously")
			}
			return assignVector(a, as.Identifiers, R, as.Modifier)
		}

		// Special case: channel scope: ⎕←C
		if c, ok := R.(apl.Channel); ok && as.Identifier == "⎕" {
			return c.Scope(a), nil
		}

		return R, assignScalar(a, as.Identifier, as.Indexes, as.Modifier, R)
	}
	return function(derived)
}

// AssignVector does a vector assignment from R to the given names.
// A modifier function may be applied.
func assignVector(a *apl.Apl, names []string, R apl.Value, mod apl.Function) (apl.Value, error) {
	var ar apl.Array
	if v, ok := R.(apl.Array); ok {
		ar = v
	} else {
		ar = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{R}}
	}

	var scalar apl.Value
	if s := ar.Shape(); len(s) != 1 {
		return nil, fmt.Errorf("vector assignment: rank of right argument must be 1")
	} else if s[0] != 1 && s[0] != len(names) {
		return nil, fmt.Errorf("vector assignment is non-conformant")
	} else if s[0] == 1 {
		if ar.Size() < 1 {
			return nil, fmt.Errorf("vector assignment: collapsed dimension")
		}
		scalar = ar.At(0)
	}

	var err error
	for i, name := range names {
		var v apl.Value
		if scalar != nil {
			v = scalar
		} else {
			if err := apl.ArrayBounds(ar, i); err != nil {
				return nil, err
			}
			v = ar.At(i)
		}
		err = assignScalar(a, name, nil, mod, v)
		if err != nil {
			return nil, err
		}
	}

	return R, nil
}

// AssignScalar assigns to a named scalar variable.
// If indexes is non-nil, it must be an IndexArray for indexed assignment.
// Mod may be a dyadic modifying function.
func assignScalar(a *apl.Apl, name string, indexes apl.Value, f apl.Function, R apl.Value) error {
	if f == nil && indexes == nil {
		return a.Assign(name, R)
	}

	w, env := a.LookupEnv(name)
	if w == nil {
		return fmt.Errorf("assign %s: modified/indexed: variable does not exist", name)
	}

	v, err := assignValue(a, w, indexes, f, R)
	if err != nil {
		return fmt.Errorf("assign %s: %s", name, err)
	}
	if v != nil {
		return a.AssignEnv(name, v, env)
	}
	return nil
}

// assignValue assigns to a given value. It may return a new value, or nil with no error.
func assignValue(a *apl.Apl, dst apl.Value, indexes apl.Value, f apl.Function, R apl.Value) (apl.Value, error) {
	// Modified assignment without indexing.
	if indexes == nil {
		if v, err := f.Call(a, dst, R); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}

	idx, ok := indexes.(apl.IntArray)
	if ok == false {
		to := ToIndexArray(nil)
		if v, ok := to.To(a, indexes); ok == false {
			return nil, fmt.Errorf("indexed assignment could not convert to IndexArray: %T", indexes)
		} else if _, ok := v.(apl.EmptyArray); ok {
			return nil, fmt.Errorf("indexed assignment could not convert to IndexArray: %T", indexes)
		} else {
			idx = v.(apl.IntArray)
		}
	}

	if t, ok := dst.(apl.Table); ok {
		return nil, assignTable(a, t, idx, f, R)
	}

	if obj, ok := dst.(apl.Object); ok {
		return nil, assignObject(a, obj, idx, f, R)
	}

	if lst, ok := dst.(apl.List); ok {
		return nil, assignList(a, lst, idx, f, R)
	}

	ar, ok := dst.(apl.ArraySetter)
	if ok == false {
		return nil, fmt.Errorf("variable is no settable array: %T", dst)
	}

	// Try to keep the original array type, upgrade only if needed.
	upgrade := func() {
		ga := apl.MixedArray{Dims: apl.CopyShape(ar)}
		ga.Values = make([]apl.Value, apl.ArraySize(ga))
		for i := range ga.Values {
			if i >= ar.Size() {
				return
			}
			ga.Values[i] = ar.At(i)
		}
		ar = ga
	}

	// modAssign assigns ar at index i with v possibly modified by f.
	modAssign := func(i int, v apl.Value) error {
		if i == -1 {
			// Index -1 is used by some indexed assignments to mark skipps.
			// E.g. replicate and compress / and \
			return nil
		}
		if f == nil {
			if err := ar.Set(i, v); err == nil {
				return nil
			}
			upgrade()
			return ar.Set(i, v)
		}
		var err error
		if err = apl.ArrayBounds(ar, i); err != nil {
			return err
		}
		v, err = f.Call(a, ar.At(i), v)
		if err != nil {
			return err
		}
		if err = ar.Set(i, v); err == nil {
			return nil
		}
		upgrade()
		return ar.Set(i, v)
	}

	var src apl.Array
	var scalar apl.Value
	if av, ok := R.(apl.Array); ok {
		src = av
		if apl.ArraySize(av) == 1 {
			scalar = av.At(0)
		}
	} else {
		scalar = R
	}
	if scalar != nil {
		// Scalar or 1-element assignment.
		for _, d := range idx.Ints {
			if err := modAssign(int(d), scalar); err != nil {
				return ar, err
			}
		}
	} else {

		// Shapes must conform. Single element axis are collapsed.
		collapse := func(s []int) []int {
			n := 0
			for _, i := range s {
				if i == 1 {
					n++
				}
			}
			if n == 0 {
				return s
			}
			r := make([]int, len(s)-n)
			k := 0
			for _, i := range s {
				if i != 1 {
					r[k] = i
					k++
				}
			}
			return r
		}
		ds := collapse(idx.Shape())
		ss := collapse(src.Shape())
		if len(ds) != len(ss) {
			return nil, fmt.Errorf("indexed assignment: arrays have different rank: %d != %d", len(ds), len(ss))
		}
		for i := range ds {
			if ss[i] != ds[i] {
				return nil, fmt.Errorf("indexed assignment: arrays are not conforming: %v != %v", ss, ds)
			}
		}
		for i, d := range idx.Ints {
			if err := apl.ArrayBounds(src, i); err != nil {
				return ar, err
			}
			if err := modAssign(int(d), src.At(i)); err != nil {
				return ar, err
			}
		}
	}
	return ar, nil
}

// assignTable updates a table.
// indexes are given in a fake IntArray. See primitives/index.go: tableSelection.
// R must be a table or array of corresponding size, an object for each row or a scalar value.
func assignTable(a *apl.Apl, t apl.Table, idx apl.IntArray, f apl.Function, R apl.Value) error {
	rows := idx.Ints[:idx.Dims[0]]
	cols := idx.Ints[idx.Dims[0]:]
	keys := make([]apl.Value, len(cols))
	for i := range keys {
		all := t.Keys()
		if i < 0 || i >= len(keys) {
			return fmt.Errorf("table-update: col idx out of range")
		}
		keys[i] = all[cols[i]]
	}

	if ar, ok := R.(apl.Array); ok {
		// convert array R to table.
		shape := ar.Shape()
		if len(shape) == 1 && ar.Size() == len(rows) {
			// Reshape column vectors.
			if rs, ok := ar.(apl.Reshaper); ok == false {
				return fmt.Errorf("table-update: cannot reshape right vector: %T", ar)
			} else {
				shape = []int{shape[0], 1}
				ar = rs.Reshape(shape).(apl.Array)
			}
		}

		if len(shape) != 2 {
			return fmt.Errorf("table-update: array on the right must have rank 2")
		}
		if shape[0] != len(rows) || shape[1] != len(cols) {
			return fmt.Errorf("table-update: array on the right has wrong shape")
		}
		m := make(map[apl.Value]apl.Value)
		for k, key := range keys {
			u := t.At(key).(apl.Uniform)
			col := u.Make([]int{len(rows)})
			to := ToType(reflect.TypeOf(u.Zero()), nil)
			for i := range rows {
				val := ar.At(i*shape[1] + k)
				v, ok := to.To(a, val)
				if ok == false {
					return fmt.Errorf("table-update: cannot convert %T to %T", val, u.Zero())
				}
				if err := col.Set(i, v); err != nil {
					return fmt.Errorf("table-update: convert array: %s", err)
				}
			}
			m[key] = col
		}
		d := apl.Dict{
			K: keys,
			M: m,
		}
		R = apl.Table{Dict: &d, Rows: shape[0]}
	} else if _, ok := R.(apl.Object); ok == false {
		// convert scalar R to dict.
		d := apl.Dict{
			K: keys, // It should be ok to share.
			M: make(map[apl.Value]apl.Value),
		}
		for _, k := range keys {
			d.M[k] = R // TODO: copy
		}
		R = &d
	}

	// TODO: this version of modified table assignment works on scalar values.
	// Maybe it should work on arrays or a subtable instead.
	update := func(col apl.ArraySetter, i int, v apl.Value) (err error) {
		if f != nil {
			ot := reflect.TypeOf(v)
			v, err = f.Call(a, col.At(i), v)
			nt := reflect.TypeOf(v)
			if ot != nt {
				return fmt.Errorf("mod changed type from %s to %s", ot, nt)
			}
			if err != nil {
				return err
			}
		}
		return col.Set(i, v)
	}

	o, ok := R.(apl.Object)
	if ok == false {
		return fmt.Errorf("table-update: illegal right argument: %T", R)
	}

	rk := o.Keys()
	if len(rk) != len(keys) {
		return fmt.Errorf("table-update: keys on the right do not match")
	}
	for i := range keys {
		if rk[i] != keys[i] {
			return fmt.Errorf("table-update: keys on the right do not match")
		}
	}
	for _, key := range keys {
		column := t.Dict.At(key)
		col, ok := column.(apl.ArraySetter)
		if ok == false {
			return fmt.Errorf("table-update: column is not settable: %T", column)
		}
		u := column.(apl.Uniform)

		if rt, ok := R.(apl.Table); ok {
			rc := rt.At(key)
			to := ToType(reflect.TypeOf(column), nil)
			nc, ok := to.To(a, rc)
			if ok == false {
				return fmt.Errorf("table-update: cannot convert %T to %T", rc, column)
			}
			ru := nc.(apl.Uniform)
			for i, n := range rows {
				if err := update(col, n, ru.At(i)); err != nil {
					return fmt.Errorf("table-update: %s", err)
				}
			}
		} else {
			z := u.Zero()
			to := ToType(reflect.TypeOf(z), nil)
			rv := o.At(key)
			v, ok := to.To(a, rv)
			if ok == false {
				return fmt.Errorf("table-update: cannot convert %T to %T", rv, z)
			}
			for _, n := range rows {
				if err := update(col, n, v); err != nil {
					return fmt.Errorf("table-update: %s", err)
				}
			}

		}
		if err := t.Dict.Set(key, column); err != nil {
			return fmt.Errorf("table-update: %s", err)
		}
	}
	return nil
}

// assignObject assigns R to index keys of a object.
func assignObject(a *apl.Apl, obj apl.Object, idx apl.IntArray, f apl.Function, R apl.Value) error {
	if f != nil {
		return fmt.Errorf("TODO: object: modified indexed assignment")
	}
	if len(idx.Ints) > 1 && idx.Ints[0] < 0 {
		return assignObjectDepth(a, obj, idx, f, R)
	} else if len(idx.Ints) == 1 && idx.Ints[0] < 0 {
		idx.Ints[0] = -1 - idx.Ints[0] + a.Origin
	}
	vectorize := false
	ar, ok := R.(apl.Array)
	if ok == true {
		if len(idx.Ints) > 1 {
			if len(idx.Ints) == ar.Size() {
				vectorize = true
			} else {
				return fmt.Errorf("assing object: assignment does not conform")
			}
		}
	}
	keys := obj.Keys()
	for i := 0; i < len(idx.Ints); i++ {
		n := int(idx.Ints[i] - a.Origin)
		if n < 0 || n >= len(keys) {
			return fmt.Errorf("assign object: index out of range")
		}
		k := keys[n]
		v := R // TODO: copy?
		if vectorize == true {
			if err := apl.ArrayBounds(ar, i); err != nil {
				return err
			}
			v = ar.At(i)
		}
		if err := obj.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func assignObjectDepth(a *apl.Apl, obj apl.Object, idx apl.IntArray, f apl.Function, R apl.Value) (err error) {
	k := -1 - idx.Ints[0]
	keys := obj.Keys()
	if k < 0 || k >= len(keys) {
		return fmt.Errorf("assign obj-depth: index out of range")
	}
	key := keys[k]
	v := obj.At(key)
	if v == nil {
		return fmt.Errorf("assign obj-depth: nil value")
	}

	ia := apl.IntArray{Dims: []int{idx.Dims[0] - 1}, Ints: idx.Ints[1:]}
	if _, ok := v.(apl.Table); ok {
		err = fmt.Errorf("assign obj-depth: tables are not supported")
	} else if o, ok := v.(apl.Object); ok {
		err = assignObject(a, o, ia, f, R)
	} else if l, ok := v.(apl.List); ok {
		err = assignList(a, l, ia, f, R)
	} else if ar, ok := v.(apl.Array); ok {
		var nv apl.Value
		nv, err = assignValue(a, ar, ia, f, R)
		if err == nil && nv != nil {
			v = nv
		}
	} else {
		err = fmt.Errorf("TODO: assign obj-depth: unsupported type: %T", v)
	}
	if err == nil {
		return obj.Set(key, v)
	}
	return err
}

// assignList assigns R to the depth index of a list.
func assignList(a *apl.Apl, l apl.List, idx apl.IntArray, f apl.Function, R apl.Value) error {
	if f != nil {
		v, err := l.GetDeep(idx.Ints)
		if err != nil {
			return err
		}
		v, err = f.Call(a, v, R)
		if err != nil {
			return err
		}
		R = v
	}
	return l.SetDeep(idx.Ints, R)
}
