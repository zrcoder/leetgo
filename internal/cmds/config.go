package cmds

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/render"
)

const (
	lanShortKey     = "l"
	projectShorkKey = "p"
)

var Config = &cli.Command{
	Name:        "config",
	Flags:       []cli.Flag{langFlag, projectFlag},
	Action:      setConfig,
	Subcommands: []*cli.Command{showCmd},
}

var langFlag = &cli.StringFlag{
	Name:    config.LangKey,
	Aliases: []string{lanShortKey},
	Value:   config.DefaultLanguage,
	Usage:   "language for the cli app",
}

var projectFlag = &cli.StringFlag{
	Name:    config.ProjectKey,
	Aliases: []string{projectShorkKey},
	Value:   config.DefaultProjectDir,
	Usage:   "project directory",
}

func setConfig(context *cli.Context) error {
	info, err := config.Get()
	if err != nil {
		return err
	}

	for _, key := range context.LocalFlagNames() {
		if _, exist := info[key]; !exist { // ignore the short flags
			continue
		}
		info[key] = context.String(key)
	}
	return config.Write(info)
}

var showCmd = &cli.Command{
	Name: "show",
	Action: func(context *cli.Context) error {
		data, err := config.Read()
		if err != nil {
			return err
		}
		res := fmt.Sprintf("```json\n%s\n```", data)
		fmt.Println(render.MarkDown(res))
		return nil
	},
}
