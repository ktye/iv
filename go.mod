module github.com/ktye/iv

go 1.12

// These dependencies are needed only for aplextra and cmd/lui
// The main package apl and commands cmd/iv or cmd/apl
// work with the standard library alone.

require (
	github.com/caio/go-tdigest v2.3.0+incompatible
	github.com/eaburns/T v0.0.0-20190217122806-dbc7887ff15c
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/ktye/plot v0.0.0-20190418145510-8e0ab7910c77
	github.com/ktye/ui v1.0.0
	github.com/leesper/go_rng v0.0.0-20171009123644-5344a9259b21 // indirect
	github.com/mattn/go-sixel v0.0.0-20190216163338-cdfbdd9946b1
	github.com/sv/kdbgo v0.11.1-0.20180806165624-70ef73f51093
	golang.org/x/image v0.0.0-20190417020941-4e30a6eb7d9a
	gonum.org/v1/gonum v0.0.0-20190424212039-2a1643c79af2 // indirect
)
