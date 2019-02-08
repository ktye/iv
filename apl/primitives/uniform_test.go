package primitives

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

// TestUnify should be in apl where Unify is defined.
// But there are no numbers available.
func TestUnify(t *testing.T) {
	testCases := []struct {
		in  apl.Array
		t1  reflect.Type
		ok1 bool
		t2  reflect.Type
		ok2 bool
	}{
		// Known uniform types.
		{
			apl.BoolArray{Dims: []int{2}, Bools: []bool{false, true}},
			reflect.TypeOf(apl.BoolArray{}), true,
			reflect.TypeOf(apl.BoolArray{}), true,
		},
		{
			apl.IndexArray{Dims: []int{2}, Ints: []int{1, 2}},
			reflect.TypeOf(apl.IndexArray{}), true,
			reflect.TypeOf(apl.IndexArray{}), true,
		},
		{
			apl.StringArray{Dims: []int{2}, Strings: []string{"a", "b"}},
			reflect.TypeOf(apl.StringArray{}), true,
			reflect.TypeOf(apl.StringArray{}), true,
		},
		{
			numbers.FloatArray{Dims: []int{2}, Floats: []float64{1.2, 3.4}},
			reflect.TypeOf(numbers.FloatArray{}), true,
			reflect.TypeOf(numbers.FloatArray{}), true,
		},
		{
			numbers.ComplexArray{Dims: []int{2}, Cmplx: []complex128{1.2, 3.4}},
			reflect.TypeOf(numbers.ComplexArray{}), true,
			reflect.TypeOf(numbers.ComplexArray{}), true,
		},
		{
			numbers.TimeArray{Dims: []int{2}, Times: []time.Time{time.Time{}, time.Time{}}},
			reflect.TypeOf(numbers.TimeArray{}), true,
			reflect.TypeOf(numbers.TimeArray{}), true,
		},

		{ // The empty array is not uniform.
			apl.EmptyArray{},
			reflect.TypeOf(apl.EmptyArray{}), false,
			reflect.TypeOf(apl.EmptyArray{}), false,
		},

		// Truely mixed arrays.
		{
			apl.MixedArray{Dims: []int{1}, Values: []apl.Value{apl.Index(0), apl.String("")}},
			reflect.TypeOf(apl.MixedArray{}), false,
			reflect.TypeOf(apl.MixedArray{}), false,
		},

		// Mixed but uptypable.
		{
			apl.MixedArray{Dims: []int{4}, Values: []apl.Value{apl.Bool(false), apl.Index(0), apl.Index(1), numbers.Float(1)}},
			reflect.TypeOf(apl.MixedArray{}), false,
			reflect.TypeOf(numbers.FloatArray{}), true,
		},

		// Mixed but uptypable.
		{
			apl.MixedArray{Dims: []int{3}, Values: []apl.Value{apl.Bool(false), apl.Index(0), apl.Index(3)}},
			reflect.TypeOf(apl.MixedArray{}), false,
			reflect.TypeOf(apl.IndexArray{}), true,
		},

		// Non-uniform list.
		{
			apl.List{apl.String(""), apl.Bool(false)},
			reflect.TypeOf(apl.List{}), false,
			reflect.TypeOf(apl.List{}), false,
		},

		// Uptypable list.
		{
			apl.List{apl.Bool(false), numbers.Complex(complex(3, 4))},
			reflect.TypeOf(apl.List{}), false,
			reflect.TypeOf(numbers.ComplexArray{}), true,
		},

		// List of lists, reports as uniform, but does not implement apl.Uniform.
		// This is a special case for tables.
		{
			apl.List{apl.List{apl.Index(1)}, apl.List{}},
			reflect.TypeOf(apl.List{}), true,
			reflect.TypeOf(apl.List{}), true,
		},
	}

	listType := reflect.TypeOf(apl.List{})
	for i, tc := range testCases {
		var buf strings.Builder
		a := apl.New(&buf)
		numbers.Register(a)

		str := tc.in.String(a)
		oks := []bool{tc.ok1, tc.ok2}
		ts := []reflect.Type{tc.t1, tc.t2}
		uptype := false

		for k := range oks {
			out, ok := a.Unify(tc.in, uptype)
			if got := out.String(a); got != str {
				if reflect.TypeOf(tc.in) == listType && uptype == true {
					// An uptyped list has a different string representation.
				} else {
					t.Fatalf("tc %d, uptype=%v expected:\n%s\ngot:\n%s", i, uptype, str, got)
				}
			}
			if ok != oks[k] {
				t.Fatalf("tc %d uptype=%v expected ok: %v got %v", i, uptype, oks[k], ok)
			}
			if got := reflect.TypeOf(out); got != ts[k] {
				t.Fatalf("tc %d uptype=%v expected type: %v got %v", i, uptype, ts[k], got)
			}
			_, uni := out.(apl.Uniform)
			if uni != ok {
				if reflect.TypeOf(out) == listType && uni == false {
					// Special case for nested lists: Unify is ok, but does not implement Uniform
				} else {
					t.Fatalf("tc %d uptype=%v, uni: %v but ok: %v", i, uptype, uni, ok)
				}
			}
			uptype = true
		}
	}
}
