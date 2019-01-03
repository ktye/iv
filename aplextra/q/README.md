# rpc interface from APL to kdb

## Examples
```
In q listen on port 1993:
    q)\p 1993
  
In APL connect:
    C←q→dial ":1993"
  
Make a function call, pass an array:
    q→call (C; "sum"; 3 3⍴⍳9;)
(1 2 3;4 5 6;7 8 9;)                 ⍝ the result is a list

Pass a defined function with an integer argument:
    q→call (C; "{n where 2 = sum 0 = n mod\:/: n:1 + til x}"; 50;)
2 3 5 7 11 13 17 19 23 29 31 37 41 43 47

Pass a dictionary
    D←`a`b`c#1 2 3
    q→call (C; "sum"; D;)
6
    q→call (C; "!"; `a`b`c ; (1;2;3;);)
a: 1				⍝ result is a dictionary
b: 2
c: 3

Pass a table
    T←⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)
    q→call (C; "sum"; T;)   ⍝ pass a table
a: 6
b: 15
c: 24

Execute a q-sql function:
    q→call (C; "{select a,c from x where b>4}"; T ;)
a c				⍝ result is a table
2 8
3 8
```
