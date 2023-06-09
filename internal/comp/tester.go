package comp

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/exec"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewTester(id string) Component {
	return &tester{
		id:      id,
		spinner: newSpinner("testing remote..."),
	}
}

type tester struct {
	id      string
	spinner *spinner.Spinner
}

func (t *tester) Run() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	// only surpport "go" to run test locally now
	if config.IsGolang(cfg) {
		fmt.Println(render.Info("begin to run tests locally"))
		err = exec.Run(local.GetDir(cfg, t.id), "go", "test", ".")
		if err != nil {
			return err
		}
	}
	t.spinner.Start()

	t.spinner.Stop()

	return nil
}
