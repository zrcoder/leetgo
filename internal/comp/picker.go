package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewPicker(id string) Component {
	return &Picker{id: id, spinner: newSpinner("Picking")}
}

type Picker struct {
	id string

	spinner *spinner.Spinner
}

func (c *Picker) Run() error {
	c.spinner.Start()
	path, question, err := mgr.Query(c.id)
	c.spinner.Stop()

	if err != nil {
		return err
	}

	fmt.Print(render.MarkDown(question.MdContent)) // question.MdContent has "\n\n" suffix
	fmt.Println(render.Info(fmt.Sprintf("  Stored in: %s\n", path)))
	return nil
}
