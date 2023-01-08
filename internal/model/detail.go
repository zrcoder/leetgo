package model

import (
	"encoding/json"
	"fmt"
	"strings"

	h2md "github.com/JohannesKaufmann/html-to-markdown"

	"github.com/zrcoder/leetgo/internal/log"
)

type Question struct {
	// redefined
	Id         string `json:"id"`
	Title      string `json:"title"`
	Referer    string `json:"referer"`
	Difficulty string `json:"difficulty"`
	MdContent  string `json:"mdContent"`

	// original
	QuestionID        string `json:"questionId"`
	Content           string `json:"content"`
	Stats             string `json:"stats"`
	CodeDefinition    string `json:"codeDefinition"`
	SampleTestCase    string `json:"sampleTestCase"`
	EnableRunCode     bool   `json:"enableRunCode"`
	MetaData          string `json:"metaData"`
	TranslatedContent string `json:"translatedContent"`
}

func (q *Question) TransformContent() error {
	var err error
	content := q.TranslatedContent
	if content == "" {
		content = q.Content
	}
	content = strings.ReplaceAll(content, "<sup>", "^")
	content = strings.ReplaceAll(content, "</sup>", "")
	content, err = h2md.NewConverter("", true, nil).ConvertString(content)
	if err != nil {
		log.Trace(err)
		return err
	}
	content = strings.ReplaceAll(content, `\[`, `[`)
	content = strings.ReplaceAll(content, `\]`, `]`)
	q.MdContent = fmt.Sprintf("## [%s. %s](%s) (%s)\n\n%s",
		q.Id, q.Title, q.Referer, q.Difficulty, content)
	q.Content = ""
	q.TranslatedContent = ""
	return nil
}

func (q *Question) ParseCodes() ([]*Code, error) {
	q.CodeDefinition = strings.ReplaceAll(q.CodeDefinition, `\\n`, `\n`)
	res := []*Code{}
	err := json.Unmarshal([]byte(q.CodeDefinition), &res)
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
