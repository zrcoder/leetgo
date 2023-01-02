package model

type Question struct {
	// redefined
	FrontendQuestionID string `json:"frontend_question_id"`
	Title              string `json:"title"`
	Referer            string `json:"referer"`
	Difficulty         string `json:"difficulty"`

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
