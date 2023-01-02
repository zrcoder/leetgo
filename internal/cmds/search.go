package cmds

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/client"
	"github.com/zrcoder/leetgo/internal/trace"
)

var Search = &cli.Command{
	Name:  "search",
	Flags: []cli.Flag{idFlag, keyFlag},
	Action: func(context *cli.Context) error {
		if context.Args().Len() == 0 {
			return trace.Wrap(errors.New("need key words"))
		}

		key := strings.Join(context.Args().Slice(), " ")
		sps, err := client.Search(key)
		if err != nil {
			return err
		}

		// TODO: print the list
		fmt.Println(sps)

		return nil
	},
}

var idFlag = &cli.StringFlag{
	Name:    "id",
	Aliases: []string{"i"},
	Usage:   "the question id",
}

var keyFlag = &cli.StringFlag{
	Name:    "key",
	Aliases: []string{"k"},
	Usage:   "key words",
}
