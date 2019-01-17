# Channels - design draft

Channels are first class types in APL\iv.

They are implemented as a pair of go channels in `apl/channel.go`
```go
type Channel [2]chan Value
```

Channels are move powerfull than file streams (descriptors) in that they cary any APL value not just bytes:
An integer, a multiprecision complex number, an n-dimensional array thereof, a list, a dictionary etc.

A channel itself is also an APL value.
It can be assigned to a variable and passed around.
Just like it is done with functions.

The primary motivation was to provide a simple rpc mechanism, but they offer more.

In Go, channels are primarily used to synchronise concurrency.

The extensions to the primitives introduced come in two forms:
- Simple ones like take and drop perform a single operation and returns. Just like most functions.
- Others create their own go routine and run concurrently

## Pairs
Channels come in pairs of two.
From the APL (client) side `Channel[0]` is the read channel and `Channel[1]` is the write channel,
when communication with an rpc server for example.

When the client want's to hang up, it never closes read `Channel[0]`, instead it closes `Channel[1]`.

This level of detail is not exposed to APL.
We take a channel pair as one value.

## Creating channels
The monadic primitive function `<` is used to create a channel.

```
	<A        return a channel, start a go routine and write A to it.
	          A may be any value (except a channel). Only once. Close afterwards.
	<[5]A     Send A 5 times.
	<[¯1]A    Send A repeatatly. As long as the channel is open.
	          Sending on a channel blocks until there is someone taking the values.
	<[0]A     return a channel that sends A once as an initial value.
	          The channel remains open.
```
Keep in mind, an expression like `A←<[5]2 3⍴⍳6` only assigns a channel to `A`.
The values are have not yet being sent anywhere.
In the background there is a sleeping go-routine waiting for a receiver, we don't have to care about.

## Atomic channel operations
Atomic read and write are implemented with take and drop:
```
	 ↑C       takes one value
	L↑C       takes multiple values (according to L) and reshapes them
	          L is similar to the left argument of ⍴
		  This is most useful, if only scalar values are send over the channel.
        C↓R       drop value R into channel C.
	          Note the anti-symmetry with take: The channel is the left argument here.
	 ↓C       Drop a channel. This closes it.
```

## Concurrent operations

When a concurrent function is called, it does not return directly with a simple value.
It creates a response channel `Z`, starts a go routine and returns the channel.

The statement evaluation can now proceed as usual and pass the channel as an input argument
to the next function. 
But the go routine is still alive.
It is just sleeping until all connections are done and the pipeline unblocks automatically.
```
	 f¨C      read a value from a channel, apply f to it and send the result to the
	          response channel.
	Lf¨C      same as above, but use L as the left argument to f on each call.
	Df¨C      D is also a channel. This is a synchronisation point.
	          The derived function waits until it has input on both channels C and D,
		  the applies the dyadic version of f to them and sends the result.
	Lf⌿C      same as Lf¨C, but skip values for which f returns an EmptyArray.
	          This can be used to implement a filter with a lambda function.
```

## Application to elementary primitive functions

Elementary primitive functions `+-×...` everything in `apl/primtives/elementary.go` could be extended
to act like being called with the each operator implicitly.

# Applications

## IO operations
```
	<`file
	!`ls
	(`wc`-l)!<`file
	`cat!<!`ls
	`dst<<`src
```

## Paste two files
```
	(<`file2)+<`file1
```

## Feedback loop
```
     ⎕←+(G)-f---(+)- <[5]1
        |         |
        +->-(F)---+
	
	f←-         ⍝ system function is negation
	F←<[0]1     ⍝ initialize feedback channel with 1
	G←f¨F+<[5]1 ⍝ main branch
	F<⎕←G       ⍝ connect feedback channel and monitor values
¯2
 1
¯2
 1
```

## Ping pong or laser
```
	⊢⎕←f⊣<1
```











