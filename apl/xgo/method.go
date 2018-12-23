package xgo

/* TODO remove
import (
	"fmt"
	"reflect"

	"github.com/ktye/iv/apl"
)

type Method struct {
	Value  reflect.Value
	Method string
}

func (m Method) String(a *apl.Apl) {
	fmt.Srintf("%v→%s", m.Value.Type(), m.Method)
}

func (m Method) Call(a *apl.Apl, L, R apl.Value) (apl.Value, error) {
	var zero reflect.Value
	if m.Value == zero {
		return nil, fmt.Errorf("method has no value")
	}
	fn := reflect.ValueOf(m.Value).MethodByName(m.Method)
	if fn == zero {
		return nil, fmt.Errorf("method %s does not exist", m.Method)
	}
	fn := Function{Name: m.Method, Fn: fn}
	return fn.Call(a, L, R)
}
*/
