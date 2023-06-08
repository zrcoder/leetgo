package cmds

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/render"
)

var Init = &cli.Command{
	Name:      "init",
	Usage:     "init the project",
	UsageText: "leetgo init",
	Action:    initAction,
	Flags: []cli.Flag{
		langFlag,
		codeLangFlag,
	},
}

func initAction(context *cli.Context) error {
	cfg := &config.Config{
		Language: context.String(langFlag.Name),
		CodeLang: context.String(codeLangFlag.Name),
	}
	err := config.Write(cfg)
	if err != nil {
		return err
	}
	fmt.Println(render.Info("Succeed"))
	return nil
}
