package apl

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// An ImageWriter is anything that can handle image output.
// Apl's stdimg device uses it.
// A single image can be written directly with WriteImage.
// An animation needs Start- and StopLoop before and after.
type ImageWriter interface {
	WriteImage(Image) error
	StartLoop()
	StopLoop()
}

// An Image is a raster image.
//
// It can by converting from a numeric array of shape HEIGHT WIDTH.
// An Image is never empty, it always has rank 2.
// It cannot be reshaped, instead reshape it's int array before creating it.
//
// Formats:
//	Img ← `img ⌶B       B BoolArray (0 White, 1 Black) // TODO or user def, or alpha?
//	Img ← `img ⌶(I;P;)  I numeric array with values in the range ⎕IO+(0..0xFF) as indexes into P
//                           P (palette) vector of shape 256 with values as below:
//	Img ← `img ⌶N       N numeric array, values between 0 and 0xFFFFFFFF (aarrggbb)
// Transparency has the value 0xFF000000, which is inverted compared to the go image library,
// to be able to specify opaque colors in the short form 0xRRGGBB.
// After creation, an image can be indexed and assigned to.
type Image struct {
	image.Image
	Dims []int
}

func (i Image) String(a *Apl) string {
	return i.toIntArray().String(a)
}
func (i Image) Copy() Value { return i } // Image is copied by reference.

func (i Image) At(k int) Value {
	ic, idx := NewIdxConverter(i.Dims)
	ic.Indexes(k, idx)
	return colorValue(i.Image.At(idx[0], idx[1]))
}
func (i Image) Shape() []int {
	return i.Dims
}
func (i Image) Size() int {
	return prod(i.Dims)
}
func (i Image) Set(k int, v Value) error {
	num, ok := v.(Number)
	if ok == false {
		return fmt.Errorf("img set: value must be a number: %T", v)
	}
	c, ok := num.ToIndex()
	if ok == false {
		return fmt.Errorf("img set: value must be an integer: %T", v)
	}
	y := k / i.Dims[1]
	x := k - y*i.Dims[1]
	r := i.Bounds()
	if d, ok := i.Image.(draw.Image); ok {
		d.Set(x+r.Min.X, y+r.Min.Y, toColor(c))
		return nil
	}
	return fmt.Errorf("img: image is not settable: %T", i.Image)
}

func (i Image) toIntArray() IntArray {
	ints := make([]int, prod(i.Dims))
	shape := make([]int, len(i.Dims))
	copy(shape, i.Dims)
	idx := 0
	r := i.Image.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			ints[idx] = int(colorValue(i.Image.At(x, y)))
			idx++
		}
	}
	return IntArray{Dims: shape, Ints: ints}
}

// ColorValue converts a Color to an Int.
func colorValue(c color.Color) Int {
	r, g, b, a := c.RGBA() // uint32 premultiplied with alpha.
	u := (255-(a>>8))<<16 | r&0xFF00<<8 | g&0xFF00 | b>>8
	return Int(u)
}
func toColor(i int) color.Color {
	u := uint32(i)
	return color.RGBA{uint8(u & 0xFF0000 >> 16), uint8(u & 0xFF00 >> 8), uint8(u & 0xFF), ^uint8(u >> 24)}
}
