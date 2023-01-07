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
	Logger = loggerDev{}
}

type loggerDev struct{}

func (l loggerDev) Errorf(s string, i ...any) {
	fmt.Print(render.Fail("[Error]: "))
	fmt.Printf(s, i...)
}

func (l loggerDev) Warningf(s string, i ...any) {
	fmt.Print(render.Fail("[Warn]: "))
	fmt.Printf(s, i...)
}

func (l loggerDev) Infof(s string, i ...any) {
	fmt.Print("[Info]: ")
	fmt.Printf(s, i...)
}

func (l loggerDev) Debugf(s string, i ...any) {
	fmt.Print(render.Debug("[Debug]: "))
	fmt.Printf(s, i...)
}
