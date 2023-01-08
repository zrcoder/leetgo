package mgr

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
)

func Query(id string) (*model.Question, error) {
	question, err := local.Read(id)
	if err == nil {
		log.Trace("got markdown data from local")
		return question, nil
	}

	if err != local.ErrNotCached {
		return nil, err
	}

	list, err := getList()
	if err != nil {
		log.Trace(err)
		return nil, err
	}

	return QueryRemote(list, id)
}

func QueryRemote(statStatusPairs map[string]model.StatStatusPair, id string) (*model.Question, error) {
	sp, ok := statStatusPairs[id]
	if !ok {
		err := fmt.Errorf("not found by id: %s", id)
		log.Trace(err)
		return nil, err
	}
	if sp.PaidOnly {
		err := fmt.Errorf("[%s. %s] is locked", sp.Stat.CalculatedID, sp.Stat.QuestionTitle)
		log.Trace(err)
		return nil, err
	}

	question, err := remote.GetQuestion(&sp)
	if err != nil {
		return nil, err
	}

	err = local.Write(sp.Stat.CalculatedID, question)
	log.Trace(err)
	return question, err
}

func Search(keyWords string) ([]model.StatStatusPair, error) {
	list, err := getList()
	if err != nil {
		return nil, err
	}

	lower := strings.ToLower(keyWords)
	var res []model.StatStatusPair
	for _, sp := range list {
		title1 := sp.Stat.CalculatedID + " " + strings.ToLower(sp.Stat.QuestionTitle)
		title2 := sp.Stat.CalculatedID + ". " + strings.ToLower(sp.Stat.QuestionTitle)
		if strings.Contains(title1, lower) || strings.Contains(title2, lower) {
			res = append(res, sp)
		}
	}
	if len(res) == 0 {
		log.Trace("no questions found")
		return nil, fmt.Errorf("no questions found by keywords: %s", keyWords)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Stat.CalculatedID < res[j].Stat.CalculatedID
	})
	return res, nil
}

func getList() (map[string]model.StatStatusPair, error) {
	res, err := local.ReadList()
	if err == nil {
		return res, nil
	}
	if err != local.ErrNotCached {
		return nil, err
	}

	var list *model.List
	list, err = remote.GetList()
	if err != nil {
		return nil, err
	}

	return local.WriteList(list)
}
