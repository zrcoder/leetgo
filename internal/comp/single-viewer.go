package comp

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/leetgo/utils/exec"
	"github.com/zrcoder/leetgo/utils/render"
)

type singleViewer struct {
	id string

	spinner *spinner.Spinner
}

func (v *singleViewer) Run() error {
	exist, err := v.checkLocal()
	if exist || err != nil {
		return err
	}

	meta, err := queryMeta(v.id)
	if err != nil {
		return err
	}
	v.id = meta.FrontendID // origin id may be "today"
	exist, err = v.checkLocal()
	if exist || err != nil {
		return err
	}

	if meta.PaidOnly {
		return errors.New("🔒 the question is locked")
	}

	v.spinner.Start()
	question, err := remote.GetQuestion(meta)
	v.spinner.Stop()
	if err != nil {
		return err
	}

	// TODO with tdoc
	fmt.Print(render.MarkDown(question.MdContent))

	err = local.Write(question)
	if err != nil {
		return err
	}

	return v.askToCode()
}

func (v *singleViewer) askToCode() error {
	prompt := &survey.Confirm{
		Message: "Solve the question now?",
		Default: true,
		Help:    "Open the local code file with your favorite editor to edit the code.",
	}

	code := true
	err := survey.AskOne(prompt, &code)
	if err != nil {
		return err
	}

	if code {
		return v.code()
	}
	return nil
}

func (v *singleViewer) checkLocal() (exist bool, err error) {
	if local.Exist(v.id) {
		exist = true
		var content []byte
		content, err = local.GetMarkdown(v.id)
		if err != nil {
			return
		}
		// TODO with tdoc
		fmt.Print(render.MarkDown(string(content)))
		return true, v.askToCode()
	}
	return
}

func (v singleViewer) code() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	codeFile := local.GetCodeFile(cfg, v.id)
	cmd, ops := config.GetEditorCmdOps(cfg.Editor)
	args := append(ops, codeFile)
	if config.IsGolang(cfg) {
		args = append(args, local.GetGoTestFile(cfg, v.id))
	}
	return exec.Run("", cmd, args...)
}
