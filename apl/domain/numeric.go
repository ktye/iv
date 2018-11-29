package domain

import (
	"reflect"

	"github.com/ktye/iv/apl"
)

// ToNumber accepts scalars and single size arrays.
// and converts them to scalars if they contain one of the types:
// apl.Bool, apl.Int, apl.Float or apl.Complex.
func ToNumber(child SingleDomain) SingleDomain {
	return number{child, true}
}

// IsNumber accepts scalars if they contain of of the types:
// apl.Bool, apl.Int, apl.Float or apl.Complex
func IsNumber(child SingleDomain) SingleDomain {
	return number{child, false}
}

type number struct {
	child   SingleDomain
	convert bool
}

func (n number) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	v := V
	if ar, ok := V.(apl.Array); ok {
		if n.convert == false {
			return V, false
		}
		if n := apl.ArraySize(ar); n != 1 {
			return V, false
		}
		v, _ = ar.At(0)
	}
	if b, ok := V.(apl.Bool); ok {
		return a.Tower.FromBool(b), true
	}
	if i, ok := V.(apl.Index); ok {
		return a.Tower.FromIndex(int(i)), true
	}
	if _, ok := a.Tower.Numbers[reflect.TypeOf(v)]; ok {
		return v, true
	}
	return V, false
}
func (n number) String(a *apl.Apl) string {
	name := "number"
	if n.convert {
		name = "tonumber"
	}
	if n.child == nil {
		return name
	}
	return name + " " + n.child.String(a)
}

// ToIndex converts the scalar to an Index.
func ToIndex(child SingleDomain) SingleDomain {
	return index{child}
}

type index struct {
	child SingleDomain
}

func (idx index) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if b, ok := V.(apl.Bool); ok {
		if b == true {
			return apl.Index(1), true
		}
		return apl.Index(0), true
	}
	if n, ok := V.(apl.Index); ok {
		return n, true
	}
	if n, ok := V.(apl.Number); ok == false {
		return V, false
	} else {
		if i, ok := n.ToIndex(); ok == false {
			return V, false
		} else {
			return propagate(a, apl.Index(i), idx.child)
		}
	}
}
func (idx index) String(a *apl.Apl) string {
	if idx.child == nil {
		return "index"
	} else {
		return "index " + idx.child.String(a)
	}
}

/* TODO: convert to new numbers package, if needed.
// ToBool accepts a number and converts it to Bool.
// If fails, if it is not 0 or 1.
func ToBool(child SingleDomain) SingleDomain {
	return toNumber{child, toBool, "bool", true}
}
func IsBool(child SingleDomain) SingleDomain {
	return toNumber{child, toBool, "bool", false}
}

// ToInt accepts a number and converts it to Int.
// by uptyping Bool and downtyping Float and Complex
// if they have no fractional or imaginary part.
func ToInt(child SingleDomain) SingleDomain {
	return toNumber{child, toInt, "int", true}
}
func IsInt(child SingleDomain) SingleDomain {
	return toNumber{child, toInt, "int", false}
}

// ToFloat accepts a number and converts it to Float
// by uptyping Bool and Int and downtyping Complex,
// if the imaginary part is zero.
func ToFloat(child SingleDomain) SingleDomain {
	return toNumber{child, toFloat, "float", true}
}
func IsFloat(child SingleDomain) SingleDomain {
	return toNumber{child, toFloat, "float", false}
}

// ToComplex accepts a number and converts it to Complex.
func ToComplex(child SingleDomain) SingleDomain {
	return toNumber{child, toComplex, "complex", true}
}
func IsComplex(child SingleDomain) SingleDomain {
	return toNumber{child, toComplex, "complex", false}
}

type toNumber struct {
	child   SingleDomain
	to      func(apl.Value, bool) (apl.Value, bool) // toBool, toInt, toFloat, toComplex
	s       string
	convert bool
}

func (t toNumber) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if v, ok := t.to(V, t.convert); ok {
		if t.child == nil {
			return v, true
		} else {
			return t.child.To(a, v)
		}
	}
	return V, false
}
func (t toNumber) String(a *apl.Apl) string {
	name := t.s
	if t.convert {
		name = "to" + t.s
	}
	if t.child == nil {
		return name
	}
	return name + " " + t.child.String(a)
}

// ToBool converts number to Bool.
// It fails if the value is not 0 or 1.
func toBool(V apl.Value, convert bool) (apl.Value, bool) {
	if _, ok := V.(apl.Bool); ok {
		return V, true
	} else if convert == false {
		return V, false
	}
	if n, ok := toInt(V, true); ok {
		if n.(apl.Int) == 0 {
			return apl.Bool(false), true
		} else if n.(apl.Int) == 1 {
			return apl.Bool(true), true
		}
	}
	return V, false
}

// ToInt converts a number to Int.
// It uptypes Bool and downtypes Float and Complex if they have no fractional
// or imaginary part.
// It fails if the conversion is not possible.
func toInt(V apl.Value, convert bool) (apl.Value, bool) {
	if _, ok := V.(apl.Int); ok {
		return V, true
	} else if convert == false {
		return V, false
	}
	switch v := V.(type) {
	case apl.Bool:
		if v {
			return apl.Int(1), true
		}
		return apl.Int(0), true
	case apl.Float:
		i := int(float64(v))
		if apl.Float(i) == v {
			return apl.Int(i), true
		}
		return V, false
	case apl.Complex:
		c := complex128(v)
		if imag(c) != 0 {
			return V, false
		}
		i := int(real(c))
		if float64(i) == real(c) {
			return apl.Int(i), true
		}
		return V, false
	}
	return V, false
}

// toFloat converts a number to Float.
// It uptypes Bool and Int and downtypes Complex, if the imaginary part is zero.
// It fails if V is not convertible.
func toFloat(V apl.Value, convert bool) (apl.Value, bool) {
	if _, ok := V.(apl.Float); ok {
		return V, true
	} else if convert == false {
		return V, false
	}
	switch v := V.(type) {
	case apl.Bool:
		if v {
			return apl.Float(1), true
		}
		return apl.Float(0), true
	case apl.Int:
		return apl.Float(v), true
	case apl.Complex:
		if imag(complex128(v)) == 0 {
			return apl.Float(real(complex128(v))), true
		}
	}
	return V, false
}

// ToComplex converts a number to Complex.
// It uptypes Bool, Int and Float.
// It fails if V is not convertible.
func toComplex(V apl.Value, convert bool) (apl.Value, bool) {
	if _, ok := V.(apl.Complex); ok {
		return V, true
	} else if convert == false {
		return V, false
	}
	switch v := V.(type) {
	case apl.Bool:
		if v {
			return apl.Complex(complex(1, 0)), true
		}
		return apl.Complex(complex(0, 0)), true
	case apl.Int:
		return apl.Complex(complex(float64(v), 0)), true
	case apl.Float:
		return apl.Complex(complex(float64(v), 0)), true
	}
	return V, false
}
*/
