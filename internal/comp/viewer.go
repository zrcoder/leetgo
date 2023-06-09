package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewViewer(id string) Component {
	return &viewer{id: id, spinner: newSpinner("Picking")}
}

type viewer struct {
	id string

	spinner *spinner.Spinner
}

func (c *viewer) Run() error {
	if local.Exist(c.id) {
		content, err := local.ReadMarkdown(c.id)
		if err != nil {
			return err
		}
		fmt.Print(render.MarkDown(string(content)))
		return nil
	}

	c.spinner.Start()
	question, err := mgr.Query(c.id)
	c.spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Print(render.MarkDown(question.MdContent))
	return local.Write(question)
}
