package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:    "⍂",
		Domain:    DyadicOp(nil),
		doc:       "axis specification",
		derived:   axis,
		selection: selection(axis),
	})
}

func axis(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		if p, ok := f.(apl.Function); ok {
			return p.Call(a, L, apl.Axis{R: R, A: g})
		} else {
			return nil, fmt.Errorf("axis: expected primitive on the left: %T", f)
		}
	}
	return function(derived)
}

// splitAxis returns ax.R and converts ax.A to []int taking account of index origin.
// It R is not an axis it returns R and nil.
func splitAxis(a *apl.Apl, R apl.Value) (apl.Value, []int, error) {
	ax, ok := R.(apl.Axis)
	if ok == false {
		return R, nil, nil
	}
	if _, ok := ax.A.(apl.EmptyArray); ok {
		return ax.R, nil, nil
	}
	to := ToIndexArray(nil)
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
