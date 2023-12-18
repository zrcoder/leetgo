package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/utils/parser"
)

type QustionsResp struct {
	Data struct {
		ProblemsetQuestionList struct {
			Questions []Meta `json:"questions"`
		} `json:"problemsetQuestionList"`
		Total int `json:"total"`
	} `json:"data"`
}

type Meta struct {
	FrontendID      string `json:"frontendQuestionId"`
	Title           string `json:"title"`
	TranslatedTitle string `json:"translatedTitle"`
	TitleSlug       string `json:"titleSlug"`
	Difficulty      string `json:"difficulty"`
	TitleCn         string `json:"titleCn"`
	Referer         string
	PaidOnly        bool `json:"paidOnly"`
}

type Question struct {
	MdContent         string `json:"mdContent"`
	QuestionID        string `json:"questionId"`
	Stats             string `json:"stats"`
	CodeDefinition    string `json:"codeDefinition"`
	SampleTestCase    string `json:"sampleTestCase"`
	Content           string `json:"content"`
	TranslatedContent string `json:"translatedContent"`
	Meta
	EnableRunCode bool `json:"enableRunCode"`
}

func (q *Meta) Transform() {
	if q.TitleCn != "" {
		q.Title = q.TitleCn
	} else if q.TranslatedTitle != "" {
		q.Title = q.TranslatedTitle
	}
}

func (q *Question) Transform(meta *Meta, refer string) error {
	q.Referer = refer
	q.FrontendID = meta.FrontendID
	q.TitleSlug = meta.TitleSlug
	q.Difficulty = meta.Difficulty
	q.Meta.Transform()
	var err error
	content := q.TranslatedContent
	if content == "" {
		content = q.Content
	}
	content, err = parser.NewWithString(content).PreRrgular().ToMarkDown().Regular().String()
	if err != nil {
		log.Debug(err)
		return err
	}
	q.MdContent = fmt.Sprintf("## [%s. %s](%s) (%s)\n\n%s\n\n",
		q.FrontendID, q.Title, q.Referer, q.Difficulty, content)
	q.Content = ""
	q.TranslatedContent = ""
	return nil
}

func (q *Question) ParseCodes() ([]*Code, error) {
	q.CodeDefinition = strings.ReplaceAll(q.CodeDefinition, `\\n`, `\n`)
	var res []*Code
	err := json.Unmarshal([]byte(q.CodeDefinition), &res)
	if err != nil {
		log.Debug(q.CodeDefinition)
		log.Debug(err)
	}
	return res, err
}

// Code the struct of leetcode codes.
type Code struct {
	Text        string `json:"text"`
	Value       string `json:"value"`
	DefaultCode string `json:"defaultCode"`
}

type GetQuestionResponseData struct {
	Question *Question `json:"question"`
}
type GetQuestionResponse struct {
	Data GetQuestionResponseData `json:"data"`
}
