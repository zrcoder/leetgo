package model

import (
	"encoding/json"
	"strings"
)

type Question struct {
	// redefined
	Id         string `json:"id"`
	Title      string `json:"title"`
	Referer    string `json:"referer"`
	Difficulty string `json:"difficulty"`

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
