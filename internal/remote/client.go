package remote

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"

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
	res := &model.List{}
	err := requests.URL(uri).
		ToJSON(res).
		Fetch(context.Background())
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
	domain := config.Domain()
	uri := fmt.Sprintf("%s/graphql", domain)
	referer := fmt.Sprintf("%s/problems/%s", domain, sp.Stat.QuestionTitleSlug)
	res := &model.GetQuestionResponse{}

	err := requests.URL(uri).
		ContentType("application/json").
		Header("Cache-Control", "no-cache").
		Header("Referer", referer).
		BodyJSON(body).
		ToJSON(res).
		Fetch(context.Background())
	if err != nil {
		log.Trace(err)
		return nil, err
	}

	question := res.Data.Question
	question.ID = sp.Stat.CalculatedID
	question.Title = sp.Stat.QuestionTitle
	question.Slug = sp.Stat.QuestionTitleSlug
	question.Referer = referer
	question.Difficulty = sp.Difficulty.String()
	err = question.TransformContent()
	return question, err
}
