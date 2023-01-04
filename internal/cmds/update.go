package cmds

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/leetgo/internal/render"
)

var Update = &cli.Command{
	Name:      "update",
	Usage:     "update question in local",
	UsageText: "pass id to update a special question, no args to update the question list",
	Action:    updateAction,
}

func updateAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return updateList()
	}
	id := strings.Join(context.Args().Slice(), " ")
	return update(id)
}

func updateList() error {
	list, err := remote.GetList()
	if err != nil {
		return err
	}

	sps, err := local.WriteList(list)
	if err != nil {
		return err
	}

	fmt.Println(render.Successf("question list updated, there are %d questions now", len(sps)))

	return nil
}

func update(id string) error {
	list, err := local.ReadList()
	if err != nil {
		return err
	}

	mdData, path, err := mgr.QueryRemote(list, id)
	if err != nil {
		return err
	}

	fmt.Println(render.MarkDown(string(mdData)))
	if path != "" {
		fmt.Println(render.Successf("Stored in %s\n", path))
	}
	return nil
}
