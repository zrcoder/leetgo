package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/remote"
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

	v.id = meta.Stat.FrontendID // origin id may be "today"
	exist, err = v.checkLocal()
	if exist || err != nil {
		return err
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
	v.printHint()
	return nil
}

func (v *singleViewer) printHint() {
	typeHint := fmt.Sprintf("Type `leetgo code %s` to solve it", v.id)
	fmt.Println(render.MarkDown(typeHint))
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
		v.printHint()
	}
	return
}
