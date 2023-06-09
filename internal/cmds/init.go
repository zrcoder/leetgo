package cmds

import (
	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/config"
)

var Init = &cli.Command{
	Name:      "init",
	Usage:     "init the project",
	UsageText: "leetgo init",
	Action:    initAction,
	Flags: []cli.Flag{
		langFlag,
		codeLangFlag,
		editorFlag,
	},
}

func initAction(context *cli.Context) error {
	cfg := &config.Config{
		Language: context.String(langFlag.Name),
		CodeLang: context.String(codeLangFlag.Name),
		Editor:   context.String(editorFlag.Name),
	}
	return writeAndShow(cfg)
}
