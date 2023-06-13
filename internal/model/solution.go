package model

import (
	"strconv"
)

type SolutionReq struct {
	ID    string
	Title string
}

type SolutionListResp interface {
	SolutionReqs() []SolutionReq
}

type SolutionListRespEN struct {
	Data struct {
		QuestionSolutions struct {
			TotalNum  int `json:"totalNum"`
			Solutions []struct {
				ID    int    `json:"id"`
				Title string `json:"title"`
			} `json:"solutions"`
		} `json:"questionSolutions"`
	} `json:"data"`
}

func (se *SolutionListRespEN) SolutionReqs() []SolutionReq {
	solutions := se.Data.QuestionSolutions.Solutions
	res := make([]SolutionReq, len(solutions))
	for i, s := range solutions {
		res[i] = SolutionReq{
			ID:    strconv.Itoa(s.ID),
			Title: s.Title,
		}
	}
	return res
}

type SolutionListRespCN struct {
	Data struct {
		QuestionSolutionArticles struct {
			Edges []struct {
				Node struct {
					Slug  string `json:"slug"`
					Title string `json:"title"`
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
			ID:    e.Node.Slug,
			Title: e.Node.Title,
		}
	}
	return res
}

type SolutionResp struct {
	Data struct {
		Topic struct { // leetcode.com
			ID     int    `json:"id"`
			Title  string `json:"title"`
			Pinned bool   `json:"pinned"`
			Post   struct {
				ID         int    `json:"id"`
				VoteCount  int    `json:"voteCount"`
				VoteStatus int    `json:"voteStatus"`
				Content    string `json:"content"`
			} `json:"post"`
		} `json:"topic"`
		SolutionArticle struct { // leetcode.cn
			Title    string `json:"title"`
			Summary  string `json:"summary"`
			Content  string `json:"content"`
			Question struct {
				QuestionTitleSlug string `json:"questionTitleSlug"`
				Typename          string `json:"__typename"`
			} `json:"question"`
		} `json:"solutionArticle"`
	} `json:"data"`
}

func (sp *SolutionResp) RegularContent() (string, error) {
	content := sp.Data.Topic.Post.Content
	if content == "" {
		content = sp.Data.SolutionArticle.Content
	}
	return Regular(content)
}
