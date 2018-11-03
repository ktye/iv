package main

import "io"

// Tabfilter returns a reader, which passes data from the
// underlying reader, only following a tab character, until a newline.
type tabfilter struct {
	r   io.Reader
	hot bool
}

func (t *tabfilter) Read(p []byte) (int, error) {
	n, err := t.r.Read(p)
	if err != nil {
		return 0, err
	}
	tail := 0
	for i := 0; i < n; i++ {
		if t.hot {
			p[tail] = p[i]
			tail++
			if p[i] == '\n' {
				t.hot = false
			}
		} else if p[i] == '\t' {
			t.hot = true
		}
	}
	return tail, nil
}
