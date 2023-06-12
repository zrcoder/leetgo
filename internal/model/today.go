package model

type Today interface {
	ToStatePair() *StatStatusPair
	FrontendID() string
}

type TodayEN struct {
	Data struct {
		ActiveDaily struct {
			Question TodayQuestion `json:"question"`
		} `json:"activeDailyCodingChallengeQuestion"`
	} `json:"data"`
}

type TodayCN struct {
	Data struct {
		TodayRecord []struct {
			Question TodayQuestion `json:"question"`
		} `json:"todayRecord"`
	} `json:"data"`
}

type TodayQuestion struct {
	Difficulty string `json:"difficulty"`
	FrontendID string `json:"frontendQuestionId"`
	PaidOnly   bool   `json:"paidOnly"`
	Title      string `json:"title"`
	TitleSlug  string `json:"titleSlug"`
}

func (t *TodayCN) FrontendID() string {
	if len(t.Data.TodayRecord) == 0 {
		return ""
	}
	return t.Data.TodayRecord[0].Question.FrontendID
}

func (t *TodayEN) FrontendID() string {
	return t.Data.ActiveDaily.Question.FrontendID
}
func (t *TodayEN) ToStatePair() *StatStatusPair {
	return t.Data.ActiveDaily.Question.ToStatePair()
}

func (t *TodayCN) ToStatePair() *StatStatusPair {
	if len(t.Data.TodayRecord) == 0 {
		return nil
	}
	return t.Data.TodayRecord[0].Question.ToStatePair()
}

func (q *TodayQuestion) ToStatePair() *StatStatusPair {
	return &StatStatusPair{
		Stat: Stat{
			QuestionTitle:     q.Title,
			QuestionTitleSlug: q.TitleSlug,
			FrontendID:        q.FrontendID,
		},
		PaidOnly:   q.PaidOnly,
		Difficulty: StrToDifficulty(q.Difficulty),
	}
}
