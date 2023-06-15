package model

type Today struct {
	Data struct {
		ActiveDaily struct { // leetcode.com
			Question Meta `json:"question"`
		} `json:"activeDailyCodingChallengeQuestion"`
		TodayRecord []struct { // leetcode.cn
			Question Meta `json:"question"`
		} `json:"todayRecord"`
	} `json:"data"`
}

func (t *Today) Meta() *Meta {
	q := t.Data.ActiveDaily.Question
	if len(t.Data.TodayRecord) > 0 {
		q = t.Data.TodayRecord[0].Question
	}
	q.Transform()
	return &q
}
