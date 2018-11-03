// Package stats contains aggregators for real and complex data for statistical analysis.
package stats

import (
	"fmt"
	"math"
)

// Real is an aggregator for float64.
type Real struct {
	Mean, Min, Max float64
	Num, Ignored   uint64
	M2             float64
}

// Push a value into the aggregator.
func (s *Real) Push(v float64) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		s.Ignored++
	}
	if s.Num == 0 {
		s.Min = v
		s.Max = v
	} else if v < s.Min {
		s.Min = v
	} else if v > s.Max {
		s.Max = v
	}
	s.Num++
	delta := v - s.Mean
	s.Mean += delta / float64(s.Num)
	s.M2 += delta * (v - s.Mean)
}

// Merge add into dst.
func (dst *Real) Merge(add *Real) {
	dst.Mean = (dst.Mean*float64(dst.Num) + add.Mean*float64(add.Num)) / float64(dst.Num+add.Num)
	dst.Num += add.Num
	dst.M2 += add.M2
	if add.Min < dst.Min {
		dst.Min = add.Min
	}
	if add.Max > dst.Max {
		dst.Max = add.Max
	}
	dst.Ignored += add.Ignored
}

// Variance returns the variance.
func (s *Real) Variance() float64 {
	return s.M2 / float64(s.Num-1)
}

// Std return the standard deviation.
func (s *Real) Std() float64 {
	return math.Sqrt(s.Variance())
}

func (s Real) String() string {
	return fmt.Sprintf("%v +/- %v [%v, %v]", s.Mean, s.Std(), s.Min, s.Max)
}
