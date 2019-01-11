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
	var shape []int = make([]int, R)
	var array []apl.Value
	var idx, max, prd int
	var pending bool

	resetShape := func() {
		for i := 0; i < R; i++ {
			shape[i] = 1
		}
		idx = 0
		max = 1
		prd = 1
		pending = false

	}
	send := func(E int) {
		shape[R-max] = len(array) / prd
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

		// An [3 2 4] array has this sequence for S
		// 1 1 1 2 1 1 1 3 1 1 1 2 1 1 1 3 1 1 1 2 1 1 ?   S
		//       4       8                            24   positions at increase

		idx++
		array = append(array, scalar)
		pending = true
		if R-S < 0 {
			send(S - R - 1)
			continue
		}
		if S > max {
			shape[R-S+1] = idx / prd
			prd *= shape[R-S+1]
			max = S
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
