package comp

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewUpdater(id string) *Updater {
	return &Updater{id: id, spinner: defaultSpinner}
}

type Updater struct {
	id      string
	success bool
	err     error
	spinner spinner.Model
}

func (u *Updater) do() tea.Msg {
	if u.id == "" {
		return u.updateList()
	}
	return u.update()
}

func (u *Updater) updateList() tea.Msg {
	list, err := remote.GetList()
	if err != nil {
		return err
	}
	sps, err := local.WriteList(list)
	if err != nil {
		return err
	}
	return render.Successf("question list updated, there are %d questions now", len(sps))
}

func (u *Updater) update() tea.Msg {
	list, err := local.ReadList()
	if err != nil {
		return err
	}
	_, err = mgr.QueryRemote(list, u.id)
	if err != nil {
		return err
	}
	return true
}

func (u *Updater) Init() tea.Cmd {
	return tea.Batch(u.spinner.Tick, u.do)
}

func (u *Updater) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		u.err = msg
		return u, tea.Quit
	case bool:
		u.success = msg
		return u, tea.Quit
	default:
		var cmd tea.Cmd
		u.spinner, cmd = u.spinner.Update(msg)
		return u, cmd
	}
}

func (u *Updater) View() string {
	if u.err != nil {
		return render.Fail(fmt.Sprintf("%s\n", u.err))
	}
	if u.success {
		return fmt.Sprintf("%s\n", render.Success("Updated!"))
	}
	return fmt.Sprintf("Updating %s", u.spinner.View())
}
