package cmds

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/render"
)

var Pick = &cli.Command{
	Name:      "pick",
	Usage:     "pick a question by id",
	UsageText: "leetgo pick 127",
	Action: func(context *cli.Context) error {
		if context.Args().Len() == 0 {
			return errors.New("need question id")
		}

		id := strings.Join(context.Args().Slice(), " ")
		mdData, path, err := mgr.Query(id)
		if err != nil {
			return err
		}

		fmt.Println(render.MarkDown(string(mdData)))
		if path != "" {
			fmt.Println(render.Successf("Stored in %s\n", path))
		}

		return nil
	},
}
