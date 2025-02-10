package model

import (
	"strconv"
	"time"

	"github.com/zrcoder/leetgo/utils/parser"
)

type SolutionReq struct {
	ID       string
	Title    string
	CreateAt time.Time
	Author   string
}

type SolutionListResp interface {
	SolutionReqs() []SolutionReq
}

type SolutionListRespEN struct {
	Data struct {
		QuestionSolutions struct {
			Solutions []struct {
				Title string `json:"title"`
				Post  struct {
					Author struct {
						Username string `json:"username"`
					} `json:"author"`
					CreationDate int `json:"creationDate"`
				} `json:"post"`
				ID int `json:"id"`
			} `json:"solutions"`
			TotalNum int `json:"totalNum"`
		} `json:"questionSolutions"`
	} `json:"data"`
}

func (se *SolutionListRespEN) SolutionReqs() []SolutionReq {
	solutions := se.Data.QuestionSolutions.Solutions
	res := make([]SolutionReq, len(solutions))
	for i, s := range solutions {
		res[i] = SolutionReq{
			ID:       strconv.Itoa(s.ID),
			Title:    s.Title,
			CreateAt: time.Unix(int64(s.Post.CreationDate), 0),
			Author:   s.Post.Author.Username,
		}
	}
	return res
}

type SolutionListRespCN struct {
	Data struct {
		QuestionSolutionArticles struct {
			Edges []struct {
				Node struct {
					Slug      string    `json:"slug"`
					Title     string    `json:"title"`
					CreatedAt time.Time `json:"createdAt"`
					Author    struct {
						Username string `json:"username"`
					} `json:"author"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"questionSolutionArticles"`
	} `json:"data"`
}

func (sc *SolutionListRespCN) SolutionReqs() []SolutionReq {
	edges := sc.Data.QuestionSolutionArticles.Edges
	res := make([]SolutionReq, len(edges))
	for i, e := range edges {
		res[i] = SolutionReq{
			ID:       e.Node.Slug,
			Title:    e.Node.Title,
			CreateAt: e.Node.CreatedAt,
			Author:   e.Node.Author.Username,
		}
	}
	return res
}

type SolutionResp struct {
	Data struct {
		SolutionArticle struct {
			Title    string `json:"title"`
			Summary  string `json:"summary"`
			Content  string `json:"content"`
			Question struct {
				QuestionTitleSlug string `json:"questionTitleSlug"`
				Typename          string `json:"__typename"`
			} `json:"question"`
		} `json:"solutionArticle"`
		Topic struct {
			Title string `json:"title"`
			Post  struct {
				Content    string `json:"content"`
				ID         int    `json:"id"`
				VoteCount  int    `json:"voteCount"`
				VoteStatus int    `json:"voteStatus"`
			} `json:"post"`
			ID     int  `json:"id"`
			Pinned bool `json:"pinned"`
		} `json:"topic"`
	} `json:"data"`
}

// RegularContent returns regulared markdown content
func (sp *SolutionResp) RegularContent() (string, error) {
	content := sp.Data.Topic.Post.Content
	if content == "" {
		content = sp.Data.SolutionArticle.Content
	}
	// content is all ready markdown
	return parser.NewWithString(content).PreRrgular().Regular().String()
}
