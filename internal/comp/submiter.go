package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/remote"
)

type submiter struct {
	id      string
	spinner *spinner.Spinner
}

func (t *submiter) Run() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	t.spinner.Start()
	defer t.spinner.Stop()

	t.id = regualarID(t.id)

	typedCode, err := local.GetTypedCode(cfg, t.id)
	if err != nil {
		return err
	}
	question, err := local.GetQuestion(cfg, t.id)
	if err != nil {
		return err
	}
	id, err := remote.Submit(question, string(typedCode), config.LeetcodeLang(cfg.CodeLang))
	if err != nil {
		return err
	}
	log.Debug("submit id:", id)

	res := &model.SubmitCheckResult{}
	err = waitToCheck(id, question, res)
	if err != nil {
		return err
	}
	fmt.Println(res.Display())
	return nil
}
