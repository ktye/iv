package p

import (
	"fmt"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/plot"
)

// plot1 returns a Plot or a PlotArray (rank-1).
// L may be nil, or an numeric array (rank-3 or smaller)
// R must be a numeric array or empty (rank-3 or smaller)
func plot1(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	if _, ok := R.(apl.EmptyArray); ok {
		return toPlot(defaultPlot()), nil
	}
	ar, rs, err := rank3(a, R)
	if err != nil {
		return nil, fmt.Errorf("p/plot: right argument %s", err)
	}
	x, err := xaxis(a, L, rs[2])
	if err != nil {
		return nil, fmt.Errorf("p/plot: left argument %s", err)
	}

	plots := make([]Plot, rs[0])
	for i := range plots {
		xoff := x.Dims[2] * x.Dims[1] * i
		if i >= x.Dims[0] {
			xoff = 0
		}
		roff := rs[2] * rs[1] * i
		p, err := plot2(a, x, xoff, ar, roff, L == nil)
		if err != nil {
			return nil, fmt.Errorf("p/plot: %s", err)
		}
		plots[i] = toPlot(p)
	}
	if len(plots) == 1 {
		return plots[0], nil
	}
	return PlotArray{plots, []int{len(plots)}}, nil
}

func plot2(a *apl.Apl, x numbers.FloatArray, xoff int, ar apl.Array, roff int, monadic bool) (*plot.Plot, error) {
	p := defaultPlot()
	rs := ar.Shape()
	np := rs[2]
	isComplex, err := complexType(ar, roff)
	if err != nil {
		return p, err
	}
	p.Type = plot.XY
	if isComplex {
		p.Type = plot.AmpAng
		if monadic {
			p.Type = plot.Polar
		}
	}
	p.Lines = make([]plot.Line, rs[1])
	for i := 0; i < rs[1]; i++ {
		l := plot.Line{
			Id: i,
			X:  make([]float64, np),
		}
		x0 := xoff + i*np
		if i >= x.Dims[1] {
			x0 = xoff
		}
		copy(l.X, x.Floats[x0:x0+np])
		if isComplex {
			l.C, err = complexVector(a, ar, roff+i*np, np)
			if err != nil {
				return p, err
			}
		} else {
			l.Y, err = realVector(a, ar, roff+i*np, np)
			if err != nil {
				return p, err
			}
		}
		p.Lines[i] = l
	}
	return p, nil
}

// realVector returns a float or a complex array for the last axis of ar at the given offset.
func realVector(a *apl.Apl, ar apl.Array, off, np int) ([]float64, error) {
	f := make([]float64, np)
	for i := 0; i < np; i++ {
		v := ar.At(i + off)
		z, ok, err := complexNumber(v)
		if err != nil {
			return nil, err
		}
		if ok {
			return nil, fmt.Errorf("expected real array, got complex")
		}
		f[i] = real(z)
	}
	return f, nil
}

// complexVector returns a float or a complex array for the last axis of ar at the given offset.
func complexVector(a *apl.Apl, ar apl.Array, off, np int) ([]complex128, error) {
	c := make([]complex128, np)
	for i := 0; i < np; i++ {
		v := ar.At(i + off)
		z, _, err := complexNumber(v)
		if err != nil {
			return nil, err
		}
		c[i] = z
	}
	return c, nil
}

// complexType returns if the subarray (last 2 dimensions) contains complex numbers.
func complexType(ar apl.Array, roff int) (bool, error) {
	rs := ar.Shape()
	for i := roff; i < roff+rs[1]*rs[2]; i++ {
		v := ar.At(i)
		_, isComplex, err := complexNumber(v)
		if isComplex || err != nil {
			return isComplex, err
		}
	}
	return false, nil
}

func complexNumber(v apl.Value) (complex128, bool, error) {
	switch x := v.(type) {
	case apl.Bool:
		if x {
			return complex(1, 0), false, nil
		}
		return complex(0, 0), false, nil
	case apl.Int:
		return complex(float64(x), 0), false, nil
	case numbers.Float:
		return complex(float64(x), 0), false, nil
	case numbers.Complex:
		return complex128(x), true, nil
	default:
		return 0, false, fmt.Errorf("unknown numeric type: %T", v)
	}
}

// rank3 makes sure the argument is an array and reshapes it to rank-3.
func rank3(a *apl.Apl, R apl.Value) (apl.Array, []int, error) {
	ar, ok := R.(apl.Array)
	if ok == false {
		return nil, nil, fmt.Errorf("not an array: %T", R)
	}
	rs := ar.Shape()
	if len(rs) > 3 {
		return nil, nil, fmt.Errorf("rank is too high: %d", len(rs))
	} else if len(rs) < 3 {
		r, ok := ar.(apl.Reshaper)
		if ok == false {
			return nil, nil, fmt.Errorf("cannot reshape to rank 3: %T", R)
		}
		shape := []int{1, 1, 1}
		copy(shape[3-len(rs):], rs) // Keep leading ones.
		v := r.Reshape(shape)
		ar = v.(apl.Array)
		rs = ar.Shape()
	}
	return ar, rs, nil
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
