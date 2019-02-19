// Package p is a plot package.
package p

import (
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

func Register(a *apl.Apl, name string) {
	if name == "" {
		name = "p"
	}
	pkg := map[string]apl.Value{
		"p": apl.ToFunction(plotf),
	}
	a.RegisterPackage(name, pkg)

	a.RegisterPrimitive("⌼", apl.ToHandler(
		plotf,
		domain.Monadic(domain.ToArray(nil)),
		"plot",
	))
	a.RegisterPrimitive("⌼", apl.ToHandler(
		plotf,
		domain.Dyadic(domain.Split(domain.ToArray(nil), domain.ToArray(nil))),
		"plot",
	))
}
