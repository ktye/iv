package apl

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"
)

// An ImageWriter is anything that can handle image output.
// Apl's stdimg device uses it.
type ImageWriter interface {
	WriteImage(Image) error
}

// An Image is one or more raster images.
// If it is more than one, it is an animation.
//
// An Image can be created by converting arrays of shape
// 	FRAMES HEIGHT WIDTH or HEIGHT WIDTH.
// An Image is never empty, it always has at least rank 2 and never more than 3.
// It cannot be reshaped, instead reshape it's int array before creating it.
//
// Formats:
//	Img ← `img ⌶B       B BoolArray (0 White, 1 Black) // TODO or user def, or alpha?
//	Img ← `img ⌶(I;P;)  I numeric array with values in the range ⎕IO+(0..0xFF) as indexes into P
//                           P (palette) vector of shape 256 with values as below:
//	Img ← `img ⌶N       N numeric array, values between 0 and 0xFFFFFFFF (rrggbbaa)
// After creation, an image can be indexed and assigned to.
type Image struct {
	Im    []image.Image
	Delay time.Duration
	Dims  []int
}

func (i Image) String(a *Apl) string {
	return i.toIntArray().String(a)
}

func (i Image) At(k int) Value {
	ic, idx := NewIdxConverter(i.Dims)
	ic.Indexes(k, idx)
	p := 0
	if len(idx) > 2 {
		p = idx[0]
		idx = idx[1:]
	}
	return colorValue(i.Im[p].At(idx[0], idx[1]))
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
	p := 0
	shape := i.Dims
	if len(i.Dims) == 3 {
		n := i.Dims[1] * i.Dims[2]
		p = k / n
		k -= p * n
	} else {
		shape = shape[1:]
	}
	y := k / shape[1]
	x := k - y*shape[1]
	m := i.Im[p]
	r := m.Bounds()
	if d, ok := m.(draw.Image); ok {
		d.Set(x+r.Min.X, y+r.Min.Y, toColor(c))
		return nil
	}
	return fmt.Errorf("img: image is not settable: %T", m)
}

func (i Image) toIntArray() IntArray {
	ints := make([]int, prod(i.Dims))
	shape := make([]int, len(i.Dims))
	copy(shape, i.Dims)
	off := 0
	for _, m := range i.Im {
		r := m.Bounds()
		for i := r.Min.Y; i < r.Max.Y; i++ {
			for k := r.Min.X; k < r.Max.X; k++ {
				ints[off] = int(colorValue(m.At(i, k)))
				off++
			}
		}
	}
	return IntArray{Dims: shape, Ints: ints}
}

// ColorValue converts a Color to an Int.
// On 32bit systems, white 0xFFFFFFFF will be -1.
func colorValue(c color.Color) Int {
	r, g, b, a := c.RGBA() // uint32 premultiplied with alpha.
	u := r&0xFF00<<16 | g&0xFF00<<8 | b&0xFF00 | a>>8
	return Int(u)
}
func toColor(i int) color.Color {
	u := uint32(i)
	return color.RGBA{uint8(u >> 24), uint8(u&0xFF0000) >> 16, uint8(u&0xFF00) >> 8, uint8(u & 0xFF)}
}
