package cmds

import (
	"errors"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var Search = &cli.Command{
	Name:      "search",
	Usage:     "search questions by keywords or id",
	UsageText: "leetgo search 227",
	Action:    searchAction,
}

func searchAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return errors.New("please pass keywords")
	}

	key := strings.Join(context.Args().Slice(), " ")
	return comp.NewSearcher(key).Run()
}
