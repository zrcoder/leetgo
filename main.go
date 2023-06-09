package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/cmds"
	"github.com/zrcoder/leetgo/internal/render"
)

func main() {
	app := &cli.App{
		Name:  "leetgo",
		Usage: "my app for Leetcode",
		Commands: []*cli.Command{
			cmds.Init,
			cmds.Config,
			cmds.Search,
			cmds.View,
			cmds.Edit,
			cmds.Test,
			cmds.Book,
		},
	}
	app.Action = func(ctx *cli.Context) error {
		md, _ := app.ToMarkdown()
		fmt.Println(render.MarkDown(md))
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(render.Error(err.Error()))
	}
}
