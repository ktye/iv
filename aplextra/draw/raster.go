package draw

import (
	"image"

	"github.com/ktye/iv/apl"
	aplimg "github.com/ktye/iv/aplextra/image"
	"gonum.org/v1/plot/vg/vgimg"
)

// Rasterize the canvas to an image.
// L: shape (width height) of the image in pixels
// R: Canvas
// The result is an image.Value.
func rasterize(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	ar, ok := r.(apl.Array)
	if ok == false {
		return false, nil, nil
	}

	c, ok := r.(Canvas)
	if ok == false {
		return false, nil, nil
	}

	shape := ar.Shape()
	if len(shape) != 1 || shape[0] != 2 {
		return false, nil, nil
	}
	a1, err := ar.At(0)
	if err != nil {
		return true, nil, err
	}
	a2, err := ar.At(2)
	if err != nil {
		return true, nil, err
	}
	width, ok := apl.ToInt(a1)
	if ok == false {
		return false, nil, nil
	}
	height, ok := apl.ToInt(a2)
	if ok == false {
		return false, nil, nil
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	dst := vgimg.NewWith(vgimg.UseImage(img))
	c.ReplayOn(dst)
	return true, aplimg.Value{img}, nil

}
