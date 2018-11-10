package draw

import (
	"fmt"

	"github.com/ktye/iv/apl"
)

func Register(a *apl.Apl) error {
	must(a.Assign("canvas", Canvas{})) // canvas 0 returns new context
	must(a.Assign("p", Path{}))        // p A or p A B or p A B C (linear, quad or qubic element)
	//must(a.Assign("l", function(polyline)))  // l A B C D E ←→ p A p B p C p D p E
	must(a.Assign("c", function(closepath)))  // c p A p B p C (close path)
	must(a.Assign("f", function(fillpath)))   // f c p A p B p C (fill path)
	must(a.Assign("s", function(strokepath))) // s c p A p B p C (stroke path)
	a.RegisterPrimitive("+", handle(translate))
	// TODO ctx x Z: multiply by complex Z (scale and rotate)
	// TODO monadic ⍒ and ⍋for push/pop?
	// TODO color, linewidth, linedash, text

	return nil
}

type handle func(*apl.Apl, apl.Value, apl.Value) (bool, apl.Value, error)

func (h handle) HandlePrimitive(a *apl.Apl, l apl.Value, r apl.Value) (bool, apl.Value, error) {
	return h(a, l, r)
}

type function func(*apl.Apl, apl.Value, apl.Value) (apl.Value, error)

func (f function) Call(a *apl.Apl, l, r apl.Value) (apl.Value, error) {
	return f(a, l, r)
}

func (f function) String(a *apl.Apl) string {
	// TODO: This prints as draw.function but not which one.
	return fmt.Sprintf("%T", f)

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
