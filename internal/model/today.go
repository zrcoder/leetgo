package model

type Today struct { // just meta info, no question content
	Data struct {
		ActiveDaily struct { // leetcode.com
			Question Question `json:"question"`
		} `json:"activeDailyCodingChallengeQuestion"`
		TodayRecord []struct { // leetcode.cn
			Question Question `json:"question"`
		} `json:"todayRecord"`
	} `json:"data"`
}

func (t *Today) Meta() *StatStatusPair {
	q := t.Data.ActiveDaily.Question
	if len(t.Data.TodayRecord) > 0 {
		q = t.Data.TodayRecord[0].Question
	}
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
