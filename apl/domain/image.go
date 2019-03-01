package domain

import (
	img "image"
	"image/color"

	"github.com/ktye/iv/apl"
)

// ToImage can convert an int array to an image.
// It also accepts a list with two int arrags: the image data as indexes into a palette, and the palette itself.
func ToImage(child SingleDomain) SingleDomain {
	return image{child, true}
}

func IsImage(child SingleDomain) SingleDomain {
	return image{child, false}
}

type image struct {
	child SingleDomain
	conv  bool
}

func (i image) String(f apl.Format) string {
	name := "image"
	if i.conv {
		name = "toimage"
	}
	if i.child == nil {
		return name
	}
	return name + " " + i.child.String(f)
}

func (im image) To(a *apl.Apl, V apl.Value) (apl.Value, bool) {
	_, ok := V.(apl.Image)
	if im.conv == false && ok == false {
		return V, false
	} else if im.conv == false && ok {
		return propagate(a, V, im.child)
	} else if ok {
		return propagate(a, V, im.child)
	}

	var pal []int
	ia, ok := indexarray{nil, true}.To(a, V)
	if ok == false {
		if l, ok := V.(apl.List); ok == false || len(l) != 2 {
			return V, false
		} else {
			idx, ok := indexarray{nil, true}.To(a, l[0])
			pa, pok := indexarray{nil, true}.To(a, l[1])
			if ok == false || pok == false {
				return V, false
			}
			ia = idx
			pal = pa.(apl.IntArray).Ints
		}
		return V, false
	}

	if m, ok := ints2image(ia.(apl.IntArray), pal); ok {
		return propagate(a, m, im.child)
	} else {
		return V, false
	}
}

func ints2image(ia apl.IntArray, pal []int) (apl.Image, bool) {
	res := apl.Image{}
	shape := ia.Shape()
	if len(shape) != 2 {
		return res, false
	}

	p := make(color.Palette, len(pal))
	for i, n := range pal {
		p[i] = toColor(n)
	}

	var r img.Rectangle
	r.Max.X = shape[1]
	r.Max.Y = shape[0]

	var pm *img.Paletted
	var im *img.RGBA
	if pal == nil {
		im = img.NewRGBA(r)
	} else {
		pm = img.NewPaletted(r, p)
	}
	i := 0
	for y := 0; y < r.Max.Y; y++ {
		for x := 0; x < r.Max.X; x++ {
			c := ia.Ints[i]
			if pal == nil {
				im.Set(x, y, toColor(c))
			} else if c >= len(pal) {
				pm.SetColorIndex(x, y, 0)
			} else {
				pm.SetColorIndex(x, y, uint8(c))
			}
			i++
		}
	}
	if pal == nil {
		res.Image = im
	} else {
		res.Image = pm
	}
	res.Dims = apl.CopyShape(ia)
	return res, true
}

// toColor is copied from apl/image.go
func toColor(i int) color.Color {
	u := uint32(i)
	return color.RGBA{uint8(u & 0xFF0000 >> 16), uint8(u & 0xFF00 >> 8), uint8(u & 0xFF), ^uint8(u >> 24)}
}
