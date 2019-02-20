package p

import (
	"fmt"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/plot"
	"github.com/ktye/plot/color"
)

// real data: plot each row in R as a line in an xy plot.
// dyadic case: L is the x axis, a vector or a conforming matrix.
//
// rank-3 input (TODO).
// create multiple plots, one for each major cell.
// Example: ⌼?10⍴10
func plot4(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	ar, rs, err := rank4(a, R)
	if err != nil {
		return nil, fmt.Errorf("p/plot: right argument %s", err)
	}
	x, err := xaxis(a, L, rs[3])
	if err != nil {
		return nil, fmt.Errorf("p/plot: left argument %s", err)
	}

	types, err := plotTypes(a, ar, L == nil)
	if err != nil {
		return nil, err
	}

	nf := rs[0] // number of frames
	np := rs[1] // number of plots
	//    rs[2] // number of lines
	//    rs[3] // number of points per line

	// Multiple frames are send over a channel (animation).
	if nf > 1 {
		plots := make([]plot.Plots, nf)
		for i := range plots {
			p, err := plot3(a, x, ar, i, types)
			if err != nil {
				return nil, err
			}
			plots[i] = p
		}

		// Each plot must have equal limits per frame (to prevent rescale during the loop).
		for n := 0; n < np; n++ {
			plts := make(plot.Plots, nf)
			for i, p := range plots {
				plts[i] = p[n]
			}
			lim, err := plts.EqualLimits()
			for i := range plots {
				if err != nil {
					plots[i][n].Limits = lim
				}
			}
		}

		frames := make([]apl.Value, nf)
		for i := range plots {
			m, err := toImage(a, plots[i], width, height)
			if err != nil {
				return nil, err
			}
			frames[i] = m
		}
		c := apl.NewChannel()
		go c.SendAll(frames)
		return c, nil
	}

	p, err := plot3(a, x, ar, 0, types)
	if err != nil {
		return nil, err
	}
	return toImage(a, p, width, height)
}

// plot3 returns plots of a single animation frame. ar is a rank 4 array.
func plot3(a *apl.Apl, al, ar apl.Array, frame int, types []string) (plot.Plots, error) {
	plts := make(plot.Plots, len(types))
	for i := range plts {
		p, err := plot2(a, al, ar, frame, i, types[i])
		if err != nil {
			return nil, fmt.Errorf("p/plot: plot %d: %s", i+1, err)
		}
		plts[i] = p
	}
	return plts, nil
}

func plot2(a *apl.Apl, al, ar apl.Array, frame, col int, plotType string) (plot.Plot, error) {
	p := plot.Plot{
		Type: plot.PlotType(plotType),
	}
	isComplex := false
	if plotType == "ampang" || plotType == "polar" {
		isComplex = true
	}
	rs := ar.Shape()
	np := rs[3]
	xf := al.(numbers.FloatArray)

	x0 := 0
	if col < xf.Dims[0] {
		x0 = col * xf.Dims[1] * xf.Dims[2]
	}
	lines := make([]plot.Line, rs[2])
	var err error
	for i := range lines {
		xoff := x0
		if i < xf.Dims[1] {
			xoff += i * np
		}
		var y []float64
		var c []complex128
		if isComplex {
			y, c, err = rcVector(a, ar, frame, col, i, true)
			if err != nil {
				return p, err
			}
		} else {
			y, c, err = rcVector(a, ar, frame, col, i, false)
			if err != nil {
				return p, err
			}
		}
		l := plot.Line{
			Id: i,
			X:  xf.Floats[xoff : xoff+np],
			Y:  y,
			C:  c,
		}
		lines[i] = l
	}
	p.Lines = lines
	p.Style.Dark = dark
	p.Style.Transparent = transparent
	p.Style.Order = color.Order(colorOrder())
	return p, nil
}

// plotTypes determines the type for the plot by looking at the data of the rank-4 array R.
// Types in the second axis may differ.
// If any number in axis 1, 3 or 4 is complex, the data is considered complex.
// Big numbers are not supported.
func plotTypes(a *apl.Apl, R apl.Array, monadic bool) ([]string, error) {
	shape := R.Shape()
	np := shape[1]
	types := make([]string, np)
	isComplex := make([]bool, np)
	idx := make([]int, len(shape))
	for i := 0; i < R.Size(); i++ {
		v := R.At(i)
		switch v.(type) {
		case apl.Bool, apl.Int, numbers.Float:
		case numbers.Complex:
			isComplex[idx[1]] = true
		default:
			return nil, fmt.Errorf("p/plot: data in R must be numeric: %T", v)
		}
		apl.IncArrayIndex(idx, shape)
	}
	for i := range types {
		if isComplex[i] {
			if monadic {
				types[i] = "polar"
			} else {
				types[i] = "ampang"
			}
		} else {
			types[i] = "xy"
		}
	}
	return types, nil
}

