package model

import (
	"strconv"
)

type All struct {
	UserName       string           `json:"user_name"`
	NumSolved      int              `json:"num_solved"`
	NumTotal       int              `json:"num_total"`
	AcEasy         int              `json:"ac_easy"`
	AcMedium       int              `json:"ac_medium"`
	AcHard         int              `json:"ac_hard"`
	StatStatusPair []StatStatusPair `json:"stat_status_pairs"`
}

type Stat struct {
	QuestionID         int    `json:"question_id"`
	QuestionTitle      string `json:"question__title"`
	QuestionTitleSlug  string `json:"question__title_slug"`
	QuestionHide       bool   `json:"question__hide"`
	TotalAcs           int    `json:"total_acs"`
	TotalSubmitted     int    `json:"total_submitted"`
	FrontendQuestionID any    `json:"frontend_question_id"` // leetcode.com is int, but leetcode.cn is string
	IsNewQuestion      bool   `json:"is_new_question"`
}

func (s *Stat) GetFrontendQuestionID() string {
	if val, ok := s.FrontendQuestionID.(string); ok {
		return val
	}
	return strconv.Itoa(int(s.FrontendQuestionID.(float64)))
}

type Difficulty struct {
	Level int `json:"level"`
}
type StatStatusPair struct {
	Stat       Stat        `json:"stat"`
	Status     interface{} `json:"status"`
	Difficulty Difficulty  `json:"difficulty"`
	PaidOnly   bool        `json:"paid_only"`
	IsFavor    bool        `json:"is_favor"`
	Frequency  int         `json:"frequency"`
	Progress   int         `json:"progress"`
}

func (d *Difficulty) String() string {
	dic := map[int]string{
		1: "Easy",
		2: "Medium",
		3: "Hard",
	}
	return dic[d.Level]
}
