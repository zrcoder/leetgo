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
	c.spinner.Start()
	question, err := mgr.Query(c.id)
	c.spinner.Stop()

	if err != nil {
		return err
	}

	fmt.Print(render.MarkDown(question.MdContent)) // question.MdContent has "\n\n" suffix

	return local.Write(question)
}
