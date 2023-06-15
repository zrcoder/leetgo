package remote

import (
	"fmt"

	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

type clientCN struct {
	*client
}

func (c *clientCN) Search(keyWords string) (*model.QustionsResp, error) {
	query := `
	query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {
		problemsetQuestionList(
		  categorySlug: $categorySlug
		  limit: $limit
		  skip: $skip
		  filters: $filters
		) {
		  total
		  questions {
			difficulty
			frontendQuestionId
			paidOnly
			title
			titleCn
			titleSlug
		  }
		}
	  }`
	return c.search(keyWords, query)
}

func (c *clientCN) GetToday() (res *model.Today, err error) {
	return c.getToday("todayRecord")
}

func (c *clientCN) GetSolutions(meta *model.Meta) (model.SolutionListResp, error) {
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
		"questionSlug": meta.TitleSlug,
		"first":        listLimit,
		"skip":         0,
		"orderBy":      "MOST_UPVOTE",
	}
	refer := fmt.Sprintf("%s/problems/%s/solutions/",
		c.domain, meta.TitleSlug)
	log.Debug("refer:", refer)
	res := &model.SolutionListRespCN{}
	err := c.graphql(operation, query, refer, vars, res)
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

func (c *clientCN) GetSolution(req *model.SolutionReq, meta *model.Meta) (*model.SolutionResp, error) {
	log.Debugf("load solution, queston slug: %s, solution slug: %s, solution title: %s", meta.TitleSlug, req.ID, req.Title)
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
	refer := fmt.Sprintf("%s/problems/%s/solution/%s/",
		c.domain,
		meta.TitleSlug,
		req.ID,
	)
	log.Debug("referer:", refer)
	res := &model.SolutionResp{}
	err := c.graphql(operation, query, refer, vars, res)
	if err != nil {
		log.Debug(err)
	}
	return res, nil
}
