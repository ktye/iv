# APL\iv interpreter and stream processor
<p align="center" >
  <img width="170" height="170" src="log.svg"><br/>
</p>

# contents
- [apl: an extendable and embeddable APL interpreter written in go](apl)
- [aplextra: extension packages to apl requiring 3rd party software](aplextra)
- Documentation
  - [REF.md](REF.md) reference of all standard primitives and operators
  - [TESTS.md](TESTS.md) output of test coverage which gives an overview of the state of affairs
  - [GOALS.md](GOALS.md) describes the target of the project
  - [DESIGN.md](DESIGN.md) description of the go implementation and how to write extra packages

# programs
- [cmd/apl](cmd/apl): APL interpreter as a command line program
- [cmd/aplui](cmd/aplui): APL gui application
- [cmd/iv](cmd/iv): a program similar to awk with an APL backend but for streaming n-dimensional data

# status
What you find in TESTS.md is what works. It has not yet been used for anything else.

Next tasks:
- [ ] assignments on tables (updates)
- [ ] decide when to copy
- [ ] make use of uniform type arrays
- [ ] *this is a stack, not a task list*
- [ ] learn APL
