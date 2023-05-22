package cmds

import (
	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
	"github.com/zrcoder/leetgo/internal/config"
)

var New = &cli.Command{
	Name:      "new",
	Usage:     "create a project to solve leetcode questions",
	UsageText: "leetgo new myLeetcode",
	Action:    newAction,
	Flags: []cli.Flag{
		langFlag,
		codeLangFlag,
		projectDirFlag,
	},
}

var (
	projectDirKey      = "directory"
	projectDirShortKey = "d"
)

var projectDirFlag = &cli.StringFlag{
	Name:     projectDirKey,
	Aliases:  []string{projectDirShortKey},
	Value:    config.DefaultCacheDir,
	Usage:    cacheDirUsage,
	Required: true,
}

func newAction(context *cli.Context) error {
	lang := context.String(langFlag.Name)
	codeLang := context.String(codeLangFlag.Name)
	name := context.String(projectDirKey)
	return comp.NewCreator(name, lang, codeLang).Run()
}
