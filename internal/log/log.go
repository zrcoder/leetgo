package log

import (
	"os"

	"github.com/charmbracelet/log"
)

var isDebug = false

func init() {
	isDebug = os.Getenv("LG_DEBUG") == "1"
	if isDebug {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}
}

func Debug(x ...any) {
	if !isDebug {
		return
	}
	log.Helper()
	log.Debug(x)
}
