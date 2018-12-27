// Package strings provides go string functions
package strings

import (
	"reflect"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/xgo"
)

func Register(a *apl.Apl) {
	pkg := map[string]apl.Value{
		"compare":        xgo.Function{Name: "Compare", Fn: reflect.ValueOf(strings.Compare)},
		"contains":       xgo.Function{Name: "Contains", Fn: reflect.ValueOf(strings.Contains)},
		"containsany":    xgo.Function{Name: "ContainsAny", Fn: reflect.ValueOf(strings.ContainsAny)},
		"containsrune":   xgo.Function{Name: "ContainsRune", Fn: reflect.ValueOf(strings.ContainsRune)},
		"count":          xgo.Function{Name: "Count", Fn: reflect.ValueOf(strings.Count)},
		"equalfold":      xgo.Function{Name: "EqualFold", Fn: reflect.ValueOf(strings.EqualFold)},
		"fields":         xgo.Function{Name: "Fields", Fn: reflect.ValueOf(strings.Fields)},
		"fieldsfunc":     xgo.Function{Name: "FieldsFunc", Fn: reflect.ValueOf(strings.FieldsFunc)},
		"hasprefix":      xgo.Function{Name: "HasPrefix", Fn: reflect.ValueOf(strings.HasPrefix)},
		"hassuffix":      xgo.Function{Name: "HasSuffix", Fn: reflect.ValueOf(strings.HasSuffix)},
		"index":          xgo.Function{Name: "Index", Fn: reflect.ValueOf(strings.Index)},
		"indexany":       xgo.Function{Name: "IndexAny", Fn: reflect.ValueOf(strings.IndexAny)},
		"indexbyte":      xgo.Function{Name: "IndexByte", Fn: reflect.ValueOf(strings.IndexByte)},
		"indexfunc":      xgo.Function{Name: "IndexFunc", Fn: reflect.ValueOf(strings.IndexFunc)},
		"indexrune":      xgo.Function{Name: "IndexRune", Fn: reflect.ValueOf(strings.IndexRune)},
		"join":           xgo.Function{Name: "Join", Fn: reflect.ValueOf(strings.Join)},
		"lastindex":      xgo.Function{Name: "LastIndex", Fn: reflect.ValueOf(strings.LastIndex)},
		"lastindexany":   xgo.Function{Name: "LastIndexAny", Fn: reflect.ValueOf(strings.LastIndexAny)},
		"lastindexbyte":  xgo.Function{Name: "LastIndexByte", Fn: reflect.ValueOf(strings.LastIndexByte)},
		"repeat":         xgo.Function{Name: "Repeat", Fn: reflect.ValueOf(strings.Repeat)},
		"replace":        xgo.Function{Name: "Replace", Fn: reflect.ValueOf(strings.Replace)},
		"split":          xgo.Function{Name: "Split", Fn: reflect.ValueOf(strings.Split)},
		"splitafter":     xgo.Function{Name: "SplitAfter", Fn: reflect.ValueOf(strings.SplitAfter)},
		"splitaftern":    xgo.Function{Name: "SplitAfterN", Fn: reflect.ValueOf(strings.SplitAfterN)},
		"splitn":         xgo.Function{Name: "SplitN", Fn: reflect.ValueOf(strings.SplitN)},
		"title":          xgo.Function{Name: "Title", Fn: reflect.ValueOf(strings.Title)},
		"tolower":        xgo.Function{Name: "ToLower", Fn: reflect.ValueOf(strings.ToLower)},
		"tolowerspecial": xgo.Function{Name: "ToLowerSpecial", Fn: reflect.ValueOf(strings.ToLowerSpecial)},
		"totitle":        xgo.Function{Name: "ToTitle", Fn: reflect.ValueOf(strings.ToTitle)},
		"totitlespecial": xgo.Function{Name: "ToTitleSpecial", Fn: reflect.ValueOf(strings.ToTitleSpecial)},
		"toupper":        xgo.Function{Name: "ToUpper", Fn: reflect.ValueOf(strings.ToUpper)},
		"toupperspecial": xgo.Function{Name: "ToUpperSpecial", Fn: reflect.ValueOf(strings.ToUpperSpecial)},
		"trim":           xgo.Function{Name: "Trim", Fn: reflect.ValueOf(strings.Trim)},
		"trimleft":       xgo.Function{Name: "TrimLeft", Fn: reflect.ValueOf(strings.TrimLeft)},
		"trimprefix":     xgo.Function{Name: "TrimPrefix", Fn: reflect.ValueOf(strings.TrimPrefix)},
		"trimright":      xgo.Function{Name: "TrimRight", Fn: reflect.ValueOf(strings.TrimRight)},
		"trimspace":      xgo.Function{Name: "TrimSpace", Fn: reflect.ValueOf(strings.TrimSpace)},
		"trimsuffix":     xgo.Function{Name: "TrimSuffix", Fn: reflect.ValueOf(strings.TrimSuffix)},
	}
	a.RegisterPackage("s", pkg)
}
