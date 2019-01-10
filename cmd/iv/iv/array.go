package iv

import (
	"fmt"
	"io"

	"github.com/ktye/iv/apl"
)

// SendArrays assembles arrays from scalars and sends them over the channel once they are complete.
// It also sends the termination level.
// SendArrays is only called for rank > 0.
func (p *InputParser) sendArrays(c apl.Channel) {
	defer close(c[0])

	var R = p.Rank
	var shape []int = make([]int, R) // shape of the current array
	var array []apl.Value            // current array
	var sp int                       // inc shape at this position
	var pending bool                 // data available to be send at eof

	resetShape := func() {
		for i := 0; i < R; i++ {
			shape[i] = 1
		}
		sp = R
		pending = false
	}
	send := func(E int) {
		size := len(array)
		if prod(shape) != size {
			panic(fmt.Errorf("array is not uniform: prod(%v) != %d", shape, size))
		}

		dims := make([]int, len(shape))
		copy(dims, shape)
		ar := apl.MixedArray{Dims: dims, Values: array} // array can be used. It's reallocated.
		c[0] <- apl.List{ar, apl.Index(E)}

		// Reset state after sending.
		resetShape()
		array = make([]apl.Value, 0, size) // Likely the next array has a similar size.
	}

	resetShape()
	for {
		// Abort request.
		select {
		case _, ok := <-c[1]:
			if ok == false {
				return
			}
		default:
		}

		scalar, S, err := p.Next()
		if err == io.EOF {
			if pending {
				send(0)
			}
			return
		} else if err != nil {
			panic(err) // TODO ?
		}

		array = append(array, scalar)
		pending = true
		if idx := R - S; idx < 0 {
			send(S - R - 1)
		} else if idx < sp {
			sp = idx
			shape[sp]++
		}
	}
}

// SendScalars is like sendArrays, for the rank=0 case.
func (p *InputParser) sendScalars(c apl.Channel) {
	defer close(c[0])
	for {
		// Abort request.
		select {
		case _, ok := <-c[1]:
			if ok == false {
				return
			}
		default:
		}

		scalar, S, err := p.Next()
		if err == io.EOF {
			return
		} else if err != nil {
			panic(err)
		}
		E := S - 1
		c[0] <- apl.List{scalar, apl.Index(E)}
	}
}

func prod(shape []int) int {
	p := 1
	for _, v := range shape {
		p *= v
	}
	return p
}
