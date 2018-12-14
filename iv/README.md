# iv stream processor
## Intro
Iv is similar to awk.
It reads data from stdin, applies a list of rules and writes the output to stdout.
```
cat data | iv rule1 rule2 ... > out
```
Instead of using a single argument for the program, like awk
```
cat data | awk 'BEGIN{...} /pattern/{rule1} /pattern2/{rule2} END{...}/'
```
each iv rule is given as a single argument, with the special `-bBEGIN` block.
The language is APL and a rule is composed of a `conditional` and an `expression`
```
cat data | iv -b 'X←0' 'N>10:X+←1' 'E=1:_'
```

## n-dimensional input data
The input data is assumed to be a stream of scalars which form an n-dimensional object.
It contains separators, if the stream reached the border of a dimension.

Consider the stream of 12 scalar values `X`, with a `|` for a separator 
```
X|X|X||X|X|X||X|X|X||X|X|X
```
It is a 2-dimensional object and with APL vocabulary, it has shape `4 3`.
Normally the first axis is very large, or endless.

It can also have a higher dimension, with multiple axes that are endless, such as `* * 2 3`.

Iv acts on this data, using the **execution rank** given as the `-r` option.
The default is `-r1`, which means iv collects all data, until a line is full,
and passes the vector to the APL program, which executes it's rules in order.

This mean: APL always sees an array of rank N, if iv is called with `-rN`

To operate on each matrix (or table), we can use
```
cat data | iv -r2 +/_
```

The input data can have more dimensions as the execution rank.
In this case, we might be interested if with the current object, a higher dimension was closed.
This information is send to APL as **termination level** in the variable `E`.
`E` is `1`, if one level above the execution rank just finished.
E.g. if iv operates in line mode (rank 1), `E` is `1`, at the end of each table,
that is if two subsequent newlines are in the stream instead of a single one.

## Variables and control flow
For each execution, the following variables are set before calling the APL rules.
They can be used in the rules or their conditionals:
- N: current record number
- E: termination level
- EOF: bool, true for last value

If the APL program wants to skip subsequent rules, it can assign to the variable `NEXT`
```
NEXT←1
```

## Backends
`iv/apl` is the default backend. But there are more given on the `-a` option:
- [Ivy](http://robpike.io/ivy) Rob Pike's big number calculator (`-ay` or `-aivy`)
- [J](http://www.jsoftware.com) (`-aj`)
- [Kona](https://github.com/kevinlawler/kona) an implementation of K (`-ak` or `-akona`)
- [Klong](http://t3x.org/klong/) (`-akg` or `-aklong`)

iv/apl and ivy are built in.

J, K and Klong are started in the background as an external process.
Their executables must be on the path, which are: `jconsole`, `k` and `kg`, respectively.

As the languages are not compatible, also the iv programs differ.
See the source for details, in `extern/`.

## Embedding
The command line program `cmd/iv` is only an example.
Iv can be embedded in any go program.
The data does not need to be a text stream. It can be connected to any data source by implementing the `iv.Nexter` interface.
```go
// Nexter can return the next scalar from the data stream.
// The call to next should returns the scalar,
// the number of separators following it,
// if it is the last value (EOF)
// and a possible error different from io.EOF.
//
// See TextStream for the default implementation.
type Nexter interface {
	Next() (Scalar, int, bool, error)
}
```

## APL Symbols
APL is a perfect fit for a command line application. It's terse. It fits before the line ends.

But how to enter the all the symbols?
¯ × ÷ ∘ ∣ ∼ ≠ ≤ ≥ ≬ ⌶ ⋆ ⌾ ⍟ ⌽ ⍉ ⍝ ⍦ ⍧ ⍪ ⍫ ⍬ ⍭ ← ↑ → ↓ ∆ ∇ ∧ ∨ ∩ ∪ ⌈ ⌊ ⊤ ⊥ ⊂ ⊃ ⌿ ⍀
⍅ ⍆ ⍏ ⍖ ⍊ ⍑ ⍋ ⍒ ⍎ ⍕ ⍱ ⍲ ○
⍳ ⍴ ⍵ ⍺ ⍶ ⍷ ⍸ ⍹ ⍘ ⍙ ⍚ ⍛ ⍜ ⍮ ¨ ⍡ ⍢ ⍣ ⍤ ⍥ ⍨ ⍩

If you don't want a special keyboard configuration, there are also [completion scripts](../complete) for `bash` and `zsh`.

## Parallel example
- TODO: 
  - embed iv in a web application
  - start multiple in parallel
  - merge results and show histogram
