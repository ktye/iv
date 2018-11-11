package operators

import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(operator{
		symbol:  ".",
		Domain:  DyadicOp(Split(Function(nil), Function(nil))),
		doc:     "inner product",
		derived: innerproduct,
	})
	register(operator{
		symbol:  ".",
		Domain:  DyadicOp(Split(primitive("+"), primitive("×"))),
		doc:     "scalar product",
		derived: scalarproduct,
	})
	register(operator{
		symbol:  ".",
		Domain:  DyadicOp(Split(primitive("∘"), Function(nil))),
		doc:     "outer product",
		derived: outer,
	})

}

func innerproduct(a *apl.Apl, f, g apl.Value) apl.Function {
	derived := func(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
		f := f.(apl.Function)
		g := g.(apl.Function)
		return inner(a, l, r, f, g)
	}
	return function(derived)
}

func scalarproduct(a *apl.Apl, f, g apl.Value) apl.Function {
	df := f.(apl.Primitive) // +
	dg := g.(apl.Primitive) // ×
	derived := func(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
		// Special case for a scalar product.
		// If both have the same type and implement a ScalarProducter, use the interface.
		if reflect.TypeOf(l) == reflect.TypeOf(r) {
			if sc, ok := l.(scalarProducter); ok {
				v, err := sc.ScalarProduct(r)
				return v, err
			}
		}
		return inner(a, l, r, df, dg)
	}
	return function(derived)
}

func inner(a *apl.Apl, l, r apl.Value, f, g apl.Function) (apl.Value, error) {
	al, lok := l.(apl.Array)
	ar, rok := r.(apl.Array)

	if lok == false && rok == false {
		// Both are scalars, compute l g r.
		v, err := g.Call(a, al, ar)
		return v, err
	}

	// If one is a scalar, convert it to a vector.
	if lok == false {
		rs := ar.Shape()
		if rs == nil || rs[0] == 0 {
			// TODO fill function?
			return nil, fmt.Errorf("inner: empty rhs array")
		}
		u := apl.GeneralArray{Dims: []int{rs[0]}}
		v := make([]apl.Value, rs[0])
		for i := range v {
			v[i] = l
		}
		u.Values = v
		al = u
	} else if rok == false {
		ls := al.Shape()
		if ls == nil || ls[0] == 0 {
			return nil, fmt.Errorf("inner: empty lhs array")
		}
		u := apl.GeneralArray{Dims: []int{ls[len(ls)-1]}}
		v := make([]apl.Value, ls[0])
		for i := range v {
			v[i] = r
		}
		u.Values = v
		ar = u
	}

	// The result is a new array with a shape of both arrays combined, without the inner dimension.
	ls := al.Shape()
	rs := ar.Shape()
	if len(ls) == 0 || len(rs) == 0 {
		return nil, fmt.Errorf("inner: empty array")
	}
	inner := ls[len(ls)-1]
	if inner != rs[0] {
		return nil, fmt.Errorf("inner dimensions must agree")
	}

	// If both arrays are vectors, compute a scalar.
	if len(ls) == 1 && len(rs) == 1 {
		var v apl.Value
		for k := inner - 1; k >= 0; k-- {
			lval, err := al.At(k)
			if err != nil {
				return nil, err
			}
			rval, err := ar.At(k)
			if err != nil {
				return nil, err
			}
			if u, err := g.Call(a, lval, rval); err != nil {
				return nil, err
			} else if k == inner-1 {
				v = u
			} else {
				if u, err := f.Call(a, u, v); err != nil {
					return nil, err
				} else {
					v = u
				}
			}
		}
		return v, nil
	}

	shape := make([]int, len(ls)+len(rs)-2)
	copy(shape, ls[:len(ls)-1])
	copy(shape[len(ls)-1:], rs[1:])
	result := apl.GeneralArray{Dims: shape}
	result.Values = make([]apl.Value, apl.ArraySize(result))

	// Iterate of all elements of the resulting array.
	idx := make([]int, len(shape))
	split := len(ls) - 1
	lidx := make([]int, len(ls))
	ridx := make([]int, len(rs))
	for i := range result.Values {
		if err := apl.ArrayIndexes(shape, idx, i); err != nil {
			return nil, err
		}
		// Split the indexes in idx into the original indexes of both arrays.
		copy(lidx, idx[:split])     // The last index is open.
		copy(ridx[1:], idx[split:]) // The first index is open.
		var v apl.Value
		for k := inner - 1; k >= 0; k-- {
			lidx[len(lidx)-1] = k
			ridx[0] = k
			lval, err := apl.ArrayAt(al, lidx)
			if err != nil {
				return nil, err
			}
			rval, err := apl.ArrayAt(ar, ridx)
			if err != nil {
				return nil, err
			}
			if u, err := g.Call(a, lval, rval); err != nil {
				return nil, err
			} else if k == inner-1 {
				v = u
			} else {
				if u, err := f.Call(a, u, v); err != nil {
					return nil, err
				} else {
					v = u
				}
			}
		}
		result.Values[i] = v
	}
	return result, nil
}

func outer(a *apl.Apl, LO, RO apl.Value) apl.Function {
	return function(func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		return nil, fmt.Errorf("TODO: outer product")
	})
}

// A scalarProducter implements a ScalarProduct which receives an argument of the same type.
// This can be implemented by matrix multiplication for special types.
type scalarProducter interface {
	ScalarProduct(interface{}) (apl.Value, error)
}
