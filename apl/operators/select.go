package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

// selection returns a derived selection function, given an operator function that creates a derived function.
func selection(op func(*apl.Apl, apl.Value, apl.Value) apl.Function) func(*apl.Apl, apl.Value, apl.Value, apl.Value, apl.Value) (apl.IntArray, error) {
	derived := func(a *apl.Apl, L, LO, RO, R apl.Value) (apl.IntArray, error) {

		// Create the derived function.
		df := op(a, LO, RO)

		// Create an index array with the shape of R.
		var ai apl.IntArray
		ar, ok := R.(apl.Array)
		if ok == false {
			return ai, fmt.Errorf("cannot select from %T", R)
		}
		ai.Dims = apl.CopyShape(ar)
		ai.Ints = make([]int, apl.Prod(ai.Dims))
		for i := range ai.Ints {
			ai.Ints[i] = i + 1
		}

		// Apply the selection function to the index array.
		v, err := df.Call(a, L, ai)
		if err != nil {
			return ai, err
		}

		to := ToIndexArray(nil)
		if av, ok := to.To(a, v); ok == false {
			return ai, fmt.Errorf("could not convert selection to index array: %T", v)
		} else {
			ai = av.(apl.IntArray)
			for i := range ai.Ints {
				ai.Ints[i]--
			}
			return ai, nil
		}
	}
	return derived
}

// selectSimple returns a general selection function for selective assignment.
// It creates an index array of the same shape of R and applies f to it.
// It is used by replicate and expand which behave like primitive functions instead of operators.
// They take only 2 arguments.
func selectSimple(f func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error)) func(*apl.Apl, apl.Value, apl.Value, apl.Value, apl.Value) (apl.IntArray, error) {
	return func(a *apl.Apl, dummyL, L apl.Value, dummyRO, R apl.Value) (apl.IntArray, error) {

		// Create an index array with the shape of R.
		var ai apl.IntArray
		ar, ok := R.(apl.Array)
		if ok == false {
			return ai, fmt.Errorf("cannot select from %T", R)
		}
		ai.Dims = apl.CopyShape(ar)
		ai.Ints = make([]int, apl.Prod(ai.Dims))
		for i := range ai.Ints {
			ai.Ints[i] = i + 1
		}

		// Apply the selection function to it.
		v, err := f(a, L, ai)
		if err != nil {
			return ai, err
		}

		to := ToIndexArray(nil)
		if av, ok := to.To(a, v); ok == false {
			return ai, fmt.Errorf("could not convert selection to index array: %T", v)
		} else {
			// Fill elements will be reported as ¯1, which the assignment should ignore.
			ai = av.(apl.IntArray)
			for i := range ai.Ints {
				ai.Ints[i]--
			}
			return ai, nil
		}
	}
}
