package cmds

import (
	"errors"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var Code = &cli.Command{
	Name:      "code",
	Usage:     "edit codes to solve the question",
	UsageText: "leetgo code 1",
	Action:    codeAction,
}

func codeAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return errors.New("please pass the question id")
	}

	id := context.Args().First()
	return comp.NewCoder(id).Run()
}
