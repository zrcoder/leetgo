package cmds

import (
	"errors"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var Submit = &cli.Command{
	Name:   "submit",
	Usage:  "submit your codes",
	Action: submitAction,
}

func submitAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return errors.New("please pass the question id")
	}
	return comp.NewSubmiter(strings.Join(context.Args().Slice(), " ")).Run()
}
