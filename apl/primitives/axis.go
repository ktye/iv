package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

// splitAxis returns ax.R and converts ax.A to []int taking account of index origin.
// It R is not an axis it returns R and nil.
func splitAxis(a *apl.Apl, R apl.Value) (apl.Value, []int, error) {
	ax, ok := R.(apl.Axis)
	if ok == false {
		return R, nil, nil
	}
	to := domain.ToIndexArray(nil)
	X, ok := to.To(a, ax.A)
	if ok == false {
		return nil, nil, fmt.Errorf("axis is not an index array")
	}
	ar := X.(apl.IndexArray)
	shape := ar.Shape()
	if len(shape) != 1 {
		return nil, nil, fmt.Errorf("axis has wrong shape: %d", len(shape))
	}
	x := make([]int, len(ar.Ints))
	for i, n := range ar.Ints {
		x[i] = n - a.Origin
	}
	return ax.R, x, nil
}
