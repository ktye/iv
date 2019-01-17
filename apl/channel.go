package apl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Channel is a pair of read and write channels.
// Channel operations:
//	↑C	take one: read one value
//	I↑C	take: take multiple values and reshape by I
//	C↓R	drop value: send value
//	↓C	crop channel: close channel
//	f/C	reduce over channel
//	f\C	scan over channel
//	[L]f¨C	each channel
type Channel [2]chan Value

func NewChannel() Channel {
	var c Channel
	c[0] = make(chan Value)
	c[1] = make(chan Value)
	return c
}

func (c Channel) String(a *Apl) string {
	return "apl.Channel"
}

// Close closes the write channel and drains the read channel.
func (c Channel) Close() {
	close(c[1])
	for range c[0] {
	}
}

// scope return a channel and copies values from R[0].
// It is called by scope assignment: ⎕←R.
func (R Channel) Scope(a *Apl) Channel {
	c := NewChannel()
	go func(r Channel) {
		defer close(c[0])
		for {
			select {
			case _, ok := <-c[1]:
				if ok == false {
					close(r[1])
					return
				}
			case v, ok := <-r[0]:
				if ok == false {
					return
				}
				fmt.Fprintf(a.stdout, "%s\n", v.String(a))
				select {
				case _, ok := <-c[1]:
					if ok == false {
						close(r[1])
						return
					}
				case c[0] <- v:
				}
			}
		}
	}(R)
	return c
}

// Apply returns a new channel.
// It reads values from c[0], applies the f the each value and writes the result to return
// returned channel.
// L (may be nil) is used as a left value for f.
// If L is also a channel, a value is read each time, before applying f.
// If filter is true, values are skipped if f returns an EmptyArray.
func (R Channel) Apply(a *Apl, f Function, L Value, filter bool) Channel {
	lv := L
	l, lc := L.(Channel)

	c := NewChannel()
	go func(r Channel) {
		defer close(c[0])
		var err error
		for {
			select {
			case _, ok := <-c[1]:
				if ok == false {
					close(r[1])
					if lc {
						close(l[1])
					}
					return
				}
			case v, ok := <-r[0]:
				if ok == false {
					if lc {
						close(l[1])
					}
					return
				}
				if lc {
					lv = <-l[0]
				}
				v, err = f.Call(a, lv, v)
				if err != nil {
					c[0] <- Error{err}
					close(r[1])
					return
				}
				if _, ok := v.(EmptyArray); filter == false || ok == false {
					c[0] <- v
				}
			}
		}
	}(R)
	return c
}

// LineReader wraps a ReadCloser with a Channel.
func LineReader(rc io.ReadCloser) Channel {
	scn := bufio.NewScanner(rc)
	c := NewChannel()
	go func(c Channel) {
		for scn.Scan() {
			line := scn.Text()
			select {
			case _, ok := <-c[1]:
				if ok == false {
					break
				}
			case c[0] <- String(line):
			}
		}
		close(c[0])
		rc.Close()
	}(c)
	return c
}

// NewChannelReader converts a channel to an io.Reader.
func NewChannelReader(a *Apl, c Channel) *ChannelReader {
	return &ChannelReader{
		a: a,
		c: c,
	}
}

// ChannelReader converts values in the channel to strings and provides an io.Reader.
// The strings are joind by newlines.
type ChannelReader struct {
	a      *Apl
	c      Channel
	buf    bytes.Buffer
	first  bool
	closed bool
}

func (r *ChannelReader) Read(p []byte) (n int, err error) {
	if r.closed {
		return r.buf.Read(p)
	}
	if r.buf.Len() < 1024 {
		select {
		case _, ok := <-r.c[1]:
			if ok == false {
				close(r.c[0])
				return 0, io.ErrClosedPipe
			}
		case v, ok := <-r.c[0]:
			if ok == false {
				r.closed = true
			} else {
				if r.first {
					r.first = false
				} else {
					r.buf.WriteRune('\n')
				}
				r.buf.WriteString(v.String(r.a)) // TODO: this could be a race when formatting.
			}
		}
	}
	return r.buf.Read(p)
}
