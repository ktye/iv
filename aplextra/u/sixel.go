package u

import (
	"fmt"
	"io"

	"github.com/ktye/iv/apl"
	sixel "github.com/mattn/go-sixel"
)

// Sxl implements an apl.ImageWriter to serve as a stdimg device in a sixel capable terminal.
type Sxl struct {
	io.Writer
	animation bool
}

func (dev *Sxl) WriteImage(m apl.Image) error {
	if dev.animation {
		fmt.Fprint(dev.Writer, "\x1b[u") // CSI u: restore position
	}
	e := sixel.NewEncoder(dev.Writer)
	return e.Encode(m.Image)
}

func (dev *Sxl) StartLoop() {
	fmt.Fprint(dev.Writer, "\x1b[s") // CSI s: save position
	dev.animation = true
}

func (dev *Sxl) StopLoop() {
	dev.animation = false
}