// rank4 makes sure the argument is an array and reshapes it to rank-4.
func rank4(a *apl.Apl, R apl.Value) (apl.Array, []int, error) {
	ar, ok := R.(apl.Array)
	if ok == false {
		return nil, nil, fmt.Errorf("not an array: %T", R)
	}
	rs := ar.Shape()
	if len(rs) > 4 {
		return nil, nil, fmt.Errorf("rank is too high: %d", len(rs))
	} else if len(rs) < 4 {
		r, ok := ar.(apl.Reshaper)
		if ok == false {
			return nil, nil, fmt.Errorf("cannot reshape to rank 4: %T", R)
		}
		shape := []int{1, 1, 1, 1}
		copy(shape[4-len(rs):], rs) // Keep leading ones.
		v := r.Reshape(shape)
		ar = v.(apl.Array)
		rs = ar.Shape()
	}
	return ar, rs, nil
}

func rcVector(a *apl.Apl, ar apl.Array, frame int, col int, line int, isComplex bool) (fs []float64, cs []complex128, err error) {
	rs := ar.Shape()
	np := rs[3]
	off := frame * rs[1] * rs[2] * rs[3]
	off += col * rs[2] * rs[3]
	off += line * rs[3]
	if isComplex {
		cs = make([]complex128, np)
	} else {
		fs = make([]float64, np)
	}
	for i := off; i < off+np; i++ {
		v := ar.At(i)
		var f float64
		var c complex128
		switch x := v.(type) {
		case apl.Bool:
			if x == true {
				f = 1.0
				c = complex(f, 0)
			}
		case apl.Int:
			f = float64(x)
			c = complex(f, 0)
		case numbers.Float:
			f = float64(x)
			c = complex(f, 0)
		case numbers.Complex:
			c = complex128(x)
			if isComplex == false {
				return nil, nil, fmt.Errorf("data is complex")
			}
		default:
			return nil, nil, fmt.Errorf("unknown numeric type: %T", v)
		}
		if isComplex {
			cs[i-off] = c
		} else {
			fs[i-off] = f
		}
	}
	return fs, cs, nil
}

// xaxis converts the left argument to an x-axis (rank-3 numbers.FloatArray).
func xaxis(a *apl.Apl, L apl.Value, lastAxis int) (numbers.FloatArray, error) {
	if L == nil {
		f := numbers.FloatArray{
			Dims:   []int{1, 1, lastAxis},
			Floats: make([]float64, lastAxis),
		}
		for i := range f.Floats {
			f.Floats[i] = float64(a.Origin + i)
		}
		return f, nil
	}
	al, ok := L.(apl.Array)
	if ok == false {
		al = apl.MixedArray{Dims: []int{1}, Values: []apl.Value{L}}
	}
	ls := al.Shape()
	if len(ls) > 3 {
		return numbers.FloatArray{}, fmt.Errorf("max rank is 3: %d", len(ls))
	}
	if n := ls[len(ls)-1]; n != lastAxis {
		return numbers.FloatArray{}, fmt.Errorf("last axis does not match last axis of R: %d != %d", n, lastAxis)
	}
	ua, ok := a.Unify(al, true)
	if ok == false {
		return numbers.FloatArray{}, fmt.Errorf("could not unify numbers")
	}
	res := numbers.FloatArray{
		Dims: []int{1, 1, lastAxis},
	}

	shape := apl.CopyShape(ua)
	copy(res.Dims[3-len(shape):], shape)

	switch ar := ua.(type) {
	case apl.BoolArray:
		res.Floats = make([]float64, len(ar.Bools))
		for i, b := range ar.Bools {
			if b {
				res.Floats[i] = 1.0
			}
		}
	case apl.IntArray:
		res.Floats = make([]float64, len(ar.Ints))
		for i, n := range ar.Ints {
			res.Floats[i] = float64(n)
		}
	case numbers.FloatArray:
		res.Floats = ar.Floats
	default:
		return numbers.FloatArray{}, fmt.Errorf("unsupported type: %T", ua)
	}
	return res, nil
}

func toImage(a *apl.Apl, plots plot.Plots, w, h int) (apl.Image, error) {
	ip, err := plots.IPlots(w, h)
	if err != nil {
		return apl.Image{}, err
	}
	return apl.Image{
		Image: plot.Image(ip, nil, w, h),
		Dims:  []int{w, h},
	}, nil
}
