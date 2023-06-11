package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/carlmjohnson/requests"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

// common config for request client
func cfg(rb *requests.Builder) {
	token, session, err := getCredentials()
	if err != nil {
		log.Debug(err)
	}
	rb.BaseURL(config.Domain()).
		ContentType("application/json").
		Cookie("LEETCODE_SESSION", session).
		Cookie("csrftoken", token).
		Header("x-csrftoken", token)

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
	const questionQueryTmp = `
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
		"query": questionQueryTmp,
		"variables": map[string]string{
			"titleSlug": sp.Stat.QuestionTitleSlug,
		},
		"operationName": "getQuestionDetail",
	}

	referer := fmt.Sprintf("%s/problems/%s",
		config.Domain(),
		sp.Stat.QuestionTitleSlug)
	res := &model.GetQuestionResponse{}
	err := requests.New(cfg).
		Path("/graphql").
		Header("Referer", referer).
		BodyJSON(body).
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

func Test(question *model.Question, typedCode, codeLang string, inputCases string) (*model.TestResult, error) {
	body := map[string]any{
		"lang":        codeLang,
		"question_id": question.QuestionID,
		"typed_code":  typedCode,
		"data_input":  inputCases,
	}
	refer := fmt.Sprintf("%s/problems/%s/",
		config.Domain(),
		question.Slug)
	res := &model.TestResult{}
	err := requests.New(cfg).
		Pathf("/problems/%s/interpret_solution/", question.Slug).
		Header("Referer", refer).
		BodyJSON(body).
		ToJSON(&res).
		Fetch(context.Background())
	return res, err
}

func CheckTest(checkID, questionSlug string) (*model.TestCheckResult, error) {
	res := &model.TestCheckResult{}
	refer := fmt.Sprintf("%s/problems/%s/",
		config.Domain(),
		questionSlug)
	err := requests.New(cfg).
		Pathf("/submissions/detail/%s/check/", checkID).
		Header("Referer", refer).
		ToJSON(&res).
		Fetch(context.Background())
	return res, err
}

func Submit(question *model.Question, typedCode, codeLang string) (string, error) {
	body := map[string]any{
		"lang":         codeLang,
		"questionSlug": question.Slug,
		"question_id":  question.QuestionID,
		"typed_code":   typedCode,
	}
	refer := fmt.Sprintf("%s/problems/%s/",
		config.Domain(),
		question.Slug)
	res := ""
	err := requests.New(cfg).
		Pathf("/problems/%s/submit/", question.Slug).
		Header("Referer", refer).
		BodyJSON(body).
		ToString(&res).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}
	type resp struct {
		SubmissionID int `json:"submission_id"`
	}
	rsp := &resp{}
	err = json.Unmarshal([]byte(res), rsp)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(rsp.SubmissionID), nil
}

func CheckSubmit(submitId, questionSlug string) (*model.SubmitCheckResult, error) {
	res := &model.SubmitCheckResult{}
	refer := fmt.Sprintf("%s/problems/%s/",
		config.Domain(),
		questionSlug)
	err := requests.New(cfg).
		Pathf("/submissions/detail/%s/check/", submitId).
		Header("Referer", refer).
		ToJSON(&res).
		Fetch(context.Background())
	return res, err
}
