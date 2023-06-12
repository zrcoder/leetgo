package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/cmds"
	"github.com/zrcoder/leetgo/utils/render"
)

func main() {
	app := &cli.App{
		Name:  "leetgo",
		Usage: "my app for Leetcode",
		Commands: []*cli.Command{
			cmds.Config,
			cmds.Search,
			cmds.View,
			cmds.Edit,
			cmds.Test,
			cmds.Submit,
			cmds.Book,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(render.Error("ERROR"), err.Error())
	}
}
