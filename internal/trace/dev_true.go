//go:build dev

package trace

import (
	"fmt"
	"runtime"

	"github.com/zrcoder/leetgo/internal/render"
)

func init() {
	Wrap = func(err error) error {
		if err == nil {
			return nil
		}

		_, filename, line, _ := runtime.Caller(1)
		return fmt.Errorf("%s:%d\n%s", filename, line, render.Fail(err.Error()))
	}
}
