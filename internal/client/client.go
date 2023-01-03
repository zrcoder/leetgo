package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

func GetAll() (*model.All, error) {
	url := fmt.Sprintf("%s/api/problems/all", config.Domain())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Dev(err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Dev(err)
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Dev(err)
		return nil, err
	}
	defer resp.Body.Close()

	res := &model.All{}
	err = json.Unmarshal(data, res)
	log.Dev(err)
	return res, err
}

func Get(sp *model.StatStatusPair) (*model.Question, error) {
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
		log.Dev(err)
		return nil, err
	}
	domain := config.Domain()
	url := fmt.Sprintf("%s/graphql", domain)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jbody))
	if err != nil {
		log.Dev(err)
		return nil, err
	}
	referer := fmt.Sprintf("%s/problems/%s", domain, sp.Stat.QuestionTitleSlug)
	req.Header.Set("Referer", referer)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Dev(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Dev("request failed")
		return nil, fmt.Errorf("GetQuestion got status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Dev(err)
		return nil, err
	}
	defer resp.Body.Close()

	res := &model.GetQuestionResponse{}
	if err = json.Unmarshal(data, res); err != nil {
		log.Dev(string(data))
		log.Dev(err)
		return nil, err
	}
	question := res.Data.Question
	question.Id = sp.Stat.CalculatedID
	question.Title = sp.Stat.QuestionTitle
	question.Referer = referer
	question.Difficulty = sp.Difficulty.String()
	return question, nil
}
