package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/trace"
)

func NewClient() *http.Client {
	const maxRedirect = 200
	return &http.Client{
		CheckRedirect: func() func(req *http.Request, via []*http.Request) error {
			redirects := 0
			return func(req *http.Request, via []*http.Request) error {
				if redirects > maxRedirect {
					err := fmt.Errorf("stopped after %d redirects", maxRedirect)
					return trace.Wrap(err)
				}
				redirects++
				return nil
			}
		}(),
	}
}

func Query(id string) (*model.Question, error) {
	all, err := GetAll()
	if err != nil {
		return nil, err
	}
	for _, sp := range all.StatStatusPair {
		if sp.Stat.GetFrontendQuestionID() == id {
			return get(&sp)
		}
	}
	err = trace.Wrap(fmt.Errorf("not found by id: %s", id))
	return nil, err
}

func Search(keyWord string) ([]*model.StatStatusPair, error) {
	all, err := GetAll()
	if err != nil {
		return nil, err
	}

	lower := strings.ToLower(keyWord)

	var res []*model.StatStatusPair
	for _, sp := range all.StatStatusPair {
		if strings.Contains(strings.ToLower(sp.Stat.QuestionTitle), lower) {
			res = append(res, &sp)
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no questions found by keyword: %s", keyWord)
	}
	return res, nil
}

func GetAll() (*model.All, error) {
	url := fmt.Sprintf("%s/api/problems/all", config.Domain())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &model.All{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return res, nil
}

func get(sp *model.StatStatusPair) (*model.Question, error) {
	query := `
query getQuestionDetail($titleSlug: String!) {
  isCurrentUserAuthenticated
  question(titleSlug: $titleSlug) {
    questionId
    content
    stats
    codeDefinition
    sampleTestCase
    enableRunCode
    metaData
    translatedContent
  }
}`

	body := map[string]any{
		"query": query,
		"variables": map[string]string{
			"titleSlug": sp.Stat.QuestionTitleSlug,
		},
		"operationName": "getQuestionDetail",
	}
	jbody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	domain := config.Domain()
	url := fmt.Sprintf("%s/graphql", domain)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jbody))
	if err != nil {
		return nil, err
	}
	referer := fmt.Sprintf("%s/problems/%s", domain, sp.Stat.QuestionTitleSlug)
	req.Header.Set("Referer", referer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetQuestion got status %d", resp.StatusCode)
	}

	res := &model.GetQuestionResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	question := res.Data.Question
	question.FrontendQuestionID = sp.Stat.GetFrontendQuestionID()
	question.Title = sp.Stat.QuestionTitle
	question.Referer = referer
	question.Difficulty = sp.Difficulty.String()
	return question, nil
}
