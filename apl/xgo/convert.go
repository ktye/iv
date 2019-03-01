package xgo

import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

// Exporter can be implemented by an apl.Value to be able convert it to a go value.
type Exporter interface {
	Export() reflect.Value
}

// export converts an apl value to a go value.
func export(v apl.Value, t reflect.Type) (reflect.Value, error) {

	if e, ok := v.(Exporter); ok {
		x := e.Export()
		if x.Type().ConvertibleTo(t) {
			return x.Convert(t), nil
		}
	}

	zero := reflect.Value{}
	switch t.Kind() {

	case reflect.Int:
		return reflect.ValueOf(int(v.(apl.Int))), nil

	case reflect.Float64:
		return reflect.ValueOf(float64(v.(numbers.Float))), nil

	case reflect.Complex128:
		return reflect.ValueOf(complex128(v.(numbers.Complex))), nil

	case reflect.String:
		return reflect.ValueOf(string(v.(apl.String))), nil

	case reflect.Slice:
		ar, ok := v.(apl.Array)
		if ok == false {
			return zero, fmt.Errorf("expected slice: %T", v)
		}
		et := t.Elem()
		n := apl.ArraySize(ar)
		s := reflect.MakeSlice(t, n, n)
		for i := 0; i < n; i++ {
			if e, err := export(ar.At(i), et); err != nil {
				return zero, err
			} else {
				se := s.Index(i)
				se.Set(e)
			}
		}
		return s, nil

	case reflect.Struct:
		if xv, ok := v.(Value); ok {
			st := reflect.Value(xv).Type()
			if st == reflect.PtrTo(t) {
				return reflect.Indirect(reflect.Value(xv)), nil
			} else if st == t {
				return reflect.Value(xv), nil
			}
		}
		return zero, fmt.Errorf("xgo: export struct: cannot convert %T to %s", v, t)

	default:
		return zero, fmt.Errorf("cannot convert to %v (%s)", t, t.Kind())
	}
}

// convert converts a go value to an apl value.
func Convert(v reflect.Value) (apl.Value, error) {
	switch v.Kind() {
	case reflect.Int:
		return apl.Int(int(v.Int())), nil

	case reflect.Uint:
		return apl.Int(int(v.Uint())), nil

	case reflect.Float64:
		return numbers.Float(v.Float()), nil

	case reflect.Complex128:
		return numbers.Complex(v.Complex()), nil

	case reflect.String:
		return apl.String(v.String()), nil

	case reflect.Slice:
		n := v.Len()
		ar := apl.MixedArray{Dims: []int{n}, Values: make([]apl.Value, n)}
		for i := range ar.Values {
			if e, err := Convert(v.Index(i)); err != nil {
				return nil, err
			} else {
				ar.Values[i] = e
			}
		}
		return ar, nil

	case reflect.Struct:
		return Value(v.Addr()), nil // TODO: populate

	default:
		return nil, fmt.Errorf("cannot convert %s to an apl value", v.Kind())
	}
}
