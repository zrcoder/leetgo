package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/carlmjohnson/requests"

	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

type client struct {
	domain string
	rb     *requests.Builder
}

func newClient(domain, token, session string) *client {
	return &client{
		domain: domain,
		rb: requests.New().BaseURL(domain).
			ContentType("application/json").
			Cookie("LEETCODE_SESSION", session).
			Cookie("csrftoken", token).
			Header("x-csrftoken", token),
	}
}

func (c *client) GetList() (*model.List, error) {
	res := &model.List{}
	err := c.rb.
		Path("/api/problems/all").
		ToJSON(res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *client) GetQuestion(sp *model.StatStatusPair) (*model.Question, error) {
	const operation = "getQuestionDetail"
	const query = `
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

	body := graphqlBody(query, operation, map[string]any{"titleSlug": sp.Stat.QuestionTitleSlug})
	refer := fmt.Sprintf("%s/problems/%s", c.domain, sp.Stat.QuestionTitleSlug)
	res := &model.GetQuestionResponse{}
	err := c.graphql(operation, refer, body, res)
	if err != nil {
		return nil, err
	}

	question := res.Data.Question
	question.FrontendID = sp.Stat.FrontendID
	question.Title = sp.Stat.QuestionTitle
	question.TitleSlug = sp.Stat.QuestionTitleSlug
	question.Referer = refer
	question.Difficulty = sp.Difficulty.String()
	err = question.TransformContent()
	return question, err
}

func (c *client) GetToday() (res *model.Today, err error) {
	return c.getToday("activeDailyCodingChallengeQuestion")
}

func (c *client) getToday(key string) (res *model.Today, err error) {
	const operation = "questionOfToday"
	const queryFmt = `
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
	body := graphqlBody(fmt.Sprintf(queryFmt, key), operation, nil)
	refer := fmt.Sprintf("%s/problemset/all/", c.domain)
	res = &model.Today{}
	err = c.graphql(operation, refer, body, res)
	if err != nil {
		log.Debug(err)
	}
	return
}

func (c *client) Test(question *model.Question, typedCode, codeLang string) (string, error) {
	body := map[string]string{
		"lang":        codeLang,
		"question_id": question.QuestionID,
		"typed_code":  typedCode,
		"data_input":  question.SampleTestCase,
	}
	type resp struct {
		InterpretId string `json:"interpret_id"`
	}
	res := &resp{}
	err := c.rb.
		Pathf("/problems/%s/interpret_solution/", question.TitleSlug).
		Header("Referer", question.Referer).
		BodyJSON(body).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
	}
	return res.InterpretId, err
}

func (c *client) Submit(question *model.Question, typedCode, codeLang string) (string, error) {
	body := map[string]any{
		"lang":         codeLang,
		"questionSlug": question.TitleSlug,
		"question_id":  question.QuestionID,
		"typed_code":   typedCode,
	}
	res := ""
	err := c.rb.
		Pathf("/problems/%s/submit/", question.TitleSlug).
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

func (c *client) CheckResult(id string, question *model.Question, res model.RunResult) error {
	err := c.rb.
		Pathf("/submissions/detail/%s/check/", id).
		Header("Referer", question.Referer).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
	}
	return err
}

func (c *client) GetSolutions(meta *model.StatStatusPair) (model.SolutionListResp, error) {
	log.Debug("query solutions for question", meta.Stat.FrontendID)
	const operation = "communitySolutions"
	const query = `
	query communitySolutions($questionSlug: String!, $skip: Int!, $first: Int!, $orderBy: TopicSortingOption) {
		questionSolutions(
		  filters: {questionSlug: $questionSlug, skip: $skip, first: $first,  orderBy: $orderBy}
		) {
		  totalNum
		  solutions {
			id
			title
		  }
		}
	  }`

	vars := map[string]any{
		"questionSlug": meta.Stat.QuestionTitleSlug,
		"skip":         0,
		"first":        solutionsLimit,
		"orderBy":      "most_votes",
	}
	body := graphqlBody(query, operation, vars)
	refer := fmt.Sprintf("%s/problems/%s/solutions/",
		c.domain, meta.Stat.QuestionTitleSlug)
	log.Debug("refer:", refer)
	res := &model.SolutionListRespEN{}
	err := c.graphql(operation, refer, body, res)
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *client) GetSolution(solution *model.SolutionReq, meta *model.StatStatusPair) (*model.SolutionResp, error) {
	const operation = "communitySolution"
	const query = `
	query communitySolution($topicId: Int!) {
		topic(id: $topicId) {
			id
			title
			pinned
			post {
				id
				content
			}
		}
	}`
	id, _ := strconv.Atoi(solution.ID)
	body := graphqlBody(query, operation, map[string]any{"topicId": id})
	refer := fmt.Sprintf("%s/problems/%s/solutions/%d/%s/",
		c.domain,
		meta.Stat.QuestionTitleSlug,
		id,
		solution.Title,
	)
	res := &model.SolutionResp{}
	err := c.graphql(operation, refer, body, res)
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *client) graphql(operation, refer string, body map[string]any, res any) error {
	log.Debug("graphql:", operation, refer)
	err := c.rb.
		Path("/graphql").
		Header("Referer", refer).
		BodyJSON(body).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
	}
	return err
}

func graphqlBody(query, operation string, variables map[string]any) map[string]any {
	return map[string]any{
		"query":         query,
		"variables":     variables,
		"operationName": operation,
	}
}
