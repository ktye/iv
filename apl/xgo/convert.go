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
func export(a *apl.Apl, v apl.Value, t reflect.Type) (reflect.Value, error) {

	if e, ok := v.(Exporter); ok {
		x := e.Export()
		if x.Type().ConvertibleTo(t) {
			return x.Convert(t), nil
		}
	}

	number := func(from apl.Value, to apl.Number) (apl.Number, error) {
		src, ok := from.(apl.Number)
		if ok == false {
			return nil, fmt.Errorf("not a number: %T", from)
		}
		n, _, err := a.Tower.SameType(src, to)
		if err != nil {
			return nil, err
		}
		if reflect.TypeOf(to) != reflect.TypeOf(n) {
			return nil, fmt.Errorf("cannot convert to %T", to)
		}
		return n, nil
	}
	zero := reflect.Value{}
	switch t.Kind() {

	case reflect.Int:
		if n, err := number(v, numbers.Integer(0)); err != nil {
			return zero, err
		} else {
			return reflect.ValueOf(int(n.(numbers.Integer))), nil
		}

	case reflect.Float64:
		if n, err := number(v, numbers.Float(0)); err != nil {
			return zero, err
		} else {
			return reflect.ValueOf(float64(n.(numbers.Float))), nil
		}

	case reflect.Complex128:
		if n, err := number(v, numbers.Complex(0)); err != nil {
			return zero, err
		} else {
			return reflect.ValueOf(complex128(n.(numbers.Complex))), nil
		}

	case reflect.String:
		if s, ok := v.(apl.String); ok == false {
			return zero, fmt.Errorf("expected string")
		} else {
			return reflect.ValueOf(string(s)), nil
		}

	case reflect.Slice:
		ar, ok := v.(apl.Array)
		if ok == false {
			return zero, fmt.Errorf("expected slice")
		}
		et := t.Elem()
		//fmt.Println(t, et)
		n := apl.ArraySize(ar)
		s := reflect.MakeSlice(t, n, n)
		for i := 0; i < n; i++ {
			if v, err := ar.At(i); err != nil {
				return zero, err
			} else {
				if e, err := export(a, v, et); err != nil {
					return zero, err
				} else {
					se := s.Index(i)
					se.Set(e)
				}
			}
		}
		return s, nil

	default:
		return zero, fmt.Errorf("cannot convert to %v", t)
	}
}

// convert converts a go value to an apl value.
func convert(a *apl.Apl, v reflect.Value) (apl.Value, error) {

	if x := a.Import(v); x != nil {
		return x, nil
	}

	switch v.Kind() {

	case reflect.Int:
		return apl.Index(int(v.Int())), nil

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
			if e, err := convert(a, v.Index(i)); err != nil {
				return nil, err
			} else {
				ar.Values[i] = e
			}
		}
		return ar, nil

	default:
		return nil, fmt.Errorf("cannot convert %s to an apl value", v.Kind())
	}
}
