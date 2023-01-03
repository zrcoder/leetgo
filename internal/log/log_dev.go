//go:build dev

package log

import (
	"fmt"
	"runtime"

	"github.com/zrcoder/leetgo/internal/render"
)

func init() {
	Dev = func(x ...any) {
		if len(x) == 0 || x[0] == nil {
			return
		}
		_, filename, line, _ := runtime.Caller(1)
		fmt.Print(render.Debug(fmt.Sprintf("[Dev] %s:%d ", filename, line)))
		fmt.Println(x...)
	}
}
