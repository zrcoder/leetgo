package log

import (
	"github.com/charmbracelet/log"

	"github.com/zrcoder/leetgo/internal/mod"
)

func init() {
	log.Default().SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}

func Debug(x ...any) {
	if !mod.IsDebug() {
		return
	}
	if len(x) == 0 || x[0] == nil {
		return
	}
	log.Helper()
	log.Debug(x)
}
