package apl

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
	"unicode"

	"github.com/ktye/iv/apl/scan"
)

// ArrayString can be used by an array implementation.
// It formats an n-dimensional array using a tabwriter for PP>=-1.
// Each dimension is terminated by k newlines, where k is the dimension index.
// For PP==-2, it uses a single line json notation with nested brackets and
// for PP==-3, it formats in a single line matlab syntax (rank <= 2).
func ArrayString(a *Apl, v Array) string {
	if a.PP == -2 {
		return jsonArray(a, v)
	} else if a.PP == -3 {
		return matArray(a, v)
	}
	shape := v.Shape()
	if len(shape) == 0 {
		return ""
	} else if len(shape) == 1 {
		s := make([]string, shape[0])
		for i := 0; i < shape[0]; i++ {
			s[i] = v.At(i).String(a)
		}
		return strings.Join(s, " ")
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
		fmt.Fprintf(tw, "%s\t", v.At(i).String(a))
		if term := inc(); term > 0 {
			for k := 0; k < term; k++ {
				fmt.Fprintln(tw)
			}
		} else if term == -1 {
			fmt.Fprintln(tw)
		}
	}
	tw.Flush()
	s := buf.String()
	if len(s) > 0 && s[len(s)-1] == '\n' {
		// Don't print the final newline.
		return s[:len(s)-1]
	}
	return s
}

// stringArray converts the array to a string array of the same shape.
// All elements are printed with the current PP.
func stringArray(a *Apl, v Array) StringArray {
	sa := StringArray{Dims: CopyShape(v), Strings: make([]string, v.Size())}
	for i := range sa.Strings {
		sa.Strings[i] = v.At(i).String(a)
	}
	return sa
}

// jsonArray is used for PP=-2
func jsonArray(a *Apl, v Array) string {
	sa := stringArray(a, v)
	var vector func(v StringArray) string
	vector = func(S StringArray) string {
		if len(S.Dims) == 1 {
			return "[" + strings.Join(S.Strings[:S.Dims[0]], ",") + "]"
		}
		vec := make([]string, S.Dims[0])
		inc := prod(S.Dims[1:])
		for i := 0; i < S.Dims[0]; i++ {
			sub := StringArray{Dims: S.Dims[1:], Strings: S.Strings[i*inc:]}
			vec[i] = vector(sub)
		}
		return "[" + strings.Join(vec, ",") + "]"
	}
	return vector(sa)
}

// matArray is used for PP=-3. It only supported for rank 1 and 2.
func matArray(a *Apl, v Array) string {
	sa := stringArray(a, v)
	if len(sa.Dims) == 1 {
		return "[" + strings.Join(sa.Strings, ",") + "]"
	} else if len(sa.Dims) != 2 {
		return "[rank error]"
	}
	var b strings.Builder
	b.WriteString("[")
	off := 0
	for i := 0; i < sa.Dims[0]; i++ {
		b.WriteString(strings.Join(sa.Strings[off:off+sa.Dims[1]], ","))
		if i < sa.Dims[0]-1 {
			b.WriteString(";")
		}
		off += sa.Dims[1]
	}
	b.WriteString("]")
	return b.String()
}

// ParseArray parses an array from a string representation.
// If the protoptye is not 0, the result will have the same type.
func (a *Apl) ParseArray(prototype Value, s string) (Value, error) {
	if prototype != nil {
		if _, ok := prototype.(Array); ok == false {
			return nil, fmt.Errorf("parse array: prototype is not an array: %T", prototype)
		}
	}
	vector := func(line string) ([]Value, bool) {
		var values []Value
		rr := strings.NewReader(line)
		for {
			r, _, err := rr.ReadRune()
			if err == io.EOF {
				return values, true
			} else if r == '"' {
				rr.UnreadRune()
				if str, err := scan.ReadString(rr); err != nil {
					return nil, false
				} else {
					values = append(values, String(str))
				}
			} else if unicode.IsSpace(r) == false {
				if num, err := scan.ScanNumber(rr); err != nil {
					return nil, false
				} else if n, err := a.Tower.Parse(num); err != nil {
					return nil, false
				} else {
					values = append(values, n)
				}
			}
		}
	}
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		s = s[1 : len(s)-1]
	}
	var values []Value
	lines := strings.Split(s, "\n")
	var shape []int
	c := 0
	for k := range lines {
		vec, ok := vector(lines[k])
		if ok == false {
			return nil, fmt.Errorf("parse array: cannot parse")
		}
		n := len(vec)
		if shape == nil && n == 0 {
			return nil, fmt.Errorf("parse array: empty first line")
		} else if shape == nil {
			shape = []int{n}
		}
		if n != 0 && n != shape[len(shape)-1] {
			return nil, fmt.Errorf("parse array: last axis is not uniform %d != %d", n, shape[len(shape)-1])
		}
		// TODO: assemble rank from the number of delimiters.
		if n == 0 {
			c++
		}
		values = append(values, vec...)
	}
	if shape == nil {
		if prototype != nil {
			if u, ok := prototype.(Uniform); ok {
				return u.Make([]int{}), nil
			}
		}
		return EmptyArray{}, nil
	}

	// TODO parse nd-arrays.
	shape = []int{len(lines) / shape[0], shape[0]}
	if prod(shape) != len(values) {
		return nil, fmt.Errorf("parse array: array is not rectangular: shape %v, size %d", shape, len(values))
	}
	A := MixedArray{
		Values: values,
		Dims:   shape,
	}
	if prototype != nil {
		if u, ok := prototype.(Uniform); ok {
			res, ok := a.Unify(A, true)
			if ok == false {
				return nil, fmt.Errorf("parse uniform array: array has no uniform type")
			}
			if reflect.TypeOf(res) != reflect.TypeOf(u) {
				return nil, fmt.Errorf("parse uniform array: result has wrong type %T != %T", res, u)
			}
			return res, nil
		}
	}
	return A, nil
}
