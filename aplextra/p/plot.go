package p

import (
	"fmt"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/plot"
)

// real data: plot each row in R as a line in an xy plot.
// dyadic case: L is the x axis, a vector or a conforming matrix.
//
// rank-3 input (TODO).
// create multiple plots, one for each major cell.
// Example: ⌼?10⍴10
func plotf(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	ar := R.(apl.Array)
	rs := ar.Shape()
	if len(rs) > 2 {
		return nil, fmt.Errorf("p/plot: TODO: rank-3 input") // multiple plots.
	} else if len(rs) < 1 {
		return nil, fmt.Errorf("p/plot: no input")
	}
	ln := rs[len(rs)-1]
	var x []float64
	var al apl.Array
	var ls = []int{ln}
	if L == nil {
		x = make([]float64, ln)
		for i := range x {
			x[i] = float64(i + a.Origin)
		}
	} else {
		al = L.(apl.Array)
		ls = al.Shape()
		if len(ls) < 1 || len(ls) > 2 {
			return nil, fmt.Errorf("p/plot: left argument does not conform. shape: %v", ls)
		}
		if l := ls[len(ls)-1]; l != ln {
			return nil, fmt.Errorf("p/plot: left and right argument have different last axis: %d != %d", l, ln)

		}
	}

	floats := func(A apl.Array, n, ln int) []float64 {
		r := make([]float64, ln)
		off := ln * n
		for i := 0; i < ln; i++ {
			v := A.At(i + off)
			num, ok := v.(apl.Number)
			if ok == false {
				continue
			}
			f, ok := num.(numbers.Float)
			if ok {
				r[i] = float64(f)
				continue
			}
			idx, ok := num.ToIndex()
			if ok {
				r[i] = float64(idx)
			}
		}
		return r
	}

	var lines []plot.Line
	n := 1
	if len(rs) > 1 {
		n = rs[len(rs)-2]
	}
	for i := 0; i < n; i++ {
		if len(ls) > 1 {
			x = floats(al, i, ln)
		}
		y := floats(ar, i, ln)
		l := plot.Line{
			Id: i,
			X:  x,
			Y:  y,
		}
		lines = append(lines, l)
	}

	// TODO: get image size from stdimg.
	w, h := 400, 400

	plt := plot.Plot{
		Type:  plot.XY,
		Lines: lines,
	}
	plots := plot.Plots{plt}
	ip, err := plots.IPlots(w, h)
	if err != nil {
		return nil, err
	}

	// TODO: implement an interactive widget in ktye/ui used by iv/aplextra/u.

	return apl.Image{
		Image: plot.Image(ip, nil, w, h),
		Dims:  []int{w, h},
	}, nil
}
