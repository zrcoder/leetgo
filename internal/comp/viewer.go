package comp

import (
	"errors"
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
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
	isToday := c.id == "today"
	var today model.Today
	var err error
	if isToday {
		today, err = remote.GetToday()
		if err != nil {
			return err
		}
		c.id = today.FrontendID()
	}

	printHint := func() {
		typeHint := fmt.Sprintf("Type `leetgo code %s` to solve it", c.id)
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

	var question *model.Question
	c.spinner.Start()
	if isToday {
		question, err = remote.GetQuestion(today.ToStatePair()) // faster than query(c.id)
	} else {
		question, err = query(c.id)
	}
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

func query(frontendID string) (*model.Question, error) {
	list, err := remote.GetList()
	if err != nil {
		return nil, err
	}
	for _, sp := range list.StatStatusPairs {
		sp.Stat.FrontendID = sp.Stat.CalFrontendID()
		if sp.Stat.FrontendID != frontendID {
			continue
		}
		if sp.PaidOnly {
			err := fmt.Errorf("[%s. %s] is locked", sp.Stat.FrontendID, sp.Stat.QuestionTitle)
			log.Debug(err)
			return nil, err
		}
		return remote.GetQuestion(&sp)
	}
	return nil, errors.New("question not found")
}
