// Package image implements an image as apl value.
package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	img "image"
	"image/png"

	"github.com/ktye/iv/apl"
)

// Value implements an apl.Value which is an image.
type Value struct {
	img.Image
}

func (v Value) String(a *apl.Apl) string {
	rect := v.Bounds()
	return fmt.Sprintf("image %dx%d", rect.Dx(), rect.Dy())
}

func (v Value) Eval(a *apl.Apl) (apl.Value, error) {
	return v, nil
}

// Encode returns the image as a base64 encoded png.
func (v Value) Encode() string {
	var buf bytes.Buffer
	if err := png.Encode(&buf, v.Image); err != nil {
		return ""
	}

	s := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + s
}
