// Package scan contains the tokenizer for iv/apl
package scan

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token struct {
	T   Type
	S   string
	Pos int
}

type Type int

const (
	Endl       Type = iota
	Symbol          // single rune APL symbol
	Number          // 1 1.23 1.234E12 -1.234 1.234@156
	String          // "quoted text", "escape "" with double quotes"
	Chars           // 'rune vector' in single quotes
	Identifier      // alpha_3
	LeftParen       // (
	RightParen      // )
	LeftBrack       // ]
	RightBrack      // ]
	LeftBrace       // {
	RightBrace      // }
	Colon           // :
	Semicolon       // ;
	Self            // ∇
	Diamond         // ⋄
)

// Scanner can split APL input into tokens.
// SetSymbols must be called before using the scanner.
type Scanner struct {
	input   string
	tokens  []Token
	symbols map[rune]string
	pos     int
	width   int
}

// SetSymbols initializes the Scanner to recognize the given APL symbols.
// In the map, both the rune and the string are the same letter.
func (s *Scanner) SetSymbols(symbols map[rune]string) {
	s.symbols = symbols
}

// Scan returns the tokens from one line of APL input.
func (s *Scanner) Scan(line string) ([]Token, error) {
	s.input = line
	s.pos = 0
	s.width = 0
	s.tokens = nil
	for {
		pos := s.pos
		if t, err := s.nextToken(); err != nil {
			return nil, err
		} else if t.T == Endl {
			break
		} else {
			t.Pos = pos
			s.tokens = append(s.tokens, t)
		}
	}
	return s.tokens, nil
}

func (t Type) String() string {
	var s string
	switch t {
	case LeftParen:
		s = "("
	case RightParen:
		s = ")"
	case LeftBrack:
		s = "["
	case RightBrack:
		s = "]"
	case LeftBrace:
		s = "{"
	case RightBrace:
		s = "}"
	default:
		// The other type are not printed.
		s = "?"
	}
	return s
}

func (t Token) String() string {
	var s string
	switch t.T {
	case Symbol:
		s = "X"
	case Number:
		s = "N"
	case String:
		s = "S"
	case Chars:
		s = "c"
	case Identifier:
		s = "I"
	case LeftParen, RightParen, LeftBrack, RightBrack, LeftBrace, RightBrace:
		s = "P"
	case Colon:
		s = ":"
	case Semicolon:
		s = ";"
	case Self:
		s = "∇"
	case Diamond:
		s = "⋄"
	case Endl:
		s = "NULL"
	default:
		panic("no stringer for token type " + strconv.Itoa(int(t.T)))
	}
	return s + "(" + t.S + ")"
}

func PrintTokens(t []Token) string {
	if len(t) == 0 {
		return "[]"
	}
	v := make([]string, len(t))
	for i := range t {
		v[i] = t[i].String()
	}
	return "[" + strings.Join(v, ",") + "]"
}

func (s *Scanner) nextRune() rune {
	if s.pos >= len(s.input) {
		return -1
	}
	r, w := utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += w
	s.width = w
	return r
}

// unreadRune can be called once to set the position to the last state.
func (s *Scanner) unreadRune() {
	s.pos -= s.width
	s.width = 0
}

// Peek is like nextRune, but does not advance the position.
func (s *Scanner) peek() rune {
	if s.pos >= len(s.input) {
		return -1
	}
	r, _ := utf8.DecodeRuneInString(s.input[s.pos:])
	return r
}

