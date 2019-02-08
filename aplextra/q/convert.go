package q

import (
	"fmt"
	"time"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	kdb "github.com/sv/kdbgo"
)

func FromAPL(a *apl.Apl, v apl.Value) (*kdb.K, error) {
	switch x := v.(type) {
	case apl.Bool:
		return &kdb.K{-kdb.KB, kdb.NONE, bool(x)}, nil
	case apl.Int:
		return kdb.Long(int64(x)), nil
	case apl.String:
		return kdb.Symbol(string(x)), nil
	case apl.BoolArray:
		if len(x.Dims) != 1 {
			return encodeArray(a, v)
		}
		return &kdb.K{kdb.KB, kdb.NONE, x.Bools}, nil
	case apl.IntArray:
		if len(x.Dims) != 1 {
			return encodeArray(a, v)
		}
		vec := make([]int64, len(x.Ints))
		for i, n := range x.Ints {
			vec[i] = int64(n)
		}
		return kdb.LongV(vec), nil
	case apl.StringArray:
		if len(x.Dims) != 1 {
			return encodeArray(a, v)
		}
		return kdb.SymbolV(x.Strings), nil
	case numbers.Float:
		return kdb.Float(float64(x)), nil
	case numbers.FloatArray:
		if len(x.Dims) != 1 {
			return encodeArray(a, v)
		}
		return kdb.FloatV(x.Floats), nil
	case numbers.Time:
		return &kdb.K{-kdb.KP, kdb.NONE, time.Time(x)}, nil
	case numbers.TimeArray:
		if len(x.Dims) != 1 {
			return encodeArray(a, v)
		}
		return &kdb.K{kdb.KP, kdb.NONE, x.Times}, nil
	case apl.List:
		return encodeList(a, v)
	case apl.Table:
		return encodeTable(a, v)
	default:
		if _, ok := v.(apl.Array); ok {
			return encodeArray(a, v)
		}
		if _, ok := v.(apl.Object); ok {
			return encodeObject(a, v)
		}
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

func encodeList(a *apl.Apl, v apl.Value) (*kdb.K, error) {
	lst := v.(apl.List)
	vec := make([]*kdb.K, len(lst))
	var err error
	for i, e := range lst {
		vec[i], err = FromAPL(a, e)
		if err != nil {
			return nil, err
		}
	}
	return kdb.NewList(vec...), nil
}
func encodeArray(a *apl.Apl, v apl.Value) (*kdb.K, error) {
	ar := v.(apl.Array)
	shape := ar.Shape()
	if len(shape) == 0 {
		return kdb.NewList(nil), nil
	}

	flat := make([]*kdb.K, ar.Size())
	for i := range flat {
		k, err := FromAPL(a, ar.At(i))
		if err != nil {
			return nil, err
		}
		flat[i] = k
	}

	for i := len(shape) - 1; i >= 0; i-- {
		w := shape[i]
		n := len(flat) / w
		lst := make([]*kdb.K, n)

		for j := 0; j < n; j++ {
			lst[j] = kdb.NewList(flat[j*w : j*w+w]...)
		}
		flat = lst
	}
	return kdb.NewList(flat...), nil
}

func encodeObject(a *apl.Apl, v apl.Value) (*kdb.K, error) {
	// kdb.Dict{ Key, Value *kdb.K}
	o := v.(apl.Object)
	keys := o.Keys()
	keystrings := make([]string, len(keys))
	values := make([]*kdb.K, len(keys))
	var err error
	for i, key := range keys {
		s, ok := key.(apl.String)
		if ok == false {
			return nil, fmt.Errorf("encode dict: only string keys are supported: %T", key)
		}
		keystrings[i] = string(s)
		v := o.At(a, key)
		if v == nil {
			return nil, fmt.Errorf("encode dict: key %s does not exist", s)
		}
		values[i], err = FromAPL(a, v)
		if err != nil {
			return nil, err
		}
	}
	return kdb.NewDict(kdb.SymbolV(keystrings), kdb.NewList(values...)), nil
}
func encodeTable(a *apl.Apl, v apl.Value) (*kdb.K, error) {
	t := v.(apl.Table)
	keys := t.Keys()
	columns := make([]string, len(keys))
	data := make([]*kdb.K, len(keys))

	var err error
	var col apl.List
	for i, k := range keys {
		if s, ok := k.(apl.String); ok == false {
			return nil, fmt.Errorf("encode table: only string headers are supported")
		} else {
			columns[i] = string(s)
		}

		col, err = toList(t.M[k])
		if err != nil {
			return nil, err
		}

		data[i], err = FromAPL(a, col)
		if err != nil {
			return nil, err
		}
	}
	return kdb.NewTable(columns, data), nil
}

func toList(v apl.Value) (apl.List, error) {
	ar, ok := v.(apl.Array)
	if ok == false {
		return nil, fmt.Errorf("cannot convert to list: %T", v)
	}
	l := make(apl.List, ar.Size())
	for i := range l {
		l[i] = ar.At(i)
	}
	return l, nil
}

func ToAPL(data *kdb.K) (apl.Value, error) {
	k := data.Data
	switch data.Type {
	case -kdb.KB:
		return apl.Bool(k.(bool)), nil
	case -kdb.KH:
		return apl.Int(k.(int16)), nil
	case -kdb.KI, -kdb.KD, -kdb.KU, -kdb.KV:
		return apl.Int(k.(int32)), nil
	case -kdb.KJ:
		return apl.Int(int(k.(int64))), nil
	case -kdb.KS:
		return apl.String(k.(string)), nil
	case -kdb.KF, -kdb.KZ:
		return numbers.Float(k.(float64)), nil
	case kdb.KB:
		v := k.([]bool)
		return apl.BoolArray{Dims: []int{len(v)}, Bools: v}, nil
	case kdb.KH:
		vec := k.([]int16)
		ints := make([]int, len(vec))
		for i := range vec {
			ints[i] = int(vec[i])
		}
		return apl.IntArray{Dims: []int{len(ints)}, Ints: ints}, nil
	case kdb.KI:
		vec := k.([]int32)
		ints := make([]int, len(vec))
		for i := range vec {
			ints[i] = int(vec[i])
		}
		return apl.IntArray{Dims: []int{len(ints)}, Ints: ints}, nil
	case kdb.KJ:
		vec := k.([]int64)
		ints := make([]int, len(vec))
		for i := range vec {
			ints[i] = int(vec[i])
		}
		return apl.IntArray{Dims: []int{len(ints)}, Ints: ints}, nil
	case kdb.KF:
		vec := k.([]float64)
		return numbers.FloatArray{Dims: []int{len(vec)}, Floats: vec}, nil
	case kdb.KS:
		vec := k.([]string)
		return apl.StringArray{Dims: []int{len(vec)}, Strings: vec}, nil
	case -kdb.KP:
		return numbers.Time(k.(time.Time)), nil
	case kdb.KP:
		vec := k.([]time.Time)
		return numbers.TimeArray{Dims: []int{len(vec)}, Times: vec}, nil
	case kdb.XD:
		return decodeDict(k)
	case kdb.XT:
		return decodeTable(k)
	case kdb.K0:
		return decodeList(k)
	default:
		return nil, fmt.Errorf("unknown K type: %d\n%v", data.Type, data)
	}
}

func decodeDict(k interface{}) (apl.Value, error) {
	xd := k.(kdb.Dict)
	kv, ok := xd.Key.Data.([]string)
	if ok == false {
		return nil, fmt.Errorf("decode dict: expected []string keys: %T", xd.Key.Data)
	}

	var d apl.Dict
	d.K = make([]apl.Value, len(kv))
	for i := range kv {
		d.K[i] = apl.String(kv[i])
	}
	d.M = make(map[apl.Value]apl.Value)

	vals, err := ToAPL(xd.Value)
	if err != nil {
		return nil, err
	}

	ar, ok := vals.(apl.Array)
	if ok == false {
		return nil, fmt.Errorf("decode dict: expected vector values: %T", vals)
	}
	if ar.Size() != len(d.K) {
		return nil, fmt.Errorf("decode dict: keys and values have differnt lengths")
	}
	for i := range d.K {
		d.M[d.K[i]] = ar.At(i)
	}

	return &d, nil
}

func decodeTable(k interface{}) (apl.Value, error) {
	kt := k.(kdb.Table)
	keys := make([]apl.Value, len(kt.Columns))
	m := make(map[apl.Value]apl.Value)
	rows := 0
	for i := range keys {
		keys[i] = apl.String(kt.Columns[i])
		v, err := ToAPL(kt.Data[i])
		if err != nil {
			return nil, err
		}
		m[keys[i]] = v
		if i == 0 {
			ar, ok := v.(apl.Array)
			if ok == false {
				return nil, fmt.Errorf("table data: expected array: %T", v)
			}
			shape := ar.Shape()
			if len(shape) != 1 {
				return nil, fmt.Errorf("table data: expected vector: %T %v", v, shape)
			}
			rows = shape[0]
		}
	}
	return apl.Table{
		Dict: &apl.Dict{
			K: keys,
			M: m,
		},
		Rows: rows,
	}, nil
}

func decodeList(k interface{}) (apl.Value, error) {
	kl := k.([]*kdb.K)
	l := make(apl.List, len(kl))
	var err error
	for i := range l {
		l[i], err = ToAPL(kl[i])
		if err != nil {
			return nil, err
		}
	}
	return l, nil
}
