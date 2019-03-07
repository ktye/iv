package primitives

import (
	"strings"

	"github.com/ktye/iv/apl"
	. "github.com/ktye/iv/apl/domain"
)

func init() {
	register(primitive{
		symbol: "⊂",
		doc:    "enclose, string catenation",
		Domain: Monadic(strvec{}),
		fn:     strcat,
	})
	register(primitive{
		symbol: "⊂",
		doc:    "join strings",
		Domain: Dyadic(Split(IsString(nil), strvec{})),
		fn:     strjoin,
	})
	register(primitive{
		symbol: "⊃",
		doc:    "split runes",
		Domain: Monadic(IsString(nil)),
		fn:     runesplit,
	})
	register(primitive{
		symbol: "⊃",
		doc:    "first",
		Domain: Monadic(IsList(nil)),
		fn:     first,
	})
	register(primitive{
		symbol: "⊃",
		doc:    "split string",
		Domain: Dyadic(Split(IsString(nil), IsString(nil))),
		fn:     strsplit,
	})
}

// strvec accepts an array if all elements are strings.
// The result is a string vector.
type strvec struct{}

func (s strvec) To(a *apl.Apl, v apl.Value) (apl.Value, bool) {
	ar, ok := v.(apl.Array)
	if ok == false {
		return v, false
	}
	n := apl.ArraySize(ar)
	vec := apl.StringArray{Dims: []int{n}, Strings: make([]string, n)}
	for i := 0; i < n; i++ {
		if s, ok := ar.At(i).(apl.String); ok {
			vec.Strings[i] = string(s)
		} else {
			return v, false
		}
	}
	return vec, true
}
func (s strvec) String(f apl.Format) string {
	return "array of strings"
}

func strcat(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	ar := R.(apl.StringArray)
	return apl.String(strings.Join(ar.Strings, "")), nil
}

func strjoin(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	v := R.(apl.StringArray)
	return apl.String(strings.Join(v.Strings, string(L.(apl.String)))), nil
}

func runesplit(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	r := []rune(string(R.(apl.String)))
	v := make([]string, len(r))
	for i, s := range r {
		v[i] = string(s)
	}
	return apl.StringArray{Dims: []int{len(v)}, Strings: v}, nil
}

// split a string to an string vector
// If L is a string, use it as the separator for strings.Split
// If L is "", use strings.Field
func strsplit(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	l := L.(apl.String)
	r := R.(apl.String)
	var v []string
	if l == "" {
		v = strings.Fields(string(r))
	} else {
		v = strings.Split(string(r), string(l))
	}
	return apl.StringArray{Dims: []int{len(v)}, Strings: v}, nil
}

func first(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	r := R.(apl.List)
	if len(r) == 0 {
		return apl.EmptyArray{}, nil
	}
	return r[0], nil
}
