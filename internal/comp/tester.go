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
	res, err := remote.Test(question, string(typedCode), config.LeetcodeLang(cfg.CodeLang), casesStr)
	if err != nil {
		log.Debug(err)
		return err
	}
	for {
		checkres, err := remote.CheckTest(res.InterpretId, question.Slug)
		if err != nil {
			log.Debug(err)
			return err
		}
		log.Debug("state:", checkres.State, checkres.StatusCode)
		if checkres.State == "FAILURE" {
			log.Debug("test failed for id:", res.InterpretId)
			data, _ := json.MarshalIndent(checkres, "", " . ")
			log.Debug(string(data))
			return errors.New("unknow internal error")
		}
		if checkres.State == "SUCCESS" {
			checkres.InputData = casesStr
			fmt.Println(checkres.Display())
			return nil
		}
		time.Sleep(time.Second)
	}
}
