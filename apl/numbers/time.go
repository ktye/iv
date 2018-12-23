package numbers

import (
	"time"

	"github.com/ktye/iv/apl"
)

var y0, y1k time.Time

func init() {
	y1k, _ = time.Parse("2006.01.02", "1000.01.01")
	y0, _ = time.Parse("15h04", "00h00")
}

type Time time.Time

func ParseTime(s string) (apl.Number, bool) {
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return Time(t), true
		}
	}
	if d, err := time.ParseDuration(s); err == nil {
		return Time(y0.Add(d)), true
	}

	return nil, false
}

// We cannot separate with colons.
// Mon Jan 2 15:04:05 -0700 MST 2006
var layouts = []string{
	"2006.01.02",
	"2006.01.02T15.04",
	"2006.01.02T15.04",
	"2006.01.02T15.04.05", // This accepts also fractional seconds.
}

func (t Time) String(a *apl.Apl) string {
	if t1 := time.Time(t); t1.Before(y1k) {
		return t1.Sub(y0).String()
	}

	format, _ := getformat(a, t, "2006.01.02T15.04.05.000")
	return time.Time(t).Format(format)
}

func (t Time) ToIndex() (int, bool) {
	return 0, false
}

// Add2 adds two times.
// At least one of the times must be a duration (before y1k).
func (t Time) Add2(R apl.Value) (apl.Value, bool) {
	t0 := time.Time(t)
	t1 := time.Time(R.(Time))
	if t0.After(y1k) && t1.After(y1k) {
		return nil, false
	} else if t0.After(y1k) {
		return Time(t0.Add(t1.Sub(y0))), true
	} else {
		return Time(t1.Add(t0.Sub(y0))), true
	}
}

func (t Time) Sub2(R apl.Value) (apl.Value, bool) {
	return Time(y0.Add(time.Time(t).Sub(time.Time(R.(Time))))), true
}

func (t Time) Less(R apl.Value) (apl.Value, bool) {
	return apl.Bool(time.Time(t).Before(time.Time(R.(Time)))), true
}
