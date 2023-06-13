package comp

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/remote"
)

type Component interface {
	Run() error
}

func newSpinner(msg string) *spinner.Spinner {
	res := spinner.New([]string{
		"▱▱▱",
		"▰▱▱",
		"▰▰▱",
		"▰▰▰",
		"▰▰▱",
		"▰▱▱",
		"▱▱▱",
	}, time.Second/7, spinner.WithHiddenCursor(true))
	res.Prefix = msg + " "
	return res
}

func regualarID(id string) string {
	if id != "today" {
		return id
	}

	spinner := newSpinner("Fetching today")
	spinner.Start()
	today, err := remote.GetToday()
	spinner.Stop()
	if err != nil {
		log.Debug(err)
		return id
	}
	return today.Question().FrontendID
}
