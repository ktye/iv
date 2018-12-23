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

func (f Function) String(a *apl.Apl) string {
	return f.Name
}

// Call a go function.
// If it requires 1 argument, that is taken from the right value.
// Two arguments may be the right and left argument or a vector of 2 arguments.
// More than two arguments must be passed in a vector of the right size.
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
		in[0], err = export(a, R, t.In(0))
		if err != nil {
			return nil, errarg(0, err)
		}
	} else if args == 2 && L != nil {
		in[0], err = export(a, R, t.In(0))
		if err != nil {
			return nil, errarg(0, err)
		}
		in[1], err = export(a, L, t.In(1))
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
		if n := apl.ArraySize(ar); n != args {
			return nil, fmt.Errorf("function %s requires %d arguments, R has size %d", f.Name, args, n)
		} else {
			for i := 0; i < args; i++ {
				v, err := ar.At(i)
				if err != nil {
					return nil, err
				}
				in[i], err = export(a, v, t.In(i))
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
		return convert(a, out[0])
	} else {
		res := apl.GeneralArray{Dims: []int{len(out)}, Values: make([]apl.Value, len(out))}
		for i := range out {
			if v, err := convert(a, out[i]); err != nil {
				return nil, err
			} else {
				res.Values[i] = v
			}
		}
		return res, nil
	}
}
