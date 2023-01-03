package render

import (
	"fmt"
	"os"

	md "github.com/charmbracelet/glamour"
	sty "github.com/charmbracelet/lipgloss"
)

func Success(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#12865f")).Render(s)
}

func Successf(format string, v ...any) string {
	s := fmt.Sprintf(format, v...)
	return Success(s)
}

func Fail(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#fd6164")).Render(s)
}

func Debug(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#ffa500")).Render(s)
}

func MarkDown(s string) string {
	res, err := md.Render(s, "auto")
	if err != nil {
		fmt.Println(Fail(err.Error()))
		os.Exit(1)
	}
	return res
}
