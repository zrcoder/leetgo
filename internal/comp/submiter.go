package comp

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/remote"
)

func NewSubmiter(id string) Component {
	return &submiter{
		id:      id,
		spinner: newSpinner("Remote testing"),
	}
}

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
		log.Debug(err)
		return err
	}
	log.Debug("submit id:", id)

	for {
		checkres, err := remote.CheckSubmit(id, question.Slug)
		if err != nil {
			log.Debug(err)
			return err
		}
		log.Debug("state:", checkres.State, checkres.StatusCode)
		if checkres.State == "FAILURE" {
			log.Debug("test failed for id:", id)
			data, _ := json.MarshalIndent(checkres, "", " . ")
			log.Debug(string(data))
			return errors.New("unknow internal error")
		}
		if checkres.State == "SUCCESS" {
			fmt.Println(checkres.Display())
			return nil
		}
		time.Sleep(time.Second)
	}
}
