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
		v, ok := f.(apl.Identifier)
		if ok == false {
			return nil, fmt.Errorf("cannot assign to %T", f)
		}
		if L != nil {
			return nil, fmt.Errorf("assign cannot be called dyadically")
		}

		return R, a.Assign(string(v), R)
	}
	return function(derived)
}

// Index returns the indexes of the index specification applied to the array.
// The indexes in the IndexArray have origin 0.
func Index(a *apl.Apl, spec apl.IdxSpec, A apl.Array) (apl.IndexArray, error) {
	shape := A.Shape()
	if len(shape) != len(spec) {
		return apl.IndexArray{}, fmt.Errorf("indexing: Array and index specification have different rank")
	}
	return apl.IndexArray{}, fmt.Errorf("TODO operators.Index")
}
