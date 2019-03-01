package numbers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ktye/iv/apl"
)

var y0, y1k time.Time

func init() {
	y1k, _ = time.Parse("2006.01.02", "1000.01.01")
	y0, _ = time.Parse("15h04", "00h00")
}

// Time holds both a time stamp and a duration in a single number type.
// Durations are identified as time stamps before year 1000 (y1k).
// The parser accepts both, durations and time stamps.
// When times and other number types are mixed, the other number types
// are identified as seconds and upgraded.
type Time time.Time

func ParseTime(s string) (apl.Number, bool) {
	s = strings.Replace(s, "¯", "-", -1)
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

func (t Time) String(f apl.Format) string {
	format, minus := getformat(f, t)
	if t1 := time.Time(t); t1.Before(y1k) {
		s := t1.Sub(y0).String()
		if minus == false {
			s = strings.Replace(s, "-", "¯", -1)
			return s
		}
	}

	if format == "W" {
		year, wk := time.Time(t).ISOWeek()
		return fmt.Sprintf("%04dw%02d", year, wk)
	} else if format == "Q" {
		_, year, quarter := roundQuarter(time.Time(t))
		return fmt.Sprintf("%04dQ%d", year, quarter)
	}

	if format == "" {
		format = "2006.01.02T15.04.05.000"
	}

	return time.Time(t).Format(format)
}
func (t Time) Copy() apl.Value { return t }

func (t Time) ToIndex() (int, bool) {
	return 0, false
}

func (t Time) Less(R apl.Value) (apl.Bool, bool) {
	return apl.Bool(time.Time(t).Before(time.Time(R.(Time)))), true
}

func (t Time) Add() (apl.Value, bool) {
	return t, true
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

func (t Time) Sub() (apl.Value, bool) {
	if t0 := time.Time(t); t0.Before(y1k) {
		d := t0.Sub(y0)
		return Time(y0.Add(-d)), true
	}
	return nil, false
}

// Sub2 returns a duration depending on it's arguments:
// If both are a duration, it is the difference.
// If both are a time, it is the difference.
// If the first is a time and the second a duration, it's a time before.
// If the first is a duration and the second a time, it is not accepted.
func (t Time) Sub2(R apl.Value) (apl.Value, bool) {
	t0 := time.Time(t)
	t1 := time.Time(R.(Time))
	if t0.After(y1k) && t1.After(y1k) {
		return Time(y0.Add(t0.Sub(t1))), true
	} else if t0.Before(y1k) && t1.Before(y1k) {
		return Time(y0.Add(t0.Sub(y0)).Add(-t1.Sub(y0))), true
	} else if t0.After(y1k) && t1.Before(y1k) {
		return Time(t0.Add(-t1.Sub(y0))), true
	}
	return nil, false
}

// Duration returns if the time value is a duration and it's value.
func (t Time) Duration() (time.Duration, bool) {
	if time.Time(t).Before(y1k) {
		return time.Time(t).Sub(y0), true
	}
	return time.Duration(0), false
}

func (t Time) Mul() (apl.Value, bool) {
	if t0 := time.Time(t); t0.Before(y0) {
		return apl.Int(-1), true
	} else if t0.After(y0) {
		return apl.Int(1), true
	}
	return apl.Int(0), true
}

// Multiplication is allowed for durations only and applied to seconds.
func (t Time) Mul2(R apl.Value) (apl.Value, bool) {
	t0 := time.Time(t)
	t1 := time.Time(R.(Time))
	if t0.After(y1k) || t1.After(y1k) {
		return nil, false
	}
	s0 := t0.Sub(y0).Seconds()
	s1 := t1.Sub(y0).Seconds()
	return Time(y0.Add(time.Duration(int64(1e9 * (s0 * s1))))), true
}

/* Does Div make any sense?
func (t Time) Div() (apl.Value, bool) {
	if t0 := time.Time(t); t0.Before(y1k) {
		s0 := t0.Sub(y0).Seconds()
		return Time(y0.Add(time.Duration(int64(1e9 / s0)))), true
	}
	return nil, false
}
*/

// Division is allowed for durations only and applied to seconds.
func (t Time) Div2(R apl.Value) (apl.Value, bool) {
	t0 := time.Time(t)
	t1 := time.Time(R.(Time))
	if t0.After(y1k) || t1.After(y1k) {
		return nil, false
	}
	s0 := t0.Sub(y0).Seconds()
	s1 := t1.Sub(y0).Seconds()
	return Time(y0.Add(time.Duration(int64(1e9 * (s0 / s1))))), true
}

func (t Time) Abs() (apl.Value, bool) {
	if t0 := time.Time(t); t0.Before(y0) {
		return Time(y0.Add(-t0.Sub(y0))), true
	}
	return t, true
}

func (t Time) Floor() (apl.Value, bool) {
	return Time(time.Time(t).Truncate(time.Second)), true
}

func (t Time) Ceil() (apl.Value, bool) {
	return Time(time.Time(t).Add(500 * time.Millisecond).Truncate(time.Second)), true
}

func (t Time) Floor2() (apl.Value, bool) {
	return Time(time.Time(t).Truncate(time.Second)), true
}

func (t Time) Ceil2() (apl.Value, bool) {
	return Time(time.Time(t).Add(500 * time.Millisecond).Truncate(time.Second)), true
}

func (t Time) Round(period string) (Time, error) {
	if d, ok := t.Duration(); ok {
		return roundDuration(d, period)
	}
	t1 := time.Time(t)
	Y, M, D := t1.Date()
	h := t1.Hour()
	m := t1.Minute()
	s := t1.Second()
	switch period {
	case "Y":
		return Time(time.Date(Y, 1, 1, 0, 0, 0, 0, time.UTC)), nil
	case "M":
		return Time(time.Date(Y, M, 1, 0, 0, 0, 0, time.UTC)), nil
	case "D":
		return Time(time.Date(Y, M, D, 0, 0, 0, 0, time.UTC)), nil
	case "h":
		return Time(time.Date(Y, M, D, h, 0, 0, 0, time.UTC)), nil
	case "m":
		return Time(time.Date(Y, M, D, h, m, 0, 0, time.UTC)), nil
	case "s":
		return Time(time.Date(Y, M, D, h, m, s, 0, time.UTC)), nil
	case "Q":
		rt, _, _ := roundQuarter(t1)
		return Time(rt), nil
	case "W":
		year, wk := t1.ISOWeek()
		return Time(roundWeek(year, wk, time.UTC)), nil
	default:
		return t, fmt.Errorf("cannot round time to %s", period)
	}

	return t, fmt.Errorf("TODO: round time")
}

func roundQuarter(t time.Time) (time.Time, int, int) {
	Y, _, _ := t.Date()
	if t.Before(time.Date(Y, time.April, 1, 0, 0, 0, 0, time.UTC)) {
		return time.Date(Y, time.January, 1, 0, 0, 0, 0, time.UTC), Y, 1
	} else if t.Before(time.Date(Y, time.July, 1, 0, 0, 0, 0, time.UTC)) {
		return time.Date(Y, time.April, 1, 0, 0, 0, 0, time.UTC), Y, 2
	} else if t.Before(time.Date(Y, time.October, 1, 0, 0, 0, 0, time.UTC)) {
		return time.Date(Y, time.July, 1, 0, 0, 0, 0, time.UTC), Y, 3
	}
	return time.Date(Y, time.October, 1, 0, 0, 0, 0, time.UTC), Y, 4
}

func roundWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year {
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week {
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}

func roundDuration(t time.Duration, s string) (Time, error) {
	var d time.Duration
	switch s {
	case "h":
		d = time.Hour
	case "m":
		d = time.Minute
	case "s":
		d = time.Second
	case "ms":
		d = time.Millisecond
	case "us", "μs":
		d = time.Microsecond
	case "ns":
		d = time.Nanosecond
	default:
		return Time(y0), fmt.Errorf("round duration: not supported: %s", s)
	}
	return Time(y0.Add(t.Round(d))), nil
}

// Not supported by elementary arithmetics on time numbers:
// - Add non-constant intervals to time, e.g. 2016.01.01 + 1 year (go: time.AddDate)
//	Currently we can only do:
//      	2015.01.01 + 365×24h
//	2016.12.31T00.00.00.000
//	       2015.01.01 + 365×24h
//	2016.01.01T00.00.00.000 // This does not take leap years into account
// - Intervals: Year, Month, Quarter, Calendar week
// 	2015Q3, 2015W12
