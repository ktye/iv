package operators

import (
	"math"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
)

// identityItem returns the identity item for the given function f, when
// the function is applied as f/⍳0.
func identityItem(f apl.Value) apl.Value {
	// Table from APL2: p 211, DyaRef p 170
	if p, ok := f.(apl.Primitive); ok {
		switch p {
		case "+", "-", "|", "∨", "<", ">", "≠", "⊤", "∪", "⌽", "⊖":
			return apl.Int(0)
		case "×", "÷", "*", "!", "^", "∧", "≤", "=", "≥", "/", "⌿", `\`, `⍀`:
			return apl.Int(1)
		case "⌊":
			return numbers.Float(-math.MaxFloat64)
		case "⌈":
			return numbers.Float(math.MaxFloat64)
		}
	}
	return nil
}
