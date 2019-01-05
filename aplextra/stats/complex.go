package stats

import (
	"math"
	"math/cmplx"
)

// Complex is an aggregator for complex128.
type Complex struct {
	Mean         complex128
	MaxAbs       float64
	Num, Ignored uint64
	Re, Im       Real
	Co           float64
}

// Push a value into the aggregator.
func (s *Complex) Push(v complex128) {
	if cmplx.IsNaN(v) || cmplx.IsInf(v) {
		s.Ignored++
		return
	}
	s.Num++

	// The order is important see en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Online
	s.Im.Push(imag(v))
	s.Co += (real(v) - s.Re.Mean) * (imag(v) - s.Im.Mean)
	s.Re.Push(real(v))

	if a := cmplx.Abs(v); a > s.MaxAbs {
		s.MaxAbs = a
	}
	s.Mean = complex(s.Re.Mean, s.Im.Mean)
}

// Merge add into dst.
func (dst *Complex) Merge(add *Complex) {
	a, b := float64(dst.Num), float64(add.Num)
	dst.Mean = (dst.Mean*complex(a, 0) + add.Mean*complex(b, 0)) / complex(a+b, 0)
	if add.MaxAbs > dst.MaxAbs {
		dst.MaxAbs = add.MaxAbs
	}
	dst.Num += add.Num
	dst.Ignored += add.Ignored
	dst.Re.Merge(&add.Re)
	dst.Im.Merge(&add.Im)
}

// Variance retuns the variance of the real and imag part as a complex number.
func (s *Complex) Variance() complex128 {
	return complex(s.Re.Variance(), s.Im.Variance())
}

// Covariance returns the sample covariance.
func (s *Complex) Covariance() float64 {
	return s.Co / float64(s.Num-1)
}

// Std retuns the standard deviation of the real and imag part as a complex number.
func (s *Complex) Std() complex128 {
	return complex(s.Re.Std(), s.Im.Std())
}

// Ellipse calculates the ellipsis of constant probability for bivariate normally distributed values.
// The returned values include the max and min standard distribution in the principal axis
// as well as the orientation of the ellipsis.
// The standard deviation is the absolute value of the values returned.
// The direction of the principal axis is the phase of the values returned.
func (c *Complex) Ellipse() (complex128, complex128) {

	// The ellipsis is calculated by a principal-axis transformation of the covariance matrix.
	varx := c.Re.Variance()
	vary := c.Im.Variance()
	cov := c.Covariance()

	// Eigenvalues and eigenvectors.
	trace := varx + vary
	det := varx*vary - cov*cov
	e1 := trace/2 + math.Sqrt(trace*trace/4-det)
	e2 := trace/2 - math.Sqrt(trace*trace/4-det)
	v1 := complex(e1-vary, cov)
	v2 := complex(e2-vary, cov)

	s1 := math.Sqrt(e1)
	s2 := math.Sqrt(e2)
	return complex(s1/cmplx.Abs(v1), 0) * v1, complex(s2/cmplx.Abs(v2), 0) * v2
}

// EllipsePath returns the elliptical path around the values for the given probability assuming bivariate normally distributed data.
func (c *Complex) EllipsePath(p float64, segments int) []complex128 {
	if segments < 4 {
		segments = 8*4 + 1
	}
	a, b := c.Ellipse()
	path := make([]complex128, segments)
	k := complex(math.Sqrt(-2*math.Log(1-p)), 0)
	for i := range path {
		phi := 2.0 * math.Pi * float64(i) / float64(len(path)-1)
		sin, cos := math.Sincos(phi)
		path[i] = c.Mean + k*(a*complex(sin, 0)+b*complex(cos, 0))
	}
	return path
}

// BinormalQuantile is an approximation of the quantile for weakly correlated normally distributed data
// for the given probability p in [0, 1].
func (s *Complex) BinormalQuantile(p float64) float64 {
	return binormalFactor(p) / binormalFactor(0.95) * s.q95()
}

// Q95 is an approximation of the 95% quantile for weakly correlated normally distributed data
func (s *Complex) q95() float64 {
	a := 1.97
	b := 3.2
	return a * math.Pow(math.Pow(s.Re.Std(), b)+math.Pow(s.Im.Std(), b), 1/b)
}

func binormalFactor(p float64) float64 {
	if p <= 0 || p >= 1 {
		return math.NaN()
	}
	return math.Sqrt(-2.0 * math.Log(1.0-float64(p)))
}
