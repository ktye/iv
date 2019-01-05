package complete

type Table struct {
	Name   string
	Symbol rune
}

// See VECTOR Vol.19 No.3
// Provisional Unicode Reference
// proposed by Adrian Smith
//
// There are some additional symbols.
// A symbol name must not be a prefix of another,
// for completion to work.

var Tab = []struct {
	Name   string
	Symbol string
}{

	// {"neg", "−"}, // U+2212, use - instead.
	// {"exp", "⋆"}, // U+22C6, use * instead.
	// {"and", "∧"},   // U+2227, use ^ instead.
	// {"tilde", "∼"}, // U+223C, use ~ instead.
	{"neg", "¯"},       // U+00AF
	{"times", "×"},     // U+00D7
	{"div", "÷"},       // U+00F7
	{"not", "~"},       // U+007E
	{"ne", "≠"},        // U+2260
	{"le", "≤"},        // U+2264
	{"ge", "≥"},        // U+2265
	{"or", "∨"},        // U+2228
	{"nor", "⍱"},       // U+2371
	{"nand", "⍲"},      // U+2372
	{"log", "⍟"},       // U+235F
	{"jot", "∘"},       // U+2218
	{"rho", "⍴"},       // U+2374
	{"iota", "⍳"},      // U+2373
	{"each", "¨"},      // U+00A8
	{"stile", "∣"},     // U+007C
	{"rev", "⌽"},       // U+233D
	{"trans", "⍉"},     // U+2349
	{"rot", "⊖"},       // U+2296
	{"avg", "ø"},       // U+00F8
	{"lamp", "⍝"},      // U+235D
	{"as", "←"},        // U+2190
	{"take", "↑"},      // U+2191
	{"ra", "→"},        // U+2192
	{"drop", "↓"},      // U+2193
	{"gu", "⍋"},        // U+234B
	{"gd", "⍒"},        // U+2352
	{"delta", "∆"},     // U+2206
	{"def", "∇"},       // U+2207
	{"intersect", "∩"}, // U+2229
	{"downsh", "∪"},    // U+222A
	{"pick", "⊃"},      // U+2283
	{"encl", "⊂"},      // U+2282
	{"max", "⌈"},       // U+2308
	{"min", "⌊"},       // U+230A
	{"enco", "⊤"},      // U+22A4
	{"dec", "⊥"},       // U+22A5
	{"rtack", "⊢"},     // U+22A2
	{"ltack", "⊣"},     // U+22A3
	{"ibeam", "⌶"},     // U+2336
	{"scan", "⌿"},      // U+233F
	{"slope", "⍀"},     // U+2340
	{"exec", "⍎"},      // U+234E
	{"format", "⍕"},    // U+2355
	{"circle", "○"},    // U+25CB
	{"circjot", "⌾"},   // U+233E
	{"diamond", "⋄"},   // U+22C4
	{"match", "≡"},     // U+2261
	{"alpha", "⍺"},     // U+237A
	{"ualpha", "⍶"},    // U+2376
	{"omega", "⍵"},     // U+2375
	{"uomega", "⍹"},    // U+2379
	{"in", "∊"},        // U+220A
	{"sigma", "σ"},     // U+03C3
	{"domino", "⌹"},    // U+2339
	{"qjot", "⌻"},      // U+233B
	{"sandwich", "⍂"},  // U+2342
	{"quad", "⎕"},      // U+2395
	{"squad", "⌷"},     // U+2337
}

/*
¯ × ÷ ∘ ∣ ∼ ≠ ≤ ≥ ≬ ⌶ ⋆ ⌾ ⍟ ⌽ ⍉ ⍝ ⍦ ⍧ ⍪ ⍫ ⍬ ⍭ ← ↑ → ↓ ∆ ∇ ∧ ∨ ∩ ∪ ⌈ ⌊ ⊤ ⊥ ⊂ ⊃ ⌿ ⍀
⍅ ⍆ ⍏ ⍖ ⍊ ⍑ ⍋ ⍒ ⍎ ⍕ ⍱ ⍲ ○
⍳ ⍴ ⍵ ⍺
⍶ ⍷ ⍸ ⍹ ⍘ ⍙ ⍚ ⍛ ⍜ ⍮
¨ ⍡ ⍢ ⍣ ⍤ ⍥ ⍨ ⍩
⎕ ⍞ ⍠ ⍯ ⍰ ⍌ ⍍ ⍐ ⍓ ⍔ ⍗ ⌷ ⌸ ⌹ ⌺ ⌻ ⌼ ⍁ ⍂ ⍃ ⍄ ⍇ ⍈
*/
