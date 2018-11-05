package image

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/ktye/iv/apl"
)

func blue(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	// we ignore both values and return an image value.
	m := image.NewRGBA(image.Rect(0, 0, 400, 400))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	return true, Value{Image: m}, nil
}
