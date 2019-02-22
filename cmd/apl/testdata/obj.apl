⍝ Create a dictionary and index-assign (change) a field
B←`a`b`c#(2 3⍴⍳6;"alpha";3;)
B[`c]←5
B[`c]=5  ⍝ out: 1

⍝ A is a nested dictionary
A←`x`y`z#(1;B;"ZZZ";)

⍝ Depth-assignment
A[`y;`c]←"@"
A[`y;`c] ⍝ out: @
A[`y;`a;2;2]←8
A[`y;`a;2;] ⍝ out: 4 8 6
A[`y;`a;2;2] ⍝ out: 8
A[`y;`c]←1 2 3
A[`y;`c] ⍝ out: 1 2 3
