package funcs

import "github.com/ktye/iv/apl"

// Register adds all functions defined in this package to the interpreter.
func Register(a *apl.Apl) {
	for _, p := range primitives {
		a.RegisterPrimitive(p.p, p.h)
	}
	for _, d := range doc {
		a.RegisterDoc(d[0], d[1])
	}
}

// Handle is both, a function but also implements a PrimitiveHandler.
// If provides a methods that calls the function value itself.
// This is used to cast a function to a handler.
type handle func(*apl.Apl, apl.Value, apl.Value) (bool, apl.Value, error)

func (h handle) HandlePrimitive(a *apl.Apl, l apl.Value, r apl.Value) (bool, apl.Value, error) {
	return h(a, l, r)
}

// Both wraps a monadic and a dyadic handle.
// The monadic handle is used if the left argument is nil,
// otherwise the dyadic is used.
func both(monadic, dyadic handle) apl.PrimitiveHandler {
	return defaultHandler{monadic, dyadic}
}

type defaultHandler struct {
	monadic, dyadic handle
}

func (d defaultHandler) HandlePrimitive(a *apl.Apl, l apl.Value, r apl.Value) (bool, apl.Value, error) {
	if l == nil {
		return d.monadic(a, l, r)
	}
	return d.dyadic(a, l, r)
}

// ArrayWrap takes a dyadic elementry handler that does not know about arrays,
// and wraps it with arrayApply.
func arrayWrap(h handle) handle {
	return func(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
		return arrayApply(a, l, r, h)
	}
}

// ArrayApply applies the given dyadic handle to l and r.
// If one of them is an array, the dyadic function is applied to each element together with the other value.
// If both are arrays, the function is applied elementwise, if the shape agrees.
// If the shape does not agree, it returns an error but accepts the handler.
// If both functions are numeric scalar, they are promoted to the same type.
func arrayApply(a *apl.Apl, l, r apl.Value, h handle) (bool, apl.Value, error) {
	// If at least one is an array, call the arrays ApplyDyadic method.
	if v, ok := l.(apl.Array); ok {
		u, err := v.ApplyDyadic(a, r, true, h)
		return true, u, err
	} else if v, ok := r.(apl.Array); ok {
		u, err := v.ApplyDyadic(a, l, false, h)
		return true, u, err
	}

	// Promote both scalars to the same type.
	l, r, err := apl.SameNumericTypes(l, r)
	if err != nil {
		return true, nil, err
	}
	return h(a, l, r)
}

type primitive struct {
	p apl.Primitive
	h apl.PrimitiveHandler
}

var primitives []primitive

func register(p apl.Primitive, h apl.PrimitiveHandler) {
	primitives = append(primitives, primitive{p, h})
}

var doc [][2]string

func addDoc(key, text string) {
	doc = append(doc, [2]string{key, text})
}
