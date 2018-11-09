// Package draw is a 2d vector graphics package for APL.
//
// The it is an interface to gonum.org/plog/vg.
package draw

import (
	"fmt"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/recorder"

	"github.com/ktye/iv/apl"
)

type Canvas struct {
	*recorder.Canvas
}

func (c Canvas) String(a *apl.Apl) string {
	return "draw.Canvas"
}

func (_ Canvas) Call(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	if l != nil {
		return nil, fmt.Errorf("ctx cannot be called dyadically")
	}
	return Canvas{&recorder.Canvas{}}, nil
}

// Translate the canvas.
// Left: canvas, Right: complex scalar.
func translate(a *apl.Apl, l, r apl.Value) (bool, apl.Value, error) {
	c, ok := l.(Canvas)
	if ok == false {
		return false, nil, nil
	}

	z, ok := apl.ToComplex(r)
	if ok == false {
		return false, nil, nil
	}

	c.Translate(vg.Point{vg.Length(real(z)), vg.Length(imag(z))})
	return true, c, nil
}

type Path vg.Path

func (p Path) String(a *apl.Apl) string {
	return "draw.Path" // TODO: return string: "c p 3@90 p 1@0 0@0"
}

func (_ Path) Call(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	if l != nil {
		return nil, fmt.Errorf("draw.p cannot be called dyadically")
	}

	// A call to path requires an array with at least 2 elements.
	ar, ok := r.(apl.Array)
	if ok == false {
		return nil, fmt.Errorf("draw.p L: L must be a vector: %T", l)
	}

	shape := ar.Shape()
	if len(shape) != 1 || shape[0] < 2 || shape[0] > 4 {
		return nil, fmt.Errorf("draw.p L: L must be a vector with 2..4 elements")
	}

	n := shape[0]
	last, err := ar.At(n - 1)
	if err != nil {
		return nil, err
	}

	var p vg.Path
	if last, ok := last.(Path); ok {
		p = vg.Path(last)
	} else if z, ok := apl.ToComplex(last); ok == false {
		return nil, fmt.Errorf("draw.p L: last vector element must be draw.Path or complex: %T", last)
	} else {
		p = vg.Path{vg.PathComp{Type: vg.MoveComp, Pos: vg.Point{X: vg.Length(real(z)), Y: vg.Length(imag(z))}}}
	}

	points := make([]vg.Point, n-1)
	for i := 0; i < n-1; i++ {
		v, err := ar.At(i)
		if err != nil {
			return nil, err
		}
		if z, ok := apl.ToComplex(v); ok == false {
			return nil, fmt.Errorf("draw.p L: element must be complex: %T", v)
		} else {
			points[i] = vg.Point{vg.Length(real(z)), vg.Length(imag(z))}
		}
	}
	switch n - 1 {
	case 1:
		p = append(p, vg.PathComp{Type: vg.LineComp, Pos: points[0]})
	//case 2:
	//	p = append(p, vg.PathComp{Type: vg.CurveComp, Pos: points[1], Control: points[0:1]})
	//case 3:
	//	p = append(p, vg.PathComp{Type: vg.CurveComp, Pos: points[2], Control: points[0:2]})
	default:
		return nil, fmt.Errorf("draw.p: wrong path length: this should not happen")
	}

	return Path(p), nil
}

/* TODO
// Polyline converts an array into linear line segements.
// The last element may be a path or a numeric value.
// Calling l V with V <- A B C D is equivalent to
// l V <--> p A p B p C D.
func polyline(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	if l != nil {
		return fmt.Errorf("draw.l cannot be called dyadically")
	}

	// A call to polyline requires an array with at least 2 elements.
	ar, ok := r.(apl.Array)
	if ok == false {
		return nil, fmt.Errorf("draw.l: L must be a vector: %T", l)
	}

	shape := ar.Shape()
	if len(shape) != 1 || shape[0] < 2 {
		return nil, fmt.Errorf("draw.l L: L must be a vector with 2..4 elements")
	}

	// The last elementy may be a number or a path.
	n := shape[0]
	pth, err := ar.At(n - 1)
	if err != nil {
		return nil, err
	}

	args := apl.GeneralArray{
		Values: make([]apl.Value, 2),
		Dims:   []int{2},
	}

	pc := Path{}
	for i := n - 2; i >= 0; i-- {
		v, err := ar.At(i)
		if err != nil {
			return nil, err
		}
		args.Values[0] = v
		args.Values[1] = pth
		pth, err = pc.Call(a, nil, args)
		if err != nil {
			return nil, err
		}
	}
	return pth, nil
}
*/

func closepath(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	if l != nil {
		return nil, fmt.Errorf("draw.c cannot be called dyadically")
	}
	p, ok := r.(Path)
	if ok == false {
		return nil, fmt.Errorf("draw.c expectes a path on the right: %T", r)
	}
	p = append(p, vg.PathComp{Type: vg.CloseComp})
	return p, nil
}

// Fillpath fills a path.
// l must be the Canvas and r the path to fill.
func fillpath(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	c, p, err := canvaspath(l, r)
	if err != nil {
		return nil, err
	}
	c.Fill(p)
	return c, nil
}

// strokepath strokes a path.
// l must be the Canvas and r the path to stroke
func strokepath(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	c, p, err := canvaspath(l, r)
	if err != nil {
		return nil, err
	}
	c.Stroke(p)
	return c, nil
}

func canvaspath(l, r apl.Value) (Canvas, vg.Path, error) {
	c, ok := l.(Canvas)
	if ok == false {
		return Canvas{}, nil, fmt.Errorf("draw.f left argument should be draw.Canvas: %T", l)
	}

	p, ok := r.(Path)
	if ok == false {
		return Canvas{}, nil, fmt.Errorf("draw.r left argument should be draw.Path: %T", r)
	}

	return c, vg.Path(p), nil
}
