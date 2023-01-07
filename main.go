package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/cmds"
	"github.com/zrcoder/leetgo/internal/render"
)

var app *cli.App

func main() {
	app = &cli.App{
		Name:   "leetgo",
		Usage:  "my app for Leetcode",
		Action: action,
		Commands: []*cli.Command{
			cmds.Config,
			cmds.Search,
			cmds.Pick,
			cmds.Update,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(render.Fail(err.Error()))
	}
}

func action(context *cli.Context) error {
	md, _ := app.ToMarkdown()
	fmt.Println(render.MarkDown(md))
	return nil
}
