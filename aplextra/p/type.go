package p

import (
	"fmt"
	"image"
	"image/draw"
	"reflect"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/xgo"
	"github.com/ktye/plot"
	"github.com/mattn/go-sixel"
)

// Plot embedds a *plot.Plot (from package ktye/plot).
// It can be manimpulated like a dictionary but has reference semantics.
type Plot struct {
	xgo.Value
}

func (p Plot) Copy() apl.Value {
	return p
}

// Line embeds a *plot.Line (from package ktye/plot).
type Line struct {
	xgo.Value
}

func (l Line) Copy() apl.Value {
	return l
}

// String overwrites the stringer if gui is 0.
func (p Plot) String(f apl.Format) string {
	if gui {
		return p.Value.String(f)
	}
	m, err := plotToImage(nil, nil, p)
	if err != nil {
		return err.Error()
	}
	return toSixel(m.(apl.Image))
}

func toSixel(m apl.Image) string {
	var buf strings.Builder
	enc := sixel.NewEncoder(&buf)
	err := enc.Encode(m.Image)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func toPlot(p *plot.Plot) Plot {
	return Plot{xgo.Value(reflect.ValueOf(p))}
}

// p converts the reflection value back to a plot pointer.
func (p Plot) p() *plot.Plot {
	return reflect.Value(p.Value).Interface().(*plot.Plot)
}

func (l Line) l() *plot.Line {
	return reflect.Value(l.Value).Interface().(*plot.Line)
}

// PlotArray implements an apl.Array. It is used to layout plots on a grid.
// Rank should not be higher than 2.
// A plot array can be reshaped to any rectangle. If the resulting size is
// larger that the original, it is filled with empty plots.
type PlotArray struct {
	p    []Plot
	Dims []int
}

func (p PlotArray) Copy() apl.Value {
	r := PlotArray{Dims: apl.CopyShape(p), p: make([]Plot, len(p.p))}
	copy(r.p, p.p)
	return r
}
func (p PlotArray) String(f apl.Format) string {
	if gui {
		return fmt.Sprintf("p.PlotArray %v", p.Dims)
	}
	m, err := plotToImage(nil, nil, p)
	if err != nil {
		return err.Error()
	}
	return toSixel(m.(apl.Image))
}
func (p PlotArray) At(i int) apl.Value {
	return p.p[i]
}
func (p PlotArray) Shape() []int {
	return p.Dims
}
func (p PlotArray) Size() int {
	return len(p.p)
}
func (p PlotArray) Reshape(shape []int) apl.Value {
	n := prod(shape)
	res := PlotArray{
		p:    make([]Plot, n),
		Dims: make([]int, len(shape)),
	}
	copy(res.Dims, shape)
	for i := 0; i < n; i++ {
		if i < len(p.p) {
			res.p[i] = p.p[i]
		} else {
			res.p[i] = toPlot(defaultPlot())
		}
	}
	return res
}
func (p PlotArray) Set(i int, v apl.Value) error {
	if i < 0 || i >= len(p.p) {
		return fmt.Errorf("index out of range")
	}
	plt, ok := v.(Plot)
	if ok == false {
		return fmt.Errorf("p.PlotArray: set: value must be Plot: %T", v)
	}
	p.p[i] = plt
	return nil
}

func prod(shape []int) int {
	if len(shape) == 0 {
		return 0
	}
	n := shape[0]
	for i := 1; i < len(shape); i++ {
		n *= shape[i]
	}
	return n
}

func plotToImage(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	var plots plot.Plots
	rows := 1
	if p, ok := R.(Plot); ok {
		plt := p.p()
		plots = plot.Plots{*plt}
	} else if pa, ok := R.(PlotArray); ok {
		plots = make(plot.Plots, len(pa.p))
		for i := range plots {
			plots[i] = *pa.p[i].p()
		}
		rs := pa.Shape()
		if len(rs) > 1 {
			rows = len(pa.p) / rs[len(rs)-1]
		}
	} else {
		return nil, fmt.Errorf("p/plot to image: argument must be Plot or PlotArray: %T", R)
	}

	var err error
	var w, h int
	if L == nil {
		w, h = width, height
	} else {
		w, h, err = read2(a, L)
		if err != nil {
			return nil, err
		}
		if w <= 0 {
			w = width
		}
		if h <= 0 {
			h = width
		}
	}
	if rows == 1 {
		return toImage(plots, w, h)
	}

	// Draw multiple rows, if the PlotArray has rank > 1.
	im := image.NewRGBA(image.Rectangle{Max: image.Point{w, h}})
	cols := len(plots) / rows
	dst := image.Rectangle{Max: image.Point{w, h / rows}}
	for i := 0; i < len(plots); i += cols {
		m, err := toImage(plots[i:i+cols], w, h/rows)
		if err != nil {
			return nil, err
		}
		draw.Draw(im, dst, m.Image, image.ZP, draw.Src)
		dst = dst.Add(image.Point{0, h / rows})
	}
	return apl.Image{im, []int{w, h}}, nil
}

func toImage(plots plot.Plots, w, h int) (apl.Image, error) {
	ip, err := plots.IPlots(w, h)
	if err != nil {
		return apl.Image{}, err
	}
	return apl.Image{
		Image: plot.Image(ip, nil, w, h),
		Dims:  []int{w, h},
	}, nil
}
