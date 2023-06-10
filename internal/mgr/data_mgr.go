package mgr

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
)

func Query(id string) (*model.Question, error) {
	list, err := remote.GetList()
	if err != nil {
		return nil, err
	}
	for _, sp := range list.StatStatusPairs {
		sp.Stat.CalculatedID = sp.Stat.GetFrontendQuestionID()
		if sp.Stat.CalculatedID != id {
			continue
		}
		if sp.PaidOnly {
			err := fmt.Errorf("[%s. %s] is locked", sp.Stat.CalculatedID, sp.Stat.QuestionTitle)
			log.Debug(err)
			return nil, err
		}
		return remote.GetQuestion(&sp)
	}
	return nil, errors.New("question not found")
}

func Search(keyWords string) ([]model.StatStatusPair, error) {
	list, err := remote.GetList()
	if err != nil {
		return nil, err
	}

	lower := strings.ToLower(keyWords)
	var res []model.StatStatusPair
	for _, sp := range list.StatStatusPairs {
		sp.Stat.CalculatedID = sp.Stat.GetFrontendQuestionID()
		oriLower := strings.ToLower(sp.Stat.QuestionTitle)
		for _, sep := range []string{" ", ". ", "."} {
			title := sp.Stat.CalculatedID + sep + oriLower
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
		return res[i].Stat.CalculatedID < res[j].Stat.CalculatedID
	})
	return res, nil
}
