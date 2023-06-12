package model

import (
	"strconv"
)

type List struct {
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
	FrontendQuestionID any    `json:"frontend_question_id"` // leetcode.com is number, but leetcode.cn is string
	FrontendID         string // caculated frontend id
}

type Difficulty struct {
	Level int `json:"level"`
}

func (s *Stat) CalFrontendID() string {
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

func StrToDifficulty(s string) Difficulty {
	dic := map[string]int{
		"Easy":   1,
		"Medium": 2,
		"Hard":   3,
	}
	return Difficulty{
		Level: dic[s],
	}
}
