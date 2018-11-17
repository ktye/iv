package numbers

import (
	"reflect"

	"github.com/ktye/iv/apl"
)

func getformat(a *apl.Apl, num apl.Number, def string) (string, bool) {
	if a == nil {
		return def, false
	}
	if n, ok := a.Tower.Numbers[reflect.TypeOf(num)]; ok == false {
		return def, false
	} else {
		f := n.Format
		if f == "" {
			return def, false
		}
		if f[0] == '-' {
			return f[1:], true
		}
		return f, false
	}
}
