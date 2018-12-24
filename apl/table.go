package apl

// A Table is a transposed dictionary where each value is a vector
// with the same number of elements and unique type.
// Tables are constructed by transposing dictionaries T←⍉D
//
// Indexing tables selects rows:
//	T[⍳5]
// returns a table with the first 5 rows.
// Right arrow indexing selects columns, just like a dict.
//	T→Col1
// Sorting by column
//	T[⍋T→Time]
// Selecting rows
//	T[⍸T→Qty>5]
type Table Dict
