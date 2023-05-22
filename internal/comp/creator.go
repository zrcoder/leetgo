package comp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewCreator(projectName, lang, codeLang string) Component {
	return &Creator{
		projectName: projectName,
		lang:        lang,
		codeLang:    codeLang,
		spinner:     newSpinner("Creating"),
	}
}

type Creator struct {
	projectName, lang, codeLang string

	spinner *spinner.Spinner
}

func (c *Creator) Run() error {
	c.spinner.Start()
	path, err := c.creatProject()
	c.spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Println(render.Infof("Your leetcode project has been generated in %s\n", path))
	return nil
}

func (c *Creator) creatProject() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path = filepath.Join(path, c.projectName)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	cfg := &config.Config{
		Language: c.lang,
		CodeLang: c.codeLang,
	}
	return path, config.WriteTo(path, cfg)
}
