package stats

import (
	"math"

	"github.com/caio/go-tdigest"
)

// Digest combines Real with a TDigest.
type Digest struct {
	Real
	tdigest.TDigest
}

// Quantile returns the quantile for the given probability [0, 1].
func (r Digest) Quantile(p float64) float64 {
	if r.Num == 0 {
		return math.NaN()
	}
	return r.TDigest.Quantile(p)
}

// Push adds a float64 to the digest.
func (r *Digest) Push(v float64) {
	r.Real.Push(v)
	if r.Num == 1 {
		t, _ := tdigest.New()
		r.TDigest = *t
	}
	r.TDigest.Add(v)
}

// Merge combines two RealDigests into dst.
func (dst *Digest) Merge(add *Digest) {
	dst.Real.Merge(&add.Real)
	add.TDigest.Merge(&dst.TDigest)
}
