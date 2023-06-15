package comp

import (
	"fmt"
	"strings"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/utils/render"
)

type searcher struct {
	key string

	spinner *spinner.Spinner
}

func (s *searcher) Run() error {
	s.spinner.Start()
	metaList, err := queryMetas(s.key)
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
	for _, meta := range metaList {
		locked := ""
		if meta.PaidOnly {
			locked = "🔒"
			lockCnt++
		}
		meta.Transform()
		row := fmt.Sprintf(rowTmp, meta.FrontendID, meta.Title, meta.Difficulty, locked)
		buf.WriteString(row)
		lastQuestionID = meta.FrontendID
	}
	buf.WriteString(fmt.Sprintf("> total: %d, locked: %d\n", len(metaList), lockCnt))
	buf.WriteString(fmt.Sprintf("> view detail? type like: `leetgo view %s`\n", lastQuestionID))
	fmt.Println(render.MarkDown(buf.String()))
	return nil
}
