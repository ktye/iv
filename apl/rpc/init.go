package rpc

import (
	"encoding/gob"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

func init() {
	// Register types for communication.
	gob.Register(apl.Bool(false))
	gob.Register(apl.Int(0))
	gob.Register(numbers.Float(0.0))
	gob.Register(numbers.Complex(0))
	gob.Register(apl.String(""))
	gob.Register(apl.List(nil))
	gob.Register(apl.MixedArray{})
	gob.Register(apl.IntArray{})
	gob.Register(apl.Bool(false))
}
