package apl

import (
	"fmt"
	"reflect"
)

// Unify tries to convert the array to a uniform array, if possible.
// If uptype is true, it uptypes numeric values, if that helps.
func (a *Apl) Unify(A Array, uptype bool) (Array, bool) {
	if _, ok := A.(EmptyArray); ok {
		return A, false // An empty array is defined to be not uniform.
	}
	if _, ok := A.(Uniform); ok {
		return A, true
	}

	boolType := reflect.TypeOf(Bool(false))
	indexType := reflect.TypeOf(Index(0))
	class := func(t reflect.Type) int {
		n, ok := a.Tower.Numbers[t]
		if ok == false {
			if t == boolType || t == indexType {
				return 0 // This should always be ok.
			}
			return -1
		}
		return n.Class
	}
	tonumber := func(v Value) Number {
		if b, ok := v.(Bool); ok {
			return a.Tower.FromBool(b)
		} else if i, ok := v.(Index); ok {
			return a.Tower.FromIndex(int(i))
		}
		return v.(Number)
	}

	// If all values of the array are the same type, it is uniform.
	// This includes a List who's primary values are lists as well.
	size := A.Size()
	if size < 1 {
		return A, true
	}
	v0 := A.At(0)
	t0 := reflect.TypeOf(v0)
	max := class(t0)
	var maxnumber Number
	if max != -1 {
		maxnumber = tonumber(v0)
	}
	sametype := true
	for i := 1; i < size; i++ {
		v := A.At(i)
		t := reflect.TypeOf(v)
		if t != t0 {
			sametype = false
			if uptype == false {
				return A, false
			}
		}
		if max != -1 {
			if c := class(t); c == -1 {
				max = -1
			} else if c > max {
				max = c
				maxnumber = tonumber(v)
			}
		}
	}

	// All values have the same type.
	// Try to convert them to a compact uniform type.
	if sametype {
		// Some uniform types are defined in the numeric implementation.
		// E.g. numbers/{FloatArray;ComplexArray;TimeArray}.
		if max != -1 {
			var values []Value
			switch v := A.(type) {
			case MixedArray:
				values = v.Values
			case List:
				values = []Value(v)
			default:
				return A, true
			}
			if u, ok := a.Tower.Uniform(values); ok == false {
				return A, true
			} else if rs, ok := u.(Reshaper); ok {
				return rs.Reshape(CopyShape(A)).(Array), true
			}
		}
		// Some uniform types are defined in array.go.
		if t0 == reflect.TypeOf(String("")) {
			ar := StringArray{}.Make(CopyShape(A))
			for i := 0; i < A.Size(); i++ {
				v := A.At(i)
				ar.(StringArray).Strings[i] = string(v.(String))
			}
			return ar, true
		} else if t0 == reflect.TypeOf(Bool(false)) {
			ar := BoolArray{}.Make(CopyShape(A))
			for i := 0; i < A.Size(); i++ {
				v := A.At(i)
				ar.(BoolArray).Bools[i] = bool(v.(Bool))
			}
			return ar, true
		} else if t0 == reflect.TypeOf(Index(0)) {
			ar := IndexArray{}.Make(CopyShape(A))
			for i := 0; i < A.Size(); i++ {
				v := A.At(i)
				ar.(IndexArray).Ints[i] = int(v.(Index))
			}
			return ar, true
		}
		// Unknown uniform type is returnd as it is.
		return A, true
	} else if max == -1 {
		// If values are not of the same type, and not identified by
		// the current tower, there is no chance to make them equal.
		return A, false
	}

	// Try to uptype all values to the same number type.
	values := make([]Value, size)
	var err error
	for i := 0; i < size; i++ {
		v := A.At(i)
		values[i], _, err = a.Tower.SameType(tonumber(v), maxnumber)
		if err != nil {
			fmt.Println(err)
			return A, false
		}
	}
	if u, ok := a.Tower.Uniform(values); ok {
		if rs, ok := u.(Reshaper); ok {
			return rs.Reshape(CopyShape(A)).(Array), true
		}
	}
	return A, false
}
