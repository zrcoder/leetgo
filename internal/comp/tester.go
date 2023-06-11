package comp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/exec"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/render"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/remote"
)

func NewTester(id string) Component {
	return &tester{
		id:      id,
		spinner: newSpinner("Remote testing"),
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
		fmt.Println("local testing...")
		err = exec.Run(local.GetDir(cfg, t.id), "go", "test", ".")
		if err != nil {
			return err
		}
		fmt.Println(render.Info("local test succeed"))
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
	cases := question.ParseDefaultTests()
	casesStr := strings.Join(cases, "\n")
	id, err := remote.Test(question, string(typedCode), config.LeetcodeLang(cfg.CodeLang), casesStr)
	if err != nil {
		log.Debug(err)
		return err
	}
	res := &model.TestCheckResult{}
	err = waitToCheck(id, question, res)
	if err != nil {
		return err
	}

	res.InputData = casesStr
	fmt.Println(res.Display())
	return nil
}

func waitToCheck(id string, question *model.Question, res model.RunResult) error {
	var err error
	for {
		err = remote.CheckResult(id, question, res)
		if err != nil {
			log.Debug(err)
			return err
		}
		state := res.Result()
		log.Debug("result:", state)
		if state == "FAILURE" {
			log.Debug("test failed for id:", id)
			data, _ := json.MarshalIndent(res, "", " . ")
			log.Debug(string(data))
			return errors.New("unknow internal error")
		}
		if state == "SUCCESS" {
			return nil
		}
		time.Sleep(time.Second)
	}
}
