# APL\iv interpreter and stream processor
<p align="center" >
  <img width="170" height="170" src="log.svg"><br/>
</p>

# contents
- [apl: an extendable and embeddable APL interpreter written in go](apl)
- Documentation
  - [REF.md](REF.md) reference of all standard primitives and operators
  - [TESTS.md](TESTS.md) output of test coverage which gives an overview of the state of affairs
  - [GOALS.md](GOALS.md) describes the target of the project
  - [DESIGN.md](DESIGN.md) description of the go implementation and how to write extra packages

# programs
- [cmd/apl](cmd/apl): APL interpreter as a command line program
- [cmd/iv](cmd/iv): a program similar to awk with an APL backend but for streaming n-dimensional data

# A random loop through pattern space
```
⎕IO←0
j←{(⍳2*⍺){⎕←(?∘⍴⌷⊢)⍸~1↓⌽0,(⍺⍴2)⊤⍵}⍣≡1⍴⍵} ⍝ not done..
```
