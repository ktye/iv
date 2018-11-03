package apl

type Value interface {
	String(*Apl) string
	Eval(*Apl) (Value, error)
	// TODO...
}
