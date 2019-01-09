package iv

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

// SendArrays assembles arrays from scalars and sends them over the channel once they are complete.
// It also sends the termination level.
func (p *InputParser) sendArrays(c apl.Channel) {
	var R = p.Rank
	var shape []int = make([]int, R) // shape of the current array
	var array []apl.Value            // current array
	var E = 0                        // termination level
	var max int                      // max level observed in input

	send := func(eof bool) {
		if eof {
			E = -1
		}
		dims := make([]int, len(shape))
		copy(dims, shape)
		ar := apl.MixedArray{Dims: dims, Values: array}
		c[0] <- apl.List{ar, apl.Index(E)}

		// Reset state after sending.
		for i := range shape {
			shape[i] = 0
		}
		array = nil
		E = 0
	}

	for i := 0; i < R-1; i++ {
		shape[i] = 1
	}
	for {
		scalar, S, eof, err := p.Next()
		if err != nil {
			panic(err) // TODO ?
		}

		E = S - R
		if S > R {
			S -= R
		}

		if scalar != nil {
			// Catenate to current array and set shape.
			array = append(array, scalar)
			if idx := R - S; idx >= 0 && idx < R {
				shape[idx]++
			}
			if prod(shape) != len(array) {
				// TODO: how to signal an error?
				panic(fmt.Errorf("array is not uniform: prod(%v) != %d", shape, len(array)))
			}
		}

		if eof {
			for i, x := range shape {
				if x == 0 {
					shape[i] = 1
				}
			}
			E = max
			send(true)
			return
		}
		if E > max {
			max = E
		}
		if E > 0 {
			send(false)
		}
	}
}

func prod(shape []int) int {
	p := 1
	for _, v := range shape {
		p *= v
	}
	return p
}
