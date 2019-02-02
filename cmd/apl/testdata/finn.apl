⍝ A : Any [Numeric, Character or Boolean]
⍝ D : Numeric
⍝ I : Integer
⍝ C : Character
⍝ B : Boolean
⍝
⍝ A number following the type indicates the rank, e.g. 
⍝ A0: Any scalar (rank 0)
⍝ A1: Any vector (rank 1)
⍝ A2: Any matrix (rank 2)


⍝⍝⍝ Grade Up ⍋
f1←{((⍴⍵)⍴⍋⍋⍵⍳⍵,⍺)⍳(⍴⍺)⍴⍋⍋⍵⍳⍺,⍵}                               ⍝ Progressive index of (without replacement)                     ⍝ ⍵←A1; ⍺←A1
f2←{⌊.5×(⍋⍋⍵)+⌽⍋⍋⌽⍵}                                           ⍝ Ascending cardinal numbers (ranking, shareable)                ⍝ ⍵←D1
f3←{⍺[A⍳⌈\A←⍋A[⍋(+\⍵)[A←⍋⍺]]]}                                 ⍝ Cumulative maxima (`⌈\`) of subvectors of ⍺ indicated by ⍵     ⍝ ⍵←B1; ⍺←D1
f4←{⍺[A⍳⌈\A←⍋A[⍋(+\⍵)[A←⍒⍺]]]}                                 ⍝ Cumulative minima (`⌊\`) of subvectors of ⍺ indicated by ⍵     ⍝ ⍵←B1; ⍺←D1
f5←{((⍋⍵⍳⍵,⍺)⍳⍳⍴⍵)⍳(⍋⍵⍳⍺,⍵)⍳⍳⍴⍺}                               ⍝ Progressive index of (without replacement)                     ⍝ ⍵←A1; ⍺←A1
f6←{⍺[⍋⍺]^.=⍵[⍋⍵]}                                             ⍝ Test if ⍵ and ⍺ are permutations of each other                 ⍝ ⍵←D1; ⍺←D1
f7←{⍵^.=⍋⍋⍵}                                                   ⍝ Test if ⍵ is a permutation vector                              ⍝ ⍵←I1
f8←{A[⍋(+\(⍳⍴⍺)∊+\⎕IO,⍵)[A←⍋⍺]]}                               ⍝ Grade up (`⍋`) for sorting subvectors of ⍺ having lengths ⍵    ⍝ ⍺←D1; ⍵←I1; (⍴⍺) ←→ +/⍵
f9←{(((1,A)/B)⌊1+⍴⍺)[(⍴⍺)↓(+\1,A←(1↓A)≠¯1↓A←A[B])[⍋B←⍋A←⍺,⍵]]} ⍝ Index of the elements of ⍵ in ⍺                                ⍝ ⍵←D1; ⍺←D1
f10←{⍺[A[⍵/⍋(+\⍵)[A←⍋⍺]]]}                                     ⍝ Minima (⌊/) of elements of subvectors of ⍺ indicated by ⍵      ⍝ ⍵←B1; ⍺←D1
f11←{A[⍋(+\⍵)[A←⍋⍺]]}                                          ⍝ Grade up (`⍋`) for sorting subvectors of ⍺ indicated by ⍵      ⍝ ⍵←B1; ⍺←D1
f12←{|-⌿(2,⍴⍵)⍴⍋⍋⍵,⍵}                                          ⍝ Occurences of the elements of ⍵                                ⍝ ⍵←D1
f13←{(⍴⍵)⍴(,⍵)[A[⍋(,⍉(⌽⍴⍵)⍴⍳1↑⍴⍵)[A←⍋,⍵]]]}                    ⍝ Sorting rows of matrix ⍵ into ascending order                  ⍝ ⍵←D2
f14←{(⍋⍋(G+1),⍳⍴⍴⍵)⍉(⍺,⍴⍵)⍴⍵}                                  ⍝ Adding a new dimension after dimension G ⍺-fold                ⍝ G←I0; ⍺←I0; ⍵←A
f15←{A←(⍋,⍵)-⎕IO ⋄ (⍴⍵)⍴(,⍵)[⎕IO+A[⍋⌊A÷¯1↑⍴⍵]]}                ⍝ Sorting rows of matrix ⍵ into ascending order                  ⍝ ⍵←D2
f16←{((⍋⍋⍵)∊⍳⍺)/⍵}                                             ⍝ ⍺ smallest elements of ⍵ in order of occurrence                ⍝ ⍵←D1, ⍺←I0
f17←{(⍵,⍺,Z)[⍋⍋G]}                                             ⍝ Merging ⍵, ⍺, Z ... under control of G (mesh)                  ⍝ ⍵←A1; ⍺←A1; Z←A1; ... ; G←I1
f18←{(⍵,⍺)[⍋⍋G]}                                               ⍝ Merging ⍵ and ⍺ under control of G (mesh)                      ⍝ ⍵←A1; ⍺←A1; G←B1
f19←{⍋⍋⍵}                                                      ⍝ Ascending cardinal numbers (ranking, all different)            ⍝ ⍵←D1
f20←{A[⍋(+\(⍳⍴⍺)∊+\⎕IO,⍵)[A←⍒⍺]]}                              ⍝ Grade down (`⍒`) for sorting subvectors of ⍺ having lengths ⍵  ⍝ ⍺←D1; ⍵←I1; (⍴⍺) ←→ +/⍵
f21←{⍺[A[⍵/⍋(+\⍵)[A←⍒⍺]]]}                                     ⍝ Maxima (⌈/) of elements of subvectors of ⍺ indicated by ⍵      ⍝ ⍵←B1; ⍺←D1
f22←{A[⍋(+\⍵)[A←⍒⍺]]}                                          ⍝ Grade down (`⍒`) for sorting subvectors of ⍺ indicated by ⍵    ⍝ ⍵←B1; ⍺←D1
f23←{((⍋⍒⍵)∊⍳⍺)/⍵}                                             ⍝ ⍺ largest elements of ⍵ in order of occurrence                 ⍝ ⍵←D1; ⍺←I0
f24←{(⍺,⍵)[⍋⍒G]}                                               ⍝ Merging ⍵ and ⍺ under control of G (mesh)                      ⍝ ⍵←A1; ⍺←A1; G←B1
f25←{⍋⍒⍵}                                                      ⍝ Descending cardinal numbers (ranking, all different)           ⍝ ⍵←D1
f26←{⍵[⍋(1+⍴⍺)⊥⍺⍳⍉⍵;]}                                         ⍝ Sorting rows of ⍵ according to key ⍺ (alphabetizing)           ⍝ ⍵←A2; ⍺←A1
f27←{(,⍵)[⍋+⌿(⍴⍵)⊤(⍳⍴,⍵)-⎕IO]}                                 ⍝ Diagonal ravel                                                 ⍝ ⍵←A
f28←{⍋⍺⍳⍵}                                                     ⍝ Grade up according to key ⍺                                    ⍝ ⍺←A1; ⍵←A1
f29←{⍵[⍋⍵]^.=⍳⍴⍵}                                              ⍝ Test if ⍵ is a permutation vector                              ⍝ ⍵←I1
f30←{⍵[⍋+⌿A<.-⍉A←⍵,0;]}                                        ⍝ Sorting a matrix into lexicographic order                      ⍝ ⍵←D2
f31←{⍵[⍋⍵+.≠' ';]}                                             ⍝ Sorting words in list ⍵ according to word length               ⍝ ⍵←C2
f32←{A[(B/C)-⍴⍺]←B/+\~B←(⍴⍺)<C←⍋⍺,⍵+A←0×⍵ ⋄ A}                 ⍝ Classification of ⍵ to classes starting with ⍺                 ⍝ ⍵←D1;⍺←D1;⍺<.≥1⌽⍺
f33←{⍺[⍋⍵++\⍵]}                                                ⍝ Rotate first elements (`1⌽`) of subvectors of ⍺ indicated by ⍵ ⍝ ⍵←B1; ⍺←A1
f34←{(⍵,'''')[(⎕IO+⍴⍵)⌊⍋(⍳⍴⍵),(''''=⍵)/⍳⍴⍵]}                   ⍝ Doubling quotes (for execution)                                ⍝ ⍵←C1
f35←{(⍵,'*')[(⎕IO+⍴⍵)⌊⍋(⍳⍴⍵),(⍺×⍴G)⍴G]}                        ⍝ Inserting ⍺ `*`'s into vector ⍵ after indices G                ⍝ ⍵←C1; ⍺←I0; G←I1
f36←{⍵[(⍋⍵)[⌈.5×⍴⍵]]}                                          ⍝ Median                                                         ⍝ ⍵←D1
f37←{¯1↑⍋⍵}                                                    ⍝ Index of last maximum element of ⍵                             ⍝ ⍵←D1
f38←{1↑⍋⍵}                                                     ⍝ Index of (first) minimum element of ⍵                          ⍝ ⍵←D1
f39←{(⍴⍵)≥⍋(⍳⍴⍵),⍺}                                            ⍝ Expansion vector with zero after indices ⍺                     ⍝ ⍵←D1; ⍺←I1
f40←{A←G×⍴,⍺ ⋄ ((A⍴H),⍵)[⍋(A⍴⍺),⍳⍴⍵]}                          ⍝ Catenating G elements H before indices ⍺ in vector ⍵           ⍝ ⍵←A1; ⍺←I1; G←I0; H←A0
f41←{A←G×⍴,⍺ ⋄ (⍵,A⍴H)[⍋(⍳⍴⍵),A⍴⍺]}                            ⍝ Catenating G elements H after indices ⍺ in vector ⍵            ⍝ ⍵←A1; ⍺←I1; G←I0; H←A0
f42←{A[⍋G]←A←⍺,⍵ ⋄ A}                                          ⍝ Merging ⍵ and ⍺ under control of G (mesh)                      ⍝ ⍵←A1; ⍺←A1; G←B1
f43←{⍵[⍋⍵[;⍺];]}                                               ⍝ Sorting a matrix according to ⍺:th column                      ⍝ ⍵←D2
f44←{⍵[⍋⍺[⍵]]}                                                 ⍝ Sorting indices ⍵ according to data ⍺                          ⍝ ⍵←I1; ⍺←D1
f45←{⍋⍵×(¯1 1)[⍺]}                                             ⍝ Choosing sorting direction during execution                    ⍝ ⍵←D1; ⍺←I0
f46←{⍺[⍋⍵]}                                                    ⍝ Sorting ⍺ according to ⍵                                       ⍝ ⍵←A1; ⍺←A1
f47←{⍵[⍋⍵]}                                                    ⍝ Sorting ⍵ into ascending order                                 ⍝ ⍵←D1
f48←{⍋⍵}                                                       ⍝ Inverting a permutation                                        ⍝ ⍵←I1

⍝⍝⍝ Grade Down ⍒
f49←{⍵[⍒⍺!⍳⍴⍵]}            ⍝ Reverse vector ⍵ on condition ⍺                          ⍝ ⍵←A1; ⍺←B0
f50←{⍵[⍒+⌿A<.-⍉A←⍵,0;]}    ⍝ Sorting a matrix into reverse lexicographic order        ⍝ ⍵←D2
f52←{⍵[⌽⍒+\(⍳⍴⍵)∊+\⎕IO,⍺]} ⍝ Reversal (`⌽`) of subvectors of ⍵ having lengths ⍺       ⍝ ⍵←D1; ⍺←I1
f53←{⍺[⌽⍒+\⍵]}             ⍝ Reversal (`⌽`) of subvectors of ⍺ indicated by ⍵         ⍝ ⍵←B1; ⍺←A1
f55←{(+/⍵)↑⍒⍵}             ⍝ Indices of ones in logical vector ⍵                      ⍝ ⍵←B1
f56←{1↑⍒⍵}                 ⍝ Index of first maximum element of ⍵                      ⍝ ⍵←D1
f57←{⍵[⍒' '≠⍵]}            ⍝ Moving all blanks to end of text                         ⍝ ⍵←C1
f58←{⍵[⍒⍵]}                ⍝ Sorting ⍵ into descending order                          ⍝ ⍵←D1
f59←{⍵[⍒⍺]}                ⍝ Moving elements satisfying condition ⍺ to the start of ⍵ ⍝ ⍵←A1; ⍺←B1

⍝⍝⍝ Matrix Inversion / Matrix Division ⌹
f60←{G⊥⍺⌹⍵∘.*⌽-⎕IO-⍳⍴⍵}               ⍝ Interpolated value of series (⍵,⍺) at G                 ⍝ ⍵←D1; ⍺←D1; G←D0
f61←{*A+.×(⍟⍺)⌹A←⍵∘.*0 1}             ⍝ Predicted values of exponential (curve) fit             ⍝ ⍵←D1; ⍺←D1
f62←{A←(⍟⍺)⌹⍵∘.*0 1 ⋄ A[1]←*A[1] ⋄ A} ⍝ Coefficients of exponential (curve) fit of points (⍵,⍺) ⍝ ⍵←D1; ⍺←D1
f63←{A+.×⍺⌹A←⍵∘.*0 1}                 ⍝ Predicted values of best linear fit (least squares)     ⍝ ⍵←D1; ⍺←D1
f64←{⌽⍺⌹⍵∘.*0,⍳G}                     ⍝ G-degree polynomial (curve) fit of points (⍵,⍺)         ⍝ ⍵←D1; ⍺←D1
f65←{⍺⌹⍵∘.*0 1}                       ⍝ Best linear fit of points (⍵,⍺) (least squares)         ⍝ ⍵←D1; ⍺←D1

⍝⍝⍝ Decode ⊥
f66←{⍕10⊥((1+⌈2⍟⌈/,⍵)⍴2)⊤⍵}                                       ⍝ Binary format of decimal number ⍵                          ⍝ ⍵←I0
f67←{' *○⍟'[⎕IO+2⊥⍵∘.≥⍳⌈/,⍵]}                                     ⍝ Barchart of two integer series (across the page)           ⍝ ⍵←I2; 1⍴⍴⍵ ←→ 2
⍝f68←{→⍺[1+2⊥⍵]}                                                   ⍝ Case structure with an encoded branch destination          ⍝ ⍝ ⍺←I1; ⍵←B1
f69←{A←⍕1000⊥3↑3↓⎕TS ⋄ A[3 6]←':' ⋄ A}                            ⍝ Representation of current time (24 hour clock)             ⍝ 
f70←{A←⍕1000⊥3↑⎕TS ⋄ A[5 8]←'-' ⋄ A}                              ⍝ Representation of current date (descending format)         ⍝ 
f71←{(1⌽,' ::',3 2⍴6 0⍕100⊥12 0 0|3↑3↓⎕TS),'AP'[1+12≤⎕TS[4]],'M'} ⍝ Representation of current time (12 hour clock)             ⍝ 
f73←{((A⍳A)=⍳⍴A←2⊥⍵^.=⍉⍵)⌿⍵}                                      ⍝ Removing duplicate rows                                    ⍝ ⍵←A2
f74←{16⊥-⎕IO-'0123456789ABCDEF'⍳⍉⍵}                               ⍝ Conversion from hexadecimal to decimal                     ⍝ ⍵←C
f75←{10⊥¯1+'0123456789'⍳⍵}                                        ⍝ Conversion of alphanumeric string into numeric             ⍝ ⍵←C1
f76←{(⍵∘.+,0)⊥⍺}                                                  ⍝ Value of polynomial with coefficients ⍺ at points ⍵        ⍝ ⍵←D1; ⍺←D1
f77←{A←(×/B←0 0+⌈/,⍵)⍴0 ⋄ A[⎕IO+B[1]⊥-⎕IO-⍵]←1 ⋄ B⍴A}             ⍝ Changing connectivity list ⍵ to a connectivity matrix      ⍝ ⍵←C2
f78←{(÷1+⍺÷100)⊥⌽⍵}                                               ⍝ Present value of cash flows ⍵ at interest rate ⍺ %         ⍝ ⍵←D1; ⍺←D0
f79←{(1-(' '=⍵)⊥1)⌽⍵}                                             ⍝ Justifying right                                           ⍝ ⍵←C
f80←{(12⍴7⍴31 30)[⍵]-0⌈¯1+2⊥(⍵=2),[.1](0≠400|⍺)-(0≠100|⍺)-0≠4|⍺}  ⍝ Number of days in month ⍵ of years ⍺ (for all leap years)  ⍝ ⍵←I0; ⍺←I
f81←{(12⍴7⍴31 30)[⍵]-0⌈¯1+2⊥(⍵=2),[.1]0≠4|⍺}                      ⍝ Number of days in month ⍵ of years ⍺ (for most leap years) ⍝ ⍵←I0; ⍺←I
f82←{100⊥100|3↑⎕TS}                                               ⍝ Encoding current date                                      ⍝ 
f83←{(1-(' '=⍵)⊥1)↓⍵}                                             ⍝ Removing trailing blanks                                   ⍝ ⍵←C1
f84←{(' '=⍵)⊥1}                                                   ⍝ Index of first non-blank, counted from the rear            ⍝ ⍵←C1
f85←{(,⍵)[⎕IO+(⍴⍵)⊥⍺-⎕IO]}                                        ⍝ Indexing scattered elements                                ⍝ ⍵←A; ⍺←I2
f86←{⎕IO+(⍴⍵)⊥⍺-⎕IO}                                              ⍝ Conversion of indices ⍺ of array ⍵ to indices of raveled ⍵ ⍝ ⍵←A; ⍺←I2
f87←{0⊥⍴⍵}                                                        ⍝ Number of columns in array ⍵ as a scalar                   ⍝ ⍵←A
f88←{(1+⍺÷100)⊥⍵}                                                 ⍝ Future value of cash flows ⍵ at interest rate ⍺ %          ⍝ ⍵←D1; ⍺←D0
f89←{1⊥⍵}                                                         ⍝ Sum of the elements of vector ⍵                            ⍝ ⍵←D1
f90←{0⊥⍵}                                                         ⍝ Last element of numeric vector ⍵ as a scalar               ⍝ ⍵←D1
f91←{0⊥⍵}                                                         ⍝ Last row of matrix ⍵ as a vector                           ⍝ ⍵←A
f92←{2⊥⍵}                                                         ⍝ Integer representation of logical vectors                  ⍝ ⍵←B
f93←{⍵⊥⍺}                                                         ⍝ Value of polynomial with coefficients ⍺ at point ⍵         ⍝ ⍵←D0; ⍺←D

⍝⍝⍝ Encode ⊤
f94←{⍉'0123456789ABCDEF'[⎕IO+((⌈⌈/16⍟,⍵)⍴16)⊤⍵]}         ⍝ Conversion from decimal to hexadecimal (`⍵=1..255`)            ⍝ ⍵←I
f94b←{⍉'0123456789ABCDEF'[⎕IO+((1+⌊16⍟⌈/X+X=0)⍴16)⊤X]}   ⍝  this alternative opens the range to 0..⌊/⍳0                   ⍝ ⍵←I                                                  ⍝ 
f95←{((⌈2⍟1+⍵)⍴2)⊤0,⍳⍵}                                  ⍝ All binary representations up to ⍵ (truth table)               ⍝ ⍵←I0
f96←{((1+⌊⍺⍟⍵)⍴⍺)⊤⍵}                                     ⍝ Representation of ⍵ in base ⍺                                  ⍝ ⍵←D0; ⍺←D0
f97←{((1+⌊10⍟⍵)⍴10)⊤⍵}                                   ⍝ Digits of ⍵ separately                                         ⍝ ⍵←I0
f98←{1 0⍕10 10⊤1-⎕IO-⍳⍵}                                 ⍝ Helps locating column positions 1..⍵                           ⍝ ⍵←I0
f99←{,' ',⍉'0123456789ABCDEF'[⎕IO+16 16⊤-⎕IO-⎕AV⍳⍵]}     ⍝ Conversion of characters to hexadecimal representation (`⎕AV`) ⍝ ⍵←C1
f100←{⌽((0,⍳⍴⍵)∘.=+⌿~A)+.×(-⍵)×.*A←((⍴⍵)⍴2)⊤¯1+⍳2*⍴⍵}    ⍝ Polynomial with roots ⍵                                        ⍝ ⍵←D1
f101←{⎕IO+(⍴⍵)⊤-⎕IO-(,(⍵=(⍴⍵)⍴⌈⌿⍵)^⍵=⍉(⌽⍴⍵)⍴⌊/⍵)/⍳×/⍴⍵}  ⍝ Index pairs of saddle points                                   ⍝ ⍵←D2
f102←{(,⍵)/1+A⊤¯1+⍳×/A←⍴⍵}                               ⍝ Changing connectivity matrix ⍵ to a connectivity list          ⍝ ⍵←C2
f103←{⎕IO+(⍴⍵)⊤(⍳×/⍴⍵)-⎕IO}                              ⍝ Matrix of all indices of ⍵                                     ⍝ ⍵←A
f104←{⍉(3⍴100)⊤⍵}                                        ⍝ Separating a date ⍺⍺MMDD to ⍺⍺, MM, DD                         ⍝ ⍵←D
f105←{⎕IO+(⍴⍵)⊤(-⎕IO)+(,⍵∊⍺)/⍳⍴,⍵}                       ⍝ Indices of elements ⍺ in array ⍵                               ⍝ ⍵←A; ⍺←A
f106←{⎕IO+(⍵,⍺)⊤(⍳⍵×⍺)-⎕IO}                              ⍝ All pairs of elements of `⍳⍵` and `⍳⍺`                         ⍝ ⍵←I0; ⍺←I0
f107←{((⍴⍵)⍴2)⊤¯1+⍳2*⍴⍵}                                 ⍝ Matrix for choosing all subsets of ⍵ (truth table)             ⍝ ⍵←A1
f108←{(⍵⍴2)⊤¯1+⍳2*⍵}                                     ⍝ All binary representations with ⍵ bits (truth table)           ⍝ ⍵←I0
f109←{1+⍺⊤⍵}                                             ⍝ Incrementing cyclic counter ⍵ with upper limit ⍺               ⍝ ⍵←D; ⍺←D0
f110←{10 100 1000⊤⍵}                                     ⍝ Decoding numeric code ABBCCC into a matrix                     ⍝ ⍵←I
f111←{0 1⊤⍵}                                             ⍝ Integer and fractional parts of positive numbers               ⍝ ⍵←D

⍝⍝⍝ Logarithm ⍟
f112←{⌊10⍟(⍎('.'≠A)/A←⍕⍵)÷⍵}        ⍝ Number of decimals of elements of ⍵                           ⍝ ⍵←D1
f113←{⌊(1+⍴⍵)⍟2*(A=¯1+A←2*⍳128)⍳1}  ⍝ Number of sortable columns at a time using `⊥` and alphabet ⍵ ⍝ ⍵←C1
f114←{,⍉(A⍴2)⍴(2*A←⌈2⍟⍵)↑⍳⍵}        ⍝ Playing order in a cup for ⍵ ranked players                   ⍝ ⍵←I0
f115←{⌊|10⍟|1-3×÷3}                 ⍝ Arithmetic precision of the system (in decimals)              ⍝ 
f116←{1+(⍵<0)+⌊10⍟|⍵+0=⍵}           ⍝ Number of digitpositions in integers in ⍵                     ⍝ ⍵←I
f117←{1+⌊10⍟(⍵=0)+⍵×(1 ¯10)[1+⍵<0]} ⍝ Number of digit positions in integers in ⍵                    ⍝ ⍵←I
f118←{1+⌊10⍟⍵+0=⍵}                  ⍝ Number of digits in positive integers in ⍵                    ⍝ ⍵←I

⍝⍝⍝ Branch →
⍝f119←{→⍺[G⍳⍵]}                 ⍝ Case structure according to key vector G                   ⍝ ⍵←A0; ⍺←I1; G←A1
⍝f120←{→⎕LC⌈⍳∨/,(⍵←⍵∨⍵∨.^⍵)≠+⍵} ⍝ Forming a transitive closure                               ⍝ ⍵←B2
⍝f121←{→⍵⌽⍺}                    ⍝ Case structure with integer switch                         ⍝ ⍵←I0; ⍺←I1
⍝f122←{→⍺⌈⍳G≥⍵←⍵+1}             ⍝ For-loop ending construct                                  ⍝ ⍵←I0; ⍺←I0; G←I0
⍝f123←{→⍺⌈⍳⍵}                   ⍝ Conditional branch to line ⍺                               ⍝ ⍵←B0; ⍺←I0; ⍺>0
⍝f124←{→0⌊⍳⍵}                   ⍝ Conditional branch out of program                          ⍝ ⍵←B0
⍝f125←{→⍺[2+×⍵]}                ⍝ Conditional branch depending on sign of ⍵                  ⍝ ⍵←I0; ⍺←I1
⍝f126←{→⍺××⍵}                   ⍝ Continuing from line ⍺ (if ⍵>0) or exit                    ⍝ ⍵←D0; ⍺←I0
⍝f127←{→(⍵≥G)/⍺}                ⍝ Case structure using levels with limits G                  ⍝ ⍵←D0; G←D1; ⍺←I1
⍝f128←{→⍵/⍺}                    ⍝ Case structure with logical switch (preferring from start) ⍝ ⍵←B1; ⍺←I1
⍝f129←{→0×⍳⍵}                   ⍝ Conditional branch out of program                          ⍝ ⍵←B0

⍝⍝⍝ Execute ⍎
f132←{⍎⍎'1','↑↓'[⎕IO+^/(⍴⍵)=⌽⍴⍵],'''0~0∊⍵=⍉⍵'''}     ⍝ Test for symmetricity of matrix ⍵                         ⍝ ⍵←A2
f133←{⍎'VAR',(⍕⍵),'←⍺'}                              ⍝ Using a variable named according to ⍵                     ⍝ ⍵←A0; ⍺←A
f134←{⍎⍕⍵}                                           ⍝ Rounding to `⎕PP` precision                               ⍝ ⍵←D1
f135←{⍎⍕⍵}                                           ⍝ Convert character or numeric data into numeric            ⍝ ⍵←A1
f136←{⍎⍕⍵}                                           ⍝ Reshaping only one-element numeric vector ⍵ into a scalar ⍝ ⍵←D1
f137←{' *'[⎕IO+(⌽(¯1+⌊/A)+⍳1+(⌈/A)-⌊/A)∘.=A←⌊.5+⍎F]} ⍝ Graph of F(⍵) at points ⍵ ('⍵'∊F)                         ⍝ F←A1; ⍵←D1
f138←{(⍵∨.≠' ')\1↓⍎'0 ',,⍵,' '}                      ⍝ Conversion of each row to a number (default zero)         ⍝ ⍵←C2
f139←{⍎(¯7*A^.=⌽A←⍴⍵)↑'0~0∊⍵=⍉⍵'}                    ⍝ Test for symmetricity of matrix ⍵                         ⍝ ⍵←A2
f140←{⍎((⍵^.=' ')/'⍺'),⍵}                            ⍝ Execution of expression ⍵ with default value ⍺            ⍝ ⍵←D1
⍝f141←{⍵←⍎,((2↑'⍵'),' ',[.5]A)[⎕IO+~' '^.=A←⍞;]}     ⍝ Changing ⍵ if a new input value is given                  ⍝ ⍵←A
f142←{A+.×⍎F,0⍴⍵←⍺[1]+(A←--/⍺÷G)×0,⍳G}               ⍝ Definite integral of F(⍵) in range ⍺ with G steps ('⍵'∊F) ⍝ F←A1; G←D0; ⍺←D1; ⍴⍺ ←→ 2
f143←{1↓⍎'0 ',(^/⍵∊' 0123456789')/⍵}                 ⍝ Test if numeric and conversion to numeric form            ⍝ ⍵←C1
f144←{(¯1↑⍵)=((~⍺∊'GIOQ')/⍺)[1+31|⍎9↑⍵]}             ⍝ Tests the social security number (Finnish)                ⍝ ⍺←'01...9ABC...Z'; 10=⍴⍵
f145←{⍎⍵/'E⍵PRESSION'}                               ⍝ Conditional execution                                     ⍝ ⍵←B0
f146←{⍎⍵/'→'}                                        ⍝ Conditional branch out of programs                        ⍝ ⍵←B0
f147←{⍎(¯3*2≠⎕NC '⍵')↑'⍵100'}                        ⍝ Using default value 100 if ⍵ does not exist               ⍝ ⍵←A
f148←{⍎⍵↓'                                           ⍝ ...'}                                                     ⍝ Conditional execution ⍝ ⍵←B0
⍝f149←{1⍴(⍎⍞,',⍳0'),⍵}                               ⍝ Giving a numeric default value for input                  ⍝ ⍵←D0
f150←{A←⍎,',','(','0','⍴',⍺,'←',⍵,')'}               ⍝ Assign values of expressions in ⍵ to variables named in ⍺ ⍝ ⍵←C2; ⍺←C2
f151←{⍎,',','(',',',⍵,')'}                           ⍝ Evaluation of several expressions; results form a vector  ⍝ ⍵←A
f152←{⍎,'+',⍵}                                       ⍝ Sum of numbers in character matrix ⍵                      ⍝ ⍵←A2
f153←{⍎'⍵[',((¯1+⍴⍴⍵)⍴';'),'⍺]'}                     ⍝ Indexing when rank is not known beforehand                ⍝ ⍵←A; ⍺←I

⍝⍝⍝ Format ⍕
f154←{(3⌽7 0⍕⍵∘.+,0),⍕⍺}                       ⍝ Numeric headers (elements of ⍵) for rows of table ⍺          ⍝ ⍵←D1; ⍺←A2
f155←{⍕⍵∘.+,0}                                 ⍝ Formatting a numerical vector to run down the page           ⍝ ⍵←D1
f156←{A←⍕⌽3↑⎕TS ⋄ A[(' '=A)/⍳⍴A]←'.' ⋄ A}      ⍝ Representation of current date (ascending format)            ⍝ 
f157←{A←⍕100|1⌽3↑⎕TS ⋄ A[(' '=A)/⍳⍴A]←'/' ⋄ A} ⍝ Representation of current date (American)                    ⍝ 
f158←{(⍴A)⍴B\(B←,('0'≠A)∨' '≠¯1⌽A)/,A←' ',⍕⍵}  ⍝ Formatting with zero values replaced with blanks             ⍝ ⍵←A
f159←{⍴⍕⍵}                                     ⍝ Number of digit positions in scalar ⍵ (depends on `⎕PP`)     ⍝ ⍵←D0
f160←{0 1↓(2↑⍺+1)⍕⍵∘.+,10*⍺}                   ⍝ Leading zeroes for ⍵ in fields of width ⍺                    ⍝ ⍵←I1; ⍺←I0; ⍵≥0
f161←{((1,G)×⍴⍵)⍴2 1 3⍉(⌽G,⍴⍵)⍴(,G,[1.1]⍺)⍕⍉⍵} ⍝ Row-by-row formatting (width G) of ⍵ with ⍺ decimals per row ⍝ ⍵←D2; ⍺←I1; G←I0
f163←{(,G,[1.1]H)⍕⍵}                           ⍝ Formatting ⍵ with H decimals in fields of width G            ⍝ ⍵←D; G←I1; H←I1

⍝⍝⍝ Roll / Deal ?
f164←{⍵[1]+?⍺⍴--/⍵}      ⍝ ⍺-shaped array of random numbers within ( ⍵[1],⍵[2] ]     ⍝ ⍵←I1; ⍺←I1
f165←{(~⍵∊' .,:;?''')/⍵} ⍝ Removing punctuation characters                           ⍝ ⍵←A1
f166←{?⍺⍴⍵}              ⍝ Choosing ⍺ objects out of `⍳⍵` with replacement (roll)    ⍝ ⍺←I; ⍵←I
f167←{⍺?⍵}               ⍝ Choosing ⍺ objects out of `⍳⍵` without replacement (deal) ⍝ ⍵←I0; ⍺←I0

⍝⍝⍝ Geometrical Functions ○
f168←{((⍵≠0)×¯3○⍺÷⍵+⍵=0)+○((⍵=0)×.5××⍺)+(⍵<0)×1-2×⍺<0} ⍝ Arctan ⍺÷⍵                                                 ⍝ ⍵←D; ⍺←D
f169←{⍵×○÷180}                                         ⍝ Conversion from degrees to radians                         ⍝ ⍵←D
f170←{⍵×180÷○1}                                        ⍝ Conversion from radians to degrees                         ⍝ ⍵←D
f171←{2 2⍴1 ¯1 1 1×2 1 1 2○⍵}                          ⍝ Rotation matrix for angle ⍵ (in radians) counter-clockwise ⍝ ⍵←D0

⍝⍝⍝ Factorial / Binomial !
f172←{(!⍺)×⍺!⍵}                        ⍝ Number of permutations of ⍵ objects taken ⍺ at a time  ⍝ ⍵←D; ⍺←D
f173←{+/⍺×(⍵*A)÷!A←¯1+⍳⍴⍺}             ⍝ Value of Taylor series with coefficients ⍺ at point ⍵  ⍝ ⍵←D0; ⍺←D1
f174←{(*-⍺)×(⍺*⍵)÷!⍵}                  ⍝ Poisson distribution of states ⍵ with average number ⍺ ⍝ ⍵←I; ⍺←D0
f175←{!⍵-1}                            ⍝ Gamma function                                         ⍝ ⍵←D0
f176←{(A!⍵)×(⍺*A)×(1-⍺)*⍵-A←-⎕IO-⍳⍵+1} ⍝ Binomial distribution of ⍵ trials with probability ⍺   ⍝ ⍵←I0; ⍺←D0
f177←{÷⍺×(⍵-1)!⍺+⍵-1}                  ⍝ Beta function                                          ⍝ ⍵←D0; ⍺←D0
f178←{⍵!⍺}                             ⍝ Selecting elements satisfying condition ⍵, others to 1 ⍝ ⍵←B; ⍺←D
f179←{⍺!⍵}                             ⍝ Number of combinations of ⍵ objects taken ⍺ at a time  ⍝ ⍵←D; ⍺←D

⍝⍝⍝ Index Of ⍳
f180←{((A⍳1)-⎕IO)↓(⎕IO-(⌽A←~⍵∊⍺)⍳1)↓⍵}                      ⍝ Removing elements ⍺ from beginning and end of vector ⍵       ⍝ ⍵←A1; ⍺←A
f181←{(G⍳⍵)<G⍳⍺}                                            ⍝ Alphabetical comparison with alphabets G                     ⍝ ⍵←A; ⍺←A
f183←{⍵+.×⍺∘.=((⍳⍴⍺)=⍺⍳⍺)/⍺}                                ⍝ Sum over elements of ⍵ determined by elements of ⍺           ⍝ ⍵←D1; ⍺←D1
f184←{(^⌿(¯1+⍳⍴⍵)⌽⍵∘.=⍺)⍳1}                                 ⍝ First occurrence of string ⍵ in string ⍺                     ⍝ ⍵←A1; ⍺←A1
f185←{((A⍳A)=⍳⍴A←⎕IO++⌿^⍀⍵∨.≠⍉⍵)⌿⍵}                         ⍝ Removing duplicate rows                                      ⍝ ⍵←A2
f186←{(⍺^.=⍵)⍳1}                                            ⍝ First occurrence of string ⍵ in matrix ⍺                     ⍝ ⍵←A1; ⍺←A2; ¯1↑⍴⍺←→⍴⍵
f187←{(+\⍵)⍳⍳+/⍵}                                           ⍝ Indices of ones in logical vector ⍵                          ⍝ ⍵←B1
f188←{(F B/⍵)[+\B←(⍵⍳⍵)=⍳⍴⍵]}                               ⍝ Executing costly monadic function F on repetitive arguments  ⍝ ⍵←A1
f189←{⍵⍳⌈/⍵}                                                ⍝ Index of (first) maximum element of ⍵                        ⍝ ⍵←D1
f190←{⌊/⍵⍳⍺}                                                ⍝ Index of first occurrence of elements of ⍺                   ⍝ ⍵←C1; ⍺←C1
f191←{⍵⍳⌊/⍵}                                                ⍝ Index of (first) minimum element of ⍵                        ⍝ ⍵←D1
f192←{^/(⍵⍳⍵)=⍳⍴⍵}                                          ⍝ Test if each element of ⍵ occurs only once                   ⍝ ⍵←A1
f193←{^/⎕IO=⍵⍳⍵}                                            ⍝ Test if all elements of vector ⍵ are equal                   ⍝ ⍵←A1
f194←{+/A×¯1*A<1⌽A←0,(1000 500 100 50 10 5 1)['MDCL⍵VI'⍳⍵]} ⍝ Interpretation of roman numbers                              ⍝ ⍵←A
f195←{(⎕IO-(~⌽⍵∊⍺)⍳1)↓⍵}                                    ⍝ Removing elements ⍺ from end of vector ⍵                     ⍝ ⍵←A1; ⍺←A
f196←{(1-(⌽' '≠⍵)⍳1)↓⍵}                                     ⍝ Removing trailing blanks                                     ⍝ ⍵←C1
f198←{((¯1 1)[2×⎕IO]+⍴⍵)-(⌽⍵)⍳⍺}                            ⍝ Index of last occurrence of ⍺ in ⍵ (`⎕IO-1` if not found)    ⍝ ⍵←A1; ⍺←A
f199←{(1+⍴⍵)-(⌽⍵)⍳⍺}                                        ⍝ Index of last occurrence of ⍺ in ⍵ (0 if not found)          ⍝ ⍵←A1; ⍺←A
f200←{(⌽⍵)⍳⍺}                                               ⍝ Index of last occurrence of ⍺ in ⍵, counted from the rear    ⍝ ⍵←A1; ⍺←A
f201←{⎕IO+(⍴⍵)|⍺+(⍺⌽⍵)⍳G}                                   ⍝ Index of first occurrence of G in ⍵ (circularly) after ⍺     ⍝ ⍵←A1; ⍺←I0; G←A
f202←{(¯1↑⍴⍺)|(,⍺)⍳⍵}                                       ⍝ Alphabetizing ⍵; equal alphabets in same column of ⍺         ⍝ ⍺←C2; ⍵←C
f203←{(1+⍴⍺)|⍺⍳⍵}                                           ⍝ Changing index of an unfound element to zero                 ⍝ ⍺←A1; ⍵←A
f204←{A[B/⍳⍴B]←⍺[(B←B≤⍴⍺)/B←⍵⍳A←,G] ⋄ (⍴G)⍴A}               ⍝ Replacing elements of G in set ⍵ with corresponding ⍺        ⍝ ⍵←A1, ⍺←A1, G←A
f205←{((⍵⍳⍵)=⍳⍴⍵)/⍵}                                        ⍝ Removing duplicate elements (nub)                            ⍝ ⍵←A1
f206←{(¯1+⍵⍳' ')↑⍵}                                         ⍝ First word in ⍵                                              ⍝ ⍵←C1
f207←{(((~⍵∊⍺)⍳1)-⎕IO)↓⍵}                                   ⍝ Removing elements ⍺ from beginning of vector ⍵               ⍝ ⍵←A1; ⍺←A
f208←{(¯1+(⍵='0')⍳0)↓⍵}                                     ⍝ Removing leading zeroes                                      ⍝ ⍵←A1
f209←{⍺+(⍺↓⍵)⍳1}                                            ⍝ Index of first one after index ⍺ in ⍵                        ⍝ G←I0; ⍵←B1
f210←{(⍵∊⍺)×⍺⍳⍵}                                            ⍝ Changing index of an unfound element to zero (not effective) ⍝ ⍵←A; ⍺←A1
f211←{(⍵⍳⍵)=⍳⍴⍵}                                            ⍝ Indicator of first occurrence of each unique element of ⍵    ⍝ ⍵←A1
f212←{⍵⍳⍳⍴⍵}                                                ⍝ Inverting a permutation                                      ⍝ ⍵←I1
f213←{(⍺≠⍵)⍳1}                                              ⍝ Index of first differing element in vectors ⍵ and ⍺          ⍝ ⍵←A1; ⍺←A1
f214←{(⎕IO+⍴⍺)=⍺⍳⍵}                                         ⍝ Which elements of ⍵ are not in set ⍺ (difference of sets)    ⍝ ⍵←A; ⍺←A1
f215←{G[⍺⍳⍵;]}                                              ⍝ Changing numeric code ⍵ into corresponding name in ⍺         ⍝ ⍵←D; ⍺←D1; G←C2
f216←{⍵⍳⍺}                                                  ⍝ Index of key ⍺ in key vector ⍵                               ⍝ ⍵←A1; ⍺←A
f217←{⎕AV⍳⍵}                                                ⍝ Conversion from characters to numeric codes                  ⍝ ⍵←A
f218←{⍵⍳1}                                                  ⍝ Index of first satisfied condition in ⍵                      ⍝ ⍵←B1

⍝⍝⍝ Outer Product ∘.! ∘.⌈ ∘.|
f219←{⍉A∘.!A←0,⍳⍵}                         ⍝ Pascal's triangle of order ⍵ (binomial coefficients) ⍝ ⍵←I0
f220←{(⍳⍵)∘.⌈⍳⍵}                           ⍝ Maximum table                                        ⍝ ⍵←I0
f221←{0+.≠(⌈(10*⍺)×10*⎕IO-⍳⍺+1)∘.|⌈⍵×10*⍺} ⍝ Number of decimals (up to ⍺) of elements of ⍵        ⍝ ⍵←D; ⍺←I0
f222←{⌈/(^/0=A∘.|⍵)/A←⍳⌊/⍵}                ⍝ Greatest common divisor of elements of ⍵             ⍝ ⍵←I1
f223←{0=(⍳⌈/⍵)∘.|⍵}                        ⍝ Divisibility table                                   ⍝ ⍵←I1
f224←{(2=+⌿0=(⍳⍵)∘.|⍳⍵)/⍳⍵}                ⍝ All primes up to ⍵                                   ⍝ ⍵←I0

⍝⍝⍝ Outer Product ∘.* ∘.× ∘.- ∘.+
f225←{⍺∘.×(1+G÷100)∘.*⍵}                          ⍝ Compound interest for principals ⍺ at rates G % in times ⍵  ⍝ ⍵←D; ⍺←D; G←D
f226←{+⌿(⎕IO-⍳⍴⍵)⌽⍵∘.×⍺,0×1↓⍵}                    ⍝ Product of two polynomials with coefficients ⍵ and ⍺        ⍝ ⍵←D1; ⍺←D1
f228←{1 2 1 2⍉⍵∘.×⍺}                              ⍝ Shur product                                                ⍝ ⍵←D2; ⍺←D2
f229←{1 3 2 4⍉⍵∘.×⍺}                              ⍝ Direct matrix product                                       ⍝ ⍵←D2; ⍺←D2
f230←{(⍳⍵)∘.×⍳⍵}                                  ⍝ Multiplication table                                        ⍝ ⍵←I0
f231←{⍵[;,(⍺⍴1)∘.×⍳(⍴⍵)[2];]}                     ⍝ Replicating a dimension of rank three array ⍵ ⍺-fold        ⍝ ⍺←I0; ⍵←A3
f232←{⍵∘.×1 ¯1}                                   ⍝ Array and its negative ('plus minus')                       ⍝ ⍵←D
f233←{1 2 1⍉⍵∘.-⌊/⍵}                              ⍝ Move set of points ⍵ into first quadrant                    ⍝ ⍵←D2
f234←{+/×⍵∘.-⍺}                                   ⍝ Test relations of elements of ⍵ to range ⍺; result in ¯2..2 ⍝ ⍵←D; ⍺←D; 2=¯1↑⍴⍺
f235←{(⍺[A∘.+¯1+⍳⍴⍵]^.=⍵)/A←(A=1↑⍵)/⍳⍴A←(1-⍴⍵)↓⍺} ⍝ Occurrences of string ⍵ in string ⍺                         ⍝ ⍵←A1; ⍺←A1
f236←{1 2 1 2⍉⍵∘.+⍺}                              ⍝ Sum of common parts of matrices (matrix sum)                ⍝ ⍵←D2; ⍺←D2
f237←{1 1 2⍉⍵∘.+⍺}                                ⍝ Adding ⍵ to each row of ⍺                                   ⍝ ⍵←D1; ⍺←D2
f238←{1 2 1⍉⍺∘.+⍵}                                ⍝ Adding ⍵ to each row of ⍺                                   ⍝ ⍵←D1; ⍺←D2
f240←{2 1 2⍉⍵∘.+⍺}                                ⍝ Adding ⍵ to each column of ⍺                                ⍝ ⍵←D1; ⍺←D2
f241←{1 2 2⍉⍺∘.+⍵}                                ⍝ Adding ⍵ to each column of ⍺                                ⍝ ⍵←D1; ⍺←D2
f242←{÷¯1+(⍳⍵)∘.+⍳⍵}                              ⍝ Hilbert matrix of order ⍵                                   ⍝ ⍵←I0
f243←{(0,⍳(⍴⍵)-⍺)∘.+⍺}                            ⍝ Moving index of width ⍺ for vector ⍵                        ⍝ ⍵←A1; ⍺←I0
f244←{⍵∘.+⍳⍺}                                     ⍝ Indices of subvectors of length ⍺ starting at ⍵+1           ⍝ ⍵←I1; ⍺←I0
f245←{⍵∘.+,0}                                     ⍝ Reshaping numeric vector ⍵ into a one-column matrix         ⍝ ⍵←D1
f246←{((⍴A)⍴⍺÷100)÷A←⍉1-(1+⍺÷100)∘.*-⍵}           ⍝ Annuity coefficient: ⍵ periods at interest rate ⍺ %         ⍝ ⍵←I; ⍺←D

⍝⍝⍝ Outer Product ∘.
f247←{⍵∘.<⌽⍳⌈/⍵}                                  ⍝ Matrix with ⍵[i] trailing zeroes on row i                 ⍝ ⍵←I1
f248←{⍵∘.<⍳⌈/⍵}                                   ⍝ Matrix with ⍵[i] leading zeroes on row i                  ⍝ ⍵←I1
f249←{+/((¯1↓⍺)∘.≤⍵)^(1↓⍺)∘.>⍵}                   ⍝ Distribution of ⍵ into intervals between ⍺                ⍝ ⍵←D; ⍺←D1
f250←{' ⎕'[⎕IO+(⌽⍳⌈/A)∘.≤A←+/(⍳1+(⌈/⍵)-⌊/⍵)∘.=⍵]} ⍝ Histogram (distribution barchart; down the page)          ⍝ ⍵←I1
f251←{' ⎕'[⎕IO+(⌽⍳⌈/⍵)∘.≤⍵]}                      ⍝ Barchart of integer values (down the page)                ⍝ ⍵←I1
f252←{^/,(0≠⍵)≤A∘.≤A←⍳1↑⍴⍵}                       ⍝ Test if ⍵ is an upper triangular matrix                   ⍝ ⍵←D2
f253←{+/A^⍉A←⍵∘.≤⍺}                               ⍝ Number of ?s intersecting ?s (⍵=starts, ⍺=stops)          ⍝ ⍵←D1; ⍺←D1
f254←{⍺[+⌿⍺∘.≤⍵]}                                 ⍝ Contour levels ⍺ at points with altitudes ⍵               ⍝ ⍵←D0; ⍺←D1
f255←{(⍳⍵)∘.≤⍳⍵}                                  ⍝ ⍵×⍵ upper triangular matrix                               ⍝ ⍵←I0
f256←{+/(A×⍵÷⌈/A←⍺-⌊/⍺)∘.≥¯1+⍳⍵}                  ⍝ Classification of elements ⍺ into ⍵ classes of equal size ⍝ ⍵←I0; ⍺←D1
f257←{⍵∘.≥⌽⍳⌈/⍵}                                  ⍝ Matrix with ⍵[i] trailing ones on row i                   ⍝ ⍵←I1
f258←{⍵∘.≥⍳⌈/⍵,0}                                 ⍝ Comparison table                                          ⍝ ⍵←I1
f259←{' ⎕'[⎕IO+⍵∘.≥(⌈/⍵)×(⍳⍺)÷⍺]}                 ⍝ Barchart of ⍵ with height ⍺ (across the page)             ⍝ ⍵←D1; ⍺←D0
f260←{' ⎕'[⎕IO+⍵∘.≥⍳⌈/⍵]}                         ⍝ Barchart of integer values (across the page)              ⍝ ⍵←I1
f261←{⍵∘.≥⍳⌈/⍵}                                   ⍝ Matrix with ⍵[i] leading ones on row i                    ⍝ ⍵←I1
f263←{^/,(0≠⍵)≤A∘.≥A←⍳1↑⍴⍵}                       ⍝ Test if ⍵ is a lower triangular matrix                    ⍝ ⍵←D2
f264←{≠/⍵∘.≥⍺}                                    ⍝ Test if ⍵ is within range [ ⍺[1],⍺[2] )                   ⍝ ⍵←D; ⍺←D1
f265←{⎕IO++/⍺∘.≥(' '=⍵)/⍳⍴⍵}                      ⍝ Ordinal numbers of words in ⍵ that indices ⍺ point to     ⍝ ⍵←C1; ⍺←I
f266←{+/⍵∘.≥0 50 100 1000}                        ⍝ Which class do elements of ⍵ belong to                    ⍝ ⍵←D
f267←{(⍳⍵)∘.≥⍳⍵}                                  ⍝ ⍵×⍵ lower triangular matrix                               ⍝ ⍵←I0
f268←{(⍴⍵)⍴(,(+/A)∘.>-⎕IO-⍳¯1↑⍴⍵)\(,A←⍵≠' ')/,⍵}  ⍝ Moving all blanks to end of each row                      ⍝ ⍵←C
f269←{(,⍺∘.>⌽(⍳G)-⎕IO)\⍵}                         ⍝ Justifying right fields of ⍵ (lengths ⍺) to length G      ⍝ ⍵←A1; ⍺←I1; G←I0
f270←{(,⍺∘.>(⍳G)-⎕IO)\⍵}                          ⍝ Justifying left fields of ⍵ (lengths ⍺) to length G       ⍝ ⍵←A1; ⍺←I1; G←I0

⍝⍝⍝ Outer Product ∘.≠ ∘.=
f271←{1++/^\1 2 1 3⍉⍺∘.≠⍵}              ⍝ Indices of elements of ⍺ in corr. rows of ⍵ (`⍵[i;]⍳⍺[i;]`)    ⍝ ⍵←A2; ⍺←A2
f273←{⍉⍵∘.=(1 1⍉<\⍵∘.=⍵)/⍵}             ⍝ Indicating equal elements of ⍵ as a logical matrix             ⍝ ⍵←A1
f275←{(1 ¯1∘.=⍉⍵)+.×⍳1↑⍴⎕←⍵}            ⍝ Changing connection matrix ⍵ (`¯1 → 1`) to a node matrix       ⍝ ⍵←I2
f276←{(G∘.=⍵)+.×⍺}                      ⍝ Sums according to codes G                                      ⍝ ⍵←A; ⍺←D; G←A
f277←{(1 1⍉<\⍵∘.=⍵)/⍵}                  ⍝ Removing duplicate elements (nub)                              ⍝ ⍵←A1
f278←{-/(⍳⌈/,⍵)∘.=⍉⍵}                   ⍝ Changing node matrix ⍵ (starts,ends) to a connection matrix    ⍝ ⍵←I2
f279←{∨/^/0 1∘.=⍵}                      ⍝ Test if all elements of vector ⍵ are equal                     ⍝ ⍵←B1
f280←{∨/1 2 1 3⍉⍵∘.=⍺}                  ⍝ Test if elements of ⍵ belong to corr. row of ⍺ (`⍵[i;]∊⍺[i;]`) ⍝ ⍵←A2; ⍺←A2; 1↑⍴⍵←→1↑⍴⍺
f281←{^/1=+⌿⍵∘.=⍳⍴⍵}                    ⍝ Test if ⍵ is a permutation vector                              ⍝ ⍵←I1
f282←{(^⌿(¯1+⍳⍴⍵)⌽(⍵∘.=⍺),0)/⍳1+⍴⍺}     ⍝ Occurrences of string ⍵ in string ⍺                            ⍝ ⍵←C1; ⍺←C1
f283←{+/(⍳⍺)∘.=⌈(⍵-G)÷H}                ⍝ Division to ⍺ classes with width H, minimum G                  ⍝ ⍵←D; ⍺←I0; G←D0; H←D0
f285←{(((¯1⌽~A)^A←(¯1↓⍵=1⌽⍵),0)/⍺)∘.=⍺} ⍝ Repeat matrix                                                  ⍝ ⍵←A1; ⍺←A1
f286←{(⍳⍵)∘.=⍳⍵}                        ⍝ ⍵×⍵ identity matrix                                            ⍝ ⍵←I0

⍝⍝⍝ Inner Product ⌈.× ⌊.× ⌊.+ ×.○ ×.* +.*
f287←{A+(⍵-A←⌊/⍵)⌈.×⍺}   ⍝ Maxima of elements of subsets of ⍵ specified by ⍺      ⍝ ⍵←A1; ⍺←B
f288←{(' '≠⍵)⌈.×⍳¯1↑⍴⍵}  ⍝ Indices of last non-blanks in rows                     ⍝ ⍵←C
f289←{⍺⌈.×⍵}             ⍝ Maximum of ⍵ with weights ⍺                            ⍝ ⍵←D1; ⍺←D1
f290←{⍺⌊.×⍵}             ⍝ Minimum of ⍵ with weights ⍺                            ⍝ ⍵←D1; ⍺←D1
f292←{⍵←⍵⌊.+⍵}           ⍝ Extending a distance table to next leg                 ⍝ ⍵←D2
f293←{1 2×.○⍵,⍺}         ⍝ A way to combine trigonometric functions (sin ⍵ cos ⍺) ⍝ ⍵←D0; ⍺←D0
f294←{(2 2⍴1 6 2 5)×.○⍵} ⍝ Sine of a complex number                               ⍝ ⍵←D; 2=1↑⍴⍵
f295←{⍵×.*⍺}             ⍝ Products over subsets of ⍵ specified by ⍺              ⍝ ⍵←A1; ⍺←B
f296←{⍵+.*2}             ⍝ Sum of squares of ⍵                                    ⍝ ⍵←D1
f297←{⎕RL←⎕TS+.*2}       ⍝ Randomizing random numbers (in `⎕L⍵` in a workspace)   ⍝ 

⍝⍝⍝ Inner Product ∨.^ .>
f298←{⍵←⍵∨.^⍵}           ⍝ Extending a transitive binary relation       ⍝ ⍵←B2
f299←{⍵<.<⍺}             ⍝ Test if ⍵ is within range [ ⍺[1;],⍺[2;] )    ⍝ ⍵←D0; ⍺←D2; 1↑⍴⍺ ←→ 2
f300←{⍵<.≤⍺}             ⍝ Test if ⍵ is within range ( ⍺[1;],⍺[2;] ]    ⍝ ⍵←D0; ⍺←D2; 1↑⍴⍺ ←→ 2
f301←{⍵<.≤⍺}             ⍝ Test if ⍵ is within range ( ⍺[1;],⍺[2;] ]    ⍝ ⍵←D; ⍺←D2; 1↑⍴⍺ ←→ 2
f302←{⍵<.≥1⌽⍵}           ⍝ Test if the elements of ⍵ are ascending      ⍝ ⍵←D1
f303←{~⍵≤.≥(⌈⍵),G,H}     ⍝ Test if ⍵ is an integer within range [ G,H ) ⍝ ⍵←I0; G←I0; H←I0
f304←{(⍵,[.1+⍴⍴⍵]⍵)>.>⍺} ⍝ Test if ⍵ is within range ( ⍺[1;],⍺[2;] ]    ⍝ ⍵←D; ⍺←D2; 1↑⍴⍺ ←→ 2

⍝⍝⍝ Inner Product ∨.≠ ^.= +.≠ +.=
f306←{(⌽∨\⌽' '∨.≠⍵)/⍵}                    ⍝ Removing trailing blank columns                           ⍝ ⍵←C2
f307←{(∨\⍵∨.≠' ')⌿⍵}                      ⍝ Removing leading blank rows                               ⍝ ⍵←C2
f308←{(∨\' '∨.≠⍵)/⍵}                      ⍝ Removing leading blank columns                            ⍝ ⍵←C2
f309←{⎕IO++⌿^⍀⍺∨.≠⍉⍵}                     ⍝ Index of first occurrences of rows of ⍵ as rows of ⍺      ⍝ ⍵←A, ⍺←A2
f310←{⎕IO++⌿^⍀⍵∨.≠⍉⍺}                     ⍝ `⍵⍳⍺` for rows of matrices                                ⍝ ⍵←A2; ⍺←A2
f311←{(A∨1↓1⌽1,A←⍵∨.≠' ')⌿⍵}              ⍝ Removing duplicate blank rows                             ⍝ ⍵←C2
f312←{(A∨1,¯1↓A←' '∨.≠⍵)/⍵}               ⍝ Removing duplicate blank columns                          ⍝ ⍵←C2
f313←{(' '∨.≠⍵)/⍵}                        ⍝ Removing blank columns                                    ⍝ ⍵←C2
f314←{(⍵∨.≠' ')⌿⍵}                        ⍝ Removing blank rows                                       ⍝ ⍵←C2
f315←{⍵∨.≠⍺}                              ⍝ Test if rows of ⍵ contain elements differing from ⍺       ⍝ ⍵←A; ⍺←A0
f316←{(-2↑+/^\⌽⍵^.=' ')↓⍵}                ⍝ Removing trailing blank rows                              ⍝ ⍵←C2
f317←{(∨⌿<\⍵^.=⍉⍵)⌿⍵}                     ⍝ Removing duplicate rows                                   ⍝ ⍵←A2
f318←{(1 1⍉<\⍵^.=⍉⍵)⌿⍵}                   ⍝ Removing duplicate rows                                   ⍝ ⍵←A2
f319←{∨/⍺^.=⍉(⍳⍴⍵)⌽(2⍴⍴⍵)⍴⍵}              ⍝ Test if circular lists are equal (excluding phase)        ⍝ ⍵←A1; ⍺←A1
f320←{⍵^.=∨/⍵}                            ⍝ Test if all elements of vector ⍵ are equal                ⍝ ⍵←B1
f321←{⍵^.=^/⍵}                            ⍝ Test if all elements of vector ⍵ are equal                ⍝ ⍵←B1
f322←{((((1↑⍴⍵),⍴⍺)↑⍵)^.=⍺)⌿⍵}            ⍝ Rows of matrix ⍵ starting with string ⍺                   ⍝ ⍵←A2; ⍺←A1
f323←{((-A)↓⍵^.=(A,1+⍴⍺)⍴⍺)/⍳(⍴⍺)+1-A←⍴⍵} ⍝ Occurrences of string ⍵ in string ⍺                       ⍝ ⍵←A1; ⍺←A1
f324←{1∊⍵^.=⍺}                            ⍝ Test if vector ⍺ is a row of array ⍵                      ⍝ ⍵←A; ⍺←A1
f325←{⍵^.=⍺}                              ⍝ Comparing vector ⍺ with rows of array ⍵                   ⍝ ⍵←A; ⍺←A1
f326←{⍵+.≠' '}                            ⍝ Word lengths of words in list ⍵                           ⍝ ⍵←C
f327←{⍵+.=,⍺}                             ⍝ Number of occurrences of scalar ⍵ in array ⍺              ⍝ ⍵←A0; ⍺←A
f328←{⍵+.=⍺}                              ⍝ Counting pairwise matches (equal elements) in two vectors ⍝ ⍵←A1; ⍺←A1

⍝⍝⍝ Inner Product -.÷ +.÷ +.×
f329←{⍺-.÷⍵}                                               ⍝ Sum of alternating reciprocal series ⍺÷⍵           ⍝ ⍵←D1; ⍺←D1
f330←{(⍵⌈1↓A)⌊1↑A←(2 2⍴¯1 1 1 ¯.1)+.×10*(-1↓⍺),-/⍺+⍺>99 0} ⍝ Limits ⍵ to fit in `⍕` field ⍺[1 2]                ⍝ ⍵←D; ⍺←I1
f331←{(⍵*¯1+⍳⍴⍺)+.×⌽⍺}                                     ⍝ Value of polynomial with coefficients ⍺ at point ⍵ ⍝ ⍵←D0; ⍺←D
f332←{(⍺+.×⍵)÷⍴⍵}                                          ⍝ Arithmetic average (mean value) of ⍵ weighted by ⍺ ⍝ ⍵←D1; ⍺←D1
f333←{⍺+.×⍵}                                               ⍝ Scalar (dot) product of vectors                    ⍝ ⍵←D1; ⍺←D1
f334←{⍵+.×⍵}                                               ⍝ Sum of squares of ⍵                                ⍝ ⍵←D1
f335←{⍵+.×⍺}                                               ⍝ Summation over subsets of ⍵ specified by ⍺         ⍝ ⍵←A1; ⍺←B
f336←{⍵+.×⍺}                                               ⍝ Matrix product                                     ⍝ ⍵←D; ⍺←D; ¯1↑⍴⍵ ←→ 1↑⍴⍺
f337←{⍺+.÷⍵}                                               ⍝ Sum of reciprocal series ⍺÷⍵                       ⍝ ⍵←D1; ⍺←D1

⍝⍝⍝ Scan ⌈\ ⌊\ ×\ -\
f338←{⍺^A=⌈\⍵×A←+\⍺>¯1↓0,⍺} ⍝ Groups of ones in ⍺ pointed to by ⍵ (or trailing parts)        ⍝ ⍵←B; ⍺←B
f339←{^/[⍺]⍵=⌈\[⍺]⍵}        ⍝ Test if ⍵ is in ascending order along direction ⍺              ⍝ ⍵←D; ⍺←I0
f340←{⍵[1⌈⌈\⍺×⍳⍴⍺]}         ⍝ Duplicating element of ⍵ belonging to `⍺,1↑⍵` until next found ⍝ ⍵←A1; ⍺←B1
f341←{^/[⍺]⍵=⌊\[⍺]⍵}        ⍝ Test if ⍵ is in descending order along direction ⍺             ⍝ ⍵←D; ⍺←I0
f342←{+/⍺××\1,⍵÷⍳¯1+⍴⍺}     ⍝ Value of Taylor series with coefficients ⍺ at point ⍵          ⍝ ⍵←D0; ⍺←D1
f343←{-\⍳⍵}                 ⍝ Alternating series (1 ¯1 2 ¯2 3 ¯3 ...)                        ⍝ ⍵←I0

⍝⍝⍝ Scan ⍲\ <\ ≤\ ≠\
f346←{(<\,(⍵=(⍴⍵)⍴⌈⌿⍵)^⍵=⍉(⌽⍴⍵)⍴⌊/⍵)/,⍵} ⍝ Value of saddle point                                          ⍝ ⍵←D2
f348←{<\⍵}                               ⍝ First one (turn off all ones after first one)                  ⍝ ⍵←B
f350←{≤\⍵}                               ⍝ Not first zero (turn on all zeroes after first zero)           ⍝ ⍵←B
f351←{≠\⍺≠⍵\A≠¯1↓0,A←⍵/≠\¯1↓0,⍺}         ⍝ Running parity (≠\) over subvectors of ⍺ indicated by ⍵        ⍝ ⍵←B1; ⍺←B1
f352←{≠\(⍳+/⍵)∊+\⎕IO,⍵}                  ⍝ Vector `(⍵[1]⍴1),(⍵[2]⍴0),(⍵[3]⍴1),...`                        ⍝ ⍵←I1
f353←{≠\(⍺∨⍵)\A≠¯1↓0,A←(⍺∨⍵)/⍺}          ⍝ Not leading zeroes(`∨\`) in each subvector of ⍺ indicated by ⍵ ⍝ ⍵←B1; ⍺←B1
f354←{~≠\(⍺≤⍵)\A≠¯1↓0,A←~(⍺≤⍵)/⍺}        ⍝ Leading ones (`^\) in each subvector of ⍺ indicated by ⍵       ⍝ ⍵←B1; ⍺←B1
f355←{A∨¯1↓0,A←≠\⍵=''''}                 ⍝ Locations of texts between and including quotes                ⍝ ⍵←C1
f356←{A^¯1↓0,A←≠\⍵=''''}                 ⍝ Locations of texts between quotes                              ⍝ ⍵←C1
f357←{⍵∨≠\⍵}                             ⍝ Joining pairs of ones                                          ⍝ ⍵←B
f358←{(~⍵)^≠\⍵}                          ⍝ Places between pairs of ones                                   ⍝ ⍵←B
f359←{≠\⍵}                               ⍝ Running parity                                                 ⍝ ⍵←B

⍝⍝⍝ Scan ∨\ ^\
f360←{((⌽∨\⌽A)^∨\A←' '≠⍵)/⍵}                         ⍝ Removing leading and trailing blanks                    ⍝ ⍵←C1
f361←{⍵^^\⍵=∨\⍵}                                     ⍝ First group of ones                                     ⍝ ⍵←B
f362←{(⌽∨\⌽∨⌿' '≠⍵)/⍵}                               ⍝ Removing trailing blank columns                         ⍝ ⍵←C2
f363←{(⌽∨\⌽' '≠⍵)/⍵}                                 ⍝ Removing trailing blanks                                ⍝ ⍵←C1
f364←{(∨\' '≠⍵)/⍵}                                   ⍝ Removing leading blanks                                 ⍝ ⍵←C1
f365←{∨\⍵}                                           ⍝ Not leading zeroes (turn on all zeroes after first one) ⍝ ⍵←B
f366←{(A-⌊0.5×(A←+/^\⌽A)++/^\A←' '=⌽⍵)⌽⍵}            ⍝ Centering character array ⍵ with ragged edges           ⍝ ⍵←C
f367←{(∨/A)⌿(⍴⍵)⍴(,A)\(,A←^\('                       ⍝'≠⍵)∨≠\⍵='''')/,⍵}                                       ⍝ Decommenting a matrix representation of a function (`⎕CR`) ⍝ ⍵←C2
f369←{(-⌊0.5×+/^\' '=⌽⍵)⌽⍵}                          ⍝ Centering character array ⍵ with only right edge ragged ⍝ ⍵←C
f370←{(-+/^\⌽' '=⍵)⌽⍵}                               ⍝ Justifying right                                        ⍝ ⍵←C
f371←{(-+/^\⌽' '=⍵)↓⍵}                               ⍝ Removing trailing blanks                                ⍝ ⍵←C1
f372←{(+/^\' '=⍵)⌽⍵}                                 ⍝ Justifying left                                         ⍝ ⍵←C
f373←{((~(⍴A↑⍵)↑'/'=⍺)/A↑⍵),(1↓A↓⍺),(A←+/^\⍺≠',')↓⍵} ⍝ Editing ⍵ with ⍺ '-wise                                 ⍝ ⍵←C1; ⍺←C1
f374←{(+/^\' '=⍵)↓⍵}                                 ⍝ Removing leading blanks                                 ⍝ ⍵←C1
f375←{⎕IO++/^\' '≠⍵}                                 ⍝ Indices of first blanks in rows of array ⍵              ⍝ ⍵←C
f377←{^\⍵}                                           ⍝ Leading ones (turn off all ones after first zero)       ⍝ ⍵←B

⍝⍝⍝ Scan +\
f378←{(⍳+/⍵,⍺)∊+\1+¯1↓0,((⍳+/⍵)∊+\⍵)\⍺}     ⍝ Vector (`⍵[1]⍴1),(⍺[1]⍴0),(⍵[2]⍴1),...`                      ⍝ ⍵←I1; ⍺←I1
f379←{((⍵≠0)/⍺)[+\¯1⌽(⍳+/⍵)∊+\⍵]}           ⍝ Replicate ⍺[i] ⍵[i] times (for all i)                        ⍝ ⍵←I1; ⍺←A1
f380←{⎕IO++\1+((⍳+/⍵)∊+\⎕IO,⍵)\⍺-¯1↓1,⍵+⍺}  ⍝ Vector (`⍺[1]+⍳⍵[1]),(⍺[2]+⍳⍵[2]),(⍺[3]+⍳⍵[3]),...`          ⍝ ⍵←I1; ⍺←I1; ⍴⍵←→⍴⍺
f381←{⍺[+\(⍳+/⍵)∊¯1↓1++\0,⍵]}               ⍝ Replicate ⍺[i] ⍵[i] times (for all i)                        ⍝ ⍵←I1; ⍺←A1
f382←{⍺[⎕IO++\(⍳+/⍵)∊⎕IO++\⍵]}              ⍝ Replicate ⍺[i] ⍵[i] times (for all i)                        ⍝ ⍵←I1; ⍺←A1
f383←{+\⍺-⍵\A-¯1↓0,A←⍵/+\¯1↓0,⍺}            ⍝ Cumulative sums (+\) over subvectors of ⍺ indicated by ⍵     ⍝ ⍵←B1; ⍺←D1
f384←{A-¯1↓0,A←(+\⍺)[+\⍵]}                  ⍝ Sums over (+/) subvectors of ⍺, lengths in ⍵                 ⍝ ⍵←I1; ⍺←D1
f386←{+\+\⍳⍵}                               ⍝ ⍵ first figurate numbers                                     ⍝ ⍵←I0
f387←{(⍳(⍴⍺)++/⍵)∊+\1+¯1↓0,(1⌽⍺)\⍵}         ⍝ Insert vector for ⍵[i] zeroes after i:th subvector           ⍝ ⍵←I1; ⍺←B1
f388←{((⍳(⍴⍺)++/⍵)∊+\1+¯1↓0,((⍳⍴⍺)∊G)\⍵)\⍺} ⍝ Open a gap of ⍵[i] after ⍺[G[i]] (for all i)                 ⍝ ⍵←I1; ⍺←A1; G←I1
f389←{((⍳(⍴⍺)++/⍵)∊+\1+((⍳⍴⍺)∊G)\⍵)\⍺}      ⍝ Open a gap of ⍵[i] before ⍺[G[i]] (for all i)                ⍝ ⍵←I1; ⍺←A1; G←I1
f390←{A←(+/⍵)⍴0 ⋄ A[+\¯1↓⎕IO,⍵]←1 ⋄ A}      ⍝ Changing lengths ⍵ of subvectors to starting indicators      ⍝ ⍵←I1
f391←{(⍳+/⍵)∊(+\⍵)-~⎕IO}                    ⍝ Changing lengths ⍵ of subvectors to ending indicators        ⍝ ⍵←I1
f392←{(⍳+/⍵)∊+\⎕IO,⍵}                       ⍝ Changing lengths ⍵ of subvectors to starting indicators      ⍝ ⍵←I1
f393←{(⍳+/A)∊+\A←1+⍵}                       ⍝ Insert vector for ⍵[i] elements before i:th element          ⍝ ⍵←I1
f394←{A-¯1↓0,A←(1⌽⍵)/+\⍺}                   ⍝ Sums over (+/) subvectors of ⍺ indicated by ⍵                ⍝ ⍵←B1; ⍺←D1
f395←{G-¯1↓0,G←0⌈(+\⍺)-⍵}                   ⍝ Fifo stock ⍺ decremented with ⍵ units                        ⍝ ⍺←D1; ⍵←D0
f396←{A∨¯1↓0,A←2|+\⍵=''''}                  ⍝ Locations of texts between and including quotes              ⍝ ⍵←C1
f397←{A^¯1↓0,A←2|+\⍵=''''}                  ⍝ Locations of texts between quotes                            ⍝ ⍵←C1
f398←{1↓(⍵=+\⍺=1↑⍺)/⍺}                      ⍝ ⍵:th subvector of ⍺ (subvectors separated by ⍺[1])           ⍝ ⍺←A1; ⍵←I0
f399←{(⍺=+\⍵=1↑⍵)/⍵}                        ⍝ Locating field number ⍺ starting with first element of ⍵     ⍝ ⍺←I0; ⍵←C1
f400←{A-¯1↓0,A←(⍺≠1↓⍺,0)/+\⍵}               ⍝ Sum elements of ⍵ marked by succeeding identicals in ⍺       ⍝ ⍵←D1; ⍺←D1
f401←{⍺^A∊(⍵^⍺)/A←+\⍺>¯1↓0,⍺}               ⍝ Groups of ones in ⍺ pointed to by ⍵                          ⍝ ⍵←B1; ⍺←B1
f402←{(+\⍵)∊⍺/⍳⍴⍺}                          ⍝ ith starting indicators ⍵                                    ⍝ ⍵←B1; ⍺←B1
f403←{(G=+\⍵)/⍺}                            ⍝ G:th subvector of ⍺ (subvectors indicated by ⍵)              ⍝ ⍵←B1; ⍺←A1; G←I0
f404←{((⍺-1)↓A)-0,(-⍺)↓A←+\⍵}               ⍝ Running sum of ⍺ consecutive elements of ⍵                   ⍝ ⍵←D1; ⍺←I0
f405←{+\('('=⍵)-¯1↓0,')'=⍵}                 ⍝ Depth of parentheses                                         ⍝ ⍵←C1
f406←{+\¯1↓⎕IO,⍵}                           ⍝ Starting positions of subvectors having lengths ⍵            ⍝ ⍵←I1
f407←{(⍳⍴⍺)∊(+\⍵)-~⎕IO}                     ⍝ Changing lengths ⍵ of subvectors of ⍺ to ending indicators   ⍝ ⍵←I1
f408←{(⍳⍴⍺)∊+\⎕IO,⍵}                        ⍝ Changing lengths ⍵ of subvectors of ⍺ to starting indicators ⍝ ⍵←I1
f409←{+\⍳⍵}                                 ⍝ ⍵ first triangular numbers                                   ⍝ ⍵←I0
f410←{+\⍵}                                  ⍝ Cumulative sum                                               ⍝ ⍵←D

⍝⍝⍝ Reduction ○/ ÷/ -/ ×/
f411←{○/¯2 1,⍵}                                     ⍝ Complementary angle (arccos sin ⍵)                        ⍝ ⍵←D0
f412←{-/×/0 1⊖⍵}                                    ⍝ Evaluating a two-row determinant                          ⍝ ⍵←D2
f413←{-/×⌿0 1⌽⍵}                                    ⍝ Evaluating a two-row determinant                          ⍝ ⍵←D2
f414←{(×/(+/⍵÷2)-0,⍵)*.5}                           ⍝ Area of triangle with side lengths in ⍵ (Heron's formula) ⍝ ⍵←D1; 3 ←→ ⍴⍵
f415←{(×⌿2 2⍴1,⍴⍵)⍴2 1 3⍉⍵}                         ⍝ Juxtapositioning planes of rank 3 array ⍵                 ⍝ ⍵←A3
f416←{×/¯1↓⍴⍵}                                      ⍝ Number of rows in array ⍵ (also of a vector)              ⍝ ⍵←A
f417←{(-⍵[2]-¯1 1×((⍵[2]*2)-×/4,⍵[1 3])*.5)÷2×⍵[1]} ⍝ (Real) solution of quadratic equation with coefficients ⍵ ⍝ ⍵←D1; 3 ←→ ⍴⍵
f418←{(×/2 2⍴1,⍴⍵)⍴⍵}                               ⍝ Reshaping planes of rank 3 array to rows of a matrix      ⍝ ⍵←A3
f419←{(×/2 2⍴(⍴⍵),1)⍴⍵}                             ⍝ Reshaping planes of rank 3 array to a matrix              ⍝ ⍵←A3
f420←{×/⍴⍵}                                         ⍝ Number of elements (also of a scalar)                     ⍝ ⍵←A
f421←{×/⍵}                                          ⍝ Product of elements of ⍵                                  ⍝ ⍵←D1
f422←{÷/⍵}                                          ⍝ Alternating product                                       ⍝ ⍵←D
f423←{⍺↑((⌊-/.5×⍺,⍴⍵)⍴' '),⍵}                       ⍝ Centering text line ⍵ into a field of width ⍺             ⍝ ⍵←C1; ⍺←I0
f424←{-/⍵}                                          ⍝ Alternating sum                                           ⍝ ⍵←D

⍝⍝⍝ Reduction ⌈/ ⌊/
f425←{(⌈/⍵)=⌊/⍵}                   ⍝ Test if all elements of vector ⍵ are equal                  ⍝ ⍵←D1
f426←{(⌈/⍵)-⌊/⍵}                   ⍝ Size of range of elements of ⍵                              ⍝ ⍵←D1
f427←{(⍳⌈/⍵)∊⍵}                    ⍝ Conversion of set of positive integers ⍵ to a mask          ⍝ ⍵←I1
f428←{⌈/⍳0}                        ⍝ Negative infinity; the smallest representable value         ⍝ 
f429←{⍵,[1+.5×⌈/(⍴⍴⍵),⍴⍴⍺]⍺}       ⍝ Vectors as column matrices in catenation beneath each other ⍝ ⍵←A1/2; ⍺←A1/2
f430←{⍵,[.5×⌈/(⍴⍴⍵),⍴⍴⍺]⍺}         ⍝ Vectors as row matrices in catenation upon each other       ⍝ ⍵←A1/2; ⍺←A1/2
f431←{A←(⌈/⍵,⍺)⍴0 ⋄ A[⍺]←1 ⋄ A[⍵]} ⍝ Quick membership (`∊`) for positive integers                ⍝ ⍵←I1; ⍺←I1
f432←{⌈/⍵,0}                       ⍝ Positive maximum, at least zero (also for empty ⍵)          ⍝ ⍵←D1
f433←{⌈/⍵}                         ⍝ Maximum of elements of ⍵                                    ⍝ ⍵←D1
f434←{⌊/⍳0}                        ⍝ Positive infinity; the largest representable value          ⍝ 
f435←{⌊/⍵}                         ⍝ Minimum of elements of ⍵                                    ⍝ ⍵←D1

⍝⍝⍝ Reduction ∨/ ⍲/ ≠/
f436←{⍲/0 1∊⍵}               ⍝ Test if all elements of vector ⍵ are equal    ⍝ ⍵←B1
f437←{(^/⍵)∨~∨/⍵}            ⍝ Test if all elements of vector ⍵ are equal    ⍝ ⍵←B1
f438←{(^/⍵)=∨/⍵}             ⍝ Test if all elements of vector ⍵ are equal    ⍝ ⍵←B1
f439←{^/⍵÷∨/⍵}               ⍝ Test if all elements of vector ⍵ are equal    ⍝ ⍵←B1
f440←{(¯1⌽1↓(∨/⍵≠¯1⊖⍵),1)⌿⍵} ⍝ Removing duplicate rows from ordered matrix ⍵ ⍝ ⍵←A2
f441←{∨/0/⍵}                 ⍝ Vector having as many ones as ⍵ has rows      ⍝ ⍵←A2
f442←{∨/⍺∊⍵}                 ⍝ Test if ⍵ and ⍺ have elements in common       ⍝ ⍵←A; ⍺←A1
f443←{~∨/⍵}                  ⍝ None, neither                                 ⍝ ⍵←B
f444←{∨/⍵}                   ⍝ Any, anyone                                   ⍝ ⍵←B
f445←{≠/0 1∊⍵}               ⍝ Test if all elements of vector ⍵ are equal    ⍝ ⍵←B1
f446←{≠/⍵}                   ⍝ Parity                                        ⍝ ⍵←B

⍝⍝⍝ Reduction ^/
f447←{+/A^⍉A←^/⍵[;A⍴1;]≤2 1 3⍉⍵[;(A←1↑⍴⍵)⍴2;]} ⍝ Number of areas intersecting areas in ⍵     ⍝ ⍵←D3 (n × 2 × dim)
f448←{^/⍵/1⌽⍵}                                 ⍝ Test if all elements of vector ⍵ are equal  ⍝ ⍵←B1
f449←{^/⍵=1⊖⍵}                                 ⍝ Comparison of successive rows               ⍝ ⍵←A2
f450←{^/⍵=1⌽⍵}                                 ⍝ Test if all elements of vector ⍵ are equal  ⍝ ⍵←A1
f451←{^/((1↑⍵)∊10↓A),⍵∊A←'0..9A..Za..z'}       ⍝ Test if ⍵ is a valid APL name               ⍝ ⍵←C1
f452←{^/⍵=1↑⍵}                                 ⍝ Test if all elements of vector ⍵ are equal  ⍝ ⍵←A1
f453←{^/(⍵∊⍺),⍺∊⍵}                             ⍝ Identity of two sets                        ⍝ ⍵←A1; ⍺←A1
f454←{^/(⍳⍴⍵)∊⍵}                               ⍝ Test if ⍵ is a permutation vector           ⍝ ⍵←I1
f455←{~^/⍵∊~⍵}                                 ⍝ Test if all elements of vector ⍵ are equal  ⍝ ⍵←B1
f456←{^/,⍵∊0 1}                                ⍝ Test if ⍵ is boolean                        ⍝ ⍵←A
f457←{^/⍺∊⍵}                                   ⍝ Test if ⍺ is a subset of ⍵ (`⍺ ⊂ ⍵`)        ⍝ ⍵←A; ⍺←A1
f458←{^/,⍵=⍺}                                  ⍝ Test if arrays of equal shape are identical ⍝ ⍵←A; ⍺←A; ⍴⍵ ←→ ⍴⍺
f459←{^/⍵=⍵[1]}                                ⍝ Test if all elements of vector ⍵ are equal  ⍝ ⍵←A1
f460←{^/' '=⍵}                                 ⍝ Blank rows                                  ⍝ ⍵←C2
f461←{^/⍵}                                     ⍝ All, both                                   ⍝ ⍵←B

⍝⍝⍝ Reduction +/
f462←{((+/(⍵-(+/⍵)÷⍴⍵)*2)÷⍴⍵)*.5} ⍝ Standard deviation of ⍵                                  ⍝ ⍵←D1
f463←{(+/(⍵-(+/⍵)÷⍴⍵)*⍺)÷⍴⍵}      ⍝ ⍺:th moment of ⍵                                         ⍝ ⍵←D1
f464←{(+/(⍵-(+/⍵)÷⍴⍵)*2)÷⍴⍵}      ⍝ Variance (dispersion) of ⍵                               ⍝ ⍵←D1
f465←{(+/,⍵)÷1⌈⍴,⍵}               ⍝ Arithmetic average (mean value), also for an empty array ⍝ ⍵←D
f466←{0=(⍴⍵)|+/⍵}                 ⍝ Test if all elements of vector ⍵ are equal               ⍝ ⍵←B1
f467←{(+⌿⍵)÷1↑(⍴⍵),1}             ⍝ Average (mean value) of columns of matrix ⍵              ⍝ ⍵←D2
f468←{(+/⍵)÷¯1↑1,⍴⍵}              ⍝ Average (mean value) of rows of matrix ⍵                 ⍝ ⍵←D2
f469←{+/⍵=,⍺}                     ⍝ Number of occurrences of scalar ⍵ in array ⍺             ⍝ ⍵←A0; ⍺←A
f470←{(+/[⍺]⍵)÷(⍴⍵)[⍺]}           ⍝ Average (mean value) of elements of ⍵ along direction ⍺  ⍝ ⍵←D; ⍺←I0
f471←{(+/⍵)÷⍴⍵}                   ⍝ Arithmetic average (mean value)                          ⍝ ⍵←D1
f472←{÷+/÷⍵}                      ⍝ Resistance of parallel resistors                         ⍝ ⍵←D1
f473←{+/⍵}                        ⍝ Sum of elements of ⍵                                     ⍝ ⍵←D1
f474←{+/⍵}                        ⍝ Row sum of a matrix                                      ⍝ ⍵←D2
f475←{+⌿⍵}                        ⍝ Column sum of a matrix                                   ⍝ ⍵←D2
f476←{+/⍵}                        ⍝ Reshaping one-element vector ⍵ into a scalar             ⍝ ⍵←A1
f477←{+/⍵}                        ⍝ Number of elements satisfying condition ⍵                ⍝ ⍵←B1

⍝⍝⍝ Reverse ⌽ ⊖
f478←{⌽⍺\⌽⍵}                               ⍝ Scan from end with function `⍺`                       ⍝ ⍵←A
f479←{A←9999⍴⎕IO+⍴⍺ ⋄ A[⌽⍺]←⌽⍳⍴⍺ ⋄ A[⍵]}   ⍝ The index of positive integers in ⍺                   ⍝ ⍵←I; ⍺←I1
f480←{((⌽A)×1,⍺)⍴2 1 3⍉(1⌽⍺,A←(⍴⍵)÷1,⍺)⍴⍵} ⍝ 'Transpose' of matrix ⍵ with column fields of width ⍺ ⍝ ⍵←A2; G←I0
f482←{⍺+⍉(⌽⍴⍺)⍴⍵}                          ⍝ Adding ⍵ to each column of ⍺                          ⍝ ⍵←D1; ⍺←D; (⍴⍵)=1↑⍴⍺
f483←{⍉(⌽⍴⍺)⍴⍵}                            ⍝ Matrix with shape of ⍺ and ⍵ as its columns           ⍝ ⍵←A1; ⍺←A2
f484←{¯1↓⍵×⌽¯1+⍳⍴⍵}                        ⍝ Derivate of polynomial ⍵                              ⍝ ⍵←D1
f485←{,⌽[⎕IO+⍺](1,⍴⍵)⍴⍵}                   ⍝ Reverse vector ⍵ on condition ⍺                       ⍝ ⍵←A1; ⍺←B0
f486←{(⌽1,⍴⍵)⍴⍵}                           ⍝ Reshaping vector ⍵ into a one-column matrix           ⍝ ⍵←A1
⍝f487←{(⌽1, ...)}                          ⍝ Avoiding parentheses with help of reversal            ⍝ 

⍝⍝⍝ Rotate ⌽ ⊖
f488←{((1⌽⍵)×¯1⌽⍺)-(¯1⌽⍵)×1⌽⍺}    ⍝ Vector (cross) product of vectors                              ⍝ ⍵←D; ⍺←D
f489←{A⊖(A←(⍳⍵)-⌈⍵÷2)⌽(⍵,⍵)⍴⍳⍵×⍵} ⍝ A magic square, side ⍵                                         ⍝ ⍵←I0; 1=2|⍵
f490←{(¯1⌽1↓(⍵≠¯1⌽⍵),1)/⍵}        ⍝ Removing duplicates from an ordered vector                     ⍝ ⍵←A1
f491←{1⌽22⍴11⍴'''1⌽22⍴11⍴'''}     ⍝ An expression giving itself                                    ⍝ 
f492←{(⍺⌽1 2)⍉⍵}                  ⍝ Transpose matrix ⍵ on condition ⍺                              ⍝ ⍵←A2; ⍺←B0
f493←{(⍵/⍺)≥A/1⌽A←(⍺∨⍵)/⍵}        ⍝ Any element true (`∨/`) on each subvector of ⍺ indicated by ⍵  ⍝ ⍵←B1; ⍺←B1
f494←{(⍵/⍺)^A/1⌽A←(⍺≤⍵)/⍵}        ⍝ All elements true (`^/`) on each subvector of ⍺ indicated by ⍵ ⍝ ⍵←B1; ⍺←B1
f495←{(1↑A)↓(A⍲1⌽A←⍺=⍵)/⍵}        ⍝ Removing leading, multiple and trailing ⍺'s                    ⍝ ⍵←A1; ⍺←A0
f496←{A-¯1↓0,A←(1⌽⍵)/⍳⍴⍵}         ⍝ Changing starting indicators ⍵ of subvectors to lengths        ⍝ ⍵←B1
f498←{(A∨1⌽A←⍵≠' ')/⍵}            ⍝ (Cyclic) compression of successive blanks                      ⍝ ⍵←C1
f499←{(1-⍳¯1↑⍴⍵)⌽⍵}               ⍝ Aligning columns of matrix ⍵ to diagonals                      ⍝ ⍵←A2
f500←{(¯1+⍳¯1↑⍴⍵)⌽⍵}              ⍝ Aligning diagonals of matrix ⍵ to columns                      ⍝ ⍵←A2
f501←{0 ¯1↓(-⍳⍴⍵)⌽((2⍴⍴⍵)⍴0),⍵}   ⍝ Diagonal matrix with elements of ⍵                             ⍝ ⍵←D1
f502←{1,1↓⍵≠¯1⌽⍵}                 ⍝ Test if elements differ from previous ones (non-empty ⍵)       ⍝ ⍵←A1
f503←{(¯1↓⍵≠1⌽⍵),1}               ⍝ Test if elements differ from next ones (non-empty ⍵)           ⍝ ⍵←A1
f504←{¯1⌽1↓⍵,⍺}                   ⍝ Replacing first element of ⍵ with ⍺                            ⍝ ⍵←A1; ⍺←A0
f505←{1⌽¯1↓⍺,⍵}                   ⍝ Replacing last element of ⍵ with ⍺                             ⍝ ⍵←A1; ⍺←A0
f506←{1⌽(⍳⍴⍵)∊⍺}                  ⍝ Ending points for ⍵ in indices pointed by ⍺                    ⍝ ⍵←A1; ⍺←I1
f507←{¯1⌽⍵}                       ⍝ Leftmost neighboring elements cyclically                       ⍝ ⍵←A
f508←{1⌽⍵}                        ⍝ Rightmost neighboring elements cyclically                      ⍝ ⍵←A

⍝⍝⍝ Transpose ⍉
f509←{⍉ ... ⍉⍵}                ⍝ Applying to columns action defined on rows               ⍝ ⍵←A1; ⍺←I0
f510←{1 1⍉⍵[⍺[1;];⍺[2;]]}      ⍝ Retrieving scattered elements ⍺ from matrix ⍵            ⍝ ⍵←A2; ⍺←I2
f511←{⍵[⍺]⍉G}                  ⍝ Successive transposes of G (⍵ after ⍺: `⍵⍉⍺⍉G`)          ⍝ ⍵←I1; ⍺←I1
f512←{(1*⍴⍵)⍉⍵}                ⍝ Major diagonal of array ⍵                                ⍝ ⍵←A
f513←{40 120⍴2 1 3⍉10 40 12⍴⍵} ⍝ Reshaping a 400×12 character matrix to fit into one page ⍝ ⍵←C2
f514←{1 3 2⍉⍵}                 ⍝ Transpose of planes of a rank three array                ⍝ ⍵←A3
f515←{1 1⍉⍵}                   ⍝ Major diagonal of matrix ⍵                               ⍝ ⍵←A2
f516←{G⍉⍵∘.f ⍺}                ⍝ Selecting specific elements from a 'large' outer product ⍝ ⍵←A; ⍺←A; G←I1
f517←{~0∊⍵=-⍉⍵}                ⍝ Test for antisymmetricity of square matrix ⍵             ⍝ ⍵←D2
f518←{~0∊⍵=⍉⍵}                 ⍝ Test for symmetricity of square matrix ⍵                 ⍝ ⍵←A2
f519←{⍉(⍵,⍴⍺)⍴⍺}               ⍝ Matrix with ⍵ columns ⍺                                  ⍝ ⍵←I0; ⍺←D1

⍝⍝⍝ Maximum ⌈ Minimum ⌊
f520←{⍺[1]⌈⍺[2]⌊⍵}                 ⍝ Limiting ⍵ between ⍺[1] and ⍺[2], inclusive          ⍝ ⍵←D; ⍺←D1
f521←{(A↑⍵),[⍳1](1↓A←(⍴⍵)⌈0,⍴⍺)↑⍺} ⍝ Inserting vector ⍺ to the end of matrix ⍵            ⍝ ⍵←A2; ⍺←A1
f522←{((0 1×⍴⍺)⌈⍴⍵)↑⍵}             ⍝ Widening matrix ⍵ to be compatible with ⍺            ⍝ ⍵←A2; ⍺←A2
f523←{((1 0×⍴⍺)⌈⍴⍵)↑⍵}             ⍝ Lengthening matrix ⍵ to be compatible with ⍺         ⍝ ⍵←A2; ⍺←A2
f524←{(1⌈¯2↑⍴⍵)⍴⍵}                 ⍝ Reshaping non-empty lower-rank array ⍵ into a matrix ⍝ ⍵←A; 2≥⍴⍴⍵
f525←{(⍵⌊⍴⍺)↑⍺}                    ⍝ Take of at most ⍵ elements from ⍺                    ⍝ ⍵←I; ⍺←A
f526←{(⍵,G)[(1+⍴⍵)⌊⍺]}             ⍝ Limiting indices and giving a default value G        ⍝ ⍵←A1; ⍺←I; G←A0

⍝⍝⍝ Ceiling ⌈ Floor ⌊
f527←{((⌈(⍴,⍵)÷⍺),⍺)⍴⍵}              ⍝ Reshaping ⍵ into a matrix of width ⍺              ⍝ ⍵←D, ⍺←I0
f528←{⌊⍵+1≤2|⍵}                      ⍝ Rounding to nearest even integer                  ⍝ ⍵←D
f529←{⌊⍵+.5×.5≠2|⍵}                  ⍝ Rounding, to nearest even integer for .5 = 1      ⍝ ⍵
f530←{⌊⍵+.5×.5≠2|⍵}                  ⍝ Rounding, to nearest even integer for .5 = 1      ⍝ ⍵
f531←{⍵+(G××⍺-⍵)×(⍳1+|⌊(⍺-⍵)÷G)-⎕IO} ⍝ Arithmetic progression from ⍵ to ⍺ with step G    ⍝ ⍵←D0; ⍺←D0; G←D0
f532←{(-⌊.5×⍺+⍴⍵)↑⍵}                 ⍝ Centering text line ⍵ into a field of width ⍺     ⍝ ⍵←C1; ⍺←I0
f533←{⍵=⌊⍵}                          ⍝ Test if integer                                   ⍝ ⍵←D
f534←{.05×⌊.5+⍵÷.05}                 ⍝ Rounding currencies to nearest 5 subunits         ⍝ ⍵←D
f535←{⌊⍵÷1000}                       ⍝ First part of numeric code ABBB                   ⍝ ⍵←I
f536←{(10*-⍵)×⌊0.5+⍺×10*⍵}           ⍝ Rounding to ⍵ decimals                            ⍝ ⍵←I; ⍺←D
f537←{0.01×⌊0.5+100×⍵}               ⍝ Rounding to nearest hundredth                     ⍝ ⍵←D
f538←{⌊0.5+⍵}                        ⍝ Rounding to nearest integer                       ⍝ ⍵←D
f539←{⌊⍵}                            ⍝ Demote floating point representations to integers ⍝ ⍵←I

⍝⍝⍝ Residue |
f540←{(0=400|⍵)∨(0≠100|⍵)^0=4|⍵}                  ⍝ Test if ⍵ is a leap year                             ⍝ ⍵←I
f541←{'_',[1]('|',⍵,'|'),[1]'¯'}                  ⍝ Framing                                              ⍝ ⍵←C2
f542←{1}                                          ⍝ Magnitude of fractional part                         ⍝ ⍵←D
f543←{(×⍵)|⍵}                                     ⍝ Fractional part with sign                            ⍝ ⍵←D
f544←{⍵,(⍺|-⍴⍵)↑0/⍵}                              ⍝ Increasing the dimension of ⍵ to multiple of ⍺       ⍝ ⍵←A1; ⍺←I0
f545←{(0≠⍺|⍳⍴⍵)/⍵}                                ⍝ Removing every ⍺:th element of ⍵                     ⍝ ⍵←A1; ⍺←I0
f546←{(0=⍺|⍳⍴⍵)/⍵}                                ⍝ Taking every ⍺:th element of ⍵                       ⍝ ⍵←A1; ⍺←I0
f547←{(0=A|⍵)/A←⍳⍵}                               ⍝ Divisors of ⍵                                        ⍝ ⍵←I0
f548←{(2|⍳⍴⍵)/⍵}                                  ⍝ Removing every second element of ⍵                   ⍝ ⍵←A1
f549←{(0=⍺|⍵)/⍵}                                  ⍝ Elements of ⍵ divisible by ⍺                         ⍝ ⍵←D1; ⍺←D0/1
f550←{(A×⍺[1]*¯1 1)⍴(A←(⍴⍵)+(⍺[1]|-1↑⍴⍵),⍺[2])↑⍵} ⍝ Ravel of a matrix to ⍺[1] columns with a gap of ⍺[2] ⍝ ⍵←A2; ⍺←I1
f551←{~2|⍵}                                       ⍝ Test if even                                         ⍝ ⍵←I
f552←{1000|⍵}                                     ⍝ Last part of numeric code ABBB                       ⍝ ⍵←I
f553←{1|⍵}                                        ⍝ Fractional part                                      ⍝ ⍵←D

⍝⍝⍝ Magnitude |, Signum ×
f554←{(×⍵)×⍺+|⍵} ⍝ Increasing absolute value without change of sign ⍝ ⍵←D; ⍺←D
f555←{⍵×⍺≤|⍵}    ⍝ Rounding to zero values of ⍵ close to zero       ⍝ ⍵←D; ⍺←D
f556←{⍵×|⍵}      ⍝ Square of elements of ⍵ without change of sign   ⍝ ⍵←D
f557←{⍺[2+×⍵]}   ⍝ Choosing according to signum                     ⍝ ⍵←D; ⍺←A1

⍝⍝⍝ Expand \ ⍀
f558←{~(B^⍵)∨(B∨⍵)\A>¯1↓0,A←(B∨⍵)/B←~⍺} ⍝ Not first zero (≤\) in each subvector of ⍺ indicated by ⍵   ⍝ ⍵←B1; ⍺←B1
f559←{(⍺^⍵)∨(⍺∨⍵)\A>¯1↓0,A←(⍺∨⍵)/⍺}     ⍝ First one (⍵←B1; ⍺←B1                                       ⍝ 
f560←{A\(A←~⍵∊⍺)/⍵}                     ⍝ Replacing elements of ⍵ in set ⍺ with blanks/zeroes         ⍝ ⍵←A0; ⍺←A1
f561←{A\(A←⍵∊⍺)/⍵}                      ⍝ Replacing elements of ⍵ not in set ⍺ with blanks/zeroes     ⍝ ⍵←A1; ⍺←A
f562←{A←G\⍵ ⋄ A[(~G)/⍳⍴G]←⍺ ⋄ A}        ⍝ Merging ⍵ and ⍺ under control of G (mesh)                   ⍝ ⍵←A1; ⍺←A1; G←B1
f563←{⍺\⍺/⍵}                            ⍝ Replacing elements of ⍵ not satisfying ⍺ with blanks/zeroes ⍝ ⍵←A; ⍺←B1
f564←{(~(⍳(⍴⍺)+1⍴⍴⍵)∊⍺+⍳⍴⍺)⍀⍵}          ⍝ Adding an empty row into ⍵ after rows ⍺                     ⍝ ⍵←A2; ⍺←I1
f565←{0∊0\0⍴⍵}                          ⍝ Test if numeric                                             ⍝ ⍵←A1
f566←{((⍺+1)≠⍳1+1⍴⍴⍵)⍀⍵}                ⍝ Adding an empty row into ⍵ after row ⍺                      ⍝ ⍵←A2; ⍺←I0
f567←{⍵,[⎕IO-.1](' '≠⍵)\'¯'}            ⍝ Underlining words                                           ⍝ ⍵←C1
f568←{(⍴⍺)⍴(,⍺)\⍵}                      ⍝ Using boolean matrix ⍺ in expanding ⍵                       ⍝ ⍵←A1; ⍺←B2
f569←{((2×⍴⍵)⍴1 0)\⍵}                   ⍝ Spacing out text                                            ⍝ ⍵←C1

⍝⍝⍝ Compress / ⌿
f570←{(A>0)/A←(1↓A)-1+¯1↓A←(~A)/⍳⍴A←0,⍵,0}     ⍝ Lengths of groups of ones in ⍵                            ⍝ ⍵←B1
f571←{(~A∊1,⍴⍵)/A←A/⍳⍴A←(1↓A,0)←~⍵∊'aeiouyÄÖ'} ⍝ Syllabization of a Finnish word ⍵                         ⍝ ⍵←A1
f572←{(G/⍵),(~G)/⍺}                            ⍝ Choosing a string according to boolean value G            ⍝ ⍵←C1; ⍺←C1; G←B0
f573←{(' '=1↑⍵)↓((1↓A,0)∨A←' '≠⍵)/⍵}           ⍝ Removing leading, multiple and trailing blanks            ⍝ ⍵←C1
f575←{(~(⍳¯1↑⍴⍵)∊⍺)/⍵}                         ⍝ Removing columns ⍺ from array ⍵                           ⍝ ⍵←A; ⍺←I1
f576←{(¯1↑(' '≠⍵)/⍳⍴⍵)⍴⍵}                      ⍝ Removing trailing blanks                                  ⍝ ⍵←C1
f577←{(1↓A)-¯1↓A←(A,1)/⍳1+⍴A←1,(1↓⍵)≠¯1↓⍵}     ⍝ Lengths of subvectors of ⍵ having equal elements          ⍝ ⍵←A1
f578←{G-¯1↓0,G←(~⎕IO)+(((1↓⍵)≠¯1↓⍵),1)/⍳⍴⍵}    ⍝ Field lengths of vector ⍵; G ←→ ending indices            ⍝ ⍵←A1; G←I1
f580←{((1↓A,0)∨A←' '≠⍵)/⍵}                     ⍝ Removing multiple and trailing blanks                     ⍝ ⍵←C1
f581←{(A∨¯1↓0,A←' '≠⍵)/⍵}                      ⍝ Removing leading and multiple blanks                      ⍝ ⍵←C1
f582←{(A∨¯1↓1,A←' '≠⍵)/⍵}                      ⍝ Removing multiple blanks                                  ⍝ ⍵←C1
f583←{(A∨¯1↓1,A←⍵≠⍺)/⍵}                        ⍝ Removing duplicate ⍺'s from vector ⍵                      ⍝ ⍵←A1; ⍺←A0
f584←{(⍵∊⍺)/⍳⍴⍵}                               ⍝ Indices of all occurrences of elements of ⍺ in ⍵          ⍝ ⍵←A1; ⍺←A
f585←{⍺,(~⍵∊⍺)/⍵}                              ⍝ Union of sets, ?                                          ⍝ ⍵←A1; ⍺←A1
f586←{(~⍵∊⍺)/⍵}                                ⍝ Elements of ⍵ not in ⍺ (difference of sets)               ⍝ ⍵←A1; ⍺←A
f587←{(⍵[;1]∊⍺)⌿⍵}                             ⍝ Rows of non-empty matrix ⍵ starting with a character in ⍺ ⍝ ⍵←A2; ⍺←A1
f588←{(⍵∊⍺)/⍵}                                 ⍝ Intersection of sets, ⍞                                   ⍝ ⍵←A1; ⍺←A
f589←{((⍴⍵)*⍺≠⍳⍴⍴⍵)⍴ ⍺/[⍺]⍵}                   ⍝ Reduction with function ⍺ in dimension ⍺, rank unchanged  ⍝ ⍺←I0; ⍵←A
f590←{A[(A=⍵)/⍳⍴A←,G]←⍺ ⋄ (⍴G)⍴A}              ⍝ Replacing all values ⍵ in G with ⍺                        ⍝ ⍵←A0; ⍺←A0; G←A
f591←{(⍺=⍵)/⍳⍴⍵}                               ⍝ Indices of all occurrences of ⍺ in ⍵                      ⍝ ⍵←A1; ⍺←A0
f592←{G[⍵/⍳⍴G]←⍺}                              ⍝ Replacing elements of G satisfying ⍵ with ⍺               ⍝ ⍺←A0; ⍵←B1; G←A1
f593←{A←9999⍴0 ⋄ A[⍵]←1 ⋄ A/⍳9999}             ⍝ Removing duplicates from positive integers                ⍝ ⍵←I1
f594←{⍵/⍳⍴⍵}                                   ⍝ Indices of ones in logical vector ⍵                       ⍝ ⍵←B1
f595←{((~⍵)/'IN'),'CORRECT'}                   ⍝ Conditional in text                                       ⍝ ⍵←B0
f596←{(' '≠⍵)/⍵}                               ⍝ Removing blanks                                           ⍝ ⍵←A1
f597←{(⍵≠⍺)/⍵}                                 ⍝ Removing elements ⍺ from vector ⍵                         ⍝ ⍵←A1; ⍺←A0
f598←{(,⍵,[1.5]1)/,⍵,[1.5]~⍵}                  ⍝ Vector to expand a new element after each one in ⍵        ⍝ ⍵←B1
f599←{⍺/,⍵}                                    ⍝ Reduction with FUNCTION `⍺` without respect to shape      ⍝ ⍵←D
f600←{1/⍵}                                     ⍝ Reshaping scalar ⍵ into a one-element vector              ⍝ ⍵←A
f601←{0⌿⍵}                                     ⍝ Empty matrix                                              ⍝ ⍵←A2
f602←{⍺/⍵}                                     ⍝ Selecting elements of ⍵ satisfying condition ⍺            ⍝ ⍵←A; ⍺←B1

⍝⍝⍝ Take ↑
f603←{⍺[⍳G;],[1]((1↓⍴⍺)↑⍵),[1](2↑G)↓⍺} ⍝ Inserting vector ⍵ into matrix ⍺ after row G                ⍝ ⍵←A1; ⍺←A2; G←I0
f604←{⍺↑⍵,⍺⍴¯1↑⍵}                      ⍝ Filling ⍵ with last element of ⍵ to length ⍺                ⍝ ⍵←A1; ⍺←I0
⍝f605←{⍵[⍺;]←(1↑⍴⍵)↑⍞}                 ⍝ Input of row ⍺ of text matrix ⍵                             ⍝ ⍵←C2; ⍺←I0
f606←{⍵>((-⍴⍴⍵)↑¯1)↓0,⍵}               ⍝ First ones in groups of ones                                ⍝ ⍵←B
f607←{(G↑⍺),⍵,G↓⍺}                     ⍝ Inserting ⍵ into ⍺ after index G                            ⍝ ⍵←A1; ⍺←A1; G←I0
f608←{⍵-((-⍴⍴⍵)↑¯1)↓0,⍵}               ⍝ Pairwise differences of successive columns (inverse of +\)  ⍝ ⍵←D
f609←{((-⍴⍴⍵)↑¯1)↓0,⍵}                 ⍝ Leftmost neighboring elements                               ⍝ ⍵←D
f610←{((-⍴⍴⍵)↑1)↓⍵,0}                  ⍝ Rightmost neighboring elements                              ⍝ ⍵←D
f611←{(-⍴⍵)↑(-⍺)↓⍵}                    ⍝ Shifting vector ⍵ right with ⍺ without rotate               ⍝ ⍵←A1; ⍺←I0
f612←{(⍴⍵)↑⍺↓⍵}                        ⍝ Shifting vector ⍵ left with ⍺ without rotate                ⍝ ⍵←A1; ⍺←I0
f613←{(2↑⍺)↓⍵}                         ⍝ Drop of ⍺ first rows from matrix ⍵                          ⍝ ⍵←A2; ⍺←I0
f614←{0∊1↑0⍴⍵}                         ⍝ Test if numeric                                             ⍝ ⍵←A
f615←{(¯2↑1 1,⍴⍵)⍴⍵}                   ⍝ Reshaping non-empty lower-rank array ⍵ into a matrix        ⍝ ⍵←A; 2≥⍴⍴⍵
⍝f616←{1↑⍞,⍵}                          ⍝ Giving a character default value for input                  ⍝ ⍵←C0
f617←{⍵+(-⍴⍵)↑⍺}                       ⍝ Adding scalar ⍺ to last element of ⍵                        ⍝ ⍵←D; ⍺←D0
f618←{1↑⍴⍵}                            ⍝ Number of rows in matrix ⍵                                  ⍝ ⍵←A2
f619←{¯1↑⍴⍵}                           ⍝ Number of columns in matrix ⍵                               ⍝ ⍵←A2
f620←{(⍵×⍺)⍴(-⍺)↑1}                    ⍝ Ending points for ⍵ fields of width ⍺                       ⍝ ⍵←I0; ⍺←I0
f621←{(⍵×⍺)⍴⍺↑1}                       ⍝ Starting points for ⍵ fields of width ⍺                     ⍝ ⍵←I0; ⍺←I0
f622←{1↑0⍴⍵}                           ⍝ Zero or space depending on the type of ⍵ (fill element)     ⍝ ⍵←A
f623←{1 80⍴80↑⍵}                       ⍝ Forming first row of a matrix to be expanded                ⍝ ⍵←A1
f624←{⍺↑⍵⍴1}                           ⍝ Vector of length ⍺ with ⍵ ones on the left, the rest zeroes ⍝ ⍵←I0; ⍺←I0
f625←{(-⍺)↑⍵}                          ⍝ Justifying text ⍵ to right edge of field of width ⍺         ⍝ ⍺←I0; ⍵←C1

⍝⍝⍝ Drop ↓
f627←{1,(1↓⍵)≠¯1↓⍵}       ⍝ Starting points of groups of equal elements (non-empty ⍵) ⍝ ⍵←A1
f628←{((1↓⍵)≠¯1↓⍵),1}     ⍝ Ending points of groups of equal elements (non-empty ⍵)   ⍝ ⍵←A1
f629←{(1↓⍵)÷¯1↓⍵}         ⍝ Pairwise ratios of successive elements of vector ⍵        ⍝ ⍵←D1
f630←{(1↓⍵)-¯1↓⍵}         ⍝ Pairwise differences of successive elements of vector ⍵   ⍝ ⍵←D1
f631←{⍵-(-⍺=⍳⍴⍴⍵)↓0,[⍺]⍵} ⍝ Differences of successive elements of ⍵ along direction ⍺ ⍝ ⍵←D; ⍺←I0
f632←{(⍺-1)↓⍳⍵}           ⍝ Ascending series of integers ⍺..⍵ (for small ⍺ and ⍵)     ⍝ ⍵←I0; ⍺←I0
f633←{⍵>¯1↓0,⍵}           ⍝ First ones in groups of ones                              ⍝ ⍵←B1
f634←{⍵>1↓⍵,0}            ⍝ Last ones in groups of ones                               ⍝ ⍵←B1
f635←{1↓,',',⍵}           ⍝ List of names in ⍵ (one per row)                          ⍝ ⍵←C2
f636←{''⍴G↓⍵,⍺}           ⍝ Selection of ⍵ or ⍺ depending on condition G              ⍝ ⍵←A0; ⍺←A0; G←B0
f637←{⍵-¯1↓0,⍵}           ⍝ Restoring argument of cumulative sum (inverse of +\)      ⍝ ⍵←D1
f638←{(⍺,0)↓⍵}            ⍝ Drop of ⍺ first rows from matrix ⍵                        ⍝ ⍵←A2; ⍺←I0
f639←{(0,⍺)↓⍵}            ⍝ Drop of ⍺ first columns from matrix ⍵                     ⍝ ⍵←A2; ⍺←I0
f640←{¯1↓⍴⍵}              ⍝ Number of rows in matrix ⍵                                ⍝ ⍵←A2
f641←{1↓⍴⍵}               ⍝ Number of columns in matrix ⍵                             ⍝ ⍵←A2
f642←{(⍺×G)↓⍵}            ⍝ Conditional drop of ⍺ elements from array ⍵               ⍝ ⍵←A; ⍺←I1; G←B1
f643←{(-⍺)↓⍵}             ⍝ Conditional drop of last element of ⍵                     ⍝ ⍵←A1; ⍺←B0

⍝⍝⍝ Member Of ∊
f644←{~(⍳(⍴⍺)+⍴⍵)∊⍺+⍳⍴⍺} ⍝ Expansion vector with zero after indices ⍺            ⍝ ⍵←A1; ⍺←I1
f645←{(~(⍳⍺)∊⍵)}         ⍝ Boolean vector of length ⍺ with zeroes in locations ⍵ ⍝ ⍵←I; ⍺←I0
f646←{(⍳⍴⍵)∊⍺}           ⍝ Starting points for ⍵ in indices pointed by ⍺         ⍝ ⍵←A1; ⍺←I1
f647←{(⍳⍺)∊⍵}            ⍝ Boolean vector of length ⍺ with ones in locations ⍵   ⍝ ⍵←I; ⍺←I0
f648←{(⍺←⎕)∊⍳⍵}          ⍝ Check for input in range 1..⍵                         ⍝ ⍵←A
f649←{~0∊⍵=⍺}            ⍝ Test if arrays are identical                          ⍝ ⍵←A; ⍺←A
f650←{⍺×~⍺∊⍵}            ⍝ Zeroing elements of ⍺ depending on their values       ⍝ ⍺←D; ⍵←D
f651←{1∊⍴,⍵}             ⍝ Test if single or scalar                              ⍝ ⍵←A
f652←{1∊⍴⍴⍵}             ⍝ Test if vector                                        ⍝ ⍵←A
f653←{0∊⍴⍵}              ⍝ Test if ⍵ is an empty array                           ⍝ ⍵←A

⍝⍝⍝ Index Generator ⍳
f654←{A←⍳⍴⍵ ⋄ A[⍵]←A ⋄ A} ⍝ Inverting a permutation                                   ⍝ ⍵←I1
f655←{⍳⍴⍴⍵}               ⍝ All axes of array ⍵                                       ⍝ ⍵←A
f656←{⍳⍴⍵}                ⍝ All indices of vector ⍵                                   ⍝ ⍵←A1
f657←{⍵+G×(⍳⍺)-⎕IO}       ⍝ Arithmetic progression of ⍺ numbers from ⍵ with step G    ⍝ ⍵←D0; ⍺←D0; G←D0
f658←{(⍵-⎕IO)+⍳1+⍺-⍵}     ⍝ Consecutive integers from ⍵ to ⍺ (arithmetic progression) ⍝ ⍵←I0; ⍺←I0
f659←{⍳0}                 ⍝ Empty numeric vector                                      ⍝ 
f660←{⍳1}                 ⍝ Index origin (⎕IO) as a vector                            ⍝ 

⍝⍝⍝ Logical Functions ~ ∨ ^ ⍱ ⍲
f661←{0∨⍵}               ⍝ Demote non-boolean representations to booleans             ⍝ ⍵←B
f662←{(⍺[1]<⍵)^⍵<⍺[2]}   ⍝ Test if ⍵ is within range ( ⍺[1],⍺[2] )                    ⍝ ⍵←D; ⍺←D1
f663←{(⍺[1]≤⍵)^(⍵≤⍺[2])} ⍝ Test if ⍵ is within range [ ⍺[1],⍺[2] ]                    ⍝ ⍵←D; ⍺←D1; 2=⍴⍺
f664←{0^⍵}               ⍝ Zeroing all boolean values                                 ⍝ ⍵←B
f666←{(⍵×G)+⍺×~G}        ⍝ Selection of elements of ⍵ and ⍺ depending on condition G  ⍝ ⍵←D; ⍺←D; G←B
f667←{(~⎕IO)+⍵}          ⍝ Changing an index origin dependent result to be as `⎕IO=1` ⍝ ⍵←I
f668←{⍺*~⍵}              ⍝ Conditional change of elements of ⍺ to one according to ⍵  ⍝ ⍺←D; ⍵←B

⍝⍝⍝ Comparison  ≠
f669←{⍵≤⍺}           ⍝ ⍵ implies ⍺                                             ⍝ ⍵←B; ⍺←B
f670←{⍵>⍺}           ⍝ ⍵ but not ⍺                                             ⍝ ⍵←B; ⍺←B
f671←{(0≠⍵)×⍺÷⍵+0=⍵} ⍝ Avoiding division by zero error (gets value zero)       ⍝ ⍵←D; ⍺←D
f672←{⍵≠⍺}           ⍝ Exclusive or                                            ⍝ ⍵←B; ⍺←B
f673←{⍵+⍺×⍵=0}       ⍝ Replacing zeroes with corresponding elements of ⍺       ⍝ ⍵←D; ⍺←D
f674←{⍺=⍵}           ⍝ Kronecker delta of ⍵ and ⍺ (element of identity matrix) ⍝ ⍵←I; ⍺←I

⍝⍝⍝ Ravel ,
f675←{,⍵,((⍴⍵),⍺)⍴G}                ⍝ Catenating ⍺ elements G after every element of ⍵  ⍝ ⍵←A1; ⍺←I0; G←A
f676←{,(((⍴⍵),⍺)⍴G),⍵}              ⍝ Catenating ⍺ elements G before every element of ⍵ ⍝ ⍵←A1; ⍺←I0; G←A0
f677←{,⍺,[⎕IO+.5]⍵}                 ⍝ Merging vectors ⍵ and ⍺ alternately               ⍝ ⍵←A1; ⍺←A1
f678←{,⍵,[1.1]⍺}                    ⍝ Inserting ⍺ after each element of ⍵               ⍝ ⍵←A1; ⍺←A0
f679←{,⍵,[1.1]' '}                  ⍝ Spacing out text                                  ⍝ ⍵←C1
f680←{(((⍴,⍵),1)×⍺*¯1 1)⍴⍵}         ⍝ Reshaping ⍵ into a matrix of width ⍺              ⍝ ⍵←D, ⍺←I0
f681←{A←⍴⍵ ⋄ ⍵←,⍵ ⋄ ⍵[G]←⍺ ⋄ ⍵←A⍴⍵} ⍝ Temporary ravel of ⍵ for indexing with G          ⍝ ⍵←A; ⍺←A; G←I
f682←{A←,⍵ ⋄ A[G]←⍺ ⋄ ⍵←(⍴⍵)⍴A}     ⍝ Temporary ravel of ⍵ for indexing with G          ⍝ ⍵←A; ⍺←A; G←I
f683←{⍵[;,1]}                       ⍝ First column as a matrix                          ⍝ ⍵←A2
f684←{⍴,⍵}                          ⍝ Number of elements (also of a scalar)             ⍝ ⍵←A

⍝⍝⍝ Catenate ,
f685←{⍵,⎕TC[2],⍺}             ⍝ Separating variable length lines                             ⍝ ⍵←A1; ⍺←A1
f686←{(⍵,⍵)⍴1,⍵⍴0}            ⍝ ⍵×⍵ identity matrix                                          ⍝ ⍵←I0
f687←{⍵,[.5+⍴⍴⍵]-⍵}           ⍝ Array and its negative ('plus minus')                        ⍝ ⍵←D
f688←{⍵,[⎕IO-.1]'¯'}          ⍝ Underlining a string                                         ⍝ ⍵←C1
f689←{⍵,[1.1]⍺}               ⍝ Forming a two-column matrix                                  ⍝ ⍵←A1; ⍺←A1
f690←{⍵,[.1]⍺}                ⍝ Forming a two-row matrix                                     ⍝ ⍵←A1; ⍺←A1
f691←{(⍵,⍺)[⎕IO+G]}           ⍝ Selection of ⍵ or ⍺ depending on condition G                 ⍝ ⍵←A0; ⍺←A0; G←B0
f692←{((((⍴⍴⍵)-⍴⍴⍺)⍴1),⍴⍺)⍴⍺} ⍝ Increasing rank of ⍺ to rank of ⍵                            ⍝ ⍵←A; ⍺←A
f693←{(⍴⍵)⍴1,0×⍵}             ⍝ Identity matrix of shape of matrix ⍵                         ⍝ ⍵←D2
f694←{((0.5×⍴⍵),2)⍴⍵}         ⍝ Reshaping vector ⍵ into a two-column matrix                  ⍝ ⍵←A1
f696←{(1,⍴⍵)⍴⍵}               ⍝ Reshaping vector ⍵ into a one-row matrix                     ⍝ ⍵←A1
f697←{((⍴⍵),1)⍴⍵}             ⍝ Reshaping vector ⍵ into a one-column matrix                  ⍝ ⍵←A1
f698←{(⍺,⍴⍵)⍴⍵}               ⍝ Forming a ⍺-row matrix with all rows alike (⍵)               ⍝ ⍵←A1; ⍺←I0
f699←{(⍴⍵)⍴ ... ,⍵}           ⍝ Handling array ⍵ temporarily as a vector                     ⍝ ⍵←A
f700←{⍺,0⍴⍵}                  ⍝ Joining sentences                                            ⍝ ⍵←A; ⍺←A1
f701←{⍵←0 2 1 2 5 8 0 4 5,⎕}  ⍝ Entering from terminal data exceeding input (printing) width ⍝ ⍵←D

⍝⍝⍝ Indexing [ ]
f702←{⍺[3]+⍵×⍺[2]+⍵×⍺[1]} ⍝ Value of fixed-degree polynomial ⍺ at points ⍵   ⍝ ⍺←D1; ⍵←D
f703←{(⍴⍵)[⍴⍴⍵]}          ⍝ Number of columns in array ⍵                     ⍝ ⍵←A
f704←{(⍴⍵)[1]}            ⍝ Number of rows in matrix ⍵                       ⍝ ⍵←A2
f705←{(⍴⍵)[2]}            ⍝ Number of columns in matrix ⍵                    ⍝ ⍵←A2
f706←{⍺×(1 ¯1)[1+⍵]}      ⍝ Conditional elementwise change of sign           ⍝ ⍺←D; ⍵←B
f707←{⍵[2×⎕IO]}           ⍝ Selection depending on index origin              ⍝ ⍵←A1
f708←{' *'[⎕IO+⍵]}        ⍝ Indexing with boolean value ⍵ (plotting a curve) ⍝ ⍵←B
f709←{⍵[⎕IO+⍺]}           ⍝ Indexing independent of index origin             ⍝ ⍵←A1; ⍺←I
f710←{⍵[1]}               ⍝ Selection depending on index origin              ⍝ ⍵←A1
f711←{⍵[]←0}              ⍝ Zeroing a vector (without change of size)        ⍝ ⍵←D1
f712←{⍵[;1]}              ⍝ First column as a vector                         ⍝ ⍵←A2

⍝⍝⍝ Shape ⍴
f713←{⍴⍴⍵}      ⍝ Rank of array ⍵                         ⍝ ⍵←A
f715←{(⍺×⍴⍵)⍴⍵} ⍝ Duplicating vector ⍵ ⍺ times            ⍝ ⍵←A1; ⍺←I0
f716←{⍺+(⍴⍺)⍴⍵} ⍝ Adding ⍵ to each row of ⍺               ⍝ ⍵←D1; ⍺←D; (⍴⍵)=¯1↑⍴⍺
f717←{(⍴⍺)⍴⍵}   ⍝ Array with shape of ⍺ and ⍵ as its rows ⍝ ⍵←A1; ⍺←A
f718←{1⍴⍴⍵}     ⍝ Number of rows in matrix ⍵              ⍝ ⍵←A2

⍝⍝⍝ Reshape ⍴
f720←{0 80⍴0} ⍝ Forming an initially empty array to be expanded ⍝ 
f721←{0⍴⍵←}   ⍝ Output of an empty line                         ⍝ ⍵←A
f722←{''⍴⍵}   ⍝ Reshaping first element of ⍵ into a scalar      ⍝ ⍵←A
f723←{1⍴⍵}    ⍝ Corner element of a (non-empty) array           ⍝ ⍵←A

⍝⍝⍝ Arithmetic + - × ÷
f724←{1+÷2+÷3+÷4+÷5+÷6}       ⍝ Continued fraction                                          ⍝ 
f725←{⍺×÷⍵}                   ⍝ Force 0÷0 into DOMAIN ERROR in division                     ⍝ ⍵←D; ⍺←D
f726←{⍵×¯1*⍺}                 ⍝ Conditional elementwise change of sign                      ⍝ ⍵←D; ⍺←B; ⍴⍵ ←→ ⍴⍺
f727←{0×⍵}                    ⍝ Zero array of shape and size of ⍵                           ⍝ ⍵←D
f728←{⍺×⍵}                    ⍝ Selecting elements satisfying condition ⍺, zeroing others   ⍝ ⍵←D; ⍺←B
f729←{1 ¯1×⍵}                 ⍝ Number and its negative ('plus minus')                      ⍝ ⍵←D0
f730←{-⎕IO-⍵}                 ⍝ Changing an index origin dependent result to be as ⎕IO=0    ⍝ ⍵←I
f731←{(⎕IO-1)+⍵}              ⍝ Changing an index origin dependent argument to act as ⎕IO=1 ⍝ ⍵←I
⍝f732←{+⍵←}                   ⍝ Output of assigned numeric value                            ⍝ ⍵←D
f733←{⎕IO+⍵}                  ⍝ Changing an index origin dependent argument to act as ⎕IO=0 ⍝ ⍵←I
f734←{⍵*⍺}                    ⍝ Selecting elements satisfying condition ⍺, others to one    ⍝ ⍵←D; ⍺←B

⍝⍝⍝ Miscellaneous
⍝f736←{⎕L⍵←⍞}                       ⍝ Setting a constant with hyphens ⍝ 
⍝f737←{⎕←⍵←}                         ⍝ Output of assigned value        ⍝ ⍵←A
⍝f738←{*}                            ⍝ Syntax error to stop execution  ⍝ 
⍝f888←{⍎⊖⍕⊃⊂|⌊-*+○⌈×÷!⌽⍉⌹~⍴⍋⍒,⍟?⍳0} ⍝ Meaning of life                 ⍝ 
