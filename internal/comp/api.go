package comp

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/leetgo/internal/config"
)

type Component interface {
	Run() error
}

func NewConfiger(cfg *config.Config, shouldWrite bool, showFunc func(*config.Config)) Component {
	return &configer{
		cfg:         cfg,
		shouldWrite: shouldWrite,
		showFunc:    showFunc,
	}
}

func NewSearcher(key string) Component {
	return &searcher{key: key, spinner: newSpinner("Searching")}
}

func NewSingleViewer(id string) Component {
	return &singleViewer{id: id, spinner: newSpinner("Picking")}
}

func NewSubmiter(id string) Component {
	return &submiter{
		id:      id,
		spinner: newSpinner("Submitting"),
	}
}

func NewTester(id string) Component {
	return &tester{
		id:      id,
		spinner: newSpinner("Remote testing"),
	}
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
