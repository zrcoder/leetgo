package mgr

import (
	"fmt"
	"sort"
	"strings"

	h2md "github.com/JohannesKaufmann/html-to-markdown"

	"github.com/zrcoder/leetgo/internal/client"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

func Query(id string) ([]byte, string, error) {
	mdData, path, err := local.Read(id)
	if err == nil {
		return mdData, path, nil
	}

	if err != local.ErrNotCached {
		return nil, "", err
	}

	all, err := getAll()
	if err != nil {
		log.Dev(err)
		return nil, "", err
	}

	sp, ok := all[id]
	if !ok {
		err = fmt.Errorf("not found by id: %s", id)
		log.Dev(err)
		return nil, "", err
	}
	if sp.PaidOnly {
		err = fmt.Errorf("[%s. %s] is locked!", sp.Stat.CalculatedID, sp.Stat.QuestionTitle)
		log.Dev(err)
		return nil, "", err
	}

	question, err := client.Get(&sp)
	if err != nil {
		return nil, "", err
	}

	data, err := parseMarkdown(question)
	if err != nil {
		return nil, "", err
	}

	path, _ = local.Write(&sp, question, []byte(data))
	return []byte(data), path, nil
}

func Search(keyWords string) ([]model.StatStatusPair, error) {
	all, err := getAll()
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

func getAll() (map[string]model.StatStatusPair, error) {
	res, err := local.ReadAll()
	if err == nil {
		return res, nil
	}
	if err != local.ErrNotCached {
		return nil, err
	}

	var all *model.All
	all, err = client.GetAll()
	if err != nil {
		return nil, err
	}

	res = make(map[string]model.StatStatusPair, len(all.StatStatusPairs))
	for _, sp := range all.StatStatusPairs {
		id := sp.Stat.GetFrontendQuestionID()
		sp.Stat.CalculatedID = id
		res[id] = sp
	}

	err = local.WriteAll(res)
	log.Dev(err)

	return res, nil
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
