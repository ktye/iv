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
			return assignVector(a, as.Identifiers, R, as.Modifier)
		}

		if as.Modifier != nil {
			return nil, fmt.Errorf("TODO: modified assignment")
		}

		return R, a.AssignIndexed(as.Identifier, as.Indexes, R)
	}
	return function(derived)
}

// AssignVector does a vector assignment from R to the given names.
// A modifier function may be applied.
func assignVector(a *apl.Apl, names []string, R apl.Value, mod apl.Value) (apl.Value, error) {
	if mod != nil {
		return nil, fmt.Errorf("TODO modified vector assignment")
	}

	var ar apl.Array
	if v, ok := R.(apl.Array); ok {
		ar = v
	} else {
		ar = apl.GeneralArray{Dims: []int{1}, Values: []apl.Value{R}}
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
		err = assignModified(a, name, v, mod)
		if err != nil {
			return nil, err
		}
	}

	return R, nil
}

func assignModified(a *apl.Apl, name string, R apl.Value, mod apl.Value) error {
	if mod != nil {
		return fmt.Errorf("TODO modified vector assignment")
	}
	return a.Assign(name, R)
}
