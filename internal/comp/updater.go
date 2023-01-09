package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewUpdater(id string) Component {
	return &Updater{id: id, spinner: newSpinner("Updating")}
}

type Updater struct {
	id      string
	spinner *spinner.Spinner
}

func (u *Updater) Run() error {
	if u.id == "" {
		return u.updateList()
	}
	return u.update()
}

func (u *Updater) updateList() error {
	u.spinner.Start()
	list, err := remote.GetList()
	if err != nil {
		u.spinner.Stop()
		return err
	}
	sps, err := local.WriteList(list)
	if err != nil {
		u.spinner.Stop()
		return err
	}
	u.spinner.Stop()
	fmt.Println(render.Successf("question list updated, there are %d questions now", len(sps)))
	return nil
}

func (u *Updater) update() error {
	u.spinner.Start()
	list, err := local.ReadList()
	if err != nil {
		u.spinner.Stop()
		return err
	}
	question, err := mgr.QueryRemote(list, u.id)
	if err != nil {
		u.spinner.Stop()
		return err
	}

	u.spinner.Stop()
	fmt.Println(render.Success("Done!"))
	fmt.Println(render.MarkDown(question.MdContent))
	return nil
}
