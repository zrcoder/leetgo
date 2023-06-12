package comp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/leetgo/utils/render"
)

func NewSearcher(key string) Component {
	return &searcher{key: key, spinner: newSpinner("Searching")}
}

type searcher struct {
	key string

	spinner *spinner.Spinner
}

func (s *searcher) Run() error {
	var sps []model.StatStatusPair
	var err error
	s.spinner.Start()
	if s.key == "today" {
		today, err := remote.GetToday()
		if err != nil {
			return err
		}
		sps = []model.StatStatusPair{*today.ToStatePair()}
	} else {
		sps, err = search(s.key)
	}

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
		row := fmt.Sprintf(rowTmp, sp.Stat.FrontendID, sp.Stat.QuestionTitle, sp.Difficulty.String(), locked)
		buf.WriteString(row)
		lastQuestionID = sp.Stat.FrontendID
	}
	buf.WriteString(fmt.Sprintf("> total: %d, locked: %d\n", len(sps), lockCnt))
	buf.WriteString(fmt.Sprintf("> view detail? type like: `leetgo view %s`\n", lastQuestionID))
	fmt.Println(render.MarkDown(buf.String()))
	return nil
}

func search(keyWords string) ([]model.StatStatusPair, error) {
	list, err := remote.GetList()
	if err != nil {
		return nil, err
	}

	lower := strings.ToLower(keyWords)
	var res []model.StatStatusPair
	for _, sp := range list.StatStatusPairs {
		sp.Stat.FrontendID = sp.Stat.CalFrontendID()
		oriLower := strings.ToLower(sp.Stat.QuestionTitle)
		for _, sep := range []string{" ", ". ", "."} {
			title := sp.Stat.FrontendID + sep + oriLower
			if strings.Contains(title, lower) {
				res = append(res, sp)
				break
			}
		}
	}
	if len(res) == 0 {
		log.Debug("no questions found")
		return nil, fmt.Errorf("no questions found for `%s`", keyWords)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Stat.FrontendID < res[j].Stat.FrontendID
	})
	return res, nil
}
