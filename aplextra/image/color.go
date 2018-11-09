package image

import (
	"image/color"

	"github.com/ktye/iv/apl"
)

type Color color.Color

// ToColor converts a value to a color.
func ToColor(a *apl.Apl, v apl.Value) (color.Color, bool) {
	if c, ok := v.(Color); ok {
		return color.Color(c), true
	}
	// TODO look up color names
	// TODO convert r g b or r g b a to a color, convert "#RRGGBBAA"
	return color.Black, true
}
