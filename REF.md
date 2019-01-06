# Reference
- [Primitive Functions](#primitive-functions)
- [Operators](#operators)

## Primitive functions
```
!                                              
   binomial                                    apl/primitives/elementary.go:74
   !R  both (table or object)                  
   binomial                                    apl/primitives/elementary.go:68
   !R  any (table or object)                   
   binomial                                    apl/primitives/elementary.go:62
   !R  arithmetic arrays with axis             
   binomial                                    apl/primitives/elementary.go:56
   !R  arithmetic arrays                       
   binomial                                    apl/primitives/elementary.go:50
   L!R  L scalar R scalar                      
   factorial                                   apl/primitives/elementary.go:44
   !R  (object or table)                       
   factorial                                   apl/primitives/elementary.go:38
   !R  array                                   
   factorial                                   apl/primitives/elementary.go:32
   !R  scalar                                  
                                               
⍉                                              
   cant, transpose, general transpose          apl/primitives/transpose.go:38
   L⍉R  L toindexarray R array                 
   cant, transpose, general transpose          apl/primitives/transpose.go:31
   L⍉R  L array R number                       
   dict from table, transpose, flip            apl/primitives/transpose.go:25
   ⍉R  table                                   
   table from object, transpose, flip          apl/primitives/transpose.go:19
   ⍉R  object                                  
   cant, transpose, reverse axes               apl/primitives/transpose.go:12
   ⍉R  array                                   
                                               
,                                              
   catenate, join along last axis              apl/primitives/comma.go:24
   ,R  any                                     
   ravel, ravel with axis                      apl/primitives/comma.go:11
   ,R  any                                     
                                               
○                                              
   circular, trigonometric                     apl/primitives/elementary.go:74
   ○R  both (table or object)                  
   circular, trigonometric                     apl/primitives/elementary.go:68
   ○R  any (table or object)                   
   circular, trigonometric                     apl/primitives/elementary.go:62
   ○R  arithmetic arrays with axis             
   circular, trigonometric                     apl/primitives/elementary.go:56
   ○R  arithmetic arrays                       
   circular, trigonometric                     apl/primitives/elementary.go:50
   L○R  L scalar R scalar                      
   pi times                                    apl/primitives/elementary.go:44
   ○R  (object or table)                       
   pi times                                    apl/primitives/elementary.go:38
   ○R  array                                   
   pi times                                    apl/primitives/elementary.go:32
   ○R  scalar                                  
                                               
↓                                              
   cut                                         apl/primitives/take.go:51
   L↓R  L toindexarray R list                  
   close channel                               apl/primitives/take.go:45
   ↓R  channel                                 
   drop to channel                             apl/primitives/take.go:39
   L↓R  L channel R any                        
   drop                                        apl/primitives/take.go:32
   L↓R  L toindexarray R any                   
                                               
?                                              
   deal                                        apl/primitives/query.go:19
   L?R  L toscalar index R toscalar index      
   roll                                        apl/primitives/query.go:13
   ?R  any                                     
                                               
⊥                                              
   decode, polynom, base value                 apl/primitives/decode.go:12
   L⊥R  L toarray R toarray                    
                                               
#                                              
   dict                                        apl/primitives/dict.go:18
   #R  any                                     
   keys, methods                               apl/primitives/dict.go:12
   #R  any                                     
                                               
÷                                              
   div, division, divide                       apl/primitives/elementary.go:74
   ÷R  both (table or object)                  
   div, division, divide                       apl/primitives/elementary.go:68
   ÷R  any (table or object)                   
   div, division, divide                       apl/primitives/elementary.go:62
   ÷R  arithmetic arrays with axis             
   div, division, divide                       apl/primitives/elementary.go:56
   ÷R  arithmetic arrays                       
   div, division, divide                       apl/primitives/elementary.go:50
   L÷R  L scalar R scalar                      
   reciprocal                                  apl/primitives/elementary.go:44
   ÷R  (object or table)                       
   reciprocal                                  apl/primitives/elementary.go:38
   ÷R  array                                   
   reciprocal                                  apl/primitives/elementary.go:32
   ÷R  scalar                                  
                                               
⊤                                              
   encode, representation                      apl/primitives/decode.go:18
   ⊤R  any                                     
                                               
=                                              
   equality                                    apl/primitives/compare.go:26
   =R  arithmetic arrays                       
   equality                                    apl/primitives/compare.go:20
   L=R  L scalar R scalar                      
                                               
⍎                                              
   execute, evaluate expression                apl/primitives/format.go:21
   ⍎R  string                                  
                                               
⍷                                              
   find                                        apl/primitives/find.go:9
   L⍷R  L toarray R toarray                    
                                               
⍕                                              
   format, convert to string                   apl/primitives/format.go:11
   ⍕R  any                                     
                                               
⍒                                              
   grade down with collating sequence          apl/primitives/grade.go:31
   L⍒R  L vector R array                       
   grade down, reverse sort index              apl/primitives/grade.go:19
   ⍒R  array                                   
                                               
⍋                                              
   grade up with collating sequence            apl/primitives/grade.go:25
   L⍋R  L vector R array                       
   grade up, sort index                        apl/primitives/grade.go:13
   ⍋R  array                                   
                                               
≥                                              
   greater or equal                            apl/primitives/compare.go:26
   ≥R  arithmetic arrays                       
   greater or equal                            apl/primitives/compare.go:20
   L≥R  L scalar R scalar                      
                                               
>                                              
   greater than                                apl/primitives/compare.go:26
   >R  arithmetic arrays                       
   greater than                                apl/primitives/compare.go:20
   L>R  L scalar R scalar                      
                                               
⍳                                              
   index of, first occurrence                  apl/primitives/iota.go:17
   L⍳R  L tovector R toarray                   
   interval, index generater, progression      apl/primitives/iota.go:11
   ⍳R  toscalar index                          
                                               
⌷                                              
   index table, []                             apl/primitives/index.go:35
   L⌷R  L [index specification] R table        
   index object, []                            apl/primitives/index.go:28
   L⌷R  L [index specification] R object       
   index list, []                              apl/primitives/index.go:21
   L⌷R  L [index specification] R list         
   index, []                                   apl/primitives/index.go:14
   L⌷R  L [index specification] R toarray      
                                               
⍸                                              
   interval index                              apl/primitives/iota.go:35
   L⍸R  L vector R array                       
   where                                       apl/primitives/iota.go:29
   ⍸R  toindexarray                            
                                               
⊂                                              
   join strings                                apl/primitives/enclose.go:17
   L⊂R  L string R array of strings            
   enclose, string catenation                  apl/primitives/enclose.go:11
   ⊂R  array of strings                        
                                               
⊣                                              
   left tack, left argument                    apl/primitives/tack.go:21
   ⊣R  any                                     
   left tack, same                             apl/primitives/tack.go:9
   ⊣R  any                                     
                                               
≤                                              
   less or equal                               apl/primitives/compare.go:26
   ≤R  arithmetic arrays                       
   less or equal                               apl/primitives/compare.go:20
   L≤R  L scalar R scalar                      
                                               
<                                              
   less that                                   apl/primitives/compare.go:26
   <R  arithmetic arrays                       
   less that                                   apl/primitives/compare.go:20
   L<R  L scalar R scalar                      
                                               
⍟                                              
   log, logarithm                              apl/primitives/elementary.go:74
   ⍟R  both (table or object)                  
   log, logarithm                              apl/primitives/elementary.go:68
   ⍟R  any (table or object)                   
   log, logarithm                              apl/primitives/elementary.go:62
   ⍟R  arithmetic arrays with axis             
   log, logarithm                              apl/primitives/elementary.go:56
   ⍟R  arithmetic arrays                       
   log, logarithm                              apl/primitives/elementary.go:50
   L⍟R  L scalar R scalar                      
   natural logarithm                           apl/primitives/elementary.go:44
   ⍟R  (object or table)                       
   natural logarithm                           apl/primitives/elementary.go:38
   ⍟R  array                                   
   natural logarithm                           apl/primitives/elementary.go:32
   ⍟R  scalar                                  
                                               
^                                              
   logical and                                 apl/primitives/boolean.go:31
   ^R  arithmetic arrays                       
   logical and                                 apl/primitives/boolean.go:25
   L^R  L scalar R scalar                      
                                               
∧                                              
   logical and                                 apl/primitives/boolean.go:31
   ∧R  arithmetic arrays                       
   logical and                                 apl/primitives/boolean.go:25
   L∧R  L scalar R scalar                      
                                               
⍲                                              
   logical nand                                apl/primitives/boolean.go:31
   ⍲R  arithmetic arrays                       
   logical nand                                apl/primitives/boolean.go:25
   L⍲R  L scalar R scalar                      
                                               
⍱                                              
   logical nor                                 apl/primitives/boolean.go:31
   ⍱R  arithmetic arrays                       
   logical nor                                 apl/primitives/boolean.go:25
   L⍱R  L scalar R scalar                      
                                               
∨                                              
   logical or                                  apl/primitives/boolean.go:31
   ∨R  arithmetic arrays                       
   logical or                                  apl/primitives/boolean.go:25
   L∨R  L scalar R scalar                      
                                               
|                                              
   magnitude, absolute value                   apl/primitives/elementary.go:74
   |R  both (table or object)                  
   magnitude, absolute value                   apl/primitives/elementary.go:68
   |R  any (table or object)                   
   magnitude, absolute value                   apl/primitives/elementary.go:62
   |R  arithmetic arrays with axis             
   magnitude, absolute value                   apl/primitives/elementary.go:56
   |R  arithmetic arrays                       
   magnitude, absolute value                   apl/primitives/elementary.go:50
   L|R  L scalar R scalar                      
   magnitude, absolute value                   apl/primitives/elementary.go:44
   |R  (object or table)                       
   magnitude, absolute value                   apl/primitives/elementary.go:38
   |R  array                                   
   magnitude, absolute value                   apl/primitives/elementary.go:32
   |R  scalar                                  
                                               
≡                                              
   match                                       apl/primitives/match.go:24
   ≡R  any                                     
   depth, level of nesting                     apl/primitives/match.go:11
   ≡R  any                                     
                                               
⌹                                              
   matrix divide, solve linear system, domino  apl/primitives/domino.go:18
   L⌹R  L toarray R toarray                    
   matrix inverse, domino                      apl/primitives/domino.go:12
   ⌹R  toarray                                 
                                               
⌈                                              
   max, maximum                                apl/primitives/elementary.go:74
   ⌈R  both (table or object)                  
   max, maximum                                apl/primitives/elementary.go:68
   ⌈R  any (table or object)                   
   max, maximum                                apl/primitives/elementary.go:62
   ⌈R  arithmetic arrays with axis             
   max, maximum                                apl/primitives/elementary.go:56
   ⌈R  arithmetic arrays                       
   max, maximum                                apl/primitives/elementary.go:50
   L⌈R  L scalar R scalar                      
   ceil                                        apl/primitives/elementary.go:44
   ⌈R  (object or table)                       
   ceil                                        apl/primitives/elementary.go:38
   ⌈R  array                                   
   ceil                                        apl/primitives/elementary.go:32
   ⌈R  scalar                                  
                                               
∊                                              
   membership                                  apl/primitives/iota.go:23
   ∊R  any                                     
   enlist                                      apl/primitives/comma.go:18
   ∊R  any                                     
                                               
⌊                                              
   min, minumum                                apl/primitives/elementary.go:74
   ⌊R  both (table or object)                  
   min, minumum                                apl/primitives/elementary.go:68
   ⌊R  any (table or object)                   
   min, minumum                                apl/primitives/elementary.go:62
   ⌊R  arithmetic arrays with axis             
   min, minumum                                apl/primitives/elementary.go:56
   ⌊R  arithmetic arrays                       
   min, minumum                                apl/primitives/elementary.go:50
   L⌊R  L scalar R scalar                      
   floor                                       apl/primitives/elementary.go:44
   ⌊R  (object or table)                       
   floor                                       apl/primitives/elementary.go:38
   ⌊R  array                                   
   floor                                       apl/primitives/elementary.go:32
   ⌊R  scalar                                  
                                               
×                                              
   multiply                                    apl/primitives/elementary.go:74
   ×R  both (table or object)                  
   multiply                                    apl/primitives/elementary.go:68
   ×R  any (table or object)                   
   multiply                                    apl/primitives/elementary.go:62
   ×R  arithmetic arrays with axis             
   multiply                                    apl/primitives/elementary.go:56
   ×R  arithmetic arrays                       
   multiply                                    apl/primitives/elementary.go:50
   L×R  L scalar R scalar                      
   signum, sign of, direction                  apl/primitives/elementary.go:44
   ×R  (object or table)                       
   signum, sign of, direction                  apl/primitives/elementary.go:38
   ×R  array                                   
   signum, sign of, direction                  apl/primitives/elementary.go:32
   ×R  scalar                                  
                                               
≠                                              
   not equal                                   apl/primitives/compare.go:26
   ≠R  arithmetic arrays                       
   not equal                                   apl/primitives/compare.go:20
   L≠R  L scalar R scalar                      
                                               
≢                                              
   not match                                   apl/primitives/match.go:30
   ≢R  any                                     
   tally, number of major cells                apl/primitives/match.go:17
   ≢R  any                                     
                                               
+                                              
   plus, addition                              apl/primitives/elementary.go:74
   +R  both (table or object)                  
   plus, addition                              apl/primitives/elementary.go:68
   +R  any (table or object)                   
   plus, addition                              apl/primitives/elementary.go:62
   +R  arithmetic arrays with axis             
   plus, addition                              apl/primitives/elementary.go:56
   +R  arithmetic arrays                       
   plus, addition                              apl/primitives/elementary.go:50
   L+R  L scalar R scalar                      
   identity, complex conjugate                 apl/primitives/elementary.go:44
   +R  (object or table)                       
   identity, complex conjugate                 apl/primitives/elementary.go:38
   +R  array                                   
   identity, complex conjugate                 apl/primitives/elementary.go:32
   +R  scalar                                  
                                               
*                                              
   power                                       apl/primitives/elementary.go:74
   *R  both (table or object)                  
   power                                       apl/primitives/elementary.go:68
   *R  any (table or object)                   
   power                                       apl/primitives/elementary.go:62
   *R  arithmetic arrays with axis             
   power                                       apl/primitives/elementary.go:56
   *R  arithmetic arrays                       
   power                                       apl/primitives/elementary.go:50
   L*R  L scalar R scalar                      
   exponential                                 apl/primitives/elementary.go:44
   *R  (object or table)                       
   exponential                                 apl/primitives/elementary.go:38
   *R  array                                   
   exponential                                 apl/primitives/elementary.go:32
   *R  scalar                                  
                                               
⍴                                              
   reshape                                     apl/primitives/rho.go:12
   L⍴R  L tovector toindexarray R toarray      
   shape                                       apl/primitives/rho.go:11
   ⍴R  any                                     
                                               
⊢                                              
   right tack, right argument                  apl/primitives/tack.go:27
   ⊢R  any                                     
   right tack, same                            apl/primitives/tack.go:15
   ⊢R  any                                     
                                               
⌽                                              
   rotate                                      apl/primitives/reverse.go:27
   L⌽R  L toindexarray R any                   
   reverse                                     apl/primitives/reverse.go:11
   ⌽R  any                                     
                                               
⊖                                              
   rotate first                                apl/primitives/reverse.go:34
   L⊖R  L toindexarray R any                   
   reverse first                               apl/primitives/reverse.go:18
   ⊖R  any                                     
                                               
-                                              
   substract, substraction                     apl/primitives/elementary.go:74
   -R  both (table or object)                  
   substract, substraction                     apl/primitives/elementary.go:68
   -R  any (table or object)                   
   substract, substraction                     apl/primitives/elementary.go:62
   -R  arithmetic arrays with axis             
   substract, substraction                     apl/primitives/elementary.go:56
   -R  arithmetic arrays                       
   substract, substraction                     apl/primitives/elementary.go:50
   L-R  L scalar R scalar                      
   reverse sign                                apl/primitives/elementary.go:44
   -R  (object or table)                       
   reverse sign                                apl/primitives/elementary.go:38
   -R  array                                   
   reverse sign                                apl/primitives/elementary.go:32
   -R  scalar                                  
                                               
⍪                                              
   table                                       apl/primitives/comma.go:36
   ⍪R  any                                     
   catenate first                              apl/primitives/comma.go:30
   ⍪R  any                                     
                                               
↑                                              
   take one from channel                       apl/primitives/take.go:26
   ↑R  channel                                 
   take from channel                           apl/primitives/take.go:20
   L↑R  L toindexarray R channel               
   take                                        apl/primitives/take.go:13
   L↑R  L toindexarray R any                   
                                               
⌶                                              
   type                                        apl/primitives/type.go:11
   ⌶R  any                                     
                                               
∪                                              
   union                                       apl/primitives/unique.go:15
   L∪R  L tovector R tovector                  
   unique                                      apl/primitives/unique.go:9
   ∪R  tovector                                
                                               
~                                              
   without, excluding                          apl/primitives/boolean.go:50
   L~R  L tovector R tovector                  
   logical not                                 apl/primitives/boolean.go:44
   ~R  array                                   
   logical not                                 apl/primitives/boolean.go:38
   ~R  scalar                                  
                                               
```
## Operators
←                                  
   assign, variable specification  apl/operators/assign.go:11
   LO←RO  LO any                   
                                   
@                                  
   at                              apl/operators/at.go:11
   @RO  any                        
                                   
⍂                                  
   axis specification              apl/operators/axis.go:11
   ⍂RO  any                        
                                   
⍨                                  
   commute, duplicate              apl/operators/commute.go:9
   LO⍨RO  LO function              
                                   
∘                                  
   compose                         apl/operators/jot.go:11
   ∘RO  L any R any                
                                   
¨                                  
   each, map                       apl/operators/each.go:11
   LO¨RO  LO function              
                                   
\                                  
   expand                          apl/operators/reduce.go:50
   LO\RO  LO toindexarray          
   scan                            apl/operators/reduce.go:23
   LO\RO  LO function              
                                   
⍀                                  
   expand first axis               apl/operators/reduce.go:59
   LO⍀RO  LO toindexarray          
   scan first axis                 apl/operators/reduce.go:29
   LO⍀RO  LO function              
                                   
⍣                                  
   power                           apl/operators/power.go:11
   ⍣RO  L function R any           
                                   
⍤                                  
   rank                            apl/operators/rank.go:11
   ⍤RO  L function R toindexarray  
                                   
/                                  
   replicate, compress             apl/operators/reduce.go:35
   LO/RO  LO toindexarray          
   reduce, n-wise reduction        apl/operators/reduce.go:11
   LO/RO  LO function              
                                   
⌿                                  
   replicate, compress first axis  apl/operators/reduce.go:44
   LO⌿RO  LO toindexarray          
   reduce first, n-wise reduction  apl/operators/reduce.go:17
   LO⌿RO  LO function              
                                   
.                                  
   scalar product                  apl/operators/dot.go:18
   .RO  L + R ×                    
   inner product                   apl/operators/dot.go:12
   .RO  L function R function      
                                   
⌺                                  
   stencil                         apl/operators/stencil.go:11
   ⌺RO  L function R toindexarray  
                                   
```
PASS
ok  	github.com/ktye/iv/apl/primitives	0.686s

generated by `go generate (apl/primitives/gen.go)` 2019-01-06 01:38:47