func (s *Scanner) nextToken() (Token, error) {
	for {
		r := s.nextRune()
		if r == -1 {
			return Token{T: Endl}, nil
		}

		if r == '"' || r == '\'' {
			return s.scanString(r)
		}

		// A number starts with [0-9] or "." or "¯".
		if (r >= '0' && r <= '9') || r == '.' || r == '¯' {
			if r == '.' {
				// If it start with . a digit must follow,
				// otherwise it could be the dot operator.
				if n := s.peek(); n >= '0' && n <= '9' {
					s.unreadRune()
					return s.scanNumber(true)
				}
			} else {
				s.unreadRune()
				return s.scanNumber(true)
			}
		}

		// Registered APL symbols.
		if r, ok := s.symbols[r]; ok {
			return Token{T: Symbol, S: r}, nil
		}

		if AllowedInVarname(r, true) {
			s.unreadRune()
			return s.scanIdentifier()
		}

		switch r {
		case '(':
			return Token{T: LeftParen, S: "("}, nil
		case ')':
			return Token{T: RightParen, S: ")"}, nil
		case '[':
			return Token{T: LeftBrack, S: "["}, nil
		case ']':
			return Token{T: RightBrack, S: "]"}, nil
		case '{':
			return Token{T: LeftBrace, S: "{"}, nil
		case '}':
			return Token{T: RightBrace, S: "}"}, nil
		case ':':
			return Token{T: Colon, S: ":"}, nil
		case ';':
			return Token{T: Semicolon, S: ";"}, nil
		case '∇':
			return Token{T: Self, S: "∇"}, nil
		case '⋄':
			return Token{T: Diamond, S: "⋄"}, nil
		case ' ', '\r', '\t':
			continue // ignore whitespace, newline should not be present.
		case '⍝':
			return Token{T: Endl}, nil
		default:
			return Token{}, fmt.Errorf("unexpected rune: %U (%d %c)", r, r, r)
		}
	}
	return Token{T: Endl}, nil
}

// ScanNumber scans the next number.
// It starts with a digit, ¯ or dot
// and stops before a character is not digit, a-zA-Z, dot or ¯.
// Valid number formats are not known to the scanner.
// Parsing is done by the parser with the current numerical tower.
func (s *Scanner) scanNumber(cmplx bool) (Token, error) {
	var buf strings.Builder
	for {
		r := s.nextRune()
		if r == -1 {
			return Token{T: Number, S: buf.String()}, nil
		} else if r >= '0' && r <= '9' {
			buf.WriteRune(r)
		} else if r >= 'a' && r <= 'z' {
			buf.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			buf.WriteRune(r)
		} else if r == '.' {
			buf.WriteRune(r)
		} else if r == '¯' {
			buf.WriteRune(r)
		} else {
			s.unreadRune()
			return Token{T: Number, S: buf.String()}, nil
		}
	}
}

// ScanString returns the next token as charstr or chars depending on the quoteChar.
// " scans the string as charstr and ' as chars.
// There is currently no way to escape newlines etc.
func (s *Scanner) scanString(quoteChar rune) (Token, error) {
	var str strings.Builder
	for {
		r := s.nextRune()
		if r == -1 {
			return Token{}, fmt.Errorf("unquoted string: %q", str.String())
		}

		if r == quoteChar {
			// Two quotes escapes a single one.
			if s.peek() == quoteChar {
				r = s.nextRune()
				str.WriteRune(r)
				continue
			} else {
				// The quotes are not part of token.s.
				if quoteChar == '"' {
					return Token{T: String, S: str.String()}, nil
				} else {
					return Token{T: Chars, S: str.String()}, nil
				}
			}
		} else {
			str.WriteRune(r)
		}
	}
}

// An identifier may start with _ or a unicode letter.
// Later characters may also be digits.
// A → may be present within an identifier.
func (s *Scanner) scanIdentifier() (Token, error) {
	var buf strings.Builder
	first := true
	arrow := false
	for {
		r := s.nextRune()
		if AllowedInVarname(r, first) {
			buf.WriteRune(r)
		} else if r == '→' && arrow == false {
			buf.WriteRune(r)
			arrow = true
		} else {
			if r != -1 {
				s.unreadRune()
			}
			if first {
				return Token{}, fmt.Errorf("cannot scan empty identifier")
			} else {
				return Token{T: Identifier, S: buf.String()}, nil
			}
		}
		first = false
	}
}

// AllowedinVarname returns true if the rune is allowed in a variable name.
func AllowedInVarname(r rune, first bool) bool {
	if first && strings.IndexRune("⎕⍺⍵", r) != -1 {
		return true
	}
	if first == false && unicode.IsNumber(r) {
		return true
	}
	return r == '_' || unicode.IsLetter(r)
}
