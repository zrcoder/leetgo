package cmds

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/comp"
)

var Search = &cli.Command{
	Name:      "search",
	Usage:     "search questions by keywords or id",
	UsageText: "leetgo search 227",
	Flags:     []cli.Flag{idFlag, keyFlag},
	Action:    searchAction,
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

func searchAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return errors.New("need key words")
	}

	key := strings.Join(context.Args().Slice(), " ")
	_, err := tea.NewProgram(comp.NewSearcher(key)).Run()
	return err
}
