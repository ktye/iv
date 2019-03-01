package domain

import "github.com/ktye/iv/apl"

// IsChannel tests if the value is a channel.
func IsChannel(child SingleDomain) SingleDomain {
	return channel{child}
}

type channel struct {
	child SingleDomain
}

func (c channel) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	if _, ok := V.(apl.Channel); ok {
		return propagate(a, V, c.child)
	}
	return V, false
}

func (c channel) String(f apl.Format) string {
	name := "channel"
	if c.child == nil {
		return name
	}
	return name + " " + c.child.String(f)
}
