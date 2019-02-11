⍝ Formatting examples for arrays, dicts and lists to csv, json, mat
¯1⍕2 3 4⍴⍳24    ⍝ display format
¯2⍕2 3 4⍴⍳24    ⍝ json
¯3⍕3⍴⍳3         ⍝ matlab vector
¯3⍕2 3⍴⍳6       ⍝ matlab vector
⍝ string arrays
S←"abc" "def" "gh\ni"
¯1⍕3 2⍴S
`json ⍕3 2⍴S
`mat ⍕3 2⍴S
⍝ floats and complex
S←1.2 3.4 5.6J¯7.8
¯1⍕3 2⍴S
`csv ⍕3 2⍴S
`json ⍕3 2⍴S
`mat ⍕3 2⍴S
⍝ dictionaries
S←`A`b`zeta#(1; "zwei";1 2 3;)
S
¯1⍕S
¯2⍕S
¯3⍕S
⍝ lists
L←(1 2;(3;(1;2;);"four";);5;)
L
¯1⍕L
¯2⍕L
⍝ Numbers
123 0123 0x1234
¯1⍕123 0123 0x123
¯8⍕123 0123 0x123
¯16⍕123 0123 0x123
`x ⍕123 0123 0x123
`x ⍕ 1.23
