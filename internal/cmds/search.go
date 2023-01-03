package cmds

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/mgr"
	"github.com/zrcoder/leetgo/internal/render"
)

var Search = &cli.Command{
	Name:      "search",
	Usage:     "search questions by keywords or id",
	UsageText: "leetgo search 227",
	Flags:     []cli.Flag{idFlag, keyFlag},
	Action:    searchAction,
}

func searchAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return errors.New("need key words")
	}

	key := strings.Join(context.Args().Slice(), " ")
	sps, err := mgr.Search(key)
	if err != nil {
		return err
	}

	buf := &strings.Builder{}
	buf.WriteString("| No. | Title | Difficulty | Locked |\n")
	buf.WriteString("| --- | ----- | ---------- | ------ |\n")
	rowTmp := "| %s  | %s    | %s         | %s     |\n"
	lockCnt := 0
	lastQustion := ""
	for _, sp := range sps {
		locked := ""
		if sp.PaidOnly {
			locked = "🔒"
			lockCnt++
		}
		row := fmt.Sprintf(rowTmp, sp.Stat.CalculatedID, sp.Stat.QuestionTitle, sp.Difficulty.String(), locked)
		buf.WriteString(row)
		lastQustion = sp.Stat.CalculatedID
	}
	buf.WriteString(fmt.Sprintf("> total: %d, locked: %d\n", len(sps), lockCnt))
	buf.WriteString(fmt.Sprintf("> pick one? type like: `leetgo pick %s`", lastQustion))

	md := buf.String()

	fmt.Println(render.MarkDown(md))

	return nil
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
