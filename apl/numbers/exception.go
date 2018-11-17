package numbers

import (
	"math"
	"math/cmplx"

	"github.com/ktye/iv/apl"
)

type exception string

func (e exception) String(a *apl.Apl) string {
	return string(e)
}

func isException(n apl.Number) (exception, bool) {
	if f, ok := n.(Float); ok {
		if math.IsNaN(float64(f)) {
			return exception("NaN"), true
		}
		if math.IsInf(float64(f), 1) {
			return exception("∞"), true
		}
		if math.IsInf(float64(f), -1) {
			return exception("-∞"), true
		}
	}
	if c, ok := n.(Complex); ok {
		if cmplx.IsNaN(complex128(c)) {
			return exception("NaN"), true
		}
		if cmplx.IsInf(complex128(c)) {
			return exception("∞"), true
		}
	}
	return "", false
}
