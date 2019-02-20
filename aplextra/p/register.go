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
		"p":           apl.ToFunction(plot4),
		"dark":        apl.ToFunction(setDark),
		"transparent": apl.ToFunction(setTransparent),
		"colors":      apl.ToFunction(setColors),
		"size":        apl.ToFunction(setSize),
		"fontsizes":   apl.ToFunction(setFontSizes),
	}
	a.RegisterPackage(name, pkg)

	a.RegisterPrimitive("⌼", apl.ToHandler(
		plot4,
		domain.Monadic(domain.ToArray(nil)),
		"plot",
	))
	a.RegisterPrimitive("⌼", apl.ToHandler(
		plot4,
		domain.Dyadic(domain.Split(domain.ToArray(nil), domain.ToArray(nil))),
		"plot",
	))
	setFontSizes(a, nil, apl.IntArray{Dims: []int{2}, Ints: []int{18, 12}})
}
