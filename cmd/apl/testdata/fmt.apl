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