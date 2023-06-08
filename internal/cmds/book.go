package cmds

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/render"
	"github.com/zrcoder/tdoc"
	"github.com/zrcoder/tdoc/docmgr"
)

var Book = &cli.Command{
	Name:   "book",
	Usage:  "generate and serve a web book for the picked questions",
	Action: bookAction,
}

func bookAction(context *cli.Context) error {
	docPath, err := local.Generate()
	if err != nil {
		return err
	}
	fmt.Println(render.Infof("Your book is generated in %s\n", docPath))

	mgr, err := docmgr.New(docPath)
	if err != nil {
		return err
	}
	return tdoc.Run(mgr.Docs())
}
