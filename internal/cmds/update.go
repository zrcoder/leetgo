package cmds

import (
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var Update = &cli.Command{
	Name:      "update",
	Usage:     "update question in local",
	UsageText: "pass id to update a special question, no args to update the question list",
	Action:    updateAction,
}

func updateAction(context *cli.Context) error {
	id := ""
	if context.Args().Len() > 0 {
		id = strings.Join(context.Args().Slice(), " ")
	}
	return comp.NewUpdater(id).Run()
}
