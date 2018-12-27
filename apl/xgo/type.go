package xgo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ktye/iv/apl"
)

// New returns an initialization function for the given type.
func New(t reflect.Type) create {
	return create{t}
}

type Value reflect.Value

func (v Value) String(a *apl.Apl) string {
	return fmt.Sprintf("%v:%v", reflect.Value(v).Type(), reflect.Value(v))
}

// Fields returns the field names, if the value is a struct.
// It does not return the method names.
// It returns nil, if the Value is not a struct.
func (v Value) Fields() []string {
	val := reflect.Value(v)
	if val.Kind() != reflect.Struct {
		return nil
	}
	t := val.Type()
	n := t.NumField()
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = t.Field(i).Name
	}
	return res
}

// Field returns the value of a field or a method with the given name.
func (v Value) Field(a *apl.Apl, name string) apl.Value {
	val := reflect.Value(v)
	var zero reflect.Value
	Name := strings.Title(name)
	m := val.MethodByName(Name)
	if m != zero {
		return Function{Name: Name, Fn: m}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	sf := val.FieldByName(Name)
	if sf == zero {
		return nil
	}
	rv, err := convert(a, sf)
	if err != nil {
		return nil
	}
	return rv
}

func (v Value) Set(a *apl.Apl, field string, fv apl.Value) error {
	val := reflect.Value(v).Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: cannot set field")
	}
	sf := val.FieldByName(field)
	var zero reflect.Value
	if sf == zero {
		return fmt.Errorf("%v: field does not exist: %s", val.Type(), field)
	}
	sv, err := export(a, fv, sf.Type())
	if err != nil {
		return err
	}
	sf.Set(sv)
	return nil
}

type create struct {
	reflect.Type
}

func (t create) String(a *apl.Apl) string {
	return fmt.Sprintf("new %v", t.Type)
}

func (t create) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	v := reflect.New(t.Type)
	return Value(v), nil
}

/*
func (v Value) MethodCall(a *apl.Apl, method string, L, R apl.Value) (apl.Value, error) {
	val := reflect.Value(v)
	var zero reflect.Value
	method = strings.Title(method)
	m := val.MethodByName(method)
	if m == zero {
		return nil, fmt.Errorf("%T has no method %s", val.Type, method)
	}
	fn := Function{Name: method, Fn: m}
	return fn.Call(a, L, R)
}
*/
