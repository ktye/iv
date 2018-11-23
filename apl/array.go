package apl

import (
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
)

// IsScalar returns true,
// if the Value can be used as a member of an array.
// It returns false for arrays, functions and identifiers.
func (a *Apl) isScalar(v Value) bool {
	if _, ok := v.(Array); ok {
		return false
	}
	if _, ok := v.(Function); ok {
		return false
	}
	if _, ok := v.(Identifier); ok {
		return false
	}
	return true
}

// array evaluates to a EmptyArray, a single Value or a GeneralArray.
type array []expr

func (ar array) Eval(a *Apl) (Value, error) {
	if len(ar) == 0 {
		return EmptyArray{}, nil
	} else if len(ar) == 1 {
		return ar[0].Eval(a)
	}

	v := make([]Value, len(ar))
	for i, x := range ar {
		e, err := x.Eval(a)
		if err != nil {
			return nil, err
		}
		if a.isScalar(e) == false {
			return nil, fmt.Errorf("vector element must be scalar: %T", e)
		}
		v[i] = e
	}
	return GeneralArray{Values: v, Dims: []int{len(ar)}}, nil
}

func (ar array) String(a *Apl) string {
	v := make([]string, len(ar))
	for i, e := range ar {
		s := e.String(a)
		if _, ok := e.(String); ok {
			s = `"` + strings.Replace(s, `"`, `""`, -1) + `"`
		}
		v[i] = s

	}
	return fmt.Sprintf("(%s)", strings.Join(v, " "))
}

// Array is the interface that an actual array must implement.
// Examples are GeneralArray, BitArray.
// Arrays can be implemented externally.
type Array interface {
	String(*Apl) string
	Eval(*Apl) (Value, error)
	At(int) (Value, error)
	Shape() []int
	//ApplyMonadic(*Apl, PrimitiveHandler) (Value, error)
	//ApplyDyadic(*Apl, Value, bool, PrimitiveHandler) (Value, error)
}

// Reshaper is an array that can reshape itself.
type Reshaper interface {
	Reshape([]int) Value
}

// ArrayMaker is an array that can allocate a new array of it's type.
// An array that implements this interface can be assumed to be able to
// create arrays of itself for shape with elements >= 0.
type ArrayMaker interface {
	MakeArray([]int) ArraySetter
}

// ArraySetter is any Array implementation that has a Set method on top.
type ArraySetter interface {
	Array
	Set(int, Value) error
}

// MakeArray creates a new array.
// It makes an array of the same type as the prototype, if it can.
// Otherwise it returns a general array.
// The prototype may be nil.
func MakeArray(prototype Array, shape []int) ArraySetter {
	var am ArrayMaker
	if prototype != nil {
		if m, ok := prototype.(ArrayMaker); ok {
			am = m
		}
	}

	if am == nil {
		g := GeneralArray{Dims: shape}
		g.Values = make([]Value, ArraySize(g))
		return g
	} else {
		return am.MakeArray(shape)
	}
}

// ArraySize returns the product of the array shape.
func ArraySize(v Array) int {
	shape := v.Shape()
	if len(shape) == 0 {
		return 0
	}
	n := 1
	for _, i := range shape {
		n *= i
	}
	return n
}

// CopyShape copies the shape of an array.
func CopyShape(v Array) []int {
	shape := v.Shape()
	newshape := make([]int, len(shape))
	copy(newshape, shape)
	return newshape
}

// ArrayIndex converts an index vector to a flat index.
func ArrayIndex(shape []int, idx []int) (int, error) {
	if len(idx) != len(shape) {
		return -1, fmt.Errorf("array index has wrong rank")
	}
	n := 0
	times := 1
	for i := len(idx) - 1; i >= 0; i-- {
		if idx[i] < 0 || idx[i] >= shape[i] {
			return -1, fmt.Errorf("index exceeds array dimension %d", i+1)
		}
		n += times * idx[i]
		times = shape[i]
	}
	return n, nil
}

// ArrayIndexes converts a flat index to an index vector.
func ArrayIndexes(shape []int, res []int, flat int) error {
	if len(shape) == 0 || len(shape) != len(res) {
		return fmt.Errorf("wrong arguments to ArrayIndexes")
	}

	for i := len(res) - 1; i >= 0; i-- {
		if shape[i] == 0 {
			return fmt.Errorf("cannot compute indexes with zero dimension")
		}
		res[i] = flat % shape[i]
		flat -= res[i]
		flat /= shape[i]
	}
	return nil
}

// IncArrayIndex increases the index vector idx for the given shape.
func IncArrayIndex(idx []int, shape []int) {
	for i := 0; i < len(idx); i++ {
		idx[i]++
		if idx[i] < shape[i] {
			break
		}
		idx[i] = 0
	}
}

// ArrayAt returns the value at the index given by a vector.
// Usually arrays are indexed by their flat index method At.
func ArrayAt(v Array, idx []int) (Value, error) {
	n, err := ArrayIndex(v.Shape(), idx)
	if err != nil {
		return nil, err
	}
	return v.At(n)
}

// ArrayString can be used by an array implementation.
// It formats an n-dimensional array using a tabwriter.
// Each dimension is terminated by k newlines, where k is the dimension index.
func ArrayString(a *Apl, v Array) string {
	shape := v.Shape()
	if len(shape) == 0 {
		return ""
	}
	size := 1
	for _, n := range shape {
		size *= n
	}

	idx := make([]int, len(shape))
	inc := func() int {
		for i := 0; i < len(idx); i++ {
			k := len(idx) - 1 - i
			idx[k]++
			if idx[k] == shape[k] {
				idx[k] = 0
			} else {
				return i
			}
		}
		return -1 // should not happen
	}
	var buf strings.Builder
	tw := tabwriter.NewWriter(&buf, 1, 0, 1, ' ', tabwriter.AlignRight) // tabwriter.AlignRight)
	for i := 0; i < size; i++ {
		e, err := v.At(i)
		if err != nil {
			fmt.Fprintf(tw, "?\t")
		} else {
			fmt.Fprintf(tw, "%s\t", e.String(a))
		}
		if term := inc(); term > 0 {
			for k := 0; k < term; k++ {
				fmt.Fprintln(tw)
			}
		} else if term == -1 {
			fmt.Fprintln(tw)
		}
	}
	tw.Flush()
	return buf.String()
}

