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
		fmt.Print(render.Trace(fmt.Sprintf("[TRACE] %s:%d ", filename, line)))
		fmt.Println(x...)
	}
	Tracef = func(format string, x ...any) {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Print(render.Trace(fmt.Sprintf("[TRACE] %s:%d ", filename, line)))
		fmt.Printf(format, x...)
	}
	Logger = loggerDev{}
}

type loggerDev struct{}

func (l loggerDev) Errorf(format string, v ...any) {
	fmt.Println(render.Errorf("[ERROR] %s", fmt.Sprintf(format, v...)))
}

func (l loggerDev) Warningf(format string, v ...any) {
	fmt.Println(render.Warnf("[ WARN] %s", fmt.Sprintf(format, v...)))
}

func (l loggerDev) Infof(format string, v ...any) {
	fmt.Println(render.Infof("[ INFO] %s", fmt.Sprintf(format, v...)))
}

func (l loggerDev) Debugf(format string, v ...any) {
	fmt.Println(render.Tracef("[DEBUG] %s", fmt.Sprintf(format, v...)))
}
