test←{X←`A ⍎⍕⍵ ⋄ X≡⍵} 
⎕PP←0
test 2 3 4⍴⍳24
⎕PP←¯1
test 2 3 4⍴⍳24
⎕PP←¯2
test 2 3 4⍴⍳24
⎕PP←¯3
test 3 4⍴⍳24