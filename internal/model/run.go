package model

import (
	"fmt"
	"strings"

	"github.com/zrcoder/leetgo/internal/render"
)

type RunResult interface {
	Display() string
	Result() string
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

func (tr *TestCheckResult) Display() string {
	stdout := ""
	if len(tr.CodeOutput) > 1 {
		stdout = "\nStdout:        " + strings.Join(tr.CodeOutput, "↩ ")
	}
	switch StatusCode(tr.StatusCode) {
	case Accepted:
		if tr.CorrectAnswer {
			return fmt.Sprintf(
				"\n%s\n%s%s%s%s%s\n",
				render.Infof("√ %s", tr.StatusMsg),
				fmt.Sprintf("\nPassed cases:  %s", formatCompare(tr.CompareResult)),
				fmt.Sprintf("\nInput:         %s", strings.ReplaceAll(tr.InputData, "\n", "↩ ")),
				fmt.Sprintf("\nOutput:        %s", strings.Join(tr.CodeAnswer, "↩ ")),
				stdout,
				fmt.Sprintf("\nExpected:      %s", strings.Join(tr.ExpectedCodeAnswer, "↩ ")),
			)
		} else {
			return fmt.Sprintf(
				"\n%s\n%s%s%s%s%s\n",
				render.Fatal(" × Wrong Answer"),
				fmt.Sprintf("\nPassed cases:  %s", formatCompare(tr.CompareResult)),
				fmt.Sprintf("\nInput:         %s", strings.ReplaceAll(tr.InputData, "\n", "↩ ")),
				fmt.Sprintf("\nOutput:        %s", strings.Join(tr.CodeAnswer, "↩ ")),
				stdout,
				fmt.Sprintf("\nExpected:      %s", strings.Join(tr.ExpectedCodeAnswer, "↩ ")),
			)
		}
	case MemoryLimitExceeded, TimeLimitExceeded, OutputLimitExceeded:
		return render.Warnf("\n × %s\n", tr.StatusMsg)
	case RuntimeError:
		return fmt.Sprintf(
			"\n%s\n%s\n\n%s\n",
			render.Errorf(" × %s", tr.StatusMsg),
			fmt.Sprintf("Passed cases:   %s", formatCompare(tr.CompareResult)),
			render.Trace(tr.FullRuntimeError),
		)
	case CompileError:
		return fmt.Sprintf(
			"\n%s\n\n%s\n",
			render.Errorf(" × %s", tr.StatusMsg),
			render.Trace(tr.FullCompileError),
		)
	default:
		return fmt.Sprintf("\n%s\n", render.Errorf(" × %s", tr.StatusMsg))
	}
}

func (tr *TestCheckResult) Result() string {
	return tr.State
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

func (sr *SubmitCheckResult) Display() string {
	stdout := ""
	if len(sr.CodeOutput) > 1 {
		stdout = "\nStdout:        " + strings.ReplaceAll(sr.StdOutput, "\n", "↩ ")
	}
	switch StatusCode(sr.StatusCode) {
	case Accepted:
		return fmt.Sprintf(
			"\n%s\n%s%s%s\n",
			render.Infof("√ %s", sr.StatusMsg),
			fmt.Sprintf("\nPassed cases:  %d/%d", sr.TotalCorrect, sr.TotalTestcases),
			fmt.Sprintf("\nRuntime:       %s, better than %.0f%%", sr.StatusRuntime, sr.RuntimePercentile),
			fmt.Sprintf("\nMemory:        %s, better than %.0f%%", sr.StatusMemory, sr.MemoryPercentile),
		)
	case WrongAnswer:
		return fmt.Sprintf(
			"\n%s\n%s%s%s%s%s\n",
			render.Errorf(" × Wrong Answer"),
			fmt.Sprintf("\nPassed cases:  %d/%d", sr.TotalCorrect, sr.TotalTestcases),
			fmt.Sprintf("\nLast case:     %s", strings.ReplaceAll(sr.LastTestcase, "\n", "↩ ")),
			fmt.Sprintf("\nOutput:        %s", strings.ReplaceAll(sr.CodeOutput, "\n", "↩ ")),
			stdout,
			fmt.Sprintf("\nExpected:      %s", strings.ReplaceAll(sr.ExpectedOutput, "\n", "↩ ")),
		)
	case MemoryLimitExceeded, TimeLimitExceeded, OutputLimitExceeded:
		return fmt.Sprintf(
			"\n%s\n%s%s\n",
			render.Warnf("\n × %s\n", sr.StatusMsg),
			fmt.Sprintf("\nPassed cases:  %d/%d", sr.TotalCorrect, sr.TotalTestcases),
			fmt.Sprintf("\nLast case:     %s", sr.LastTestcase),
		)
	case RuntimeError:
		return fmt.Sprintf(
			"\n%s\n%s\n\n%s\n",
			render.Errorf(" × %s", sr.StatusMsg),
			fmt.Sprintf("Passed cases:   %s", formatCompare(sr.CompareResult)),
			render.Tracef(sr.FullRuntimeError),
		)
	case CompileError:
		return fmt.Sprintf(
			"\n%s\n\n%s\n",
			render.Errorf(" × %s", sr.StatusMsg),
			render.Tracef(sr.FullCompileError),
		)
	default:
		return fmt.Sprintf("\n%s\n", render.Errorf(" × %s", sr.StatusMsg))
	}
}

func (sr *SubmitCheckResult) Result() string {
	return sr.State
}
