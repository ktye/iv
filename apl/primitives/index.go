package primitives

import (
	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
	"github.com/ktye/iv/apl/operators"
)

func init() {
	// An expression such as A[1;2;] is translated by the parser to
	//	[1;2;] ⌷ A
	// ⌷ cannot be used directly, as an index specification is converted by the parser.
	register(primitive{
		symbol: "⌷",
		doc:    "index, []",
		Domain: Dyadic(Split(indexSpec{}, ToArray(nil))),
		fn:     index,
	})
}

// indexSpec is the domain type for an index specification.
type indexSpec struct{}

func (i indexSpec) To(a *apl.Apl, v apl.Value) (apl.Value, bool) {
	if _, ok := v.(apl.IdxSpec); ok {
		return v, true
	}
	return v, false
}
func (i indexSpec) String(a *apl.Apl) string {
	return "[index specification]"
}

func index(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	spec := L.(apl.IdxSpec)
	ar := R.(apl.Array)

	idx, err := operators.Index(a, spec, ar)
	if err != nil {
		return nil, err
	}

	res := apl.GeneralArray{
		Dims:   apl.CopyShape(idx),
		Values: make([]apl.Value, apl.ArraySize(idx)),
	}
	for i, n := range idx.Ints {
		v, err := ar.At(n)
		if err != nil {
			return nil, err
		}
		res.Values[i] = v // TODO copy?
	}
	return res, nil
}
