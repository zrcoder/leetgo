//go:build dev

package log

import (
	"fmt"
	"runtime"

	"github.com/zrcoder/leetgo/internal/render"
)

func init() {
	Trace = func(x ...any) {
		if len(x) == 0 || x[0] == nil {
			return
		}
		_, filename, line, _ := runtime.Caller(1)
		fmt.Print(render.Trace(fmt.Sprintf("\n[TRACE] %s:%d ", filename, line)))
		fmt.Println(x...)
	}
	Tracef = func(format string, x ...any) {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Print(render.Trace(fmt.Sprintf("\n[TRACE] %s:%d ", filename, line)))
		fmt.Printf(format, x...)
	}
}
