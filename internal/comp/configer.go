package comp

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/zrcoder/leetgo/internal/config"
)

type configer struct {
	cfg         *config.Config
	shouldWrite bool
	showFunc    func(*config.Config)
}

func (c *configer) Run() error {
	_, err := config.Get()
	if err != nil {
		if err != config.ErrConfigNotExist {
			return err
		}
		init := false
		prompt := &survey.Confirm{
			Message: "No config found in the current directory, initial?",
			Default: true,
			Help:    "Initial the current directory as your leetgo project.",
		}
		if err = survey.AskOne(prompt, &init); err != nil {
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
