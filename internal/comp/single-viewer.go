package comp

import (
	"errors"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/tdoc"
	tmodel "github.com/zrcoder/tdoc/model"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/remote"
)

type singleViewer struct {
	spinner  *spinner.Spinner
	id       string
	showGlow bool
}

func (v *singleViewer) Run() error {
	exist, err := v.localAction()
	if exist || err != nil {
		return err
	}

	meta, err := queryMeta(v.id)
	if err != nil {
		return err
	}
	v.id = meta.FrontendID // origin id may be "today"
	exist, err = v.localAction()
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

	err = local.Write(question)
	if err != nil {
		return err
	}

	if v.showGlow {
		err = v.show()
		if err != nil {
			return err
		}
		return askToCode(v.id)
	}
	return code(v.id)
}

func (v *singleViewer) localAction() (exist bool, err error) {
	if local.Exist(v.id) {
		if v.showGlow {
			err = v.show()
			if err != nil {
				return
			}
			return true, askToCode(v.id)
		}
		return true, code(v.id)
	}
	return
}

func (v *singleViewer) show() error {
	doc, err := getDocFromLocal(v.id)
	if err != nil {
		return err
	}
	return tdoc.Run([]*tmodel.DocInfo{doc})
}
