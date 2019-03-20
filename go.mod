module github.com/ktye/iv

go 1.12

// These dependencies are needed only for aplextra and cmd/lui
// The main package apl and commands cmd/iv or cmd/apl
// work with the standard library alone.

require (
	github.com/atotto/clipboard v0.1.2 // indirect
	github.com/caio/go-tdigest v2.3.0+incompatible
	github.com/eaburns/T v0.0.0-20190217122806-dbc7887ff15c
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/ktye/plot v0.0.0-20190313115457-f0482d904546
	github.com/ktye/ui v0.0.0-20190210201213-3c0b096c17bb
	github.com/mattn/go-sixel v0.0.0-20190216163338-cdfbdd9946b1
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/soniakeys/quant v1.0.0 // indirect
	github.com/sv/kdbgo v0.11.1-0.20180806165624-70ef73f51093
	golang.org/x/exp v0.0.0-20190316020145-860388717186 // indirect
	golang.org/x/image v0.0.0-20190227222117-0694c2d4d067
	golang.org/x/mobile v0.0.0-20190319155245-9487ef54b94a // indirect
)
