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
	a         *apl.Apl
	Separator byte
	Rank      int
}

func (ts *InputParser) Next() (apl.Value, int, bool, error) {
	var eof bool

	s, err := ts.Reader.ReadString(ts.Separator)
	if err != nil && err != io.EOF {
		return nil, 0, false, err
	} else if err == io.EOF {
		eof = true
	} else {
		s = s[:len(s)-1] // remove delimiter
	}

	numexp, err := ts.a.Tower.Parse(s)
	if err != nil {
		return nil, 0, false, err
	}
	v := numexp.Number

	if eof {
		return v, ts.Rank, true, nil
	}

	separators := 1
	for {
		if b, err := ts.Reader.ReadByte(); err == io.EOF {
			return v, ts.Rank, true, nil
		} else if err != nil {
			return nil, 0, false, err
		} else if b == ts.Separator {
			separators++
		} else {
			if err := ts.Reader.UnreadByte(); err != nil {
				return nil, 0, false, err
			}
			return v, separators, false, nil
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
