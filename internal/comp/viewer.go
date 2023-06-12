package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/utils/render"
)

func NewViewer(id string) Component {
	return &viewer{id: id, spinner: newSpinner("Picking")}
}

type viewer struct {
	id string

	spinner *spinner.Spinner
}

func (c *viewer) Run() error {
	printHint := func() {
		typeHint := fmt.Sprintf("Type `leetgo edit %s` to solve it", c.id)
		fmt.Println(render.MarkDown(typeHint))
	}
	if local.Exist(c.id) {
		content, err := local.GetMarkdown(c.id)
		if err != nil {
			return err
		}
		fmt.Print(render.MarkDown(string(content)))
		printHint()
		return nil
	}

	c.spinner.Start()
	question, err := mgr.Query(c.id)
	c.spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Print(render.MarkDown(question.MdContent))

	err = local.Write(question)
	if err != nil {
		return err
	}
	printHint()
	return nil
}
