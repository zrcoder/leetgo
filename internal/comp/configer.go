package comp

import (
	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/utils/huh"
)

type configer struct {
	cfg         *config.Config
	showFunc    func(*config.Config)
	shouldWrite bool
}

func (c *configer) Run() error {
	_, err := config.Get()
	if err != nil {
		if err != config.ErrConfigNotExist {
			return err
		}
		init := true
		err = huh.NewConfirm("No config found, initialize?", "", &init).Run()
		if err != nil {
			return err
		}
		if !init {
			return nil
		}
		c.shouldWrite = true
	}
	if c.shouldWrite {
		c.cfg, err = config.Write(c.cfg)
	} else {
		c.cfg, err = config.Get()
	}
	if err != nil {
		return err
	}
	c.showFunc(c.cfg)
	return nil
}
