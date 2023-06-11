package model

import (
	"fmt"
	"strings"

	"github.com/zrcoder/leetgo/internal/render"
)

type StatusCode int

const (
	Accepted            StatusCode = 10
	WrongAnswer         StatusCode = 11
	MemoryLimitExceeded StatusCode = 12
	OutputLimitExceeded StatusCode = 13
	TimeLimitExceeded   StatusCode = 14
	RuntimeError        StatusCode = 15
	CompileError        StatusCode = 20
)

type TestResult struct {
	InterpretExpectedId string `json:"interpret_expected_id"`
	InterpretId         string `json:"interpret_id"`
	TestCase            string `json:"test_case"`
}

type CheckResult interface {
	Display(q *Question) string
	GetState() string
}

type TestCheckResult struct {
	InputData          string
	State              string   `json:"state"` // STARTED, SUCCESS
	StatusCode         int      `json:"status_code"`
	StatusMsg          string   `json:"status_msg"`  // Accepted, Wrong Answer, Time Limit Exceeded, Memory Limit Exceeded, Runtime Error, Compile Error, Output Limit Exceeded, Unknown Error
	CodeAnswer         []string `json:"code_answer"` // return values of our code
	FullCompileError   string   `json:"full_compile_error"`
	FullRuntimeError   string   `json:"full_runtime_error"`
	CompareResult      string   `json:"compare_result"` // "111", 1 means correct, 0 means wrong
	CorrectAnswer      bool     `json:"correct_answer"` // true means all passed
	CodeOutput         []string `json:"code_output"`    // output to stdout of our code
	ExpectedCodeAnswer []string `json:"expected_code_answer"`
}

func (r *TestCheckResult) Display() string {
	stdout := ""
	if len(r.CodeOutput) > 1 {
		stdout = "\nStdout:        " + strings.Join(r.CodeOutput, "↩ ")
	}
	switch StatusCode(r.StatusCode) {
	case Accepted:
		if r.CorrectAnswer {
			return fmt.Sprintf(
				"\n%s\n%s%s%s%s%s\n",
				render.Infof("√ %s", r.StatusMsg),
				fmt.Sprintf("\nPassed cases:  %s", formatCompare(r.CompareResult)),
				fmt.Sprintf("\nInput:         %s", strings.ReplaceAll(r.InputData, "\n", "↩ ")),
				fmt.Sprintf("\nOutput:        %s", strings.Join(r.CodeAnswer, "↩ ")),
				stdout,
				fmt.Sprintf("\nExpected:      %s", strings.Join(r.ExpectedCodeAnswer, "↩ ")),
			)
		} else {
			return fmt.Sprintf(
				"\n%s\n%s%s%s%s%s\n",
				render.Fatal(" × Wrong Answer"),
				fmt.Sprintf("\nPassed cases:  %s", formatCompare(r.CompareResult)),
				fmt.Sprintf("\nInput:         %s", strings.ReplaceAll(r.InputData, "\n", "↩ ")),
				fmt.Sprintf("\nOutput:        %s", strings.Join(r.CodeAnswer, "↩ ")),
				stdout,
				fmt.Sprintf("\nExpected:      %s", strings.Join(r.ExpectedCodeAnswer, "↩ ")),
			)
		}
	case MemoryLimitExceeded, TimeLimitExceeded, OutputLimitExceeded:
		return render.Warnf("\n × %s\n", r.StatusMsg)
	case RuntimeError:
		return fmt.Sprintf(
			"\n%s\n%s\n\n%s\n",
			render.Errorf(" × %s", r.StatusMsg),
			fmt.Sprintf("Passed cases:   %s", formatCompare(r.CompareResult)),
			render.Trace(r.FullRuntimeError),
		)
	case CompileError:
		return fmt.Sprintf(
			"\n%s\n\n%s\n",
			render.Errorf(" × %s", r.StatusMsg),
			render.Trace(r.FullCompileError),
		)
	default:
		return fmt.Sprintf("\n%s\n", render.Errorf(" × %s", r.StatusMsg))
	}
}

