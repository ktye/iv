package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/funcs"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/aplextra/image"
)

func main() {
	http.Handle("/", file("index.html"))
	http.Handle("/logo.svg", file("logo.svg"))
	http.Handle("/apl", aplHandler{})
	log.Fatal(http.ListenAndServe(":80", nil))
}

type file string

func (f file) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, string(f))
}

type aplHandler struct {
}

func (h aplHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rc := r.Body
	defer rc.Close()

	var buf bytes.Buffer
	a := apl.New(&buf)
	funcs.Register(a)
	operators.Register(a)
	image.Register(a)

	var res Response
	scn := bufio.NewScanner(rc)
	for scn.Scan() {
		line := scn.Text()
		if strings.HasPrefix(line, "\t") {
			fmt.Printf("apl: %s\n", line)
			fmt.Fprintln(&buf, line)
			prog, err := a.Parse(line)
			if err != nil {
				fmt.Fprintln(&buf, err)
				break
			}
			vals, err := a.EvalProgram(prog)
			if err != nil {
				fmt.Fprintln(&buf, err)
				break
			}
			if len(vals) > 0 {
				v := vals[0]
				if img, ok := v.(image.Value); ok {
					if s := img.Encode(); len(s) > 0 {
						res.Image = s
					}
				} else {
					fmt.Fprintln(&buf, v.String(a))
				}
			}
		}
	}

	res.Text = string(buf.Bytes())
	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		log.Print(err)
	}
}

type Response struct {
	Text  string
	Image string
}
