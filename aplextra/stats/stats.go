// Package stats is a statistical aggregator value for apl.
//
// min, max, std, var, histogram, tdigest, quantile,
// ascii and image plots...
// TODO
//
// Stats uses the average symbol ø as a primitive.
// This is not a standard APL symbol.
//
// 	S ← ø0    ⍝ initialize a stats.Value
//	⍝ overloaded operators on stats.Value
//	S +← 1   ⍝ push a value, Float or Complex or a Vector
//	S ← 2 3 ⍴ø0 ⍝ initialize a stats.Value array
//	S1 +← S2  ⍝ merge two stats.Values
//	⍝ after aggregation, get results:
//	⌈S	⍝ max value, etc...
//	⍝ quantile, stddev, etc...?
//	⍝ histogram text or image?
//	⍝ plots for CDF, Weibull, Whisker...?
//	⍝ serialization?
//	⍝ register function variables on S for quantile, etc...
//	Z ← 0.95 quantile S
//
//	⍝ maybe also use ø on normal types:
//	Z ← øV    ⍝ monadic use: average value
//	Z ← NøV   ⍝ dyadic use: ?? quantile, histogram??
//
//	⍝ random number generator
//	0⋏1	⍝ rectangular random variable [0,1]
//	1⋏0	⍝ normally distributed random variable
//	2⋏0	⍝ bivariate distribution (Complex)
//	⍝ maybe more distributions, maybe other parameters
//	0⋏ MIN MAX
//	1⋏ MEAN VAR or STD
//	2⋏ CMEAN CVAR ROTATE or CORR
package stats

// Value is a scalar value which collects statitical information
// about Float and Complex values, which are added to it.
// It overloads the + Primitive.
type Value struct {
	// TODO
	// use package github.com/ktye/iv/stats
	// Collect all: Real, Complex, Histogram, Digest
}
