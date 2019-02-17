package u

import (
	"io"

	"github.com/ktye/iv/apl"
	sixel "github.com/mattn/go-sixel"
)

// Sxl implements an apl.ImageWriter to serve as a stdimg device in a sixel capable terminal.
type Sxl struct {
	io.Writer
}

func (dev Sxl) WriteImage(m apl.Image) error {
	e := sixel.NewEncoder(dev.Writer)
	return e.Encode(m.Image)
}
