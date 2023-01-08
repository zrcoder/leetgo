package comp

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewPicker(id string) *Picker {
	return &Picker{id: id, spinner: defaultSpinner}
}

type Picker struct {
	id string

	question *model.Question
	err      error
	spinner  spinner.Model
}

func (c *Picker) do() tea.Msg {
	question, err := mgr.Query(c.id)
	if err != nil {
		return err
	}
	return question
}

func (c *Picker) Init() tea.Cmd {
	return tea.Batch(c.spinner.Tick, c.do)
}

func (c *Picker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		c.err = msg
		return c, tea.Quit
	case *model.Question:
		c.question = msg
		return c, tea.Quit
	default:
		var cmd tea.Cmd
		c.spinner, cmd = c.spinner.Update(msg)
		return c, cmd
	}
}

func (c *Picker) View() string {
	if c.err != nil {
		// must append \n, or not show
		return render.Fail(fmt.Sprintf("%s\n", c.err))
	}
	if c.question != nil {
		return render.MarkDown(c.question.MdContent)
	}
	return fmt.Sprintf("Picking %s", c.spinner.View())
}
