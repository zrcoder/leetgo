package comp

import (
	"fmt"
	"strings"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewSearcher(key string) Component {
	return &searcher{key: key, spinner: newSpinner("Searching")}
}

type searcher struct {
	key string

	spinner *spinner.Spinner
}

func (s *searcher) Run() error {
	s.spinner.Start()
	sps, err := mgr.Search(s.key)
	s.spinner.Stop()

	if err != nil {
		return err
	}

	buf := &strings.Builder{}
	buf.WriteString("| No. | Title | Difficulty | Locked |\n")
	buf.WriteString("| --- | ----- | ---------- | ------ |\n")
	rowTmp := "| %s  | %s    | %s         | %s     |\n"
	lockCnt := 0
	lastQuestionID := ""
	for _, sp := range sps {
		locked := ""
		if sp.PaidOnly {
			locked = "🔒"
			lockCnt++
		}
		row := fmt.Sprintf(rowTmp, sp.Stat.CalculatedID, sp.Stat.QuestionTitle, sp.Difficulty.String(), locked)
		buf.WriteString(row)
		lastQuestionID = sp.Stat.CalculatedID
	}
	buf.WriteString(fmt.Sprintf("> total: %d, locked: %d\n", len(sps), lockCnt))
	buf.WriteString(fmt.Sprintf("> view detail? type like: `leetgo view %s`\n", lastQuestionID))
	fmt.Println(render.MarkDown(buf.String()))
	return nil
}
