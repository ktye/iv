package iv

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// TabularText translates blanks in src to newlines.
// Multile blanks are treated as a single one.
func TabularText(src io.Reader) io.Reader {
	r, w := io.Pipe()

	go translateText(src, w)
	return r
}

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
	// can we do anything about scanner.Err()?
	w.Close()
}
