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
	token, session, _ := getCredentials()
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
	if err != nil {
		log.Debug(err)
	}
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
	question.FrontendID = sp.Stat.FrontendID
	question.Title = sp.Stat.QuestionTitle
	question.Slug = sp.Stat.QuestionTitleSlug
	question.Referer = referer
	question.Difficulty = sp.Difficulty.String()
	err = question.TransformContent()
	return question, err
}

func GetToday() (res model.Today, err error) {
	isEn := config.IsDefaultLang()
	queryTmp := `
    query questionOfToday {
        %s {
            question {
				difficulty
				frontendQuestionId: questionFrontendId
				paidOnly: isPaidOnly
				title
				titleSlug
            }
        }
    }`
	queryKey := "todayRecord"
	if isEn {
		queryKey = "activeDailyCodingChallengeQuestion"
	}

	body := map[string]any{
		"query":         fmt.Sprintf(queryTmp, queryKey),
		"operationName": "questionOfToday",
	}
	referer := fmt.Sprintf("%s/problemset/all/", config.Domain())
	if isEn {
		res = &model.TodayEN{}
	} else {
		res = &model.TodayCN{}
	}
	err = requests.New(cfg).
		Path("/graphql").
		Header("Referer", referer).
		BodyJSON(body).
		ToJSON(res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
	}
	return
}

func Test(question *model.Question, typedCode, codeLang string) (string, error) {
	body := map[string]any{
		"lang":        codeLang,
		"question_id": question.QuestionID,
		"typed_code":  typedCode,
		"data_input":  question.SampleTestCase,
	}
	type resp struct {
		InterpretId string `json:"interpret_id"`
	}
	res := &resp{}
	err := requests.New(cfg).
		Pathf("/problems/%s/interpret_solution/", question.Slug).
		Header("Referer", question.Referer).
		BodyJSON(body).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
	}
	return res.InterpretId, err
}

func Submit(question *model.Question, typedCode, codeLang string) (string, error) {
	body := map[string]any{
		"lang":         codeLang,
		"questionSlug": question.Slug,
		"question_id":  question.QuestionID,
		"typed_code":   typedCode,
	}
	res := ""
	err := requests.New(cfg).
		Pathf("/problems/%s/submit/", question.Slug).
		Header("Referer", question.Referer).
		BodyJSON(body).
		ToString(&res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
		return "", err
	}
	type resp struct {
		SubmissionID int `json:"submission_id"`
	}
	rsp := &resp{}
	err = json.Unmarshal([]byte(res), rsp)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	return strconv.Itoa(rsp.SubmissionID), nil
}

func CheckResult(id string, question *model.Question, res model.RunResult) error {
	err := requests.New(cfg).
		Pathf("/submissions/detail/%s/check/", id).
		Header("Referer", question.Referer).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
	}
	return err
}
