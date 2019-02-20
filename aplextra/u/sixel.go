package u

import (
	"fmt"
	"io"
	"os"

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

// Device size query only works on stdout.
// Stdout and stdin must be connected to the terminal.
// If this fails, it is in bad state.
// TODO: this does not work. It hangs.
func (dev *Sxl) DeviceSize() (width, height int) {
	width = 800
	height = 400

	// Request size of text area in cells and pixels.
	os.Stdout.WriteString("\x1b[18t\x1b[14t")
	os.Stdout.Sync()
	var rows, cols, x, y int
	if n, err := fmt.Scanf("\x1b[8;%d;%d;t", &rows, &cols); n != 2 || err != nil {
		return
	}
	if n, err := fmt.Scanf("\x1b[4;%d;%d;t", &y, &x); n != 2 || err != nil {
		return
	}
	rowheight := y / height

	// Substract 2 lines to make the first and last line visible.
	height = y - 2*rowheight
	width = x
	if width <= 0 || height <= 0 {
		return 800, 400
	}
	return
}
