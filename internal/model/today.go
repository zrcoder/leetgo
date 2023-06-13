package model

type Today struct {
	Data struct {
		ActiveDaily struct { // leetcode.com
			Question Question `json:"question"`
		} `json:"activeDailyCodingChallengeQuestion"`
		TodayRecord []struct { // leetcode.cn
			Question Question `json:"question"`
		} `json:"todayRecord"`
	} `json:"data"`
}

func (t *Today) Question() Question {
	if len(t.Data.TodayRecord) == 0 {
		return t.Data.ActiveDaily.Question
	}
	return t.Data.TodayRecord[0].Question
}
func (q Question) StatePair() *StatStatusPair {
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
