package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

// Selection returns a general selection function for selective assignment.
// It creates an index array of the same shape of R and applies f to it.
func selection(f func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error)) func(*apl.Apl, apl.Value, apl.Value) (apl.IndexArray, error) {
	return func(a *apl.Apl, L apl.Value, R apl.Value) (apl.IndexArray, error) {
		// Create an index array with the shape of R.
		var ai apl.IndexArray
		ar, ok := R.(apl.Array)
		if ok == false {
			return ai, fmt.Errorf("cannot select from %T", R)
		}
		ai.Dims = apl.CopyShape(ar)
		ai.Ints = make([]int, apl.ArraySize(ai))
		for i := range ai.Ints {
			ai.Ints[i] = i
		}

		// Apply the selection function to it.
		v, err := f(a, L, ai)
		if err != nil {
			return ai, err
		}

		to := domain.ToIndexArray(nil)
		if av, ok := to.To(a, v); ok == false {
			return ai, fmt.Errorf("could not convert selection to index array: %T", v)
		} else {
			return av.(apl.IndexArray), nil
		}
	}
}
