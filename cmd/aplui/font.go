package main

import (
	"fmt"

	"github.com/golang/freetype/truetype"
	duitdraw "github.com/ktye/duitdraw"
	"github.com/ktye/iv/font"
)

// registerFont parses APL385.ttf and registers it in duitdraw under the name APL385.
// This must be called before starting duit.
func registerFont(size int) {
	f, err := truetype.Parse(font.APL385())
	if err != nil {
		fmt.Println(err)
		return
	}

	// It only works if the DPI value is not changed
	// from the built-in default.
	opt := truetype.Options{
		Size: float64(size),
		DPI:  float64(duitdraw.DefaultDPI),
	}
	face := truetype.NewFace(f, &opt)

	id := duitdraw.FaceID{
		Name: "APL385",
		Size: size,
		DPI:  duitdraw.DefaultDPI,
	}

	duitdraw.RegisterFont(id, face)
}
