# APL\iv
APL\iv is an implementation of APL written in go, depending only on the standard library.

As such, it is a package that can be included into any go project.
This may be as small as `cmd/apl` or a huge application where APL simply acts as a hidden debugging interface with all it's expressive powers.

- extendable and embeddable
  - extend APL by implementing extensions in go
  - embed APL by spreading it over already existing go programs
  - *If that sounds too complicated and uses too many words:* `(APL\iv÷Go)←→Lua÷C`
- compatibility and portability: it runs everywhere go runs
  - cross compiling from x to y out of the box, e.g. `GOOS=linux GOARCH=mipsle go install` working on any host including windows is a matter of course

# Goals
- #1: Build the smallest APL machine. Small with respect to m³

# Compatibility
- The compatibility goal is to be mostly conforming to APL2/Dyalog core language substracting nested arrays
- The parser adds some more restrictions
  - function variables have to be lowercase: `f←+/`, nouns are uppercase
  - lambdas (dfns) exist but no user defined operators directly in APL (go extension can define operators)
  - minor issues:
    - `/\ etc` are implemented as operators. These are really nasty.
    - assignment is also implemented as an operator. But `{indexed, modified, selective}` assignment should work.

# Non-compatibility
- Workspace, Namespace, Quadfunctions, I-Beams, user functions others than lambdas.
All this does not exist, is or will be done differently.

# Additions
- Replaceable numeric tower
  - bool and int (index) is always there as `apl.Bool` and `apl.Index`
  - the default tower is compiled in with `numbers.Register(a)` and provides `int64, float64, complex128, time (time.Time + time.Duration)`
  - big numbers can be used with `big.Register(a)` which adds two more towers:
    - big: bigint, bigrat
    - precise: bigfloat, bigcomplex
  - towers are separated, they cannot be mixed within an expression, but can be changed at runtime (if they are compiled in)
    - `big→set 1`
- Overloading primitives and operators
  - adding the packages `primitives.Register(a)` and `operators.Register(a)` adds the default implementation to the interpreter.
  - a user package (e.g. to interface any application dependend go type) may overload each primitive function or operator.
- Lists
  - Instead of nested arrays, there is a list type similar to K's.
  - A list needs a terminating semicolon: `(1;2;3;)` or nested `(1;(2;3;);4;5;)`
- Dicts and tables
  - Also influence by K are dictionaries and tables (and the time type mentioned above).
  - A `Dict` is a special implementation of the more general `Object`
- Go interface
  - go structs are mapped into APL by implementing an `Object`. Working on the from APL side feels like working on a Dict.
  - struct methods are mapped to fields of the dict that have function type.
  - package xgo `xgo.Register(a)` adds helpers to translate from go to apl types.
- Streaming
  - the type `Channel` combines two go channels for sequential reading and writing of any `apl.Value` to a concurrent process, rpc call or go routine 
  
# Speed and size
Of course it should be fast and compact. However the primary goal is implementation speed. This is hard enough.
  
# Disclaimer
The author has never used APL. It's a chicken and egg problem.
Primary source for the implementation is the APL2 Language Reference, Dyalog 17 Language Reference guide and on some occasions the ISO spec. All testing has been done with tryapl.org. Thank you a lot, Dyalog for this.
