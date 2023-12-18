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
	FrontendQuestionID any    `json:"frontend_question_id"`
	QuestionTitle      string `json:"question__title"`
	QuestionTitleSlug  string `json:"question__title_slug"`
}

type Difficulty struct {
	Level int `json:"level"`
}

func (sp *StatStatusPair) Meta() Meta {
	return Meta{
		FrontendID: sp.Stat.getFrontendID(),
		TitleSlug:  sp.Stat.QuestionTitleSlug,
		Title:      sp.Stat.QuestionTitle,
		PaidOnly:   sp.PaidOnly,
		Difficulty: sp.Difficulty.String(),
	}
}

func (s *Stat) getFrontendID() string {
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
