package model

type InterpretSolutionResult struct {
	InterpretExpectedId string `json:"interpret_expected_id"`
	InterpretId         string `json:"interpret_id"`
	TestCase            string `json:"test_case"`
}

// type CheckResult interface {
// 	Display(q *Question) string
// 	GetState() string
// }

type TestCheckResult struct {
	InputData              string
	State                  string   `json:"state"` // STARTED, SUCCESS
	StatusCode             int      `json:"status_code"`
	StatusMsg              string   `json:"status_msg"`         // Accepted, Wrong Answer, Time Limit Exceeded, Memory Limit Exceeded, Runtime Error, Compile Error, Output Limit Exceeded, Unknown Error
	Memory                 int      `json:"memory"`             // 内存消耗 in bytes
	StatusMemory           string   `json:"status_memory"`      // 内存消耗
	MemoryPercentile       float64  `json:"memory_percentile"`  // 内存消耗击败百分比
	StatusRuntime          string   `json:"status_runtime"`     // 执行用时
	RuntimePercentile      float64  `json:"runtime_percentile"` // 用时击败百分比
	Lang                   string   `json:"lang"`
	PrettyLang             string   `json:"pretty_lang"`
	CodeAnswer             []string `json:"code_answer"`   // return values of our code
	CompileError           string   `json:"compile_error"` //
	FullCompileError       string   `json:"full_compile_error"`
	FullRuntimeError       string   `json:"full_runtime_error"`
	CompareResult          string   `json:"compare_result"`  // "111", 1 means correct, 0 means wrong
	CorrectAnswer          bool     `json:"correct_answer"`  // true means all passed
	CodeOutput             []string `json:"code_output"`     // output to stdout of our code
	StdOutputList          []string `json:"std_output_list"` // list of output to stdout, same as code_output
	TaskName               string   `json:"task_name"`
	TotalCorrect           int      `json:"total_correct"`   // number of correct answers
	TotalTestcases         int      `json:"total_testcases"` // number of test cases
	ElapsedTime            int      `json:"elapsed_time"`
	TaskFinishTime         int      `json:"task_finish_time"`
	RunSuccess             bool     `json:"run_success"` // true if run success
	FastSubmit             bool     `json:"fast_submit"`
	Finished               bool     `json:"finished"`
	ExpectedOutput         string   `json:"expected_output"`
	ExpectedCodeAnswer     []string `json:"expected_code_answer"`
	ExpectedCodeOutput     []string `json:"expected_code_output"`
	ExpectedElapsedTime    int      `json:"expected_elapsed_time"`
	ExpectedLang           string   `json:"expected_lang"`
	ExpectedMemory         int      `json:"expected_memory"`
	ExpectedRunSuccess     bool     `json:"expected_run_success"`
	ExpectedStatusCode     int      `json:"expected_status_code"`
	ExpectedStatusRuntime  string   `json:"expected_status_runtime"`
	ExpectedStdOutputList  []string `json:"expected_std_output_list"`
	ExpectedTaskFinishTime int      `json:"expected_task_finish_time"`
	ExpectedTaskName       string   `json:"expected_task_name"`
}

// func (r *TestCheckResult) Display(q *QuestionData) string {
// 	stdout := ""
// 	if len(r.CodeOutput) > 1 {
// 		stdout = "\nStdout:        " + strings.Join(r.CodeOutput, "↩ ")
// 	}
// 	switch StatusCode(r.StatusCode) {
// 	case Accepted:
// 		if r.CorrectAnswer {
// 			return fmt.Sprintf(
// 				"\n%s\n%s%s%s%s%s\n",
// 				colorGreen.Sprintf("√ %s", r.StatusMsg),
// 				fmt.Sprintf("\nPassed cases:  %s", formatCompare(r.CompareResult)),
// 				fmt.Sprintf("\nInput:         %s", strings.ReplaceAll(r.InputData, "\n", "↩ ")),
// 				fmt.Sprintf("\nOutput:        %s", strings.Join(r.CodeAnswer, "↩ ")),
// 				stdout,
// 				fmt.Sprintf("\nExpected:      %s", strings.Join(r.ExpectedCodeAnswer, "↩ ")),
// 			)
// 		} else {
// 			return fmt.Sprintf(
// 				"\n%s\n%s%s%s%s%s\n",
// 				colorRed.Sprint(" × Wrong Answer"),
// 				fmt.Sprintf("\nPassed cases:  %s", formatCompare(r.CompareResult)),
// 				fmt.Sprintf("\nInput:         %s", strings.ReplaceAll(r.InputData, "\n", "↩ ")),
// 				fmt.Sprintf("\nOutput:        %s", strings.Join(r.CodeAnswer, "↩ ")),
// 				stdout,
// 				fmt.Sprintf("\nExpected:      %s", strings.Join(r.ExpectedCodeAnswer, "↩ ")),
// 			)
// 		}
// 	case MemoryLimitExceeded, TimeLimitExceeded, OutputLimitExceeded:
// 		return colorYellow.Sprintf("\n × %s\n", r.StatusMsg)
// 	case RuntimeError:
// 		return fmt.Sprintf(
// 			"\n%s\n%s\n\n%s\n",
// 			colorRed.Sprintf(" × %s", r.StatusMsg),
// 			fmt.Sprintf("Passed cases:   %s", formatCompare(r.CompareResult)),
// 			colorFaint.Sprint(r.FullRuntimeError),
// 		)
// 	case CompileError:
// 		return fmt.Sprintf(
// 			"\n%s\n\n%s\n",
// 			colorRed.Sprintf(" × %s", r.StatusMsg),
// 			colorFaint.Sprint(r.FullCompileError),
// 		)
// 	default:
// 		return fmt.Sprintf("\n%s\n", colorRed.Sprintf(" × %s", r.StatusMsg))
// 	}
// }

// func (r *TestCheckResult) GetState() string {
// 	return r.State
// }
