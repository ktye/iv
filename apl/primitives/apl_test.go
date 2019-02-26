package primitives

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/big"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	aplstrings "github.com/ktye/iv/apl/strings"
	"github.com/ktye/iv/apl/xgo"
)

//go:generate go run gen.go

var testCases = []struct {
	in, exp string
	flag    int
}{
	// {"T←⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)⋄T,(+⌿÷≢)T", "a b c\n1 4 7\n2 5 8\n3 6 9\n2 5 8", 0}, // table with avg value.

	{"⍝ Basic numbers and arithmetics", "", 0},
	{"1", "1", 0},
	{"1b", "1", 0},
	{"1+1", "2", 0},
	{"1-2", "¯1", 0}, // negative number
	{"¯1", "¯1", 0},
	{"1-¯2", "3", 0},
	{"1a90", "0J1", float}, // a complex number
	{"1a60+1a300", "1J0", float},
	{"1J1", "1J1", float},

	{"⍝ Vectors", "", 0},
	{"1 2 3", "1 2 3", 0},
	{"1+1 2 3", "2 3 4", 0},
	{"1 2 3+¯1", "0 1 2", 0},
	{"1 2 3+4 5 6", "5 7 9", 0},

	{"⍝ Braces", "apl/parse.go", 0},
	{"1 2+3 4", "4 6", 0},
	{"(1 2)+3 4", "4 6", 0},
	{"1×2+3×4", "14", 0},
	{"1×(2+3)×4", "20", 0},
	{"(3×2)+3×4", "18", 0},
	{"3×2+3×4", "42", 0},
	// {"1 (2+3) 4", "1 5 4", 0}, not supported
	// {"1 2 (+/1 2 3) 4 5", "1 2 6 4 5", 0},

	{"⍝ Comparison", "apl/primitives/compare.go", 0},
	{"1 2 3 4 5 > 2", "0 0 1 1 1", 0},         // greater than
	{"1 2 3 4 5 ≥ 3", "0 0 1 1 1", 0},         // greater or equal
	{"2 4 6 8 10<6", "1 1 0 0 0", 0},          // less than
	{"2 4 6 8 10≤6", "1 1 1 0 0", 0},          // less or equal
	{"1 2 3 ≠ 1.1 2 3", "1 0 0", 0},           // not equal
	{"3=3.1 3 ¯2 ¯3 3J0", "0 1 0 0 1", float}, // equal
	{"2+2=2", "3", 0},                         // calculating with boolean values
	{"2×1 2 3=4 2 1", "0 2 0", 0},             // dyadic array
	{"-3<4", "¯1", 0},                         // monadic scalar
	{"-1 2 3=0 2 3", "0 ¯1 ¯1", 0},            // monadic array
	{"⍝ TODO Comparison tolerance is not implemented.", "", 0},

	{"⍝ Boolean, logical", "apl/primitives/boolean.go", 0},
	{"0 1 0 1 ^ 0 0 1 1", "0 0 0 1", 0}, // and
	{"0 1 0 1 ∧ 0 0 1 1", "0 0 0 1", 0}, // accept both ^ and ∧
	{"0^0 0 1 1", "0 0 0 0", 0},         // or
	{"0 0 1 1∨0 1 0 1", "0 1 1 1", 0},   // or
	{"1∨0 1 0 1", "1 1 1 1", 0},         // or
	{"0 0 1 1⍱0 1 0 1", "1 0 0 0", 0},   // nor
	{"0 0 1 1⍲0 1 0 1", "1 1 1 0", 0},   // nand
	{"~0", "1", 0},                      // scalar not
	{"~1.0", "0", 0},                    // scalar not
	{"~0 1", "1 0", 0},                  // array not

	{"⍝ Least common multiple, greatest common divisor", "apl/primitives/boolean.go", 0},
	{"30^36", "180", small},                     // lcm
	{"0^3", "0", 0},                             // lcm with 0
	{"3^0", "0", 0},                             // lcm with 0
	{"15 1 2 7 ^ 35 1 4 0", "105 1 4 0", small}, // least common multiple
	{"30∨36", "6", small},                       // gcm
	{"15 1 2 7 ∨ 35 1 4 0", "5 1 2 7", small},   // greatest common divisor
	{"0∨3", "3", 0},                             // gcm with 0
	{"3∨0", "3", 0},                             // gcm with 0
	{"3^3.6", "18", small},                      // lcm
	//{"¯29J53^¯1J107", "¯853J¯329", 0},          // lcm
	//{"2 3 4 ∧ 0j1 1j2 2j3", "0J2 3J6 8J12", 0}, // least common multiple
	//{"2j2 2j4 ∧ 5j5 4j4", "10J10 ¯4J12", 0},    // least common multiple
	{"3∨3.6", "0.6", small}, // gcm
	//{"¯29J53∨¯1J107", "7J1", 0},                // gcm
	{"⍝ TODO: lcm and gcm of float and complex", "", 0},

	{"⍝ Multiple expressions", "apl/parse.go", 0},
	{"1⋄2⋄3", "1\n2\n3", 0},
	{"1⋄2", "1\n2", 0},
	{"1 2⋄3 4", "1 2\n3 4", 0},
	{"X←3 ⋄ Y←4", "", 0},

	{"⍝ Index origin, print precision", "apl/var.go", 0},
	{"⎕IO←0 ⋄ ⍳3", "0 1 2", 0},
	{"⎕IO", "1", 0},
	{"⎕IO←0 ⋄ ⎕IO", "0", 0},
	{"⎕PP←1 ⋄ ⎕PP", "1", 0},
	{"⎕PP←0 ⋄ 1.23456789", "1.23457", small},
	{"⎕PP←¯1 ⋄ 1.23456789", "1.23456789", small},
	{"⎕PP←1 ⋄ 1.23456789", "1", small},
	{"⎕PP←3 ⋄ 1.23456789", "1.23", small},

	{"⍝ Type, typeof", "apl/primitives/type.go", 0},
	{"⌶'a'", "apl.String", 0},

	{"⍝ Bracket indexing", "apl/primitives/index.go", 0},
	{"A←⍳6 ⋄ A[1]", "1", 0},
	{"A←2 3⍴⍳6 ⋄ A[1;] ⋄ ⍴A[1;]", "1 2 3\n3", 0},
	{"A←2 3⍴⍳6 ⋄ A[2;3]", "6", 0},
	{"A←2 3⍴⍳6 ⋄ A[2;2 3]", "5 6", 0},
	{"A←2 3⍴⍳6 ⋄ ⍴⍴A[2;3]", "0", 0},
	{"A←2 3 4 ⋄ A[]", "2 3 4", 0},
	{"⎕IO←0 ⋄ A←2 3⍴⍳6 ⋄ A[1;2]", "5", 0},
	{"5 6 7[2+1]", "7", 0},
	{"(2×⍳3)[2]", "4", 0},
	{"A←2 3 ⍴⍳6⋄A[A[1;1]+1;]", "4 5 6", 0},
	{"A←1 2 3⋄A[3]+1", "4", 0},
	{"A←1 2 3⋄1+A[3]", "4", 0},

	{"⍝ Scalar primitives with axis", "apl/primitives/array.go", 0},
	{"(2 3⍴⍳6)+[2]1 2 3", "2 4 6\n5 7 9", 0},
	{"1 2 3 +[2] 2 3⍴⍳6", "2 4 6\n5 7 9", 0},
	{"K←2 3⍴.1×⍳6⋄J←2 3 4⍴⍳24⋄N←J+[1 2]K⋄⍴N⋄N[1;2;3]⋄N[2;3;4]", "2 3 4\n7.2\n24.6", small},

	{"⍝ Iota", "apl/primitives/iota.go", 0},
	{"⍳5", "1 2 3 4 5", 0}, // index generation
	{"⍳0", "", 0},          // empty array

	{"⍝ Rho, reshape", "apl/primitives/rho.go", 0},
	{"⍴⍳5", "5", 0},              // shape
	{"⍴5", "", 0},                // shape of scalar is empty
	{"⍴⍴5", "0", 0},              // shape of empty is 0
	{"⍴⍳0", "0", 0},              // empty array has zero dimensions
	{"⍴⍴⍳0", "1", 0},             // rank of empty array is 1
	{"2 3⍴1", "1 1 1\n1 1 1", 0}, // shape
	{"3⍴⍳0", "0 0 0", 0},         // reshape empty array
	{"⍴0 2⍴⍳0", "0 2", 0},        // reshape empty array
	{"⍴3 0⍴⍳0", "3 0", 0},        // reshape empty array
	{"⍴3 0⍴3", "3 0", 0},         // reshape empty array
	{"⍳'a'", "fail: strings are not in the input domain of ⍳", 0},

	{"⍝ Where, interval index", "apl/primitives/iota.go", 0},
	{"⍸1 0 1 0 0 0 0 1 0", "1 3 8", 0},
	{"⍸'e'='Pete'", "2 4", 0},
	{"⍸1=1", "1", 0},
	{"10 20 30⍸11 1 31 21", "1 0 3 2", 0},
	{"'AEIOU'⍸'DYALOG'", "1 5 1 3 4 2", 0},
	{"0.8 2 3.3⍸1.3 1.9 0.7 4 .6 3.2", "1 1 0 3 0 2", 0},

	{"⍝ Membership", "apl/primitives/iota.go", 0},
	{"'BANANA'∊'AN'", "0 1 1 1 1 1", 0},
	{"5 1 2∊6 5 4 1 9", "1 1 0", 0},
	{"(2 3⍴8 3 5 8 4 8)∊1 8 9 3", "1 1 0\n1 0 1", 0},
	{"8 9 7 3∊⍳0", "0 0 0 0", 0},
	{"3.1 5.1 7.1∊2 2⍴1.1 3.1 5.1 4.1", "1 1 0", 0},
	{"19∊'CLUB'", "0", 0},
	{"'BE'∊'BOF'", "1 0", 0},
	{"'NADA'∊⍳0", "0 0 0 0", 0},
	{"(⌈/⍳0)∊⌊/⍳0", "0", 0},
	{"5 10 15∊⍳10", "1 1 0", 0},

	{"⍝ Without", "apl/primitives/boolean.go", 0},
	{"1 2 3 4 5~2 3 4", "1 5", 0},
	{"'RHYME'~'MYTH'", "R E", 0},
	{"1 2~⍳0", "1 2", 0},
	{"1~3", "1", 0},
	{"3~3", "", 0},
	{"⍴⍳0~1 2", "0", 0},
	{"5 10 15~⍳10", "15", 0},
	{"3 1 4 1 5 5~3 1 4 1 5 5~4 2 5 2 6", "4 5 5", 0}, // intersection

	{"⍝ Unique, union", "apl/primitives/unique.go", 0},
	{"∪3", "3", 0},
	{"⍴∪3", "1", 0},
	{"∪ 22 10 22 22 21 10 5 10", "22 10 21 5", 0},
	{"∪2 7 1 8 2 8 1 8 2 8 4 5 9 0 4 4 9", "2 7 1 8 4 5 9 0", 0},
	{"∪'MISSISSIPPI'", "M I S P", 0},
	{"⍴∪⍳0", "0", 0},
	{"∪⍳0", "", 0},
	{"3∪3", "3", 0},
	{"⍴3∪3", "1", 0},
	{"3∪⍳0", "3", 0},
	{"(⍳0)∪3", "3", 0},
	{"⍴(⍳0)∪⍳0", "0", 0},
	{"1 2 3∪5 3 2 1 4", "1 2 3 5 4", 0},
	{"5 6 7∪1 2 3", "5 6 7 1 2 3", 0},

	{"⍝ Find", "apl/primitives/find.go", 0},
	{"'AN'⍷'BANANA'", "0 1 0 1 0 0", 0},
	{"'ANA'⍷'BANANA'", "0 1 0 1 0 0", 0},
	{"(2 2⍴1)⍷1 2 3", "0 0 0", 0},
	{"(2 2⍴5 6 8 9)⍷3 3⍴⍳9", "0 0 0\n0 1 0\n0 0 0", 0},
	{"4 5 6⍷3 3⍴⍳9", "0 0 0\n1 0 0\n0 0 0", 0},

	{"⍝ Magnitude, Residue, Ceil, Floor, Min, Max", "apl/primitives/elementary.go", 0},
	{"|1 ¯2 ¯3.2 2.2a20", "1 2 3.2 2.2", float},                  // magnitude
	{"3 3 ¯3 ¯3|¯5 5 ¯4 4", "1 2 ¯1 ¯2", 0},                      // residue
	{"0.5|3.12 ¯1 ¯0.6", "0.12 0 0.4", small},                    // residue
	{"¯1 0 1|¯5.25 0 2.41", "¯0.25 0 0.41", small},               // residue
	{"1j2|2j3 3j4 5j6", "1J1 ¯1J1 0J1", float},                   // complex residue
	{"4J6|7J10", "3J4", float},                                   // complex residue
	{"¯10 7J10 .3|17 5 10", "¯3 ¯5J7 0.1", float},                // residue
	{"⌊¯2.3 0.1 100 3.3", "¯3 0 100 3", 0},                       // floor
	{"⌊0.5 + 0.4 0.5 0.6", "0 1 1", 0},                           // floor
	{"⌊1j3.2 3.3j2.5 ¯3.3j¯2.5", "1J3 3J2 ¯3J¯3", float},         // complex floor
	{"⌊1.5J2.5", "2J2", float},                                   // complex floor
	{"⌊1J2 1.2J2.5 ¯1.2J¯2.5", "1J2 1J2 ¯1J¯3", float},           // complex floor
	{"⌈¯2.7 3 .5", "¯2 3 1", 0},                                  // ceil
	{"⌈1.5J2.5", "1J3", float},                                   // complex ceil
	{"⌈1J2 1.2J2.5 ¯1.2J¯2.5", "1J2 1J3 ¯1J¯2", float},           // complex ceil
	{"⌈¯2.3 0.1 100 3.3", "¯2 1 100 4", 0},                       // ceil
	{"⌈1.2j2.5 1.2j¯2.5", "1J3 1J¯2", float},                     // ceil
	{"5⌊4 5 7", "4 5 5", 0},                                      // min
	{"¯2⌊¯3", "¯3", 0},                                           // min
	{"3.3 0 ¯6.7⌊3.1 ¯4 ¯5", "3.1 ¯4 ¯6.7", small},               // min
	{"¯2.1 0.1 15.3 ⌊ ¯3.2 1 22", "¯3.2 0.1 15.3", small},        // min
	{"5⌈4 5 7", "5 5 7", 0},                                      // max
	{"¯2⌈¯3", "¯2", 0},                                           // max
	{"3.3 0 ¯6.7⌈3.1 ¯4 ¯5", "3.3 0 ¯5", small},                  // max
	{"¯2.01 0.1 15.3 ⌈ ¯3.2 ¯1.1 22.7", "¯2.01 0.1 22.7", small}, // max

	{"⍝ Factorial, gamma, binomial", "apl/primitives/elementary.go", 0},
	{"!4", "24", small},                                   // factorial
	{"!1 2 3 4 5", "1 2 6 24 120", small},                 // factorial
	{"!3J2", "¯3.01154J1.77017", small},                   // complex gamma
	{"!.5 ¯.05", "0.886227 1.03145", small},               // real gamma (APL2 doc: "0.0735042656 1.031453317"?)
	{"2!5", "10", small},                                  // binomial
	{"3.2!5.2", "10.92", small},                           // binomial, floats with beta function
	{"3!¯2", "¯4", small},                                 // binomial, negative R
	{"¯6!¯3", "¯10", small},                               // binomial negative L and R
	{"2 3 4!6 18 24", "15 816 10626", small},              // binomial
	{"3!.05 2.5 ¯3.6", "0.0154375 0.3125 ¯15.456", small}, // binomial
	{"0 1 2 3!3", "1 3 3 1", small},                       // binomial coefficients
	{"2!3J2", "1J5", small},                               // binomial complex

	{"⍝ Match, Not match, tally, depth", "apl/primitives/match.go", 0},
	{"≡5", "0", 0},          // depth
	{"≡⍳0", "1", 0},         // depth for empty array
	{`≡"alpha"`, "0", 0},    // a string is a scalarin APLv.
	{"≢2 3 4⍴⍳10", "2", 0},  // tally
	{"≢2", "1", 0},          // tally
	{"≢⍳0", "0", 0},         // tally
	{"1 2 3≡1 2 3", "1", 0}, // match
	{"3≡1⍴3", "0", 0},       // match shape
	{`""≡⍳0`, "0", 0},       // match empty string
	{`''≡⍳0`, "1", 0},       // this is false in other APLs (here '' is an empty array).
	{"2.0-1.0≡1>0", "1", 0}, // compare numbers of different type
	{"1≢2", "1", 0},         // not match
	{"1≢1", "0", 0},         // not match
	{"3≢1⍴3", "1", 0},       // not match
	{`""≢⍳0`, "1", 0},       // not match

	{"⍝ Left tack, right tack", "apl/primitives/tack.go", 0},
	{"⊣1 2 3", "1 2 3", 0},      // monadic left: same
	{"3 2 1⊣1 2 3", "3 2 1", 0}, // dyadic left
	{"1 2 3⊢3 2 1", "3 2 1", 0}, // dyadic right
	{"⊢4", "4", 0},              // monadic right: same
	{"⊣/1 2 3", "1", 0},         // ⊣ reduction selects the first sub array
	{"⊢/1 2 3", "3", 0},         // ⊢ reduction selects the last sub array
	{"⊣/2 3⍴⍳6", "1 4", 0},      // ⊣ reduction over array
	{"⊢/2 3⍴⍳6", "3 6", 0},      // ⊢ reduction over array

	{"⍝ Array expressions", "apl/primitives/array.go", 0},
	{"-⍳3", "¯1 ¯2 ¯3", 0},

	{"⍝ Ravel, enlist, catenate, join", "apl/primitives/comma.go", 0},
	{",2 3⍴⍳6", "1 2 3 4 5 6", 0},     // ravel
	{"⍴,3", "1", 0},                   // scalar ravel
	{"⍴,⍳0", "0", 0},                  // ravel empty array
	{"1 2 3,4 5 6", "1 2 3 4 5 6", 0}, // catenate
	{`"abc",1 2`, `abc 1 2`, 0},
	{"(2 3⍴⍳6),2 2⍴7 8 9 10", "1 2 3 7 8\n4 5 6 9 10", 0},
	{"2 3≡2,3", "1", 0},                       // catenate vector result
	{"(1 2 3,4 5 6)≡⍳6", "1", 0},              // catenate vector result
	{"0,2 3⍴1", "0 1 1 1\n0 1 1 1", 0},        // catenate scalar and array
	{"0,[1]2 3⍴⍳6", "0 0 0\n1 2 3\n4 5 6", 0}, // catenate with axis
	{"(2 3⍴⍳6),[1]0", "1 2 3\n4 5 6\n0 0 0", 0},
	{"(2 3⍴⍳6),[1]5 4 3", "1 2 3\n4 5 6\n5 4 3", 0},
	{"⍴(3 5⍴⍳15),[1]3 3 5⍴-⍳45", "4 3 5", 0},
	{"⍴(3 5⍴⍳15),[2]3 3 5⍴-⍳45", "3 4 5", 0},

	{"⍝ Ravel with axis", "apl/primitives/comma.go", 0},
	{",[0.5]1 2 3", "1 2 3", 0},
	{"⍴,[0.5]1 2 3", "1 3", 0},
	{",[1.5]1 2 3", "1\n2\n3", 0},
	{"⍴,[1.5]1 2 3", "3 1", 0},
	{"A←3 4⍴⍳12⋄⍴,[0.5]A", "1 3 4", 0},
	{"A←3 4⍴⍳12⋄⍴,[1.5]A", "3 1 4", 0},
	{"A←3 4⍴⍳12⋄⍴,[2.5]A", "3 4 1", 0},
	{"A←2 3⍴⍳6⋄⍴,[.1]A", "1 2 3", 0},
	{"A←2 3⍴⍳6⋄⍴,[1.1]A", "2 1 3", 0},
	{"A←2 3⍴⍳6⋄⍴,[2.1]A", "2 3 1", 0},
	{",[1.1]5 6 7", "5\n6\n7", 0},
	{"A←2 3 4⍴⍳24⋄A←,[1 2]A⋄⍴A⋄A[5;3]", "6 4\n19", 0},
	{"A←2 3 4⍴⍳24⋄⍴,[2 3]A", "2 12", 0},
	{"A←3 2 4⍴⍳24⋄⍴,[2 3]A", "3 8", 0},
	{"A←3 2 4⍴⍳24⋄⍴,[1 2]A", "6 4", 0},
	{"⍴,[⍳0]1 2 3", "3 1", 0},
	{"⍴,[⍳0]2 3⍴⍳6", "2 3 1", 0},
	{"A←3 2 5⍴⍳30⋄⍴,[⍳⍴⍴A],[.5]A", "6 5", 0}, // Turn array into matrix
	{"A←2 3 4⍴⍳24⋄(,[2 3]A)←2 12⍴-⍳24⋄⍴A⋄A[1;3;4]", "2 3 4\n¯12", 0},

	{"⍝ Laminate", "apl/primitives/comma.go", 0},
	{"1 2 3,[0.5]4", "1 2 3\n4 4 4", 0},
	{"1 2 3,[1.5]4", "1 4\n2 4\n3 4", 0},
	{"⎕IO←0⋄1 2 3,[¯0.5]4", "1 2 3\n4 4 4", 0},
	{"'FOR',[.5]'AXE'", "F O R\nA X E", 0},
	{"'FOR',[1.1]'AXE'", "F A\nO X\nR E", 0},

	{"⍝ Table, catenate first", "apl/primitives/comma.go", 0},
	{"⍪0", "0", 0},
	{"⍴⍪0", "1 1", 0},
	{"⍪⍳4", "1\n2\n3\n4", 0},
	{"⍪2 2⍴⍳4", "1 2\n3 4", 0},
	{"⍪2 2 2⍴⍳8", "1 2 3 4\n5 6 7 8", 0},
	{"10 20⍪2 2⍴⍳4", "10 20\n1 2\n3 4", 0},

	{"⍝ Decode", "apl/primitives/decode.go", 0},
	{"3⊥1 2 1", "16", 0},
	{"3⊥4 3 2 1", "142", 0},
	{"2⊥1 1 1 1", "15", 0},
	{"1 2 3⊥3 2 1", "25", 0},
	{"1J1⊥1 2 3 4", "5J9", float},
	{"24 60 60⊥2 23 12", "8592", 0},
	{"(2 1⍴2 10)⊥3 2⍴ 1 4 0 3 1 2", "5 24\n101 432", 0},

	{"⍝ Encode, representation", "apl/primitives/decode.go", 0},
	{"2 2 2 2⊤15", "1 1 1 1", 0},
	{"10⊤5 15 125", "5 5 5", 0},
	{"⍴10⊤5 15 125", "3", 0},
	{"⍴(1 1⍴10)⊤5 15 125", "1 1 3", 0},
	{"0 10⊤5 15 125", "0 1 12\n5 5 5", 0},
	{"0 1⊤1.25 10.5", "1 10\n0.25 0.5", small},
	{"24 60 60⊤8592", "2 23 12", 0},
	{"2 2 2 2 2⊤15", "0 1 1 1 1", 0},
	{"2 2 2⊤15", "1 1 1", 0},
	{"4 5 6⊤⍳0", "", 0},
	{"⍴4 5 6⊤⍳0", "3 0", 0},
	{"⍴(⍳0)⊤4 5 6", "0 3", 0},
	{"((⌊1+2⍟135)⍴2)⊤135", "1 0 0 0 0 1 1 1", float},
	{"24 60 60⊤162507", "21 8 27", 0},
	{"0 24 60 60⊤162507", "1 21 8 27", 0},
	{"10 10 10⊤215 345 7", "2 3 0\n1 4 0\n5 5 7", 0},
	{"(4 2⍴8 2)⊤15", "0 1\n0 1\n1 1\n7 1", 0},
	{"3 2J3⊤2", "0J2 ¯1J2", float},
	{"0 2J3⊤2", "0J¯1 ¯1J2", float},
	{"3 2J3⊤2", "0J2 ¯1J2", float},
	{"3 2J3⊤2 1", "0J2 0J2\n¯1J2 ¯2J2", float},
	{"10⊥2 2 2 2⊤15", "1111", 0},
	{"10 10 10⊤123", "1 2 3", 0},
	{"10 10 10⊤123 456", "1 4\n2 5\n3 6", 0},
	{"2 2 2⊤¯1", "1 1 1", 0},
	{"0 2 2⊤¯1", "¯1 1 1", 0},
	{"0 1⊤3.75 ¯3.75", "3 ¯4\n0.75 0.25", small},
	{"1 0⊤0", "0 0", 0},
	{"0⊤0", "0", 0},
	{"0⊤0 0", "0 0", 0},
	{"0 0⊤0", "0 0", 0},
	{"1 0⊤234", "0 234", 0},

	{"⍝ Reduce, reduce first, reduce with axis", "apl/operators/reduce.go", 0},
	{"+/1 2 3", "6", 0},
	{"+⌿1 2 3", "6", 0},
	{"+/2 3 1 ⍴⍳6", "1 2 3\n4 5 6", 0},
	{"⍴+/3", "", 0},
	{"⍴+/1 1⍴3", "1", 0},
	{"+/2 3⍴⍳6", "6 15", 0},
	{"+⌿2 3⍴⍳6", "5 7 9", 0},
	{"+/⍳0", "0", 0},
	{"+/1", "1", 0},
	{"+/1⍴1", "1", 0},
	{"-/1⍴1", "1", 0},
	{"+/[1]2 3⍴⍳6", "5 7 9", 0},
	{"+/[1]3 4⍴⍳12", "15 18 21 24", 0},
	{"+/[2]3 4⍴⍳12", "10 26 42", 0},
	{"×/[1]3 4 ⍴⍳12", "45 120 231 384", 0},
	{"÷/[2]2 1 4⍴2×⍳8", "2 4 6 8\n10 12 14 16", 0},
	{"÷/[2]2 0 3⍴0", "1 1 1\n1 1 1", 0},

	{"⍝ N-wise reduction", "apl/operators/reduce.go", 0},
	{"6+/⍳6", "21", 0},
	{"4+/⍳6", "10 14 18", 0},
	{"5+/⍳6", "15 20", 0},
	{"3+/⍳6", "6 9 12 15", 0},
	{"1+/⍳6", "1 2 3 4 5 6", 0},
	{"0+/⍳0", "0", 0},
	{"⍴0+/⍳0", "1", 0},
	{"1+/⍳0", "", 0},
	{"¯1+/⍳0", "", 0},
	{"⍴4+/2 3⍴⍳6", "2 0", 0},
	{"2+/3 4⍴⍳12", "3 5 7\n11 13 15\n19 21 23", 0},
	{"¯2-/1 4 9 16 25", "3 5 7 9", 0},
	{"2-/1 4 9 16 25", "¯3 ¯5 ¯7 ¯9", 0},
	{"3×/⍳6", "6 24 60 120", 0},
	{"¯3×/⍳6", "6 24 60 120", 0},
	{"0×/⍳5", "1 1 1 1 1 1", 0},
	{"4+/[1]4 3⍴⍳12", "22 26 30", 0},
	{"3+/[1]4 3⍴⍳12", "12 15 18\n21 24 27", 0},
	{"2+/[1]4 3⍴⍳12", "5 7 9\n11 13 15\n17 19 21", 0},
	{"0×/[1]2 3⍴⍳12", "1 1 1\n1 1 1\n1 1 1", 0},
	{"1+/⍳6", "1 2 3 4 5 6", 0},
	{`+/1000+/⍳10000`, "45009500500", small},

	{"⍝ Scan, scan first, scan with axis", "apl/operators/reduce.go", 0},
	{`+\1 2 3 4 5`, "1 3 6 10 15", 0},
	{`+\2 3⍴⍳6`, "1 3 6\n4 9 15", 0},
	{`+⍀2 3⍴⍳6`, "1 2 3\n5 7 9", 0},
	{`-\1 2 3`, "1 ¯1 2", 0},
	{"∨/0 0 1 0 0 1 0", "1", 0},
	{`^\1 1 1 0 1 1 1`, "1 1 1 0 0 0 0", 0},
	{`+\1 2 3 4 5`, "1 3 6 10 15", 0},
	{`+\[1]2 3⍴⍳6`, "1 2 3\n5 7 9", 0},

	{"⍝ Replicate, compress", "apl/operators/reduce.go", 0},
	{"1 1 0 0 1/'STRAY'", "S T Y", 0},
	{"1 0 1 0/3 4⍴⍳12", "1 3\n5 7\n9 11", 0},
	{"1 0 1/1 2 3", "1 3", 0},
	{"1/1 2 3", "1 2 3", 0},
	{"3 2 1/1 2 3", "1 1 1 2 2 3", 0},
	{"1 0 1/2", "2 2", 0},
	{"⍴1/1", "1", 0},
	{"⍴⍴(,1)/2", "1", 0},
	{"3 4/1 2", "1 1 1 2 2 2 2", 0},
	{"1 0 1 0 1/⍳5", "1 3 5", 0},
	{"1 ¯2 3 ¯4 5/⍳5", "1 0 0 3 3 3 0 0 0 0 5 5 5 5 5", 0},
	{"2 0 1/2 3⍴⍳6", "1 1 3\n4 4 6", 0},
	{"0 1⌿2 3⍴⍳6", "4 5 6", 0},
	{"0 1⌿⍴⍳6", "6", 0},
	{"1 0 1/4", "4 4", 0},
	{"1 0 1/,3", "3 3", 0},
	{"1 0 1/1 1⍴5", "5 5", 0},
	{"1 2/[2]2 2 1⍴⍳4", "1\n2\n2\n\n3\n4\n4", 0},
	{"A←2 ¯1 1/[1]3 2 4⍴⍳24⋄⍴A⋄+/+/A", "4 2 4\n36 36 0 164", 0},
	{"⍴2/[2]3 2 4⍴⍳24", "3 4 4", 0},
	{"⍴¯1 1/[2]3 1 4⍴⍳12", "3 2 4", 0},
	{"⍴1 0 2 ¯1⌿[2]3 4⍴⍳12", "3 4", 0},
	{"0 1/[1]2 3⍴⍳6", "4 5 6", 0},
	{"B←2 2⍴'ABCD'⋄A←3 2⍴⍳6⋄(1 0 1/[1]A)←B⋄A", "A B\n3 4\nC D", 0},

	{"⍝ Expand, expand first", "apl/operators/reduce.go", 0},
	{`1 0 1 0 0 1\1 2 3`, "1 0 2 0 0 3", 0},
	{`1 0 0\5`, "5 0 0", 0},
	{`0 1 0\3 1⍴7 8 9`, "0 7 0\n0 8 0\n0 9 0", 0},
	{`1 0 0 1 0 1\7 8 9`, "7 0 0 8 0 9", 0},
	{`⍴(⍳0)\3`, "0", 0},
	{`⍴(⍳0)\2 0⍴3`, "2 0", 0},
	{`⍴1 0 1\0 2⍴0`, "0 3", 0},
	{`0 0 0\2 0⍴0`, "0 0 0\n0 0 0", 0},
	{`1 0 1⍀2 3⍴⍳6`, "1 2 3\n0 0 0\n4 5 6", 0},
	{`0\⍳0`, "0", 0},
	{`1 ¯2 3 ¯4 5\3`, "3 0 0 3 3 3 0 0 0 0 3 3 3 3 3", 0},
	{`1 0 1\1 3`, "1 0 3", 0},
	{`1 0 1\2`, "2 0 2", 0},
	{`1 0 1 1\1 2 3`, "1 0 2 3", 0},
	{`1 0 1 1⍀3`, "3 0 3 3", 0},
	{`0 1\3 1⍴3 2 4`, "0 3\n0 2\n0 4", 0},
	{`0 0\5`, "0 0", 0},
	{`1 0 1⍀2 3⍴⍳6`, "1 2 3\n0 0 0\n4 5 6", 0},
	{`1 0 1\3 2⍴⍳6`, "1 0 2\n3 0 4\n5 0 6", 0},
	{`1 0 1 1\2 3⍴⍳6`, "1 0 2 3\n4 0 5 6", 0},
	{`1 0 1\[1]2 3⍴⍳6`, "1 2 3\n0 0 0\n4 5 6", 0},
	{"⍝ TODO expand with selective specification", "", 0},

	{"⍝ Pi times, circular, trigonometric", "apl/primitives/elementary.go", 0},
	{"○0 1 2", "0 3.14159 6.28319", small},                  // pi times
	{"1E¯12>|1+*○0J1", "1", small},                          // Euler identity
	{"0 ¯1 ○ 1", "0 1.5708", small},                         //
	{"1○(○1)÷2 3 4", "1 0.866025 0.707107", small},          //
	{"2○(○1)÷3", "0.5", small},                              //
	{"9 11○3.5J¯1.2", "3.5 ¯1.2", small},                    //
	{"9 11∘.○3.5J¯1.2 2J3 3J4", "3.5 2 3\n¯1.2 3 4", small}, //
	{"¯4○¯1", "0", small},                                   //
	{"3○2", "¯2.18504", small},                              //
	{"2○1", "0.540302", small},                              //
	{"÷3○2", "¯0.457658", small},                            //
	{"1○○30÷180", "0.5", small},
	{"2○○45÷180", "0.707107", small},
	{"¯1○1", "1.5708", small},
	{"¯2○.54032023059", "0.999979", small},
	{"(¯1○.5)×180÷○1", "30", small},
	{"(¯3○1)×180÷○1", "45", small},
	{"5○1", "1.1752", small},
	{"6○1", "1.54308", small},
	{"¯5○1.175201194", "1", small},
	{"¯6○1.543080635", "1", small},

	{"⍝ Take, drop", "apl/primitives/take.go", 0}, // Monadic First and split are not implemented.
	{"5↑'ABCDEF'", "A B C D E", 0},
	{"5↑1 2 3", "1 2 3 0 0", 0},
	{"¯5↑1 2 3", "0 0 1 2 3", 0},
	{"2 3↑2 4⍴⍳8", "1 2 3\n5 6 7", 0},
	{"¯1 ¯2↑2 4⍴⍳8", "7 8", 0},
	{"1↑2", "2", 0},
	{"⍴1↑2", "1", 0},
	{"1 1 1↑2", "2", 0},
	{"⍴1 1 1↑2", "1 1 1", 0},
	{"(⍳0)↑2", "2", 0},
	{"⍴(⍳0)↑2", "", 0},
	{"2↑⍳0", "0 0", 0},
	{"2 3↑2", "2 0 0\n0 0 0", 0},
	{"4↓'OVERBOARD'", "B O A R D", 0},
	{"¯5↓'OVERBOARD'", "O V E R", 0},
	{"⍴10↓'OVERBOARD'", "0", 0},
	{"0 ¯2↓3 3⍴⍳9", "1\n4\n7", 0},
	{"¯2 ¯1↓3 3⍴⍳9", "1 2", 0},
	{"1↓3 3⍴⍳9", "4 5 6\n7 8 9", 0},
	{"1 1↓2 3 4⍴⍳24", "17 18 19 20\n21 22 23 24", 0},
	{"¯1 ¯1↓2 3 4⍴⍳24", "1 2 3 4\n5 6 7 8", 0},
	{"3↓12 31 45 10 57", "10 57", 0},
	{"¯3↓12 31 45 10 57", "12 31", 0},
	{"0 2↓3 5⍴⍳15", "3 4 5\n8 9 10\n13 14 15", 0},
	{"⍴3 1↓2 3⍴'ABCDEF'", "0 2", 0},
	{"⍴2 3↓2 3⍴'ABCDEF'", "0 0", 0},
	{"0↓4", "4", 0},
	{"⍴0↓4", "1", 0},
	{"0 0 0↓4", "4", 0},
	{"⍴0 0 0↓4", "1 1 1", 0},
	{"⍴1↓5", "0", 0},
	{"⍴0↓5", "1", 0},
	{"⍴1 2 3↓4", "0 0 0", 0},
	{"''↓5", "5", 0},
	{"⍴⍴''↓5", "0", 0},
	{"1↑2 3⍴⍳6", "1 2 3", 0},
	{"1↑[1]2 3⍴⍳6", "1 2 3", 0},
	{"1 3↑[1 2]2 3⍴⍳6", "1 2 3", 0},
	{"2↑[1]3 5⍴'GIANTSTORETRAIL'", "G I A N T\nS T O R E", 0},
	{"¯3↑[2]3 5⍴'GIANTSTORETRAIL'", "A N T\nO R E\nA I L", 0},
	{"3↑[1]2 3⍴⍳6", "1 2 3\n4 5 6\n0 0 0", 0},
	{"¯4↑[1]2 3⍴⍳6", "0 0 0\n0 0 0\n1 2 3\n4 5 6", 0},
	{"¯1 3↑[1 3]3 3 4⍴'HEROSHEDDIMESODABOARPARTLAMBTOTODAMP'", "L A M\nT O T\nD A M", 0},
	{"2↑[2]2 3 4⍴⍳24", "1 2 3 4\n5 6 7 8\n\n13 14 15 16\n17 18 19 20", 0},
	{"2↑[3]2 3 4⍴⍳24", "1 2\n5 6\n9 10\n\n13 14\n17 18\n21 22", 0},
	{"2 ¯2↑[3 2]2 3 4⍴⍳24", "5 6\n9 10\n\n17 18\n21 22", 0},
	{"2 ¯2↑[2 3]2 3 4⍴⍳24", "3 4\n7 8\n\n15 16\n19 20", 0},
	{"1↓[1]3 4⍴'FOLDBEATRODE'", "B E A T\nR O D E", 0},
	{"1↓[2]3 4⍴'FOLDBEATRODE'", "O L D\nE A T\nO D E", 0},
	{"A←3 4⍴'FOLDBEATRODE'⋄(1↓[1]A)≡1 0↓A", "1", 0},
	{"A←3 4⍴'FOLDBEATRODE'⋄(1↓[2]A)≡0 1↓A", "1", 0},
	{"A←3 2 4⍴⍳24⋄1 ¯1↓[2 3]A", "5 6 7\n\n13 14 15\n\n21 22 23", 0},
	{"A←3 2 4⍴⍳24⋄1 ¯1↓[3 2]A", "2 3 4\n\n10 11 12\n\n18 19 20", 0},
	{"A←2 3 4⍴⍳24⋄⍴1↓[2]A", "2 2 4", 0},
	{"A←2 3 4⍴⍳24⋄2↓[3]A", "3 4\n7 8\n11 12\n\n15 16\n19 20\n23 24", 0},
	{"A←2 3 4⍴⍳24⋄2 1↓[3 2]A", "7 8\n11 12\n\n19 20\n23 24", 0},

	{"⍝ Format as a string, Execute", "apl/primitives/format.go", 0},

	{"⍕10", "10", 0},                                  // format as string
	{"⍕10.1", "10.1", small},                          // format as string
	{"⍕123.45678901234", "123.457", small},            // format as string
	{"4⍕123.45678901234", "123.5", small},             // format with precision
	{"`%.3f@%.1f ⍕1J2", "2.236@63.4", small},          // format with string
	{"`%.3f ⍕¯1.23456", "¯1.235", small},              // format with string
	{"`-%.3f ⍕¯1.23456", "-1.235", small},             // format with string (normal minus sign)
	{`⍕"alpha"`, `alpha`, 0},                          // format with default stringer
	{`¯1⍕"alpha"`, `"alpha"`, 0},                      // format with text marshaler
	{`¯1⍕"al\npha"`, `"al\npha"`, 0},                  // format with text marshaler
	{"`csv ⍕2 3⍴⍳6", "1,2,3\n4,5,6", 0},               // format as csv
	{"`csv ⍕2 2⍴`a`b`c\"t`d", "a,b\n\"c\"\"t\",d", 0}, // format as csv
	{`⍎"1+1"`, "2", 0},                                // evaluate expression
	{"⍝ TODO: dyadic format with specification.", "", 0},
	{"⍝ TODO: dyadic execute with namespace.", "", 0},

	{"⍝ Grade up, grade down, sort", "apl/primitives/grade.go", 0},
	{"⍋23 11 13 31 12", "2 5 3 1 4", 0},                             // grade up
	{"⍋23 14 23 12 14", "4 2 5 1 3", 0},                             // identical subarrays
	{"⍋5 3⍴4 16 37 2 9 26 5 11 63 3 18 45 5 11 54", "2 4 1 5 3", 0}, // grade up rank 2
	{"⍋22.5 1 15 3 ¯4", "5 2 4 3 1", 0},                             // grade up
	{"⍒33 11 44 66 22", "4 3 1 5 2", 0},                             // grade down
	{"⍋'alpha'", "1 5 4 2 3", 0},                                    // strings grade up
	{"'ABCDE'⍒'BEAD'", "2 4 1 3", 0},                                // grade down with collating sequence
	{"⍝ TODO dyadic grade up/down is only implemented for vector L", "", 0},
	{"A←23 11 13 31 12⋄A[⍋A]", "11 12 13 23 31", 0}, // sort

	{"⍝ Reverse, revere first", "apl/primitives/reverse.go", 0},
	{"⌽1 2 3 4 5", "5 4 3 2 1", 0}, // reverse vector
	{"⌽2 3⍴⍳6", "3 2 1\n6 5 4", 0}, // reverse matrix
	{"⊖2 3⍴⍳6", "4 5 6\n1 2 3", 0}, // reverse first
	{"⌽[1]2 3⍴⍳6", "4 5 6\n1 2 3", 0},
	{"⊖[2]2 3⍴⍳6", "3 2 1\n6 5 4", 0},
	{"A←2 3⍴⍳12 ⋄ (⌽[1]A)←2 3⍴-⍳6⋄A", "¯4 ¯5 ¯6\n¯1 ¯2 ¯3", 0},
	{"⌽'DESSERTS'", "S T R E S S E D", 0}, // reverse strings
	{"⍝ Rotate", "", 0},
	{"1⌽1 2 3 4", "2 3 4 1", 0},                                                     // rotate vector
	{"10⌽1 2 3 4", "3 4 1 2", 0},                                                    // rotate vector
	{"¯1⌽1 2 3 4", "4 1 2 3", 0},                                                    // rotate vector negative
	{"(-7)⌽1 2 3 4", "2 3 4 1", 0},                                                  // rotate vector negative
	{"1 2⌽2 3⍴⍳6", "2 3 1\n6 4 5", 0},                                               // rotate array
	{"(2 2⍴2 ¯3 3 ¯2)⌽2 2 4⍴⍳16", "3 4 1 2\n6 7 8 5\n\n12 9 10 11\n15 16 13 14", 0}, // rotate array
	{"(2 3⍴2 ¯3 3 ¯2 1 2)⊖2 2 3⍴⍳12", "1 8 9\n4 11 6\n\n7 2 3\n10 5 12", 0},         // rotate array
	{"(2 4⍴0 1 ¯1 0 0 3 2 1)⌽[2]2 2 4⍴⍳16", "1 6 7 4\n5 2 3 8\n\n9 14 11 16\n13 10 15 12", 0},
	{"A←3 4⍴⍳12⋄(1 ¯1 2 ¯2⌽[1]A)←3 4⍴'ABCDEFGHIJKL'⋄A", "I F G L\nA J K D\nE B C H", 0},

	{"⍝ Transpose", "apl/primitives/transpose.go", 0},
	{"1 2 1⍉2 3 4⍴⍳6", "1 5 3\n2 6 4", 0},
	{"⍉3 1⍴1 2 3", "1 2 3", 0},
	{"⍴⍉2 3⍴⍳6", "3 2", 0},
	{"+/+/1 3 2⍉2 3 4⍴⍳24", "78 222", 0},
	{"+/+/3 2 1⍉2 3 4⍴⍳24", "66 72 78 84", 0},
	{"+/+/2 1 3⍉2 3 4⍴⍳24", "68 100 132", 0},
	{"1 1 1⍉2 3 3⍴⍳18", "1 14", 0},
	{"1 1 1⍉2 3 4⍴'ABCDEFGHIJKL',⍳12", "A 6", 0},
	{"1 1 2⍉2 3 4⍴'ABCDEFGHIJKL',⍳12", "A B C D\n5 6 7 8", 0},
	{"2 2 1⍉2 3 4⍴'ABCDEFGHIJKL',⍳12", "A 5\nB 6\nC 7\nD 8", 0},
	{"1 2 2⍉2 3 4⍴'ABCDEFGHIJKL',⍳12", "A F K\n1 6 11", 0},
	{"1 2 1⍉2 3 4⍴'ABCDEFGHIJKL',⍳12", "A E I\n2 6 10", 0},
	{"⍴⍴(⍳0)⍉5", "0", 0},
	{"⍴2 1 3⍉3 2 4⍴⍳24", "2 3 4", 0},
	{"⎕IO←0⋄⍴1 0 2⍉3 2 4⍴⍳24", "2 3 4", 0},
	{"A←3 3⍴⍳9⋄(1 1⍉A)←10 20 30⋄A", "10 2 3\n4 20 6\n7 8 30", 0},

	{"⍝ Enclose, string catenation, join strings, disclose, split", "apl/primitives/enclose.go", 0},
	{`⊂'alpha'`, "alpha", 0},
	{`"+"⊂'alpha'`, "a+l+p+h+a", 0},
	{`"\n"⊂"alpha" "beta" "gamma"`, "alpha\nbeta\ngamma", 0},
	{"`alpha`beta`gamma", "alpha beta gamma", 0},
	{"(`alpha`beta`gamma)", "alpha beta gamma", 0},
	{"`alpha`beta`gamma⋄", "alpha beta gamma", 0},
	{`⊃"alpha"`, "a l p h a", 0},
	{`'p'⊃"alpha"`, "al ha", 0},
	{`⍴','⊃",a,,b,c"`, "5", 0},
	{`⍴""⊃" a  b c\tc "`, "4", 0},

	{"⍝ Domino, solve linear system", "apl/primitives/domino.go", 0},
	{"⌹2 2⍴2 0 0 1", "0.5 0\n0 1", small},
	// TODO: this fails for big.Float. Remove sfloat and debug
	{"(1 ¯2 0)⌹3 3⍴3 2 ¯1 2 ¯2 4 ¯1 .5 ¯1", "1\n¯2\n¯2", small},
	// A←2a30
	// B←1a10
	// RHS←A+B**(¯1+⍳6)×○1÷3
	// S←⍉2 6⍴(6⍴1),*0J1×(¯1+⍳6)×○1÷3
	// ⍉RHS⌹S
	// With rational numbers:
	// A←3 3⍴9?100
	// B←3 3⍴9?100
	// 0=⌈/⌈/|B-A+.×B⌹A

	{"⍝ Dates, Times and durations", "apl/numbers/time.go", small},
	{"2018.12.23", "2018.12.23T00.00.00.000", small},       // Parse a time
	{"2018.12.23+12s", "2018.12.23T00.00.12.000", small},   // Add a duration to a time
	{"2018.12.24<2018.12.23", "0", small},                  // Times are comparable
	{"⌊/3s 2s 10s 4s", "2s", small},                        // Durations are comparable
	{"2018.12.23-1s", "2018.12.22T23.59.59.000", small},    // Substract a duration from a time
	{"2017.03.01-2017.02.28", "24h0m0s", small},            // Substract two times returns a duration
	{"2016.03.01-2016.02.28", "48h0m0s", small},            // Leap years are covered
	{"3m-62s", "1m58s", small},                             // Substract two durations
	{"-3s", "¯3s", small},                                  // Negate a duration
	{"×¯3h 0s 2m 2015.01.02", "¯1 0 1 1", small},           // Signum
	{"(|¯1s)+|1s", "2s", small},                            // Absolute value of a duration
	{"3×1h", "3h0m0s", small},                              // Uptype numbers to seconds and multiply durations
	{"1m × ⍳5", "1m0s 2m0s 3m0s 4m0s 5m0s", small},         // Generate a duration vector
	{"⍴⍪2018.12.23 + 1h×(¯1+⍳24)", "24 1", small},          // Table with all starting hours in a day
	{"4m×42.195", "2h48m46.8s", small},                     //
	{"⌈2018.12.23+3.5s", "2018.12.23T00.00.04.000", small}, // Ceil rounds to seconds
	{"⌊3h÷42.195", "4m15s", small},                         // Floor truncates seconds.

	{"⍝ Basic operators", "apl/operators/", 0},
	{"+/1 2 3", "6", 0},                            // plus reduce
	{"1 2 3 +.× 4 3 2", "16", 0},                   // scalar product
	{"(2 3⍴⍳6) +.× 3 2⍴5+⍳6", "52 58\n124 139", 0}, // matrix multiplication
	{`-\×\+\1 2 3`, "1 ¯2 16", 0},                  // chained monadic operators
	{"+/+/+/+/1 2 3", "6", 0},
	{`+.×/2 3 4`, "24", 0},
	// {`S←0.0 n→f "%.0f"⋄ +.×.*/2 3 4`, "2417851639229258349412352", 0},
	{`+.×.*/2 3 4`, "2.41785E+24", small},
	{`+.*.×/2 3 4`, "24", 0},

	{"⍝ Identify item for reduction over empty array", "apl/operators/identity.go", 0},
	{"+/⍳0", "0", 0},
	{"-/⍳0", "0", 0},
	{"×/⍳0", "1", 0},
	{"÷/⍳0", "1", 0},
	{"|/⍳0", "0", 0},
	{"⌊/⍳0", "¯1.79769E+308", small},
	{"⌈/⍳0", "1.79769E+308", small},
	{"*/⍳0", "1", 0},
	{"!/⍳0", "1", 0},
	{"^/⍳0", "1", 0},
	{"∧/⍳0", "1", 0},
	{"∨/⍳0", "0", 0},
	{"</⍳0", "0", 0},
	{"≤/⍳0", "1", 0},
	{"=/⍳0", "1", 0},
	{"≥/⍳0", "1", 0},
	{">/⍳0", "0", 0},
	{"≠/⍳0", "0", 0},
	{"⊤/⍳0", "0", 0},
	{"⌽/⍳0", "0", 0},
	{"⊖/⍳0", "0", 0},
	{"∨/0 3⍴ 1", "", 0},
	{"∨/3 3⍴ ⍳0", "0 0 0", 0},
	{"∪/⍳0", "0", 0},
	// These are implemented as operators and do not parse.
	// {"//⍳0", "0", 0},
	// {"⌿/⍳0", "0", 0},
	// {`\/⍳0`, "0", 0},
	// {`⍀/⍳0`, "0", 0},

	{"⍝ Outer product", "apl/operators/dot.go", 0},
	{"10 20 30∘.+1 2 3", "11 12 13\n21 22 23\n31 32 33", 0},
	{"(⍳3)∘.=⍳3", "1 0 0\n0 1 0\n0 0 1", 0},
	{"1 2 3∘.×4 5 6", "4 5 6\n8 10 12\n12 15 18", 0},

	{"⍝ Each", "apl/operators/each.go", 0},
	{"-¨1 2 3", "¯1 ¯2 ¯3", 0},   // monadic each
	{"1+¨1 2 3", "2 3 4", 0},     // dyadic each
	{"1 2 3+¨1", "2 3 4", 0},     // dyadic each
	{"1 2 3+¨4 5 6", "5 7 9", 0}, // dyadic each
	{"1+¨1", "2", 0},             // dyadic each

	{"⍝ Commute, duplicate", "apl/operators/commute.go", 0},
	{"∘.≤⍨1 2 3", "1 1 1\n0 1 1\n0 0 1", 0},
	{"+/∘(÷∘⍴⍨)⍳10", "5.5", small}, // mean value
	{"⍴⍨3", "3 3 3", 0},
	{"3-⍨4", "1", 0},
	{"+/2*⍨2 2⍴4 7 1 8", "65 65", 0},
	{"3-⍨4", "1", 0},

	{"⍝ Composition", "apl/operators/jot.go", 0},
	{"+/∘⍳¨2 4 6", "3 10 21", 0}, // Form I
	{"1∘○ 10 20 30", "¯0.544021 0.912945 ¯0.988032", small},
	{"+∘÷/40⍴1", "1.61803", small},     // Form IV, golden ratio (continuous-fraction)
	{"(*∘0.5)4 16 25", "2 4 5", float}, // Form III

	{"⍝ Power operator", "apl/operators/power.go", 0},
	{"⍟⍣2 +2 3 4", "¯0.366513 0.0940478 0.326634", float}, // log log
	// TODO: 1+∘÷⍣=1 oscillates for big.Float.
	// TODO: Add comparison tolerance and remove sfloat.
	{"1+∘÷⍣=1", "1.61803", small}, // fixed point iteration golden ratio
	{"⍝ TODO: function inverse", "", 0},

	{"⍝ Rank operator", "apl/operators/rank.go", 0},
	{`+\⍤0 +2 3⍴1`, "1 1 1\n1 1 1", 0},
	{`+\⍤1 +2 3⍴1`, "1 2 3\n1 2 3", 0},
	{"⍴⍤1 +2 3⍴1", "3\n3", 0},
	{"⍴⍤2 +2 3 5⍴1", "3 5\n3 5", 0},
	{"4 5+⍤1 0 2 +2 2⍴7 8 9 10", "11 12\n13 14\n\n12 13\n14 15", 0},
	{"⍉2 2 2⊤⍤1 0 ⍳5", "0 0 0 1 1\n0 1 1 0 0\n1 0 1 0 1", 0},
	{"⍳⍤1 +3 1⍴⍳3", "1 0 0\n1 2 0\n1 2 3", 0},

	{"⍝ At", "apl/operators/at.go", 0},
	{"(10 20@2 4)⍳5", "1 10 3 20 5", 0},
	{"10 20@2 4⍳5", "1 10 3 20 5", 0},
	{"((2 3⍴10 20)@2 4)4 3⍴⍳12", "1 2 3\n10 20 10\n7 8 9\n20 10 20", 0},
	{"⍴@(0.5∘<)3 3⍴1 4 0.2 0.3 0.3 4", "5 5 0.2\n0.3 0.3 5\n5 5 0.2", small},
	{"÷@2 4 ⍳5", "1 0.5 3 0.25 5", small},
	{"⌽@2 4 ⍳5", "1 4 3 2 5", 0},
	{"10×@2 4⍳5", "1 20 3 40 5", 0},
	{`(+\@2 4)4 3⍴⍳12`, "1 2 3\n4 9 15\n7 8 9\n10 21 33", 0},
	{"0@(2∘|)⍳5", "0 2 0 4 0", 0},
	{"÷@(2∘|)⍳5", "1 2 0.333333 4 0.2", small},
	{"⌽@(2∘|)⍳5", "5 2 3 4 1", 0},

	{"⍝ Stencil", "apl/operators/stencil.go", 0},
	{"{⌈/⌈/⍵}⌺(3 3) ⊢3 3⍴⍳25", "5 6 6\n8 9 9\n8 9 9", 0},

	{"⍝ Assignment, specification", "apl/operators/assign.go", 0},
	{"X←3", "", 0},              // assign a number
	{"-X←3", "¯3", 0},           // assign a value and use it
	{"X←3⋄X←4", "", 0},          // assign and overwrite
	{"X←3⋄⎕←X", "3", 0},         // assign and check
	{"f←+", "", 0},              // assign a function
	{"f←+⋄⎕←3 f 3", "6", 0},     // assign a function and apply
	{"X←4⋄⎕←÷X", "0.25", small}, // assign and use it in another expr
	{"A←2 3 ⋄ A", "2 3", 0},     // assign a vector

	{"⍝ Indexed assignment", "apl/operators/assign.go", 0},
	{"A←2 3 4 ⋄ A[1]←1 ⋄ A", "1 3 4", 0},
	{"A←2 2⍴⍳4 ⋄ +A[1;1]←3 ⋄ A", "3\n3 2\n3 4", 0},
	{"A←⍳5 ⋄ A[2 3]←10 ⋄ A", "1 10 10 4 5", 0},
	{"A←2 3⍴⍳6 ⋄ A[;2 3]←2 2⍴⍳4 ⋄ A", "1 1 2\n4 3 4", 0},
	{"⍝ TODO: choose/reach indexed assignment", "", 0},

	{"⍝ Multiple assignment", "apl/operators/assign.go", 0},
	{"A←B←C←D←1 ⋄ A B C D", "1 1 1 1", 0},

	{"⍝ Vector assignment", "apl/operators/assign.go", 0},
	{"(A B C)←2 3 4 ⋄ A ⋄ B ⋄ C ", "2\n3\n4", 0},
	{"-A B C←1 2 3 ⋄ A B C", "¯1 ¯2 ¯3\n1 2 3", 0},

	{"⍝ Modified assignment", "apl/operators/assign.go", 0},
	{"A←1 ⋄ A+←1 ⋄ A", "2", 0},
	{"A←1 2⋄ A+←1 ⋄ A", "2 3", 0},
	{"A←1 2 ⋄ A+←3 4 ⋄ A", "4 6", 0},
	{"A←1 2 ⋄ A{⍺+⍵}←3 ⋄ A", "4 5", 0},
	{"A B C←1 2 3 ⋄ A B C +← 4 5 6 ⋄ A B C", "5 7 9", 0},

	// Selective specification APL2 p.41, DyaRef p.21
	{"⍝ Selective assignment/specification", "apl/operators/assign.go", 0},
	{"A←10 20 30 40 ⋄ (2↑A)←100 200 ⋄ A", "100 200 30 40", 0},
	{"A←'ABCD' ⋄ (3↑A)←1 2 3 ⋄ A", "1 2 3 D", 0},
	{"A←1 2 3 ⋄ ((⍳0)↑A)←4 ⋄ A", "4 4 4", 0},
	//{"A←1 2 3 ⋄ (4↑A)←4 ⋄ A", "4 4 4", 0}, // overtake is ignored
	{"A←2 3⍴⍳6 ⋄ (,A)←2×⍳6 ⋄ A", "2 4 6\n8 10 12", 0},
	{"A←3 4⍴⍳12 ⋄ (4↑,⍉A)←10 20 30 40 ⋄ ,A ", "10 40 3 4 20 6 7 8 30 10 11 12", 0},
	{"A←2 3⍴'ABCDEF' ⋄ A[1;1 3]←8 9 ⋄ A", "8 B 9\nD E F", 0},
	{"A←2 3 4 ⋄ A[]←9 ⋄ A", "9 9 9", 0},
	{"A←4 3⍴⍳12 ⋄ (1 0 0/A)←1 4⍴⍳4 ⋄ A[3;1]", "3", 0}, // single element axis are collapsed
	{"A←3 2⍴⍳6 ⋄ (1 0/A)←'ABC' ⋄ A", "A 2\nB 4\nC 6", 0},
	{"A←4 5 6 ⋄ (1 ¯1  1/A)←7 8 9 ⋄ A", "7 5 9", 0},
	{"A←3 2⍴⍳6 ⋄ B←2 2⍴'ABCD' ⋄ (1 0 1/[1]A)←B ⋄ A", "A B\n3 4\nC D", 0},
	{"A←5 6 7 8 9 ⋄ (2↓A)←⍳3 ⋄ A", "5 6 1 2 3", 0},
	{"A←3 4⍴'ABCDEFGHIJKL' ⋄ (1 ¯1↓A)←2 3⍴⍳6 ⋄ A", "A B C D\n1 2 3 H\n4 5 6 L", 0},
	{"A←2 3⍴⍳6 ⋄ (1↓[1]A)←9 8 7 ⋄ A", "1 2 3\n9 8 7", 0},
	{"A←2 3 4⍴⍳12⋄(¯1 2↓[3 2]A)←0⋄A", "1 2 3 4\n5 6 7 8\n0 0 0 12\n\n1 2 3 4\n5 6 7 8\n0 0 0 12", 0},
	{`A←'ABC' ⋄ (1 0 1 0 1\A)←⍳5 ⋄ A`, "1 3 5", 0},
	{`A←2 3⍴⍳6 ⋄ (1 0 1 1\A)←10×2 4⍴⍳8 ⋄ A`, "10 30 40\n50 70 80", 0},
	{`A←3 2⍴⍳6 ⋄ (1 1 0 0 1\[1]A)←5 2⍴-⍳10 ⋄ A`, "¯1 ¯2\n¯3 ¯4\n¯9 ¯10", 0},
	{"A←2 3⍴⍳6 ⋄ (,A)←10×⍳6 ⋄ A", "10 20 30\n40 50 60", 0},
	{"A←2 3 4⍴⍳24 ⋄ (,[2 3]A)←2 12⍴-⍳24⋄⍴A⋄A[2;3;]", "2 3 4\n¯21 ¯22 ¯23 ¯24", 0},
	{"A←'GROWTH' ⋄ (2 3⍴A)←2 3⍴-⍳6 ⋄ (4⍴A)←⍳4 ⋄ A", "1 2 3 4 ¯5 ¯6", 0},
	{"A←3 4⍴⍳12 ⋄ (⌽A)←3 4⍴'STOPSPINODER' ⋄ A", "P O T S\nN I P S\nR E D O", 0},
	{"A←2 3⍴⍳6 ⋄ (⌽[1]A)←2 3⍴-⍳6 ⋄ A", "¯4 ¯5 ¯6\n¯1 ¯2 ¯3", 0},
	{"A←⍳6 ⋄ (2⌽A)←10×⍳6 ⋄ A", "50 60 10 20 30 40", 0},
	{"A←3 4⍴⍳12 ⋄ (1 ¯1 2 ¯2⊖A)←3 4⍴4×⍳12 ⋄ A", "36 24 28 48\n4 40 44 16\n20 8 12 32", 0},
	{"A←3 4⍴⍳12 ⋄ (1 ¯1 2 ¯2⌽[1]A)←3 4⍴4×⍳12 ⋄ A", "36 24 28 48\n4 40 44 16\n20 8 12 32", 0},
	{"A←⍳5 ⋄ (2↑A)← 10 20 ⋄ A", "10 20 3 4 5", 0},
	{"A←2 3⍴⍳6 ⋄ (¯2↑[2]A)←2 2⍴10×⍳4 ⋄ A", "1 10 20\n4 30 40", 0},
	{"A←3 3⍴⍳9 ⋄ (1 1⍉A)←10 20 30 ⋄ A", "10 2 3\n4 20 6\n7 8 30", 0},
	{"A←3 3⍴'STYPIEANT' ⋄ (⍉A)←3 3⍴⍳9 ⋄ A", "1 4 7\n2 5 8\n3 6 9", 0},
	{"⍝ TODO: First (↓) and Pick (⊃) are not implemented", "", 0},

	{"⍝ Functional selective specification", "apl/operators/assign.go", 0},    // iv extension
	{"A←3 3⍴⍳9 ⋄ A[{⍺[2]>⍺[1]}]←0 ⋄ A", "1 0 0\n4 5 0\n7 8 9", 0},             // ⍺ is the index vector
	{"A←10×3 3⍴⍳9 ⋄ A[{(⍵>30)^⍵<60}]←0 ⋄ A", "10 20 30\n0 0 60\n70 80 90", 0}, // ⍵ is the scalar value

	{"⍝ Lambda expressions", "apl/lambda.go", 0},
	{"{2×⍵}3", "6", 0},           // lambda in monadic context
	{"2{⍺+3{⍺×⍵}⍵+2}2", "14", 0}, // nested lambas
	{"2{(⍺+3){⍺×⍵}⍵+⍺{⍺+1+⍵}1+2}2", "40", 0},
	{"1{1+⍺{1+⍺{1+⍺+⍵}1+⍵}1+⍵}1", "7", 0},
	{"2{}4", "", 0}, // empty lambda expression ignores arguments
	{"{⍺×⍵}/2 3 4", "24", 0},
	{"A←1⋄{A+←1⋄A>0:B←A⋄B}0", "2", 0}, // continue if guarded expr is an assignment (differs from Dyalog)
	{`{1:1+2⋄{1:1+⍵}3}4`, "3", 0},

	{"⍝ Evaluation order", "apl/function.go", 0},
	{"A←1⋄A+(A←2)", "4", 0},
	{"A+A←3", "6", 0},
	{"A←1⋄A{(⍺ ⍵)}A+←10", "11 10", 0},

	{"⍝ Lexical scoping", "apl/lambda.go", 0},
	{"A←1⋄{A←2⋄A}0⋄A", "2\n1", 0},
	{"X←{A←3⋄B←4⋄0:ignored⋄42}0⋄X⋄A⋄B", "42\nA\nB", 0},
	{"{A←1⋄{A←⍵}⍵+1}1", "2", 0},
	{"A←1⋄S←{A←2}0⋄A", "1", 0},
	{"A←1⋄S←{A⊢←2}0⋄A", "2", 0}, // overwrite a global
	{"A←1⍴1⋄S←{A[1]←2}0⋄A", "2", 0},
	{"A←1⋄{A+←1⋄A}0⋄A", "2\n2", 0},
	{"+X←{A←3⋄B←4}0", "4", 0},

	{"⍝ Default left argument", "apl/lambda.go", 0},
	{"f←{⍺←3⋄⍺+⍵}⋄ f 4 ⋄ 1 f 4", "7\n5", 0},

	{"⍝ Recursion", "apl/lambda.go", 0},
	{`f←{⍵≤1: 1 ⋄ ⍵×∇⍵-1} ⋄ f 10`, "3628800", small},
	{"S←0{⍺>20:⍺⋄⍵∇⎕←⍺+⍵}1", "1\n2\n3\n5\n8\n13\n21\n34", 0},

	{"⍝ Tail call", "apl/lambda.go", 0},
	{"{⍵>1000:⍵⋄∇⍵+1}1", "1001", 0},

	{"⍝ Trains, forks, atops", "apl/train.go", 0},
	{"-,÷ 5", "¯0.2", float},
	{"(-,÷)5", "¯5 0.2", float},
	{"3(+×-)1", "8", 0},
	{"(+⌿÷≢)3+⍳13", "10", 0},
	{"(⍳{⍺/⍵}⍳)3", "1 2 2 3 3 3", 0},
	{"(2/⍳)3", "1 1 2 2 3 3", 0},
	{"6(+,-,×,÷)2", "8 4 12 3", 0},
	{"6(⌽+,-,×,÷)2", "3 12 4 8", 0},
	{"(⍳12) (⍳∘1 ≥)9", "9", 0},
	{"(*-)1", "0.367879", float},
	{"2(*-)1", "2.71828", float},
	{"1(*-)2", "0.367879", float},
	{"3(÷+×-)1", "0.125", float},
	{"(÷+×-)4", "¯0.0625", float},
	{"(⌊÷+×-)4", "¯0.25", float},
	{"6(⌊÷+×-)4", "0.2", float},
	{"(3+*)4", "57.5982", float}, // Agh fork
	//{"(⍳(/∘⊢)⍳)3", "1 2 2 3 3 3", 0}, // The hybrid token does not parse.

	{"⍝ Go interface package strings", "apl/strings/register.go", 0},
	{`u←s→toupper ⋄ u "alpha"`, "ALPHA", 0},
	{`";" s→join "alpha" "beta" `, "alpha;beta", 0},

	{"⍝ Lists", "apl/list.go", 0},
	{"(1;2;)", "(1;2;)", 0},
	{"(1 5 9;(2;3+4;);)", "(1 5 9;(2;7;);)", 0},
	{"(+;1;2;)", "(+;1;2;)", 0},
	{"(/;+;1;2;)", "(/;+;1;2;)", 0},
	{"(.;+;×;1 2;3 4;)", "(.;+;×;1 2;3 4;)", 0},
	{"1 2 3+(4;5;6;)", "5 7 9", 0},
	{"+/(1;2;3;)", "6", 0},
	{"+/(1;2;(3;4;);)", "6 7", 0},

	{"⍝ Lists catenate, enlist, cut, each", "apl/primitives/comma.go", 0},
	{"1,(2;3;)", "(1;2;3;)", 0},
	{"(1;2;),3", "(1;2;3;)", 0},
	{"(1;2;),(3;4;)", "(1;2;3;4;)", 0},
	{"((1;2;);(3;4;);),(5;6;)", "((1;2;);(3;4;);5;6;)", 0},
	{"∊3", "(3;)", 0},
	{"∊⍳0", "(;)", 0},
	{"∊(1;2;3;)", "(1;2;3;)", 0},
	{"∊(1;(2;3;);(4;(5;6;););7 8 9;)", "(1;2;3;4;5;6;7 8 9;)", 0},
	{"1 3↓(1;2;3;)", "((1;2;);(3;);)", 0},
	{"(1;2;(3;4;);)+¨(1;2;(3;4;);)", "(2;4;6 8;)", 0},
	{"≢¨(1;2;(3;4;);)", "(1;1;2;)", 0},

	{"⍝ List indexing", "apl/primitives/index.go", 0},
	{"L←(1;2;)⋄L[2]", "2", 0},
	{"L←(1;(2;3;);4;)⋄L[2;1]", "2", 0},
	{"L←(1;(2;3;);4;)⋄L[0]", "4", 0},
	{"L←(1;(2;3;);4;)⋄L[2;0]", "3", 0},
	{"L←(1;(2;3;);4;)⋄L[2]", "(2;3;)", 0},
	{"L←(1;(2;3;);4;)⋄L[2][2]", "3", 0},
	{"⍝ Indexing with lists is not supported", "", 0},

	{"⍝ List indexed assignment", "apl/primitives/index.go", 0},
	{"L←(1;2;)⋄L[1]←3⋄L", "(3;2;)", 0},
	{"L←(1;(2;3;);4;)⋄L[2;1]←5⋄L", "(1;(5;3;);4;)", 0},
	{"L←(1;(2;3;);4;)⋄L[2;0]←5⋄L", "(1;(2;5;);4;)", 0},
	{"L←(1;(2;3;);4;)⋄L[2;¯1]×←5⋄L", "(1;(10;3;);4;)", 0},

	{"⍝ Dictionaries", "apl/object.go", 0},
	{"D←`alpha#1 2 3⋄D[`alpha]←`xyz⋄D", "alpha: xyz", 0},
	{"D←`alpha#1⋄D[`alpha`beta]←3 4⋄D", "alpha: 3\nbeta: 4", 0},
	{"D←`a`b`c#1⋄D⋄#D", "a: 1\nb: 1\nc: 1\na b c", 0},
	{"D←`a`b`c#1 2 3⋄G←D[`a`c]⋄G", "a: 1\nc: 3", 0},
	{"D←`a`b#(1;(`c`d#`F`G);)⋄D[`b;`d]←123⋄D[`b]", "c: F\nd: 123", 0},

	{"⍝ Table, transpose a dict to create a table", "apl/primitives/transpose.go", 0},
	{"⍉`a`b#1 2", "a b\n1 2", 0},
	{"⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)", "a b c\n1 4 7\n2 5 8\n3 6 9", small},
	{"⍉⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)", "a: 1 2 3\nb: 4 5 6\nc: 7 8 9", small},
	{"⍴`a`b#(1 2 3;4 5 6;)", "2", 0},
	{"⍴⍉`a`b#(1 2 3;4 5 6;)", "3 2", small},

	{"⍝ Indexing tables", "apl/primitives/index.go", 0},                              // see cmd/apl/testdata/table.apl for queries
	{"T←⍉`a`b#1 2⋄T[1]", "a: 1\nb: 2", small},                                        // single row as a dict
	{"T←⍉`a`b#(1;3;)⋄T[1]", "a: 1\nb: 3", small},                                     // single row as a dict
	{"T←⍉`a`b#1 2⋄T[1⍴1]", "a b\n1 2", small},                                        // single row as a table
	{"T←⍉`a`b#(1 2 3;3 4 5;)⋄T[1]", "a: 1\nb: 3", small},                             // single row as a dict
	{"T←⍉`a`b#(1 2 3;3 4 5;)⋄T[2;`b]", "4", small},                                   // scalar value
	{"T←⍉`a`b#(1 2 3;3 4 5;)⋄T[1 3]", "a b\n1 3\n3 5", small},                        // multiple rows as a table
	{"T←⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)⋄T[;`b]", "4 5 6", small},                        // single column as a vector
	{"T←⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)⋄T[`b]", "4 5 6", small},                         // single column as a vector (string key)
	{"T←⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)⋄T[;1⍴`b]", "b\n4\n5\n6", small},                 // single column table
	{"T←⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)⋄T[1 2;`b]", "b\n4\n5", small},                   // subtable if any index is multiple
	{"T←⍉`a`b`c#(1 2 3;4.1 5.2 6.3;7 8 9;)⋄T[]", "1 4.1 7\n2 5.2 8\n3 6.3 9", small}, // empty index converts to array
	{"T←⍉`a`b#(1 2 3;3 4 5;)⋄T[{⍺>2}]", "a b\n3 5", small},                           // functional row index
	{"T←⍉`A`B#(1 2 3;3 4 5;)⋄T[{6=A+B};`B]", "B\n4", small},                          // functional row index with column variable

	{"⍝ Table updates", "apl/operators/assign.go", 0},
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T", "a b\n1 3\n2 2\n3 1", small},
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[1 3]←0 ⋄ T", "a b\n0 0\n2 2\n0 0", small},                    // update with scalar
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[1 3]←10×2 2⍴⍳4 ⋄ T", "a b\n10 20\n2 2\n30 40", small},        // update with array
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[1 3]←`a`b#8 9 ⋄ T", "a b\n8 9\n2 2\n8 9", small},             // update with object
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[1 3]←⍉`a`b#(8 9;10 11;) ⋄ T", "a b\n8 10\n2 2\n9 11", small}, // update with table
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[;`b]←5 6 7 ⋄ T", "a b\n1 5\n2 6\n3 7", small},                // update column
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[{⍺<3};`b]←9 ⋄ T", "a b\n1 9\n2 9\n3 1", small},               // update column with row selection function
	{"T←⍉`A`B#(⍳3;4-⍳3;) ⋄ T[{B<3};`A]←9 ⋄ T", "A B\n1 3\n9 2\n9 1", small},               // update column with row selection function using a key value
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[1 3]+←1 ⋄ T", "a b\n2 4\n2 2\n4 2", small},                   // update with modification function
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[`a]←1 ⋄ T", "a b\n1 3\n1 2\n1 1", small},                     // column name are given as first index
	{"T←⍉`a`b#(⍳3;4-⍳3;) ⋄ T[`a`b]←1 ⋄ T", "a b\n1 1\n1 1\n1 1", small},                   // column names are given as first index

	{"⍝ Elementary functions on dicts and tables", "apl/primitives/elementary.go", 0},
	{"A←`a`b#(1 2;3 4;)⋄-A", "a: ¯1 ¯2\nb: ¯3 ¯4", small},
	{"A←⍉`a`b#(1 2;3 4;)⋄-A", "a b\n¯1 ¯3\n¯2 ¯4", small},
	{"A←`a`b#(1 2;3 4;)⋄B←`a`b#(9 8;7 6;)⋄B-A", "a: 8 6\nb: 4 2", small},
	{"A←`a`b#(1 2;3 4;)⋄B←`b`c#(9 8;7 6;)⋄B-A", "b: 6 4\nc: 7 6\na: ¯1 ¯2", small},
	{"A←⍉`a`b#(1 2;3 4;)⋄B←⍉`b`c#(9 8;7 6;)⋄B-A", "b c a\n6 7 ¯1\n4 6 ¯2", small},
	{"A←⍉`a`b#(1 2;3 4;)⋄A-3", "a b\n¯2 0\n¯1 1", small},
	{"A←`a`b#(1 2;3 4;)⋄A-5 7", "a: ¯4 ¯5\nb: ¯2 ¯3", small},
	{"A←`a`b#(1 2;3 4;)⋄3-A", "a: 2 1\nb: 0 ¯1", small},

	{"⍝ Catenate tables or objects", "apl/primitives/comma.go", 0},
	{"A←`a`b#(1 2;3 4;)⋄B←`a`b#(5 6;7 8;)⋄A,B", "a: 1 2 5 6\nb: 3 4 7 8", small},
	{"A←`a`b#(1 2;3 4;)⋄B←`b`c#(5 6;7 8;)⋄A,B", "a: 1 2\nb: 3 4 5 6\nc: 7 8", small},
	{"A←`a`b#(1 2;3 4;)⋄B←`a`b#(5 6;7 8;)⋄A⍪B", "a: 5 6\nb: 7 8", small},
	{"A←`a`b#(1 2;3 4;)⋄B←`b`c#(5 6;7 8;)⋄A⍪B", "a: 1 2\nb: 5 6\nc: 7 8", small},
	{"A←⍉`a`b#(1 2;3 4;)⋄B←⍉`a`b#(5 6;7 8;)⋄A,B", "a b\n5 7\n6 8", small},
	{"A←⍉`a`b#(1 2;3 4;)⋄B←⍉`b`c#(5 6;7 8;)⋄A,B", "a b c\n1 5 7\n2 6 8", small},
	{"A←⍉`a`b#(1 2;3 4;)⋄B←⍉`a`b#(5 6;7 8;)⋄A⍪B", "a b\n1 3\n2 4\n5 7\n6 8", small},
	{"T←⍉`a`b#(1 2;3 4;)⋄T,←⍉`c#5 6⋄T", "a b c\n1 3 5\n2 4 6", small}, // catenate a column
	{"A←⍉`a`b#(1 2;3 4;)⋄A⍪5", "a b\n1 3\n2 4\n5 5", small},
	{"A←⍉`a`b#(1 2;3 4;)⋄A,5 6", "a b\n1 3\n2 4\n5 5\n6 6", small},
	{"A←⍉`a`b#(1 2;3 4;)⋄5 6,A", "a b\n5 5\n6 6\n1 3\n2 4", small},
	{"A←`a`b#(1 2;3 4;)⋄A,5", "a: 1 2 5\nb: 3 4 5", small},
	{"A←`a`b#(1 2;3 4;)⋄5 6⍪A", "a: 5 6 1 2\nb: 5 6 3 4", small},

	{"⍝ Reduction over objects and tables", "apl/operators/reduce.go", 0},
	{"+/`a`b`c#(1 2 3;4 6;7;)", "a: 6\nb: 10\nc: 7", small},
	{"+\\`a`b`c#(1 2 3;4 6;7;)", "a: 1 3 6\nb: 4 10\nc: 7", small},
	{"+/⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)", "a b c\n6 15 24", small},
	{"+\\⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)", "a b c\n1 4 7\n3 9 15\n6 15 24", small},
	{"2+/`a`b#(1 2 3;4 6 7;)", "a: 3 5\nb: 10 13", small},
	{"2+/⍉`a`b#(1 2 3;4 6 7;)", "a b\n3 10\n5 13", small},
	{"T←⍉`a`b`c#(1 2 3;4 5 6;7 8 9;)⋄T⍪(+⌿÷≢)T", "a b c\n1 4 7\n2 5 8\n3 6 9\n2 5 8", small},

	{"⍝ Object, go example", "apl/xgo/register.go", 0},
	{"X←go→t 0⋄X[`V]←`a`b⋄X[`V]", "a b", 0},
	{"X←go→t 0⋄X[`I]←55⋄X[`inc]⍨0⋄X[`I]", "56", small},
	{"X←go→t 0⋄X[`V]←'abcd'⋄X[`join]⍨'+'", "(4;a+b+c+d;)", small},
	{"S←go→s 0⋄#[1]S", "sum", 0},
	{"T←go→t 0⋄T[`S;`A]←3⋄T[`S;`V]←2 3⋄T[`S]", "A: 3\nB: 0\nV: 2 3", 0},

	{"⍝ Channels read, write and close", "apl/primitives/take.go", 0},
	{"C←go→source 6⋄2 3↑C", "0 1 2\n3 4 5", 0},
	{"C←go→source 6⋄↑C⋄↑C⋄↓C", "0\n1\n1", 0},

	{"⍝ Reduce, scan and each over channel", "apl/operators/reduce.go", 0},
	{"C←go→source 6⋄+/C", "15", 0},
	{`C←go→source 6⋄+\C`, "0 1 3 6 10 15", 0},
	{`C←go→source 6⋄⊢\C`, "0 1 2 3 4 5", 0},
	{`C←go→source 4⋄5+¨C`, "5\n6\n7\n8", 0},
	{"C←go→source 3⋄C", "0\n1\n2", 0},
	{"C←go→source 3⋄-¨C", "0\n¯1\n¯2", 0},
	{"<¨⍳3", "1\n2\n3", 0},                                 // channel-each
	{"(<⍤2)2 2 3⍴⍳12", "1 2 3\n4 5 6\n7 8 9\n10 11 12", 0}, // channel-rank

	{"⍝ Communicate over a channel", "apl/channel.go", 0},
	{`C←go→echo"?"⋄C↓'a'⋄C↓'b'⋄2↑C⋄↓C`, "a\nb\n?a ?b\n1", 0},

	{"⍝ Primes", "", 0},
	{"f←{(2=+⌿0=X∘.|X)⌿X←⍳⍵} ⋄ f 42", "2 3 5 7 11 13 17 19 23 29 31 37 41", 0},        // 01-primes
	{"⎕IO←0 ⋄ f←{(~X∊X∘.×X)⌿X←2↓⍳⍵} ⋄ f 42", "2 3 5 7 11 13 17 19 23 29 31 37 41", 0}, // 01-primes

	{"⍝ π", "", 0},
	{".5*⍨6×+/÷2*⍨⍳1000", "3.14064", float},
	{"4×-/÷¯1+2×⍳100", "3.13159", float},
	{"4×+/{(⍵ ⍴ 1 0 ¯1 0)÷⍳⍵}100", "3.12159", float},

	{"⍝ Conway-completeness", "", 0},
	{"A←5 5⍴(23⍴2)⊤1215488⋄l←{3=S-⍵∧4=S←({+/,⍵}⌺3 3)⍵}⋄(l⍣8)A", "0 0 0 0 0\n0 0 0 0 0\n0 0 0 0 1\n0 0 1 0 1\n0 0 0 1 1", 0},
	// A←5 5⍴(23⍴2)⊤1215488⋄l←3=s-⊢∧4=s←{+/,⍵}⌺3 3⋄(l⍣8)A // TODO: This cannot be reduced.
	{"A←5 5⍴(23⍴2)⊤1215488 ⋄ s←{+/,⍵}⌺3 3 ⋄ l←{(3=s-⊢∧(4=s))⍵} ⋄ (l⍣8)A", "0 0 0 0 0\n0 0 0 0 0\n0 0 0 0 1\n0 0 1 0 1\n0 0 0 1 1", 0},
	{"A←5 5⍴(23⍴2)⊤1215488 ⋄ l←{(⊢(3=⊢-⊣∧(4=⊢)){+/,⍵}⌺3 3)⍵} ⋄ (l⍣8)A", "0 0 0 0 0\n0 0 0 0 0\n0 0 0 0 1\n0 0 1 0 1\n0 0 0 1 1", 0},
	// {≢⍸⍵}⌺3 3∊¨3+0,¨⊢ // needs nested arrays.
	// life2←{3=s-⍵∧4=s←{+/,⍵}⌺3 3⊢⍵} // Dya: works without braces.

	// github.com/DhavalDalal/APL-For-FP-Programmers
	// filter←{(⍺⍺¨⍵)⌿⍵} // 01-primes
	// ⎕IO←0 ⋄ sieve ← {⍸⊃{~⍵[⍺]:⍵ ⋄ 0@(⍺×2↓⍳⌈(≢⍵)÷⍺)⊢⍵}/⌽(⊂0 0,(⍵-2)⍴1),⍳⍵} // 02-sieve
	// ⎕IO←0 ⋄ triples←{{⍵/⍨(2⌷x)=+⌿2↑x←×⍨⍵}⍉↑,1+⍳⍵ ⍵ ⍵}// 03-pythagoreans
	// ⎕IO←0 ⋄ '-:'⊣@(' '=⊢)¨(14⍴(4⍴1),0)(17⍴1 1 0)\¨⊂⍉(⎕D,6↑⎕A)[(12⍴16)⊤?10⍴2*48] // 04-MacAddress
	// life←{⊃1 ⍵∨.∧3 4=+⌿,1 0 ¯1∘.⊖1 0 ¯1⌽¨⊂⍵} // 05-life
	// life2←{3=s-⍵∧4=s←{+/,⍵}⌺3 3⊢⍵} // 05-life

	// Trees: https://youtu.be/hzPd3umu78g

	//https://github.com/theaplroom/apl-sound-wave/blob/master/src/DSP.dyalog

}

func testCompare(got, exp string) bool {
	got = strings.TrimSpace(got)
	gotlines := strings.Split(got, "\n")
	explines := strings.Split(exp, "\n")
	if len(gotlines) != len(explines) {
		return false
	}
	for i, g := range gotlines {
		e := explines[i]
		gf := strings.Fields(g)
		ef := strings.Fields(e)
		if len(gf) != len(ef) {
			return false
		}
		for k := range gf {
			if gf[k] != ef[k] {
				return false
			}
		}
	}
	return true
}

const (
	float int = 1 << iota // only for floating point towers
	small                 // normal tower only
)

func TestNormal(t *testing.T) {
	testApl(t, nil, 0)
}

func TestBig(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	testApl(t, big.SetBigTower, small|float)
}

func TestPrecise(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	testApl(t, func(a *apl.Apl) { big.SetPreciseTower(a, 256) }, small)
}

func testApl(t *testing.T, tower func(*apl.Apl), skip int) {
	log := func(v ...interface{}) {
		if testing.Short() {
			t.Log(v...)
		}
	}
	logf := func(f string, v ...interface{}) {
		if testing.Short() {
			t.Logf(f, v...)
		}
	}

	// Print table of contents
	if testing.Short() {
		log("# Test results")
		logf("%s from `apl/primitives/gen.go` on %s\n", `Generated by [apl_test](apl/primitives/apl_test.go)`, time.Now().Format("2006-01-02 15:04:05"))
		for _, tc := range testCases {
			if strings.HasPrefix(tc.in, "⍝") {
				if strings.HasPrefix(tc.in, "⍝ TODO") {
					continue
				}
				s := strings.TrimPrefix(tc.in, "⍝ ")
				a := strings.ToLower(s)
				a = strings.Replace(a, " ", "-", -1)
				logf("- [%s](#%s)\n", s, a)
			}
		}
		logf("\n```apl\n")
	}

	// Compare result with expectation but ignores differences in whitespace.
	for i, tc := range testCases {

		if strings.HasPrefix(tc.in, "⍝") {
			if strings.HasPrefix(tc.in, "⍝ TODO") {
				log(tc.in)
			} else {
				s := strings.TrimPrefix(tc.in, "⍝ ")
				log("```")
				logf("## %s\n", s)
				if tc.exp != "" {
					logf("[→%s](%s)\n", tc.exp, tc.exp)
				}
				logf("\n```apl\n")
			}
			continue
		}

		// Skip tests for unsupported numberic types
		if skip&tc.flag != 0 {
			continue
		}

		var buf strings.Builder
		a := apl.New(&buf)
		numbers.Register(a)
		if tower != nil {
			tower(a)
		}
		Register(a)
		operators.Register(a)
		aplstrings.Register(a, "s")
		xgo.Register(a, "go")

		mustfail := strings.HasPrefix(tc.exp, "fail:")
		lines := strings.Split(tc.in, "\n")
		for k, s := range lines {
			logf("\t%s", s)
			err := a.ParseAndEval(s)
			if err != nil && mustfail == false {
				t.Fatalf("tc%d:%d: %s: %s\n", i+1, k+1, tc.in, err)
			} else if err == nil && mustfail == true {
				t.Fatalf("tc%d:%d: %s: should fail but did not", i+1, k+1, tc.in)
			}
		}
		if mustfail {
			log("Must", tc.exp) // This prints: "Must fail: ..."
			continue
		}

		got := buf.String()
		log(got)

		g := got
		g = spaces.ReplaceAllString(g, " ")
		g = newline.ReplaceAllString(g, "\n")
		g = strings.TrimSpace(g)
		if g != tc.exp {
			fmt.Printf("%q != %q\n", g, tc.exp)
			t.Fatalf("tc%d:\nin>\n%s\ngot>\n%s\nexpected>\n%s", i+1, tc.in, got, tc.exp)
		}
	}
}

var rat0, _ = big.ParseRat("0")
var spaces = regexp.MustCompile(`  *`)
var newline = regexp.MustCompile(`\n *`)
