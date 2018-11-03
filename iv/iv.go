// Iv is a stream processor with an APL backend.
package iv

// It reads text from the input and writes to the output.
// Input is parsed as n-dimensional data, which is send
// to an APL backend repeatedly for fixed chunks of data.
//
// The data chunks depend on the rank parameter for iv.
// For rank 0, it calls APL for each scalar read.
// For rank 1 for each vector and so on.
//
// It also recognizes input data dimensions which are higher
// than the rank setting. These are send to APL together
// with the current array as the termination level E.
//
// rank0 (-r=0) execute each scalar
//   R=0                                  ?:   parseNext
//          |e    |e    |e    |e    |e    e:   error or EOF
// -->o-S++-?-E++-?-E++-?-E++-?-...-?-... E++: increase termination level
//   [X]    |     |     |     |     |     X:   execute and reset
//    +--<--+-----+-----+-----+-----+-...
//
//
// rank1        execute each vector       Initialize shape with [0]
//         R=1                            C:   append scalar to vector;
//          |e    |e    |e    |e    |e         increase shape[0]
// -->o-S++-?-E++-?-E++-?-E++-?-...-?-...
//    |     |     |     |     |     |
//    +--<--C1    |     |     |     |
//   [X]          |     |     |     |
//    +--<--------+-----+-----+-...-+-...
//
//
// rank2        execute each matrix       Initialize shape with [0, 0]
//               R=2                      CS:  append scalar to vector;
//          |e    |e    |e    |e    |e         increase shape[R-S] for CS;
// -->o-S++-?-S++-?-E++-?-E++-?-...-?-...      check uniformity at C2
//    |     |     |     |     |     |     S++: increase separator count
//    +--<--C1    |     |     |     |          S=1: a single separator follows
//    +--<--------C2    |     |     |               the scalar
//   [X]                |     |     |          S=2: two separators follow
//    +--<--------------+-----+-...-+-...           the scalar
//
//
// rankR   execute each R-dim array       Initialize shape with zeros of length R.
//          1     2     3     R           Each next scalar read will be catenated,
//          |e    |e    |e    |e          if S <= R.
// -->o-S++-?-S++-?-S++-?-....?.E++.?...  At a catenation CS, the value of the shape
//    |     |     |     |     |     |     is increased at index R-S.
//    +--<--C1    |     |     |     |     Uniformity is checked for CN, with N > 1.
//    +--<--------C2    |     |     |     If the number of separators is > R,
//    +--<--------------C3    |     |     An apl execution step is done, sending
//    +-...                   |     |     the aggregated vector, it's shape
//   [X]                      |     |     and the termination level E, which indicates
//    +--<----------------...-+-...-+-... the dimension that terminates with the current
//                                        object above rank R.
// At EOF: If the shape vector containts zeros,
// then the current object rank is smaller than expected,
// in this case, set the zeros to ones.
// The data may have missing closing separators, e.g. a scalar followed
// immediatedly by an EOF, or a table ending with a single newline.
// For convenience, the level at EOF will always be the maximum level observed.
//
// Uniformity: Input data is expected to be uniform upto rank R by default.
// E.g. in matrix mode (rank 2), each row must have the same number of columns,
// but subsequent table are allowed to have different matrix sizes.
// This is done by checking if the shape product equals the length of the
// aggregation vector.
// As a consequence, the apl side receives only simple arrays, not nested ones.
//
// To simplify the case, when apl expects each object to have idential shape
// as the previous without requireing a special rule, the option -u can be set.

import (
	"bufio"
	"fmt"
	"io"
)

// Iv is an instance of the stream interpreter. Call New to create one.
type Iv struct {
	apl     Apl
	rank    int
	uniform bool
	next    Nexter
}

// New creates a new stream interpreter for the given rank and an apl backend.
// For the default backend, a may be nil.
func New(rank int, uniform bool, apl string) (*Iv, error) {
	create := backends[apl]
	if create == nil {
		return nil, fmt.Errorf("unknown apl backend: %s", apl)
	}
	a, err := create(rank)
	if err != nil {
		return nil, err
	}
	return &Iv{
		apl:     a,
		rank:    rank,
		uniform: uniform,
	}, nil
}

// SetNext can be called to set a non-standard nexter.
// By default, the next value is read from a text format on stdin.
func (v *Iv) SetNext(n Nexter) {
	v.next = n
}

