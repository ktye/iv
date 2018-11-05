package image

import "github.com/ktye/iv/apl"

func Register(a *apl.Apl) {
	// This is just a dummy.
	a.RegisterPrimitive("⍞", blue)
}
