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

func (c *client) GetAll() (*model.All, error) {
	res := &model.All{}
	err := c.rb.
		Path("/api/problems/all").
		ToJSON(res).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *client) Search(keyWords string) (*model.QustionsResp, error) {
	query := `
	query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
		problemsetQuestionList: questionList(
		  categorySlug: $categorySlug
		  limit: $limit
		  skip: $skip
		  filters: $filters
		) {
		  total: totalNum
		  questions: data {
			difficulty
			frontendQuestionId: questionFrontendId
			paidOnly: isPaidOnly
			title
			titleSlug
		  }
		}
	  }`
	return c.search(keyWords, query)
}

func (c *client) search(keyWords, query string) (*model.QustionsResp, error) {
	log.Debug("remote search:", keyWords)
	const operation = "problemsetQuestionList"
	vars := map[string]any{
		"skip":         0,
		"limit":        listLimit,
		"categorySlug": "",
		"filters": map[string]string{
			"searchKeywords": keyWords,
		},
	}
	refer := fmt.Sprintf("%s/problemset/all/?search=%s&page=%d", c.domain, keyWords, 1)
	res := &model.QustionsResp{}
	err := c.graphql(operation, query, refer, vars, res)
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *client) GetQuestion(meta *model.Meta) (*model.Question, error) {
	log.Debug("get question:", meta.FrontendID, meta.Title, meta.TitleSlug)
	const operation = "getQuestionDetail"
	const query = `
	query getQuestionDetail($titleSlug: String!) {
	  isCurrentUserAuthenticated
	  question(titleSlug: $titleSlug) {
		questionId
		stats
		codeDefinition
		sampleTestCase
		enableRunCode
		title
		content
 		translatedTitle
		translatedContent
	  }
	}`

	refer := fmt.Sprintf("%s/problems/%s", c.domain, meta.TitleSlug)
	res := &model.GetQuestionResponse{}
	err := c.graphql(operation, query, refer, map[string]any{"titleSlug": meta.TitleSlug}, res)
	if err != nil {
		return nil, err
	}

	question := res.Data.Question
	err = question.Transform(meta, refer)

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
				titleCn: translatedTitle
				titleSlug
            }
        }
    }`
	refer := fmt.Sprintf("%s/problemset/all/", c.domain)
	res = &model.Today{}
	err = c.graphql(operation, fmt.Sprintf(queryFmt, key), refer, nil, res)
	if err != nil {
		log.Debug(err)
	}
	return
}

func (c *client) Test(question *model.Question, typedCode, codeLang string) (string, error) {
	log.Debugf("remote test, slug(%s), frontendID(%s), titile(%s)",
		question.TitleSlug, question.FrontendID, question.Title)
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
	log.Debugf("submit, slug(%s), frontendID(%s), titile(%s)",
		question.TitleSlug, question.FrontendID, question.Title)
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

func (c *client) GetSolutions(meta *model.Meta) (model.SolutionListResp, error) {
	log.Debug("query solutions for question", meta.FrontendID, meta.Title, meta.TitleSlug)
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
			post {
				creationDate
				author {
					username
				}
			}
		  }
		}
	  }`

	vars := map[string]any{
		"questionSlug": meta.TitleSlug,
		"skip":         0,
		"first":        listLimit,
		"orderBy":      "most_votes",
	}
	refer := fmt.Sprintf("%s/problems/%s/solutions/",
		c.domain, meta.TitleSlug)
	log.Debug("refer:", refer)
	res := &model.SolutionListRespEN{}
	err := c.graphql(operation, query, refer, vars, res)
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *client) GetSolution(solution *model.SolutionReq, meta *model.Meta) (*model.SolutionResp, error) {
	log.Debug("get solution", solution.ID, solution.Title, meta.TitleSlug)
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
	refer := fmt.Sprintf("%s/problems/%s/solutions/%d/%s/",
		c.domain,
		meta.TitleSlug,
		id,
		solution.Title,
	)
	res := &model.SolutionResp{}
	err := c.graphql(operation, query, refer, map[string]any{"topicId": id}, res)
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *client) graphql(operation, query, refer string, vars map[string]any, res any) error {
	log.Debug("graphql:", operation, refer)
	body := map[string]any{
		"operationName": operation,
		"query":         query,
		"variables":     vars,
	}
	errJson := map[string]any{}
	err := c.rb.
		Path("/graphql").
		Header("Referer", refer).
		BodyJSON(body).
		ToJSON(&res).
		ErrorJSON(&errJson).
		Fetch(context.Background())
	if err != nil {
		log.Debug(err)
		log.Debug(errJson)
	}
	return err
}
