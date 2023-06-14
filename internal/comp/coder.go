package comp

import (
	"fmt"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/utils/exec"
)

type coder struct {
	id string
}

func (c *coder) Run() error {
	c.id = regualarID(c.id)

	if !local.Exist(c.id) {
		return fmt.Errorf("not picked yet, type `leetgo view %s`", c.id)
	}

	cfg, err := config.Get()
	if err != nil {
		return err
	}

	codeFile := local.GetCodeFile(cfg, c.id)
	cmd, ops := config.GetEditorCmdOps(cfg.Editor)
	args := append(ops, codeFile)
	if config.IsGolang(cfg) {
		args = append(args, local.GetGoTestFile(cfg, c.id))
	}
	return exec.Run("", cmd, args...)
}
