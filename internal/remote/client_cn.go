package remote

import (
	"fmt"

	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

type clientCN struct {
	*client
}

func newClientCN(domain, token, session string) *clientCN {
	return &clientCN{
		client: newClient(domain, token, session),
	}
}

func (c *clientCN) GetToday() (res *model.Today, err error) {
	return c.getToday("todayRecord")
}

func (c *clientCN) GetSolutions(question *model.Question) (model.SolutionListResp, error) {
	log.Debug("query solutions for question", question.FrontendID)
	const operation = "questionSolutionArticles"
	const query = `
	query questionSolutionArticles($questionSlug: String!, $skip: Int, $first: Int, $orderBy: SolutionArticleOrderBy) {
		questionSolutionArticles(questionSlug: $questionSlug, skip: $skip, first: $first, orderBy: $orderBy) {
		  edges {
			node {
			  ...solutionArticle
			}
		  }
		}
	  }
	  fragment solutionArticle on SolutionArticleNode {
		title
		slug
	  }`

	vars := map[string]any{
		"questionSlug": question.TitleSlug,
		"first":        solutionsLimit,
		"skip":         0,
		"orderBy":      "MOST_UPVOTE",
	}
	body := graphqlBody(query, operation, vars)
	refer := fmt.Sprintf("%s/problems/%s/solutions/",
		c.domain, question.TitleSlug)
	log.Debug("refer:", refer)
	res := &model.SolutionListRespCN{}
	err := c.graphql(operation, refer, body, res)
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *clientCN) GetSolution(req *model.SolutionReq, question *model.Question) (*model.SolutionResp, error) {
	log.Debugf("load solution, queston slug: %s, solution slug: %s, solution title: %s", question.TitleSlug, req.ID, req.Title)
	const operation = "solutionDetailArticle"
	const query = `
	query solutionDetailArticle($slug: String!, $orderBy: SolutionArticleOrderBy!) {
		solutionArticle(slug: $slug, orderBy: $orderBy) {
			...solutionArticle
			content
		}
	}
	fragment solutionArticle on SolutionArticleNode {
		title
		slug
	}`
	vars := map[string]any{
		"slug":    req.ID,
		"orderBy": "MOST_UPVOTE",
	}
	body := graphqlBody(query, operation, vars)
	refer := fmt.Sprintf("%s/problems/%s/solution/%s/",
		c.domain,
		question.TitleSlug,
		req.ID,
	)
	log.Debug("referer:", refer)
	res := &model.SolutionResp{}
	err := c.graphql(operation, refer, body, res)
	if err != nil {
		log.Debug(err)
	}
	return res, nil
}
