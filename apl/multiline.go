package apl

import (
	"fmt"
	"strings"

	"github.com/ktye/iv/apl/scan"
)

// ParseLines parses multiple lines separated by newline, that may contain continuation lines.
// Continuation lines are only allowed for lambda functions.
func (a *Apl) ParseLines(lines string) (Program, error) {
	b := NewLineBuffer(a)
	v := strings.Split(line, "\n")
	for i, s := range v {
		if ok, err := b.Add(s); err != nil {
			return nil, err
		} else if i == len(v)-1 {
			if ok == false {
				return nil, fmt.Errorf("unbalanced {")
			}
			return lb.Parse()
		}
	}
	return nil, nil
}

// LineBuffer buffers multiline statements for lambda functions.
type LineBuffer struct {
	a      *Apl
	tokens []scan.Token
	level  int
}

func NewLineBuffer(a *Apl) *LineBuffer {
	return &LineBuffer{a: a}
}

// Add a line to the buffer.
// The function returns ok, if the line is complete and can be parsed.
func (b *LineBuffer) Add(line string) (bool, error) {
	if b.a == nil {
		return false, fmt.Errorf("linebuffer is not initialized (no APL)")
	}
	tokens, err := b.a.Scan(line)
	if err != nil {
		b.reset()
		return false, err
	}
	if len(tokens) == 0 {
		return false, nil
	}

	// Join with diamonds. Ommit the diamond if the last token is LeftBrace
	// or the next token is a RightBrace.
	diamond := true
	if len(b.tokens) == 0 {
		diamond = false
	} else if b.tokens[len(b.tokens)-1].T == scan.LeftBrace {
		diamond = false
	}
	if len(tokens) > 0 && diamond == true {
		b.tokens = append(b.tokens, scan.Token{T: scan.Diamond, S: "⋄"})
	}
	b.tokens = append(b.tokens, tokens...)

	for _, t := range tokens {
		if t.T == scan.LeftBrace {
			b.level++
		} else if t.T == scan.RightBrace {
			b.level--
			if b.level < 0 {
				b.reset()
				return false, fmt.Errorf("too many }")
			}
		}
	}
	if b.level == 0 {
		return true, nil
	}
	return false, nil
}

// Parse parses the tokens in the buffer.
// Lines must be pushed to the buffer with Add and Parse should only be called if Add returned true.
func (b *LineBuffer) Parse() (Program, error) {
	defer b.reset()
	return b.a.parse(b.tokens)
}

func (b *LineBuffer) reset() {
	b.level = 0
	if len(b.tokens) > 0 {
		b.tokens = b.tokens[:0]
	}
}
