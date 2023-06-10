package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/exec"
	"github.com/zrcoder/leetgo/internal/render"

	//	"github.com/zrcoder/leetgo/internal/exec"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/remote"
	//"github.com/zrcoder/leetgo/internal/render"
)

func NewTester(id string) Component {
	return &tester{
		id:      id,
		spinner: newSpinner("remote testing..."),
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
		fmt.Println(render.Info("local testing..."))
		err = exec.Run(local.GetDir(cfg, t.id), "go", "test", ".")
		if err != nil {
			return err
		}
		fmt.Println(render.Info("succeed"))
	}

	t.spinner.Start() // remote testing...
	defer t.spinner.Stop()

	typedCode, err := local.GetTypedCode(cfg, t.id)
	if err != nil {
		return err
	}
	question, err := local.GetQuestion(cfg, t.id)
	if err != nil {
		return err
	}
	resp, err := remote.Test(question, string(typedCode), cfg.CodeLang, question.ParseDefaultTests())
	if err != nil {
		return err
	}
	for {
		res, err := remote.Check(resp.InterpretId)
		if err != nil {
			return err
		}
		fmt.Println(res)
		// if res.State == "SUCCESS" {
		// todo
		return nil
		// }
		// time.Sleep(time.Second)
	}
}
