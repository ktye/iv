package p

import (
	"fmt"
	"strings"

	"golang.org/x/image/font"

	"github.com/golang/freetype/truetype"
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/domain"
	aplfont "github.com/ktye/iv/cmd/lui/font"
	"github.com/ktye/plot"
)

var dark bool = true
var transparent bool
var colors []int
var width int = 800
var height int = 400

func setDark(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	dark = boolean(a, R)
	return apl.EmptyArray{}, nil
}

func setTransparent(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	transparent = boolean(a, R)
	return apl.EmptyArray{}, nil
}

func setColors(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.EmptyArray); ok {
		colors = nil
		return R, nil
	}
	to := domain.ToIndexArray(nil)
	v, ok := to.To(a, R)
	if ok == false {
		return nil, fmt.Errorf("p/colors: right argument must be integers: %T", R)
	}
	ia := v.(apl.IntArray)
	if ia.Size() < 1 {
		colors = nil
	} else {
		colors = make([]int, len(ia.Ints))
		copy(colors, ia.Ints)
	}
	return apl.EmptyArray{}, nil
}

func setSize(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	w, h, err := read2(a, R)
	if err != nil {
		return nil, fmt.Errorf("p/size: %s", err)
	}
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("p/size: WIDTH HEIGHT must be > 0")
	}
	width = w
	height = h
	return apl.EmptyArray{}, nil
}

func setFontSizes(a *apl.Apl, _, R apl.Value) (apl.Value, error) {
	big, small, err := read2(a, R)
	if err != nil {
		return nil, fmt.Errorf("p/fontsizes: %s", err)
	}
	if big <= 0 || small <= 0 {
		return nil, fmt.Errorf("p/fontsizes: BIG SMALL must be > 0")
	}
	var sizes [2]float64
	sizes[0] = float64(big)
	sizes[1] = float64(small)

	ttf := aplfont.APL385()
	f, err := truetype.Parse(ttf)
	if err != nil {
		return nil, err
	}

	var faces [2]font.Face
	for i := range sizes {
		opt := truetype.Options{
			Size: sizes[i],
			DPI:  72,
		}
		faces[i] = truetype.NewFace(f, &opt)
	}
	plot.SetFonts(faces[0], faces[1])
	return apl.EmptyArray{}, nil
}

func read2(a *apl.Apl, R apl.Value) (int, int, error) {
	to := domain.ToIndexArray(nil)
	ia, ok := to.To(a, R)
	if ok == false {
		return 0, 0, fmt.Errorf("argument must be 2 ints: %T", R)
	}
	var i [2]int
	ints := ia.(apl.IntArray).Ints
	if len(ints) > 0 {
		i[0] = ints[0]
	}
	if len(ints) > 1 {
		i[1] = ints[1]
	}
	return i[0], i[1], nil
}

func colorOrder() string {
	if len(colors) == 0 {
		return ""
	}
	v := make([]string, len(colors))
	for i, n := range v {
		v[i] = fmt.Sprintf("#%06X\n", n)
	}
	return strings.Join(v, ",")
}

func boolean(a *apl.Apl, v apl.Value) bool {
	num, ok := v.(apl.Number)
	if ok == false {
		return false
	}
	b, ok := a.Tower.ToBool(num)
	if ok == false {
		return false
	}
	return bool(b)
}
