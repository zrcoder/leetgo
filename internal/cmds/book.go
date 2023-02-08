package cmds

import (
	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var Book = &cli.Command{
	Name:   "book",
	Usage:  "generate and serve a web book for the picked questions",
	Action: bookAction,
}

func bookAction(context *cli.Context) error {
	return comp.NewBook().Run()
}
