package render

import (
	"fmt"
	"os"

	md "github.com/charmbracelet/glamour"
	sty "github.com/charmbracelet/lipgloss"
)

func Trace(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#3e94d2")).Render(s)
}

func Tracef(format string, v ...any) string {
	return Trace(fmt.Sprintf(format, v...))
}

func Info(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#5b9033")).Render(s)
}

func Infof(format string, v ...any) string {
	return Info(fmt.Sprintf(format, v...))
}

func Warn(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#a58821")).Render(s)
}

func Warnf(format string, v ...any) string {
	return Warn(fmt.Sprintf(format, v...))
}

func Error(s string) string {
	return sty.NewStyle().Foreground(sty.Color("#fc4355")).Render(s)
}

func Errorf(format string, v ...any) string {
	return Error(fmt.Sprintf(format, v...))
}

func Fatal(s string) string {
	return sty.NewStyle().Bold(true).Foreground(sty.Color("#fc4355")).Render(s)
}

func MarkDown(s string) string {
	res, err := md.Render(s, "auto")
	if err != nil {
		fmt.Println(Fatal(err.Error()))
		os.Exit(1)
	}
	return res
}