func (v *Iv) SetOut(stdout, stderr io.Writer) {
	v.apl.SetOut(stdout, stderr)
}

// Run is the iv main loop which is called after parsing APL rules.
func (v *Iv) Run() error {
	var R = v.rank                   // rank
	var shape []int = make([]int, R) // shape of the current array
	var array []Scalar               // current array
	var N = 1                        // number of executions, like awk NR.
	var E = 0                        // termination level
	var max int                      // max level observed in input
	var lastShape []int              // used for hard uniformity check

	if v.next == nil {
		return fmt.Errorf("iv needs to be initialized with SetNext")
	}

	execute := func(eof bool) error {
		defer func() {
			// reset state after execution
			copy(lastShape, shape)
			for i := range shape {
				shape[i] = 0
			}
			array = nil
			N++
			E = 0
		}()
		if lastShape == nil {
			lastShape = make([]int, R)
		} else if v.uniform {
			for i := range shape {
				if shape[i] != lastShape[i] {
					return fmt.Errorf("array #%d has different shape than previous: %v != %v", N, shape, lastShape)
				}
			}
		}
		return v.apl.Execute(shape, array, E, eof)
	}

	product := func() int { // 7 lines for a */ρ data
		p := 1
		for _, v := range shape {
			p *= v
		}
		return p
	}

	for i := 0; i < R-1; i++ {
		shape[i] = 1
	}
	for {
		scalar, S, eof, err := v.next.Next()
		if err != nil {
			return err
		}

		E = S - R
		if S > R {
			S -= R
		}

		if scalar != nil {
			// Catenate to current array and set shape.
			array = append(array, scalar)
			if idx := R - S; idx >= 0 && idx < R {
				shape[idx]++
			}

			// Check data uniformity.
			if product() != len(array) {
				return fmt.Errorf("array #%d is not uniform: prod(%v) != %d", N, shape, len(array))
			}
		}

		if eof {
			for i, x := range shape {
				if x == 0 {
					shape[i] = 1
				}
			}
			E = max
			return execute(true)
		}

		if E > max {
			max = E
		}
		if E > 0 {
			if err := execute(false); err != nil {
				return err
			}
		}
	}
}

// Nexter can return the next scalar from the data stream.
// The call to next should returns the scalar,
// the number of separators following it,
// if it is the last value (EOF)
// and a possible error different from io.EOF.
//
// See TextStream for the default implementation.
type Nexter interface {
	Next() (Scalar, int, bool, error)
}

// TextStream is a Nexter which read scalars in text format.
type TextStream struct {
	*bufio.Reader
	Parse     func(string) (Scalar, error) // Parse is a backend dependend scalar parser.
	Separator byte
	Rank      int
}

func (ts *TextStream) Next() (Scalar, int, bool, error) {
	var eof bool

	s, err := ts.Reader.ReadString(ts.Separator)
	if err != nil && err != io.EOF {
		return nil, 0, false, err
	} else if err == io.EOF {
		eof = true
	} else {
		s = s[:len(s)-1] // remove delimiter
	}

	a, err := ts.Parse(s)
	if err != nil {
		return nil, 0, false, err
	}

	if eof {
		return a, ts.Rank, true, nil
	}

	separators := 1
	for {
		if b, err := ts.Reader.ReadByte(); err == io.EOF {
			return a, ts.Rank, true, nil
		} else if err != nil {
			return nil, 0, false, err
		} else if b == ts.Separator {
			separators++
		} else {
			if err := ts.Reader.UnreadByte(); err != nil {
				return nil, 0, false, err
			}
			return a, separators, false, nil
		}
	}
}

// TextParser returns the ParseScalar method of the connected backend.
func (v *Iv) TextParser() func(string) (Scalar, error) {
	return v.apl.ParseScalar
}

// AddRule adds a new rule.
func (v *Iv) AddRule(conditional, statement string) error {
	return v.apl.AddRule(conditional, statement)
}

// Parse parses and executes the given content from a begin block
// or a library file. It is called before adding rules and
// reading input data.
func (v *Iv) Parse(s, name string) error {
	return v.apl.Parse(s, name)
}

// Register registers an APL backend under a given name.
func Register(name string, create func(int) (Apl, error)) {
	backends[name] = create
}

// Backends contains all registered backends.
var backends map[string]func(int) (Apl, error)

func init() {
	backends = make(map[string]func(int) (Apl, error))
	Register("", newAplIv) // Default backend: iv/apl.
}
