package comp

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewSearcher(key string) *Searcher {
	return &Searcher{key: key, spinner: defaultSpinner}
}

type Searcher struct {
	key string

	err     error
	sps     []model.StatStatusPair
	spinner spinner.Model
}

func (s *Searcher) do() tea.Msg {
	sps, err := mgr.Search(s.key)
	if err != nil {
		return err
	}
	return sps
}

func (s *Searcher) Init() tea.Cmd {
	return tea.Batch(s.spinner.Tick, s.do)
}

func (s *Searcher) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		s.err = msg
		return s, tea.Quit
	case []model.StatStatusPair:
		s.sps = msg
		return s, tea.Quit
	default:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}
}

func (s *Searcher) View() string {
	if s.err != nil {
		return render.Fail(fmt.Sprintf("%s\n", s.err))
	}
	if s.sps != nil {
		buf := &strings.Builder{}
		buf.WriteString("| No. | Title | Difficulty | Locked |\n")
		buf.WriteString("| --- | ----- | ---------- | ------ |\n")
		rowTmp := "| %s  | %s    | %s         | %s     |\n"
		lockCnt := 0
		lastQuestionID := ""
		for _, sp := range s.sps {
			locked := ""
			if sp.PaidOnly {
				locked = "🔒"
				lockCnt++
			}
			row := fmt.Sprintf(rowTmp, sp.Stat.CalculatedID, sp.Stat.QuestionTitle, sp.Difficulty.String(), locked)
			buf.WriteString(row)
			lastQuestionID = sp.Stat.CalculatedID
		}
		buf.WriteString(fmt.Sprintf("> total: %d, locked: %d\n\n", len(s.sps), lockCnt))
		buf.WriteString(fmt.Sprintf("pick one? type like `leetgo pick %s`\n", lastQuestionID))
		return render.MarkDown(buf.String())
	}
	return fmt.Sprintf("Searching %s", s.spinner.View())
}
