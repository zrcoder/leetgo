package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/utils/parser"
)

type Meta struct {
	FrontendID string `json:"frontendQuestionId"`
	Title      string `json:"title"`
	Referer    string
	TitleSlug  string `json:"titleSlug"`
	Difficulty string `json:"difficulty"`
	PaidOnly   bool   `json:"paidOnly"`
}

type Question struct {
	Meta
	MdContent string `json:"mdContent"`

	// original
	QuestionID        string `json:"questionId"`
	Content           string `json:"content"`
	Stats             string `json:"stats"`
	CodeDefinition    string `json:"codeDefinition"`
	SampleTestCase    string `json:"sampleTestCase"`
	EnableRunCode     bool   `json:"enableRunCode"`
	TranslatedContent string `json:"translatedContent"`
}

func (q *Question) Transform(meta *StatStatusPair, refer string) error {
	q.FrontendID = meta.Stat.FrontendID
	q.Title = meta.Stat.QuestionTitle
	q.TitleSlug = meta.Stat.QuestionTitleSlug
	q.Difficulty = meta.Difficulty.String()
	q.Referer = refer
	return q.transformContent()
}

func (q *Question) transformContent() error {
	var err error
	content := q.TranslatedContent
	if content == "" {
		content = q.Content
	}
	content, err = parser.NewWithString(content).PreRrgular().ToMarkDown().Regular().String()
	if err != nil {
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
	log.Debug(err)
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
