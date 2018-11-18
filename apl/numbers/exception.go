package numbers

import (
	"math"
	"math/cmplx"

	"github.com/ktye/iv/apl"
)

const (
	NaN    exception = "NaN"
	Inf    exception = "∞"
	NegInf exception = "¯∞"
)

type exception string

func (e exception) String(a *apl.Apl) string {
	return string(e)
}

func isException(n apl.Number) (exception, bool) {
	if f, ok := n.(Float); ok {
		if math.IsNaN(float64(f)) {
			return NaN, true
		}
		if math.IsInf(float64(f), 1) {
			return Inf, true
		}
		if math.IsInf(float64(f), -1) {
			return NegInf, true
		}
	}
	if c, ok := n.(Complex); ok {
		if cmplx.IsNaN(complex128(c)) {
			return NaN, true
		}
		if cmplx.IsInf(complex128(c)) {
			return Inf, true
		}
	}
	return "", false
}
