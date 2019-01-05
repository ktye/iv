package stats

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// Histogram collects the distribution of real values over discrete bins.
// A variable of the type can be used directly or with specified bins.
// If Bins are nil on the first call to Push, the histogram will automatically
// redistribute the bins as new data arrives.
// MaxBins can be used to control how many bins will be created.
// If Inclusive is true, values that exceed the first or last bin are included in the edge bins.
// Bins are expected to be monotonically increasing and have equal width.
// EqualBins can be used to create bins.
type Histogram struct {
	Bins      []Bin
	MaxBins   int
	Inclusive bool
	auto      bool
	ignored   uint64
}

// Bin is an interval and a counter of how many values are collected.
type Bin struct {
	Min, Max float64
	N        uint64
}

// EqualBins returns N equally distributed Bins from min to max.
// If centered is true, the bins are centered around min and max,
// otherwise min is the left edge and max the right edge of the first and last
// bin respectively.
//	centered=false  centered=true
//        +-+-+-+-+       +-+-+-+-+
//	  | | | | |       |^| | |^|
//	 min     max      min   max
func EqualBins(min, max float64, N int, centered bool) []Bin {
	if N < 2 {
		N = 30
	}
	if min >= max {
		return nil
	}

	width := (max - min) / float64(N)
	if centered {
		width = (max - min) / float64(N-1)
	}
	w2 := width / 2.0

	bins := make([]Bin, N)
	for i := range bins {
		if centered {
			x := scale(float64(i), 0, float64(len(bins)-1), min, max)
			bins[i].Min = x - w2
			bins[i].Max = x + w2
		} else {
			x := scale(float64(i), 0, float64(len(bins)), min, max)
			bins[i].Min = x
			bins[i].Max = x + width
		}
	}
	return bins
}

// Push adds a value to the histogram.
func (h *Histogram) Push(v float64) {
	h.push(v, 1)
}

func (h *Histogram) push(v float64, weight uint64) {
	if len(h.Bins) == 0 {
		h.auto = true
		h.Bins = []Bin{Bin{Min: v, Max: v, N: weight}}
		if h.MaxBins < 2 {
			h.MaxBins = 30
		}
		return
	}
	if h.auto == false {
		if idx, ok := h.index(v); ok {
			h.Bins[idx].N += weight
		} else {
			h.ignored += weight
		}
		return
	} else {
		// The first Bin has zero width.
		// Subsequent numbers could have the same values.
		if len(h.Bins) == 1 {
			if v >= h.Bins[0].Min && v <= h.Bins[0].Max {
				h.Bins[0].N += weight
				return
			}
		}
		extend := false
		min := h.Bins[0].Min
		max := h.Bins[len(h.Bins)-1].Max
		if v < min {
			min = v
			extend = true
		} else if v > max {
			max = v
			extend = true
		}
		if extend {
			var width float64
			old := h.Bins
			min, max, width = niceLimits(min, max, h.MaxBins)
			n := math.Round((max - min) / width)
			h.Bins = EqualBins(min, max, int(n), false)
			for _, b := range old {
				v := (b.Max + b.Min) / 2.0
				if idx, ok := h.index(v); ok {
					h.Bins[idx].N += b.N
				} else {
					h.ignored += b.N
				}
			}
		}
		if idx, ok := h.index(v); ok {
			h.Bins[idx].N += weight
		} else {
			h.ignored += weight
		}
	}
}

// Merge adds the histogram in src to the receiver.
func (dst *Histogram) Merge(src Histogram) {
	for _, b := range src.Bins {
		dst.push((b.Min+b.Max)/2.0, b.N)
	}
}

// Sparkline returns a small string visualization of the histogram using unicode block characters.
// Not all fonts display these characters nicely.
func (h *Histogram) Sparkline() string {
	if len(h.Bins) == 0 {
		return ""
	}

	max := h.Bins[0].N
	for _, b := range h.Bins {
		if b.N > max {
			max = b.N
		}
	}

	levels := []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	var buf strings.Builder
	for _, b := range h.Bins {
		idx := int(scale(float64(b.N), 0, float64(max), 0, 8))
		if idx < 0 {
			idx = 0
		} else if idx > 8 {
			idx = 8
		}
		buf.WriteRune(levels[idx])
	}
	return buf.String()
}

