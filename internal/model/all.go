package model

import (
	"strconv"
)

type All struct {
	StatStatusPairs []StatStatusPair `json:"stat_status_pairs"`
}

type StatStatusPair struct {
	Stat       Stat       `json:"stat"`
	Difficulty Difficulty `json:"difficulty"`
	PaidOnly   bool       `json:"paid_only"`
}

type Stat struct {
	QuestionTitle      string `json:"question__title"`
	QuestionTitleSlug  string `json:"question__title_slug"`
	FrontendQuestionID any    `json:"frontend_question_id"` // leetcode.com is int, but leetcode.cn is string
	CalculatedID       string `json:"calculated_id"`
}

type Difficulty struct {
	Level int `json:"level"`
}

func (s *Stat) GetFrontendQuestionID() string {
	if val, ok := s.FrontendQuestionID.(string); ok {
		return val
	}
	return strconv.Itoa(int(s.FrontendQuestionID.(float64)))
}

func (d *Difficulty) String() string {
	dic := map[int]string{
		1: "Easy",
		2: "Medium",
		3: "Hard",
	}
	return dic[d.Level]
}
