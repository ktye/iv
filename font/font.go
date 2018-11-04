// Package font embedds APL385 unicode.
package font

import (
	"archive/zip"
	"io/ioutil"
	"strings"
)

// APL385 returns the font in ttf format.
func APL385() []byte {

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
