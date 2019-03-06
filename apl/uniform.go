package apl

import (
	"fmt"
	"reflect"
)

// Unify tries to convert the array to a uniform array, if possible.
// If uptype is true, it uptypes numeric values, if that helps.
func (a *Apl) Unify(A Array, uptype bool) (resultarray Array, resultok bool) {
	if _, ok := A.(EmptyArray); ok {
		return A, false // An empty array is defined to be not uniform.
	}
	if _, ok := A.(Uniform); ok {
		return A, true
	}

	noNumber := -10
	boolType := reflect.TypeOf(Bool(false))
	indexType := reflect.TypeOf(Int(0))
	class := func(t reflect.Type) int {
		n, ok := a.Tower.Numbers[t]
		if ok == false {
			if t == indexType {
				return -1
			} else if t == boolType {
				return -2
			}
			return noNumber
		}
		return n.Class
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
	if max != noNumber {
		maxnumber = v0.(Number)
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
		if max != noNumber {
			if c := class(t); c == noNumber {
				max = noNumber
			} else if c > max {
				max = c
				maxnumber = v.(Number)
			}
		}
	}

	// All values have the same type.
	// Try to convert them to a compact uniform type.
	if sametype {
		// Some uniform types are defined in the numeric implementation.
		// E.g. numbers/{FloatArray;ComplexArray;TimeArray}.
		if max != noNumber {
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
		} else if t0 == reflect.TypeOf(Int(0)) {
			ar := IntArray{}.Make(CopyShape(A))
			for i := 0; i < A.Size(); i++ {
				v := A.At(i)
				ar.(IntArray).Ints[i] = int(v.(Int))
			}
			return ar, true
		}
		// Unknown uniform type is returnd as it is.
		return A, true
	} else if max == noNumber {
		// If values are not of the same type, and not identified by
		// the current tower, there is no chance to make them equal.
		return A, false
	}

	// Try to uptype all values to the same number type.
	values := make([]Value, size)
	var err error
	for i := 0; i < size; i++ {
		v := A.At(i)
		values[i], _, err = a.Tower.SameType(v.(Number), maxnumber)
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

// UnifyArray tries to unify the input array without uptyping.
func (a *Apl) UnifyArray(A Array) Array {
	u, _ := a.Unify(A, false)
	return u
}
