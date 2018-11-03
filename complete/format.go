package complete

import "strings"

// Format replaces APL symbol names with their symbols and removes space around them.
// Format respectes double quoted strings.
func Format(s string) string {
	r := strings.NewReader(s)

	type token struct {
		s string
		t int
	}
	var vector []token
	var buf strings.Builder

	var space, text, other, match = 0, 1, 2, 3
	quoted := false
	which := func(c rune) int {
		if quoted {
			if c != '"' {
				return other
			} else {
				quoted = false
				return other
			}
		} else if c == '"' {
			quoted = true
			return other
		}

		if c == ' ' {
			return space
		} else if c >= 'a' && c <= 'z' {
			return text
		}
		return other
	}

	// Tokenize string into vector.
	// A token has on of the types: space, text, other.
	var state int = -1
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			break
		}

		v := which(c)
		if v == state {
			buf.WriteRune(c)
		} else if state != -1 {
			vector = append(vector, token{buf.String(), state})
			buf.Reset()
			buf.WriteRune(c)
			state = v
		} else {
			state = v
			buf.WriteRune(c)
		}
	}
	if buf.Len() > 0 {
		vector = append(vector, token{buf.String(), state})
	}

	tab := make(map[string]string)
	for _, e := range Tab {
		tab[e.Name] = e.Symbol
	}

	// Iterate over the vector and replace matched text with it's APL symbol.
	// Mark the token as a match.
	for i, t := range vector {
		if t.t == text {
			if s, ok := tab[t.s]; ok {
				vector[i].t = match
				vector[i].s = s
			}
		}
	}

	// Write tokens to out, but ommit all space before and after a match.
	var out strings.Builder
	for i, t := range vector {
		if t.t == space {
			if i > 0 && vector[i-1].t == match {
				continue
			}
			if i < len(vector)-1 && vector[i+1].t == match {
				continue
			}
		}
		out.WriteString(t.s)
	}
	return out.String()
}
