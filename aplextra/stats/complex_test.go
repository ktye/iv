package stats

import (
	"fmt"
	"testing"
)

func TestComplexCovariance(t *testing.T) {
	data := []complex128{-2, -1i, 2 + 1i}

	var s Complex
	for _, z := range data {
		s.Push(z)
	}
	if s.Covariance() != 1 {
		t.Fatal(fmt.Errorf("covariance should be 1"))
	}
}
