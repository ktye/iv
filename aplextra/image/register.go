package image

import "github.com/ktye/iv/apl"

func Register(a *apl.Apl) {
	// This is just a dummy.
	a.RegisterPrimitive("⍞", handle(blue))
}

// Handle is both, a function but also implements a PrimitiveHandler.
// If provides a methods that calls the function value itself.
// This is used to cast a function to a handler.
type handle func(*apl.Apl, apl.Value, apl.Value) (bool, apl.Value, error)

func (h handle) HandlePrimitive(a *apl.Apl, l apl.Value, r apl.Value) (bool, apl.Value, error) {
	return h(a, l, r)
}
