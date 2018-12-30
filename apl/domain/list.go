package domain

import "github.com/ktye/iv/apl"

// IsList tests if the value is a list.
func IsList(child SingleDomain) SingleDomain {
	return list{child, false}
}

func ToList(child SingleDomain) SingleDomain {
	return list{child, true}
}

type list struct {
	child SingleDomain
	conv  bool
}

func (v list) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	_, ok := V.(apl.List)
	if v.conv == false && ok == false {
		return V, false
	} else if ok == true {
		return propagate(a, V, v.child)
	}

	// enlist
	return propagate(a, apl.List{N: 1, L: []apl.Value{V}}, v.child)
}
func (v list) String(a *apl.Apl) string {
	name := "list"
	if v.conv {
		name = "tolist"
	}
	if v.child == nil {
		return name
	}
	return name + " " + v.child.String(a)
}
