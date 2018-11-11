package domain

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"reflect"

	"github.com/ktye/iv/apl"
)

var ErrNaN error

type ErrInf int

func (e ErrInf) Error() string {
	if e < 0 {
		return "-Inf"
	} else {
		return "Inf"
	}
}

func init() {
	ErrNaN = errors.New("NaN")
}

// IsFloatErr returns ErrInf(1), ErrInf(-1) or ErrNaN, if the
// Float or Complex argument is not finite.
func IsFloatErr(f apl.Value) error {
	switch t := f.(type) {
	case apl.Float:
		v := float64(t)
		if math.IsInf(v, 1) {
			return ErrInf(1)
		} else if math.IsInf(v, -1) {
			return ErrInf(-1)
		} else if math.IsNaN(v) {
			return ErrNaN
		}
		return nil
	case apl.Complex:
		c := complex128(t)
		if cmplx.IsInf(c) {
			return ErrInf(1)
		} else if cmplx.IsNaN(c) {
			return ErrNaN
		}
		return nil
	default:
		return fmt.Errorf("IsFloatErr expected Float or Complex: %T", f)
	}
}

// CompareScalars compares and b and returns
// a == b, a < b.
// Non-numeric values are checked for equality.
// If they have the same type and are a Lesser, they
// are also checke for a < b.
// Numeric types are converted to the same type before comparison.
func CompareScalars(a, b apl.Value) (bool, bool, error) {
	// Equal types and equal.
	if a == b {
		return true, false, nil
	}
	// EqualTypes have lesser.
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		if al, ok := a.(Lesser); ok {
			le := al.Less(b)
			return false, le, nil
		}
	}
	// Convert both to numbers.
	num := number{nil, true}
	an, aok := num.To(nil, a)
	bn, bok := num.To(nil, b)
	if aok == false || bok == false {
		return false, false, fmt.Errorf("cannot compare %T with %T", a, b)
	}

	an, bn = MustPromote(an, bn)
	switch an := an.(type) {
	case apl.Bool:
		if an == bn.(apl.Bool) {
			return true, false, nil
		} else if bool(an) == false && bool(bn.(apl.Bool)) == true {
			return false, true, nil
		}
		return false, false, nil
	case apl.Int:
		i := bn.(apl.Int)
		return an == i, an < i, nil
	case apl.Float:
		f := bn.(apl.Float)
		if math.IsNaN(float64(f)) || math.IsNaN(float64(an)) {
			return false, false, ErrNaN
		}
		return an == bn, an < f, nil
	case apl.Complex:
		c := bn.(apl.Complex)
		if cmplx.IsNaN(complex128(an)) || cmplx.IsNaN(complex128(c)) {
			return false, false, ErrNaN
		}
		if an == c {
			return true, false, nil
		}
		if imag(an) != 0 || imag(c) != 0 {
			return false, false, fmt.Errorf("cannot compare complex numbers")
		}
		return real(an) == real(c), real(an) < real(c), nil
	}
	return false, false, fmt.Errorf("cannot compare %T with %T", a, b)
}

// Lesser is used for scalar comparison.
// Custom types may implement it.
// It must be called only after checking that the types are identical.
// It returns true, if the receiver is less that b.
type Lesser interface {
	Less(b interface{}) bool
}

// MustPromote returns two numbers of the same type, by possibly uptyping one of them.
// mustPromote can only be called with input types apl.Bool, apl.Int, apl.Float and apl.Complex.
func MustPromote(L, R apl.Value) (apl.Value, apl.Value) {
	to := ToBool(nil)
	if _, ok := L.(apl.Complex); ok {
		to = ToComplex(nil)
	} else if _, ok := R.(apl.Complex); ok {
		to = ToComplex(nil)
	} else if _, ok := L.(apl.Float); ok {
		to = ToFloat(nil)
	} else if _, ok := R.(apl.Float); ok {
		to = ToFloat(nil)
	} else if _, ok := L.(apl.Int); ok {
		to = ToInt(nil)
	} else if _, ok := R.(apl.Int); ok {
		to = ToInt(nil)
	}
	L, _ = to.To(nil, L)
	R, _ = to.To(nil, R)
	return L, R
}
