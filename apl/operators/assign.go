package operators

import (
	"fmt"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  "←",
		Domain:  MonadicOp(nil),
		doc:     "assign, variable assignment, specification, copula",
		derived: assign,
	})
}

func assign(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		as, ok := f.(*apl.Assignment)
		if ok == false {
			return nil, fmt.Errorf("cannot assign to %T", f)
		}
		if L != nil {
			return nil, fmt.Errorf("assign cannot be called dyadically")
		}

		if as.Identifiers != nil {
			if as.Indexes != nil {
				return nil, fmt.Errorf("vector and indexed assignment cannot exist simulaneously")
			}
			return assignVector(a, as.Identifiers, R, as.Modifier)
		}

		return R, assignScalar(a, as.Identifier, as.Indexes, as.Modifier, R)
	}
	return function(derived)
}

// AssignVector does a vector assignment from R to the given names.
// A modifier function may be applied.
func assignVector(a *apl.Apl, names []string, R apl.Value, mod apl.Value) (apl.Value, error) {
	var ar apl.Array
	if v, ok := R.(apl.Array); ok {
		ar = v
	} else {
		ar = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{R}}
	}

	var scalar apl.Value
	if s := ar.Shape(); len(s) != 1 {
		return nil, fmt.Errorf("vector assignment: rank of right argument must be 1")
	} else if s[0] != 1 && s[0] != len(names) {
		return nil, fmt.Errorf("vector assignment is non-conformant")
	} else if s[0] == 1 {
		if v, err := ar.At(0); err != nil {
			return nil, err
		} else {
			scalar = v
		}
	}

	var err error
	for i, name := range names {
		var v apl.Value
		if scalar != nil {
			v = scalar
		} else {
			v, err = ar.At(i)
			if err != nil {
				return nil, err
			}
		}
		err = assignScalar(a, name, nil, mod, v)
		if err != nil {
			return nil, err
		}
	}

	return R, nil
}

// AssignScalar assigns to a single numeric variable.
// If indexes is non-nil, it must be an IndexArray for indexed assignment.
// Mod may be a dyadic modifying function.
func assignScalar(a *apl.Apl, name string, indexes apl.Value, mod apl.Value, R apl.Value) error {
	if mod == nil && indexes == nil {
		return a.Assign(name, R)
	}

	w, env := a.LookupEnv(name)
	if w == nil {
		return fmt.Errorf("modified/indexed assignment to non-existing variable %s", name)
	}

	var f apl.Function
	if mod != nil {
		if fn, ok := mod.(apl.Function); ok == false {
			return fmt.Errorf("modified assignment needs a function: %T", mod)
		} else {
			f = fn
		}
	}

	// Modified assignment without indexing.
	if indexes == nil {
		if v, err := f.Call(a, w, R); err != nil {
			return err
		} else {
			return a.AssignEnv(name, v, env)
		}
	}

	idx, ok := indexes.(apl.IndexArray)
	if ok == false {
		to := ToIndexArray(nil)
		if v, ok := to.To(a, indexes); ok == false {
			return fmt.Errorf("indexed assignment could not convert to IndexArray: %T", indexes)
		} else {
			idx = v.(apl.IndexArray)
		}
	}

	ar, ok := w.(apl.ArraySetter)
	if ok == false {
		return fmt.Errorf("variable %s is no settable array: %T", name, w)
	}

	// Try to keep the original array type, upgrade only if needed.
	upgrade := func() {
		ga := apl.MixedArray{Dims: apl.CopyShape(ar)}
		ga.Values = make([]apl.Value, apl.ArraySize(ga))
		for i := range ga.Values {
			v, err := ar.At(i)
			if err != nil {
				return
			}
			ga.Values[i] = v
		}
		ar = ga
	}

	// modAssign assigns ar at index i with v possibly modified by f.
	modAssign := func(i int, v apl.Value) error {
		if i == -1 {
			// Index -1 is used by some indexed assignments to mark skipps.
			// E.g. replicate and compress / and \
			return nil
		}
		if f == nil {
			if err := ar.Set(i, v); err == nil {
				return nil
			}
			upgrade()
			return ar.Set(i, v)
		}
		u, err := ar.At(i)
		if err != nil {
			return err
		}
		v, err = f.Call(a, u, v)
		if err != nil {
			return err
		}
		if err := ar.Set(i, v); err == nil {
			return nil
		}
		upgrade()
		return ar.Set(i, v)
	}

	var src apl.Array
	var scalar apl.Value
	if av, ok := R.(apl.Array); ok {
		src = av
		if apl.ArraySize(av) == 1 {
			scalar, _ = av.At(0)
		}
	} else {
		scalar = R
	}
	if scalar != nil {
		// Scalar or 1-element assignment.
		for _, d := range idx.Ints {
			if err := modAssign(d, scalar); err != nil {
				return err
			}
		}
	} else {

		// Shapes must conform. Single element axis are collapsed.
		collapse := func(s []int) []int {
			n := 0
			for _, i := range s {
				if i == 1 {
					n++
				}
			}
			if n == 0 {
				return s
			}
			r := make([]int, len(s)-n)
			k := 0
			for _, i := range s {
				if i != 1 {
					r[k] = i
					k++
				}
			}
			return r
		}
		ds := collapse(idx.Shape())
		ss := collapse(src.Shape())
		if len(ds) != len(ss) {
			return fmt.Errorf("indexed assignment: arrays have different rank: %d != %d", len(ds), len(ss))
		}
		for i := range ds {
			if ss[i] != ds[i] {
				return fmt.Errorf("indexed assignment: arrays are not conforming: %v != %v", ss, ds)
			}
		}
		for i, d := range idx.Ints {
			if v, err := src.At(i); err != nil {
				return err
			} else {
				if err := modAssign(d, v); err != nil {
					return err
				}
			}
		}
	}
	return a.AssignEnv(name, ar, env)
}
