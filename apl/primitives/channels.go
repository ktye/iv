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
		// go sendInitial(c, r)
		// return c, nil
		return nil, fmt.Errorf("TODO: sendInitial")
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

/*
// sendInitial sends the initial value to c[0], then it read repeatedly from c[0] and
// send the value back to c[0] until c[1] is closed.
func sendInitial(c apl.Channel, initial apl.Value) {
	c[0] <- initial
	var v apl.Value
	var ok bool
	for {
	???

		select {
		case _, ok = <-c[1]:
			if ok == false {
				return
			}
		case v, ok = <-c[0]:
			if ok == false {

			select {
			case _, ok := <-c[1]:
				if ok == false {
					close(c[1])
					return
				}
			case c[0] <- v:
			}
		}
	}
}
*/
