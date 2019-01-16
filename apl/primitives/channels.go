package primitives

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

// primitive < is defined in compare.go

// channelSource sends any value R over a channel.
// <[axis] R
func channelSource(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	r, ax, err := splitAxis(a, R)
	if err != nil {
		return nil, fmt.Errorf("channel send: %s", err)
	}
	if len(ax) > 1 {
		return nil, fmt.Errorf("channel send: axis must be scalar")
	}

	if _, ok := r.(apl.Channel); ok {
		return nil, fmt.Errorf("channel send: right argument is a channel")
	}

	n := 1
	if ax != nil {
		n = ax[0]
	}
	n += a.Origin // splitAxis substracts the origin.

	c := apl.NewChannel()
	if n == 0 {
		// Send only once, but do not close any channels.
		go func(v apl.Value) {
			c[0] <- v
		}(r)
		return c, nil
	}

	// Send v n times. If n is negative send until c[1] is closed.
	go func(v apl.Value, n int) {
		defer close(c[0])
		i := 0
		for {
			select {
			case _, ok := <-c[1]:
				if ok == false {
					return
				}
			case c[0] <- v:
				i++
				if n > 0 && i >= n {
					return
				}
			}
		}
	}(r, n)
	return c, nil
}
