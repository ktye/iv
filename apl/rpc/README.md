# rpc package

The package provides client and server interfaces for communication between distributed
APL instances.

## Server
An server could be implemented as:

```go
package main

import (
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	"github.com/ktye/iv/apl/rpc"
)

func main() {
	a := apl.New(nil)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)
	rpc.Register(a, "")
	
	rpc.ListenAndServe(a, ":1966")
}

```
When running, it listens on port 1966 for connections.

## Client
On a different process, run a normal APL session:

```
	C←rpc→dial ":1966"
	rpc→call (C; "+/"; 5; ⍳10;)
15 20 25 30 35 40
```

The rpc call evaluates the function string in the remote environment
and calls it with the local values on the remote process.

APL value that should be transfered over the wire need to be
registerd to the gob package.
See `init.go` for values that are already registerd.
