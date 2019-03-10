package xgo

import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
)

type Function struct {
	Name string
	Fn   reflect.Value
}

func (f Function) String(af apl.Format) string {
	return f.Name
}
func (f Function) Copy() apl.Value { return f }

// Call a go function.
// If it requires 1 argument, that is taken from the right value.
// Two arguments may be the right and left argument or a vector of 2 arguments.
// More than two arguments must be passed in a vector of the right size.
// If the function returns an error as the last value, it is checked and returned.
// Otherwise, or if the error is nil the result is converted and returned.
// More than one result will be returned as a List.
func (f Function) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	errarg := func(i int, err error) error {
		return fmt.Errorf("function %s argument %d: %s", f.Name, i+1, err)
	}
	t := f.Fn.Type()
	args := t.NumIn()
	in := make([]reflect.Value, args)
	var err error
	if args == 0 {
	} else if args == 1 {
		in[0], err = export(R, t.In(0))
		if err != nil {
			return nil, errarg(0, err)
		}
	} else if args == 2 && L != nil {
		in[0], err = export(R, t.In(0))
		if err != nil {
			return nil, errarg(0, err)
		}
		in[1], err = export(L, t.In(1))
		if err != nil {
			return nil, errarg(1, err)
		}
	} else if L == nil {
		ar, ok := R.(apl.Array)
		if ok == false {
			return nil, fmt.Errorf("function %s requires %d arguments", f.Name, args)
		}
		rs := ar.Shape()
		if len(rs) > 1 {
			return nil, fmt.Errorf("function argument has rank %d", len(rs))
		}
		if n := ar.Size(); n != args {
			return nil, fmt.Errorf("function %s requires %d arguments, R has size %d", f.Name, args, n)
		} else {
			for i := 0; i < args; i++ {
				in[i], err = export(ar.At(i), t.In(i))
				if err != nil {
					return nil, errarg(i, err)
				}
			}
		}
	}
	out := f.Fn.Call(in)

	// Test if the last output value is an error, check and remove it.
	if len(out) > 0 {
		if last := out[len(out)-1]; last.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			if last.IsNil() == false {
				return nil, last.Interface().(error)
			} else {
				out = out[:len(out)-1]
			}
		}
	}

	if len(out) == 0 {
		return apl.EmptyArray{}, nil
	} else if len(out) == 1 {
		return Convert(out[0])
	} else {
		res := make(apl.List, len(out))
		for i := range out {
			if v, err := Convert(out[i]); err != nil {
				return nil, err
			} else {
				res[i] = v
			}
		}
		return res, nil
	}
}
