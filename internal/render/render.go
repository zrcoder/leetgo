package render

import (
	"fmt"
	"os"

	md "github.com/charmbracelet/glamour"
	sty "github.com/charmbracelet/lipgloss"
)

func Success(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#3cbb33")).Render(s)
}

func Fail(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#fd6164")).Render(s)
}

func MarkDown(s string) string {
	res, err := md.Render(s, "auto")
	if err != nil {
		fmt.Println(Fail(err.Error()))
		os.Exit(1)
	}
	return res
}