// WriteHorizontal writes a horizontal histogram using unicode block characters up to the given height.
func (h *Histogram) WriteHorizontal(w io.Writer, height int) {
	if len(h.Bins) == 0 {
		return
	}
	if height < 2 {
		height = 10
	}

	max := h.Bins[0].N
	for _, b := range h.Bins {
		if b.N > max {
			max = b.N
		}
	}

	levels := []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇'}
	for i := height; i > 0; i-- {
		for _, b := range h.Bins {
			y := int((b.N * uint64(8*height)) / max)
			lines := y / 8
			rem := y % 8
			if i == lines+1 {
				fmt.Fprintf(w, "%c", levels[rem])
			} else if i > lines {
				fmt.Fprintf(w, ` `)
			} else {
				fmt.Fprintf(w, `█`)
			}
		}
		fmt.Fprintln(w)
	}
}

// WriteVertical writes a vertical histogram using unicode block characters up to the given width.
func (h *Histogram) WriteVertical(w io.Writer, width int) {
	if len(h.Bins) == 0 {
		return
	}

	width *= 8
	max := h.Bins[0].N
	for _, b := range h.Bins {
		if b.N > max {
			max = b.N
		}
	}

	for _, b := range h.Bins {
		x := int((b.N * uint64(width)) / max)
		rem := x % 8
		for i := 0; i < x/8; i++ {
			fmt.Fprintf(w, `█`)
		}
		levels := []rune{'▏', '▎', '▍', '▌', '▋', '▊', '▉'}
		if rem > 0 {
			fmt.Fprintf(w, "%c", levels[rem-1])
		}
		fmt.Fprintln(w)
	}
}

// scale interpolates linearly and returns y.
func scale(x, xmin, xmax, ymin, ymax float64) float64 {
	return ymin + (x-xmin)*(ymax-ymin)/(xmax-xmin)
}

// Index returns the bin index for the given value and if it is valid.
// Index assumes bins to be monotonically increasing and equally spaced.
func (h *Histogram) index(v float64) (int, bool) {
	if len(h.Bins) < 2 {
		return 0, false
	}
	left := h.Bins[0].Min
	right := h.Bins[len(h.Bins)-1].Max
	fidx := scale(v, left, right, 0, float64(len(h.Bins)))
	if fidx < 0 {
		fidx = -1
	} else if fidx > float64(len(h.Bins)) {
		fidx = float64(len(h.Bins))
	}
	idx := int(fidx)
	if idx < 0 {
		if h.Inclusive {
			return 0, true
		} else {
			return -1, false
		}
	} else if idx >= len(h.Bins) {
		if h.Inclusive {
			return len(h.Bins) - 1, false
		} else {
			return 1, false
		}
	}
	return idx, true
}

// NiceLimits returns limits which include data values in [min, max] with nice numbers.
// Reference: Heckbert algorithm (Nice numbers for graph labels)
func niceLimits(min, max float64, maxTicks int) (niceMin, niceMax, spacing float64) {

	extent := niceNumber(max-min, false)
	spacing = niceNumber(extent/float64(maxTicks-1), true)
	niceMin = math.Floor(min/spacing) * spacing
	niceMax = math.Ceil(max/spacing) * spacing

	return niceMin, niceMax, spacing
}

// niceNumber only returns number like 0.01, 0.05, 0.1, 0.5, 1, 1.5, ...
// Twos are excluded, because these cannot be merged when rescaling the histogram.
func niceNumber(extent float64, round bool) float64 {
	exponent := math.Floor(math.Log10(extent))
	fraction := extent / math.Pow10(int(exponent))
	var niceFraction float64
	if round {
		if fraction < 3 {
			niceFraction = 1
		} else if fraction < 7 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	} else {
		if fraction <= 2 {
			niceFraction = 1
		} else if fraction <= 5 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	}
	return niceFraction * math.Pow10(int(exponent))
}
