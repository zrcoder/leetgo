package cmds

import (
	"errors"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var Test = &cli.Command{
	Name:   "test",
	Usage:  "test your code locally and remotely",
	Action: testAction,
}

func testAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return errors.New("please pass the question id")
	}
	return comp.NewTester(context.Args().First()).Run()
}
