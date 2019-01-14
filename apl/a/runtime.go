package a

import (
	"reflect"
	"runtime"

	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/xgo"
)

// Memstats returns runtime.MemStats as an object.
func Memstats(p *apl.Apl, L, R apl.Value) (apl.Value, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return xgo.Convert(p, reflect.ValueOf(m))
}

func cpus(p *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return apl.Index(runtime.NumCPU()), nil
}

func goroutines(p *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return apl.Index(runtime.NumGoroutine()), nil
}

func goversion(p *apl.Apl, _, R apl.Value) (apl.Value, error) {
	return apl.String(runtime.Version()), nil
}
