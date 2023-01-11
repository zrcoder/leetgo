package remote

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

const (
	questionQueryTmp = `
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
)

func GetList() (*model.List, error) {
	uri := fmt.Sprintf("%s/api/problems/all", config.Domain())
	resp, err := http.Get(uri)
	if err != nil {
		log.Trace(err)
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Trace(err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	res := &model.List{}
	err = json.Unmarshal(data, res)
	log.Trace(err)
	return res, err
}

func GetQuestion(sp *model.StatStatusPair) (*model.Question, error) {
	body := map[string]any{
		"query": questionQueryTmp,
		"variables": map[string]string{
			"titleSlug": sp.Stat.QuestionTitleSlug,
		},
		"operationName": "getQuestionDetail",
	}
	reqBody, _ := json.Marshal(body)
	domain := config.Domain()
	uri := fmt.Sprintf("%s/graphql", domain)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Trace(err)
		return nil, err
	}
	referer := fmt.Sprintf("%s/problems/%s", domain, sp.Stat.QuestionTitleSlug)
	setCommonHeaders(req, referer)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Trace(err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("querey question failed, response status: %s", resp.Status)
		log.Trace(err)
		return nil, err
	}

	res := &model.GetQuestionResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		log.Trace(err)
		return nil, err
	}
	question := res.Data.Question
	question.ID = sp.Stat.CalculatedID
	question.Title = sp.Stat.QuestionTitle
	question.Referer = referer
	question.Difficulty = sp.Difficulty.String()
	err = question.TransformContent()
	return question, err
}

func setCommonHeaders(req *http.Request, referer string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Add("Referer", referer)
}
