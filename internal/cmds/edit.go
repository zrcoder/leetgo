package cmds

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/local"
)

var Edit = &cli.Command{
	Name:      "edit",
	Usage:     "edit codes to solve the question",
	UsageText: "leetgo edit 1",
	Action:    editAction,
}

func editAction(context *cli.Context) error {
	if context.Args().Len() == 0 {
		return errors.New("please pass the question id")
	}

	cfg, err := config.Get()
	if err != nil {
		return err
	}

	id := context.Args().First()
	if !local.Exist(id) {
		return fmt.Errorf("Not picked yet, type `leetgo view %s`", id)
	}

	codeFile, markdownFile := local.GetCodeFile(cfg, id), local.GetMarkdownFile(cfg, id)
	cmd := exec.Command(config.GetEditorCmd(cfg.Editor), "-p", codeFile, markdownFile)
	if config.IsGolang(cfg) {
		cmd.Args = append(cmd.Args, local.GetGoTestFile(cfg, id))
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
