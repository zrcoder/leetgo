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
		Name:   "leetgo",
		Action: action,
		Commands: []*cli.Command{
			cmds.Config,
			cmds.Search,
			cmds.Pick,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(render.Fail(err.Error()))
	}
}

func action(context *cli.Context) error {
	// TODO: help info
	md := `
## leetgo

My cli app to interact with Leetcode.
`
	fmt.Println(render.MarkDown(md))
	return nil
}
