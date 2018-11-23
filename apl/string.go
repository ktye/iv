package apl

type String string

// String formats s with %s by default.
// The format can be changed in Format.String.
func (s String) String(a *Apl) string {
	return string(s)
}

func (s String) Eval(a *Apl) (Value, error) {
	return s, nil
}

// Less implements primitives.lesser to be used for comparison and sorting.
func (s String) Less(r Value) (Bool, bool) {
	b, ok := r.(String)
	if ok == false {
		return false, false
	}
	return s < b, true
}
