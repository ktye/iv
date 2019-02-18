package primitives

import (
	"fmt"
	"time"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
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
	n += int(a.Origin) // splitAxis substracts the origin.

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

// channelCopy connects two channels. It writes to L what it reads from R.
// The function returns the number of values copied.
func channelCopy(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	l := L.(apl.Channel)
	r := R.(apl.Channel)
	// Should this run in a go-routine and return directly?
	defer close(l[0])
	ret := apl.EmptyArray{}
	for {
		select {
		case _, ok := <-l[1]:
			if ok == false {
				close(r[1])
				return ret, nil
			}
		case v, ok := <-r[0]:
			if ok == false {
				return ret, nil
			}
			select {
			case _, ok := <-l[1]:
				if ok == false {
					close(r[1])
					return ret, nil
				}
			case l[0] <- v:
			}
		}
	}
}

// channelDelay returns a channel that sends at fixed intervals what it receives.
func channelDelay(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	d, ok := L.(numbers.Time).Duration()
	if ok == false {
		return nil, fmt.Errorf("channel delay: left argument is not a duration: %T", L)
	}
	in := R.(apl.Channel)
	out := in.Apply(a, Delay(d), nil, false)
	return out, nil
}

// Delay is a function that pauses execution for a given duration.
// It is currently not bound to a primitive and only used by channel-delay.
type Delay time.Duration

func (d Delay) Call(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	time.Sleep(time.Duration(d))
	return R, nil
}

// channel1 applies the monadic elementary function to each value in a channel.
func channel1(symbol string, fn func(*apl.Apl, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	return func(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
		c := R.(apl.Channel)
		return c.Apply(a, apl.Primitive(symbol), nil, false), nil
	}
}

// channel2 applies the dyadic elementary function to each value in a channel.
func channel2(symbol string, fn func(*apl.Apl, apl.Value, apl.Value) (apl.Value, bool)) func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error) {
	return func(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
		c := R.(apl.Channel)
		return c.Apply(a, apl.Primitive(symbol), L, false), nil
	}
}
