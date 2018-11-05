package funcs

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

func init() {
	register("←", both(sink, assign))
	addDoc("←", `← primitive function: assign, sink
←R: sink R
	returns the empty array to suppress printing
Z←R: Z: identifier
	assign R to the identifier Z
`)
}

// Sink converts the left argument to an empty array to suppress printing.
func sink(a *apl.Apl, ignored, v apl.Value) (bool, apl.Value, error) {
	return true, apl.EmptyArray{}, nil
}

// Assign the right value to the identifier given on the left.
func assign(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	if id, ok := l.(apl.Identifier); ok == false {
		return true, nil, fmt.Errorf("variable assignment needs an identifier on the left, not %T", l)
	} else {
		if err := a.Assign(string(id), r); err != nil {
			return true, nil, err
		}
	}
	return true, r, nil
}
