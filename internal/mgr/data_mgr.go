package mgr

import (
	"fmt"
	"sort"
	"strings"

	h2md "github.com/JohannesKaufmann/html-to-markdown"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
)

func Query(id string) ([]byte, string, error) {
	mdData, path, err := local.Read(id)
	if err == nil {
		log.Dev("got markdown data from local")
		return mdData, path, nil
	}

	if err != local.ErrNotCached {
		return nil, "", err
	}

	list, err := getList()
	if err != nil {
		log.Dev(err)
		return nil, "", err
	}

	return QueryRemote(list, id)
}

func QueryRemote(statStatusPairs map[string]model.StatStatusPair, id string) ([]byte, string, error) {
	sp, ok := statStatusPairs[id]
	if !ok {
		err := fmt.Errorf("not found by id: %s", id)
		log.Dev(err)
		return nil, "", err
	}
	if sp.PaidOnly {
		err := fmt.Errorf("[%s. %s] is locked", sp.Stat.CalculatedID, sp.Stat.QuestionTitle)
		log.Dev(err)
		return nil, "", err
	}

	question, err := remote.GetQuestion(&sp)
	if err != nil {
		return nil, "", err
	}

	data, err := parseMarkdown(question)
	if err != nil {
		return nil, "", err
	}

	path, err := local.Write(&sp, question, []byte(data))
	log.Dev(err)
	return []byte(data), path, err
}

func Search(keyWords string) ([]model.StatStatusPair, error) {
	all, err := getList()
	if err != nil {
		return nil, err
	}
	lower := strings.ToLower(keyWords)

	var res []model.StatStatusPair
	for _, sp := range all {
		title1 := sp.Stat.CalculatedID + " " + strings.ToLower(sp.Stat.QuestionTitle)
		title2 := sp.Stat.CalculatedID + ". " + strings.ToLower(sp.Stat.QuestionTitle)
		if strings.Contains(title1, lower) || strings.Contains(title2, lower) {
			res = append(res, sp)
		}
	}
	if len(res) == 0 {
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

func parseMarkdown(question *model.Question) (string, error) {
	content := question.TranslatedContent
	if content == "" {
		content = question.Content
	}
	content = strings.ReplaceAll(content, "<sup>", "^")
	content = strings.ReplaceAll(content, "</sup>", "")
	content, err := h2md.NewConverter("", true, nil).ConvertString(content)
	if err != nil {
		log.Dev(err)
		return "", err
	}
	content = strings.ReplaceAll(content, `\[`, `[`)
	content = strings.ReplaceAll(content, `\]`, `]`)
	md := fmt.Sprintf("## [%s. %s](%s) (%s)\n\n%s",
		question.Id, question.Title, question.Referer, question.Difficulty, content)
	return md, nil
}
