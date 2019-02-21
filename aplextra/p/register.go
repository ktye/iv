// Package p is a plot package.
package p

import (
	"reflect"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
)

func Register(a *apl.Apl, name string) {
	if name == "" {
		name = "p"
	}
	pkg := map[string]apl.Value{
		"p":           apl.ToFunction(plot1),
		"l":           apl.ToFunction(newLine),
		"dark":        apl.ToFunction(setDark),
		"transparent": apl.ToFunction(setTransparent),
		"colors":      apl.ToFunction(setColors),
		"size":        apl.ToFunction(setSize),
		"fontsizes":   apl.ToFunction(setFontSizes),
		"gui":         apl.ToFunction(setGui),
	}
	a.RegisterPackage(name, pkg)

	a.RegisterPrimitive("⌼", apl.ToHandler(
		plot1,
		domain.Monadic(domain.ToArray(nil)),
		"plot",
	))
	a.RegisterPrimitive("⌼", apl.ToHandler(
		plot1,
		domain.Dyadic(domain.Split(domain.ToArray(nil), domain.ToArray(nil))),
		"plot",
	))
	a.RegisterPrimitive("⌼", apl.ToHandler(
		plotToImage,
		domain.Monadic(
			domain.Or(
				domain.IsType(reflect.TypeOf(Plot{}), nil),
				domain.IsType(reflect.TypeOf(PlotArray{}), nil),
			),
		),
		"plot to image",
	))
	a.RegisterPrimitive("⌼", apl.ToHandler(
		plotToImage,
		domain.Dyadic(
			domain.Split(domain.ToIndexArray(nil),
				domain.Or(
					domain.IsType(reflect.TypeOf(Plot{}), nil),
					domain.IsType(reflect.TypeOf(PlotArray{}), nil),
				),
			),
		),
		"plot to image",
	))
	a.RegisterPrimitive("+", apl.ToHandler(
		plotAddLine,
		domain.Dyadic(domain.Split(
			domain.IsType(reflect.TypeOf(Line{}), nil),
			domain.IsType(reflect.TypeOf(Plot{}), nil),
		)),
		"plot add line",
	))
	setFontSizes(a, nil, apl.IntArray{Dims: []int{2}, Ints: []int{18, 12}})
}
