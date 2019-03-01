package apl

// Error carries an error value.
// It is used by go routines to signal errors.
// To send err over Channel c, use: c[0]<-Error{e}
type Error struct {
	E error
}

func (e Error) String(f Format) string {
	if e.E == nil {
		return "<nil error>"
	}
	return e.E.Error()
}
func (e Error) Copy() Value { return e }
