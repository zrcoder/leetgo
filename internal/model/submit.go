package model

type SubmitRequest struct {
	Lang       string `json:"lang"`
	QuestionID string `json:"question_id"`
	TestMode   string `json:"test_mode"`
	Name       string `json:"name"`
	JudgeType  string `json:"judge_type"`
	TypedCode  string `json:"typed_code"`
}

type SubmitResult struct {
	SubmissionID int `json:"submission_id"`
}
