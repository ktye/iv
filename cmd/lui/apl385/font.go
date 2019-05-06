// Package apl385 embedds APL385 unicode truetype font.
package apl385

import (
	"archive/zip"
	"io/ioutil"
	"strings"
)

// TTF returns the APL385-Unicode font in trutype format.
func TTF() []byte {
	r := strings.NewReader(data)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		panic(err)
	}
	for _, f := range zr.File {
		if f.Name == "APL385.ttf" {
			rc, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer rc.Close()

			b, err := ioutil.ReadAll(rc)
			if err != nil {
				panic(err)
			}
			return b
		}
	}
	panic("cannot find APL385.ttf in embedded zip")
	return nil
}
