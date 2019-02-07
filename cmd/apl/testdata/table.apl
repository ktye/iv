⍝ Create a table
D ← 2019.02.07T13.25 + ⍳3
N ← `Peter`Jack`Thomas
B ← 1 0 1
I ← 100 × ⍳3
F ← 100÷⍳3
C ← 1+0J1×⍳3
M ← ((1.2 2.1 3.0;);(;);(7.8912345678;);) ⍝ Multiple values; not allowed in K

⍝ A table is created by transposing a dictionary.
⍝ A dictionary is created with the # primitive
D ← `Time`Name`Mark`Count`Number`Comp`Mult#(D;N;B;I;F;C;M;)
T← ⍉D

⍝ Print the table in stanard form, which is equal to ⍕T
"Default table format:"
T

"PP set to 2:"
⎕PP←2
T
⎕PP←⍳0 ⍝ reset PP

⍝ Print the table in parsable form
"Parsable table format:"
¯1⍕T

⍝ Print table in csv format
"csv format:"
`csv ⍕T

⍝ Custom formatting is created in a dictionary with the corresponding field names
"custom format:"
F←`Time#`alpha
F[`Time]←`2006-01-02T15:04 ⍝ To format dates, use the desired result for the prototype date Jan 2, 2006-01-02 15:04:05
F[`Count]←`0x%x       ⍝ Hexadecimal
F[`Comp]←`%.3f@%.1f ⍝ Amplitude @ angle in degree
F
F⍕T

⍝ Custom format in csv output
"custom format with csv:"
F[`CSV]←1 ⍝ Add the special key CSV with value 1
F⍕T
