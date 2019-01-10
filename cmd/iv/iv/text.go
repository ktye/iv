package iv

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ktye/iv/apl"
)

// Iv is also tested from primitives/apl_test.go

// InputParser parses the input stream into apl.Numbers.
type InputParser struct {
	*bufio.Reader
	Apl       *apl.Apl
	Separator byte
	Rank      int
	state
}

type state struct {
	init bool
	eof  bool
	max  int
	v    apl.Value
	n    int
}

func (s *state) push(v apl.Value, n int) {
	if n > s.max {
		s.max = n
	}
	s.v = v
	s.n = n
}

// Next returns the next scalar apl.Value from the stream.
// The value is converted from a string representation using the current apl.Tower.
// It returns the scalar and the number of separators following it.
// EOF is reported as io.EOF and a nil Value.
// The number of separators following the last valid value is always the max
// number of separators observed over all calls, independend on the actual value.
func (ts *InputParser) Next() (apl.Value, int, error) {
	if ts.state.eof == true {
		return nil, 0, io.EOF
	} else if ts.state.init == false {
		ts.state.init = true
		v, n, e := ts.next()
		if e != nil {
			return v, n, e
		}
		ts.state.push(v, n)
	}
	v, n, e := ts.next()
	if e == io.EOF {
		ts.eof = true
		return ts.state.v, ts.state.max, nil
	} else if e != nil {
		return nil, 0, e
	}
	ar, nr := ts.state.v, ts.state.n
	ts.state.push(v, n)
	return ar, nr, nil
}
func (ts *InputParser) next() (apl.Value, int, error) {
	s, err := ts.Reader.ReadString(ts.Separator)
	if err != nil {
		return nil, 0, err
	} else {
		s = s[:len(s)-1] // remove delimiter
	}

	numexp, err := ts.Apl.Tower.Parse(s)
	if err != nil {
		return nil, 0, err
	}
	v := numexp.Number

	separators := 1
	for {
		if b, err := ts.Reader.ReadByte(); err == io.EOF {
			return v, 0, nil
		} else if err != nil {
			return nil, 0, err
		} else if b == ts.Separator {
			separators++
		} else {
			if err := ts.Reader.UnreadByte(); err != nil {
				return nil, 0, err
			}
			return v, separators, nil
		}
	}
}

// tabularText is used to convert text in column layout to a stream of
// scalar strings, separated by newlines.
// Multiple spaces are treated as a single one.
// Single newlines are translated into double newlines.
// "a b  c\nd e f" -> "a\nb\nc\n\nd\ne\nf".
func tabularText(src io.Reader) io.Reader {
	r, w := io.Pipe()
	go translateText(src, w)
	return r
}

// TabularText translates blanks in src to newlines.
func translateText(r io.Reader, w *io.PipeWriter) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if t == "" {
			fmt.Fprintln(w)
		} else {
			if strings.Contains(t, `"`) || strings.Contains(t, "'") {
				// TODO: handle quoted strings.
				// Do not split inside those.
			}
			fmt.Fprintln(w, strings.Join(strings.Fields(t), "\n")+"\n")
		}
	}
	w.Close()
}
