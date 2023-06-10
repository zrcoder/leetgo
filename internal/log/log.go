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
	log.Helper()
	log.Debug(x)
}
