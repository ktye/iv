package apl

type Object interface {
	Field(*Apl, string) Value
	Set(*Apl, string, Value) error
}
