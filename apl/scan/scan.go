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
	T Type
	S string
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
		if t, err := s.nextToken(); err != nil {
			return nil, err
		} else if t.T == Endl {
			break
		} else {
			s.tokens = append(s.tokens, t)
		}
	}
	return s.tokens, nil
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

/* TODO remove
// ScanNumber scans the next number.
// It may include a minus sign ¯, an exponential part (E or e) and an @ for complex angle.
// Complex numbers are also accepted by real and imag parts as in 1J2.
// Numbers may start with a dot.
func (s *Scanner) scanNumber(cmplx bool) (Token, error) {
	var buf strings.Builder
	valid := false
	exp := 0
	dot := 0
	acceptMinus := true
	if cmplx == false {
		// The angle part cannot start with a minus: 1.234@-234 is not allowed.
		acceptMinus = false
	}
	for {
		r := s.nextRune()
		if r == -1 {
			if valid {
				return Token{T: Number, S: buf.String()}, nil
			}
			return Token{}, fmt.Errorf("cannot scan number: %s", buf.String())
		}

		// A minus ¯ is accepted at the beginning, or after the exponential part.
		// We replace it with -.
		if acceptMinus {
			if r == '¯' {
				buf.WriteRune('-')
				acceptMinus = false
				continue
			}
		}

		if r >= '0' && r <= '9' {
			buf.WriteRune(r)
			acceptMinus = false
			if dot == -1 {
				dot = 0
			}
			valid = true
			continue
		}

		if r == 'e' || r == 'E' {
			if exp == 0 {
				exp++
				dot = -1
			} else {
				return Token{}, fmt.Errorf("cannot scan number: %s%c", buf.String(), r)
			}
			valid = false
			acceptMinus = true
			buf.WriteRune(r)
			continue
		}

		// A dot is only accepted if dot==0.
		if r == '.' {
			if dot == 0 {
				dot = 1
				buf.WriteRune(r)
				valid = false
				continue
			} else {
				return Token{}, fmt.Errorf("cannot scan number: %s%c", buf.String(), r)
			}
		}

		if r == '@' || r == 'J' {
			if cmplx == false {
				return Token{}, fmt.Errorf("cannot scan number: %s%c", buf.String(), r)
			} else {
				buf.WriteRune(r)
				if t, err := s.scanNumber(false); err != nil {
					return t, err
				} else {
					return Token{T: Number, S: buf.String() + t.S}, nil
				}
			}
		}

		if valid {
			s.unreadRune()
			return Token{T: Number, S: buf.String()}, nil
		} else {
			return Token{}, fmt.Errorf("cannot scan number: %s%c", buf.String(), r)
		}
	}
}
*/

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
func (s *Scanner) scanIdentifier() (Token, error) {
	var buf strings.Builder
	first := true
	for {
		r := s.nextRune()
		if AllowedInVarname(r, first) {
			if IsSpecial(r) {
				// Special can only be the first one.
				return Token{T: Identifier, S: string(r)}, nil
			}
			buf.WriteRune(r)
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
	if first && IsSpecial(r) {
		return true
	}
	if first == false && unicode.IsNumber(r) {
		return true
	}
	return r == '_' || unicode.IsLetter(r)
}

// IsSpecial checks for special single-rune variables.
// ⎕ is special because assigning to it, prints the value but does not assign it.
// ⍺ and ⍵ are special, because they are assigned by lambda functions.
// _ could be used as the last value for an interactive interpreter and is used by iv.
func IsSpecial(r rune) bool {
	for _, s := range "⎕⍺⍵_" {
		if r == s {
			return true
		}
	}
	return false
}