func formatCompare(s string) string {
	var sb strings.Builder
	for _, c := range s {
		if c == '1' {
			sb.WriteString(render.Info("√"))
		} else {
			sb.WriteString(render.Error("×"))
		}
	}
	return sb.String()
}

type SubmitCheckResult struct {
	CodeOutput        string  `json:"code_output"` // answers of our code
	CompareResult     string  `json:"compare_result"`
	ExpectedOutput    string  `json:"expected_output"`
	LastTestcase      string  `json:"last_testcase"`
	MemoryPercentile  float64 `json:"memory_percentile"`
	RuntimePercentile float64 `json:"runtime_percentile"`
	State             string  `json:"state"`
	StatusCode        int     `json:"status_code"`
	StatusMemory      string  `json:"status_memory"`
	StatusMsg         string  `json:"status_msg"`
	StatusRuntime     string  `json:"status_runtime"`
	StdOutput         string  `json:"std_output"`
	TotalCorrect      int     `json:"total_correct"`
	TotalTestcases    int     `json:"total_testcases"`
	FullCompileError  string  `json:"full_compile_error"`
	FullRuntimeError  string  `json:"full_runtime_error"`
}

func (r *SubmitCheckResult) Display() string {
	stdout := ""
	if len(r.CodeOutput) > 1 {
		stdout = "\nStdout:        " + strings.ReplaceAll(r.StdOutput, "\n", "↩ ")
	}
	switch StatusCode(r.StatusCode) {
	case Accepted:
		return fmt.Sprintf(
			"\n%s\n%s%s%s\n",
			render.Infof("√ %s", r.StatusMsg),
			fmt.Sprintf("\nPassed cases:  %d/%d", r.TotalCorrect, r.TotalTestcases),
			fmt.Sprintf("\nRuntime:       %s, better than %.0f%%", r.StatusRuntime, r.RuntimePercentile),
			fmt.Sprintf("\nMemory:        %s, better than %.0f%%", r.StatusMemory, r.MemoryPercentile),
		)
	case WrongAnswer:
		return fmt.Sprintf(
			"\n%s\n%s%s%s%s%s\n",
			render.Errorf(" × Wrong Answer"),
			fmt.Sprintf("\nPassed cases:  %d/%d", r.TotalCorrect, r.TotalTestcases),
			fmt.Sprintf("\nLast case:     %s", strings.ReplaceAll(r.LastTestcase, "\n", "↩ ")),
			fmt.Sprintf("\nOutput:        %s", strings.ReplaceAll(r.CodeOutput, "\n", "↩ ")),
			stdout,
			fmt.Sprintf("\nExpected:      %s", strings.ReplaceAll(r.ExpectedOutput, "\n", "↩ ")),
		)
	case MemoryLimitExceeded, TimeLimitExceeded, OutputLimitExceeded:
		return fmt.Sprintf(
			"\n%s\n%s%s\n",
			render.Warnf("\n × %s\n", r.StatusMsg),
			fmt.Sprintf("\nPassed cases:  %d/%d", r.TotalCorrect, r.TotalTestcases),
			fmt.Sprintf("\nLast case:     %s", r.LastTestcase),
		)
	case RuntimeError:
		return fmt.Sprintf(
			"\n%s\n%s\n\n%s\n",
			render.Errorf(" × %s", r.StatusMsg),
			fmt.Sprintf("Passed cases:   %s", formatCompare(r.CompareResult)),
			render.Tracef(r.FullRuntimeError),
		)
	case CompileError:
		return fmt.Sprintf(
			"\n%s\n\n%s\n",
			render.Errorf(" × %s", r.StatusMsg),
			render.Tracef(r.FullCompileError),
		)
	default:
		return fmt.Sprintf("\n%s\n", render.Errorf(" × %s", r.StatusMsg))
	}
}
