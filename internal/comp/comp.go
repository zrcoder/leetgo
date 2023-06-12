package comp

import (
	"time"

	"github.com/briandowns/spinner"
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
