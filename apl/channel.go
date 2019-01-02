package apl

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
