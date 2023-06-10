package remote

import (
	"context"
	"fmt"
	"strings"

	"github.com/carlmjohnson/requests"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/mod"
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

// common config for request client
func cfg(rb *requests.Builder) {
	rb.BaseURL(config.Domain())

	if mod.IsDebug() {
		u, err := rb.URL()
		if err == nil {
			log.Debug("request url:", u.String())
		} else {
			log.Debug(err)
		}
	}
}

func GetList() (*model.List, error) {
	res := &model.List{}
	err := requests.New(cfg).
		Path("/api/problems/all").
		ToJSON(res).
		Fetch(context.Background())
	log.Debug(err)
	return res, err
}

func GetQuestion(sp *model.StatStatusPair) (*model.Question, error) {
	referer := fmt.Sprintf("%s/problems/%s",
		config.Domain(),
		sp.Stat.QuestionTitleSlug)
	res := &model.GetQuestionResponse{}
	err := requests.New(cfg).
		Path("/graphql").
		ContentType("application/json").
		Header("Cache-Control", "no-cache").
		Header("Referer", referer).
		BodyJSON(map[string]any{
			"query": questionQueryTmp,
			"variables": map[string]string{
				"titleSlug": sp.Stat.QuestionTitleSlug,
			},
			"operationName": "getQuestionDetail",
		}).
		ToJSON(res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
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

func Test(question *model.Question, typedCode, codeLang string, tests []string) (*model.InterpretSolutionResult, error) {
	res := &model.InterpretSolutionResult{}
	builder := requests.New(cfg).
		Pathf("/problems/%s/interpret_solution", question.Slug).
		BodyJSON(map[string]any{
			"lang":        codeLang,
			"question_id": question.QuestionID,
			"typed_code":  typedCode,
			"data_input":  strings.Join(tests, "\n"),
		}).
		ToJSON(res)

	err := builder.Fetch(context.Background())
	return res, err
}

func Check(submitId string) (string, error) {
	path := fmt.Sprintf("%s/submissions/detail/%s/check", config.Domain(), submitId)
	res := ""
	err := requests.URL(path).ToString(&res).Fetch(context.Background())
	return res, err
}
