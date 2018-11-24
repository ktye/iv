package apl

import "testing"

func TestIdx(t *testing.T) {
	testCases := []struct {
		shape []int
		idx   []int
		n     int
	}{
		{[]int{2}, []int{1}, 1},
		{[]int{2, 2, 4}, []int{1, 1, 3}, 15},
		{[]int{3, 4, 2}, []int{1, 3, 0}, 14},
	}

	// Test IdxConverter
	for _, tc := range testCases {
		ic, idx := NewIdxConverter(tc.shape)
		ic.Indexes(tc.n, idx)
		for i := range idx {
			if idx[i] != tc.idx[i] {
				t.Fatalf("expected %v got %v", tc.idx, idx)
			}
		}
		n := ic.Index(idx)
		if n != tc.n {
			t.Fatalf("expected %d got %d", tc.n, n)
		}
	}

	// Test ArraySize and IncArrayIndex
	for _, tc := range testCases {
		ar := GeneralArray{Dims: tc.shape}
		ic, idx := NewIdxConverter(tc.shape)
		for i := 0; i < ArraySize(ar); i++ {
			n := ic.Index(idx)
			if n != i {
				t.Fatalf("expected %d got %d", i, n)
			}
			IncArrayIndex(idx, tc.shape)
		}
		for i := range idx {
			if idx[i] != 0 {
				t.Fatalf("idx should be zeros: %v", idx)
			}
		}
	}
}