// GeneralArray is an n-dimensional array that can hold any Value.
type GeneralArray struct {
	Values []Value
	Dims   []int
}

func (v GeneralArray) Eval(a *Apl) (Value, error) {
	return v, nil
}

func (v GeneralArray) String(a *Apl) string {
	return ArrayString(a, v)
}

func (v GeneralArray) At(i int) (Value, error) {
	if i >= 0 && i < len(v.Values) {
		return v.Values[i], nil
	}
	return nil, fmt.Errorf("array index out of range")
}

func (v GeneralArray) Set(i int, e Value) error {
	if i < 0 || i >= len(v.Values) {
		return fmt.Errorf("index out of range")
	}
	v.Values[i] = e
	return nil
}

func (v GeneralArray) Shape() []int {
	return v.Dims
}

func (v GeneralArray) Reshape(shape []int) Value {
	if len(v.Values) == 0 {
		return EmptyArray{}
	}
	size := 1
	for _, k := range shape {
		size *= k
	}
	if size == 0 {
		return EmptyArray{}
	}
	rv := GeneralArray{
		Values: make([]Value, size),
		Dims:   shape,
	}
	k := 0
	for i := range rv.Values {
		rv.Values[i] = v.Values[k]
		k++
		if k == len(v.Values) {
			k = 0
		}
	}
	return rv
}

type EmptyArray struct{}

func (e EmptyArray) String(a *Apl) string       { return "" }
func (e EmptyArray) Eval(a *Apl) (Value, error) { return e, nil }
func (e EmptyArray) At(i int) (Value, error)    { return nil, fmt.Errorf("index out of range") }
func (e EmptyArray) Shape() []int               { return nil }
func (e EmptyArray) Reshape(s []int) Value      { return e }

/* TODO remove
// Bitarray is an array implementation which has only boolean values.
type Bitarray struct {
	Bits []Bool
	Dims []int
}

func (b Bitarray) String(a *Apl) string {
	return ArrayString(a, b)
}

func (b Bitarray) Eval(a *Apl) (Value, error) {
	return b, nil
}

func (b Bitarray) At(i int) (Value, error) {
	if i < 0 || i >= len(b.Bits) {
		return nil, fmt.Errorf("index exceeds array dimensions")
	}
	return b.Bits[i], nil
}

func (b Bitarray) Shape() []int {
	return b.Dims
}

func (b Bitarray) Reshape(shape []int) Value {
	if len(b.Bits) == 0 {
		return EmptyArray{}
	}
	size := 1
	for _, k := range shape {
		size *= k
	}
	if size == 0 {
		return EmptyArray{}
	}
	v := Bitarray{
		Bits: make([]Bool, size),
		Dims: shape,
	}
	k := 0
	for i := range v.Bits {
		v.Bits[i] = b.Bits[k]
		k++
		if k == len(b.Bits) {
			k = 0
		}
	}
	return v
}

// IntArray converts a Bitarray to a GeneralArray containing integers.
func (b Bitarray) IntArray() GeneralArray {
	ints := make([]Value, len(b.Bits))
	for i, v := range b.Bits {
		if v {
			ints[i] = Int(1)
		} else {
			ints[i] = Int(0)
		}
	}
	dims := make([]int, len(b.Dims))
	copy(dims, b.Dims)
	return GeneralArray{Dims: dims, Values: ints}
}
*/

type Bool bool

func (b Bool) String(a *Apl) string {
	if b {
		return "1"
	}
	return "0"
}

func (i Bool) Less(v Value) (Bool, bool) {
	j, ok := v.(Bool)
	if ok == false {
		return false, false
	}
	return i == false && j == true, true
}

type Index int

func (i Index) String(a *Apl) string {
	s := strconv.Itoa(int(i))
	return strings.Replace(s, "-", "¯", 1)
}

func (i Index) ToIndex() (int, bool) {
	return int(i), true
}

func (i Index) Less(v Value) (Bool, bool) {
	j, ok := v.(Index)
	if ok == false {
		return false, false
	}
	return i < j, true
}

// IntdexArray is an array implementation which has only int values.
type IndexArray struct {
	Ints []int
	Dims []int
}

func (ar IndexArray) String(a *Apl) string {
	return ArrayString(a, ar)
}

func (ar IndexArray) Eval(a *Apl) (Value, error) {
	return ar, nil
}

func (ar IndexArray) At(i int) (Value, error) {
	if i < 0 || i >= len(ar.Ints) {
		return nil, fmt.Errorf("index exceeds array dimensions")
	}
	return Index(ar.Ints[i]), nil
}

func (ar IndexArray) Shape() []int {
	return ar.Dims
}

func (ar IndexArray) Reshape(shape []int) Value {
	if len(ar.Ints) == 0 {
		return EmptyArray{}
	}
	size := 1
	for _, k := range shape {
		size *= k
	}
	if size == 0 {
		return EmptyArray{}
	}
	rv := IndexArray{
		Ints: make([]int, size),
		Dims: shape,
	}
	k := 0
	for i := range rv.Ints {
		rv.Ints[i] = ar.Ints[k]
		k++
		if k == len(ar.Ints) {
			k = 0
		}
	}
	return rv
}
