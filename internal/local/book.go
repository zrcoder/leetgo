package local

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	rdbook "github.com/zrcoder/rdbook"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
)

const (
	bookMarkdownName = "markdown"
	bookWebName      = "web"
)

func WriteDocs() (string, error) {
	cfg, err := config.Get()
	if err != nil {
		return "", err
	}

	doc := &rdbook.Doc{}
	fmt.Println(doc)

	dir, err := getBookMarkdownDir(cfg)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		log.Trace(err)
		return "", err
	}

	return "", nil
}

func ReadAnswerCode(id string) ([]byte, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	path, err := getCodeFilePath(cfg, id)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Trace(err)
		return nil, err
	}

	index := bytes.Index(data, []byte(codeStartFlag))
	if index == -1 {
		return nil, newHintError(fmt.Sprintf("start flag %s not found", codeStartFlag), path)
	}
	data = data[index+len(codeStartFlag):]
	index = bytes.Index(data, []byte(codeEndFlag))
	if index == -1 {
		return nil, newHintError(fmt.Sprintf("end flag %s not found", codeEndFlag), path)
	}
	return data[:index], nil
}

func getBookMarkdownDir(cfg *config.Config) (string, error) {
	return filepath.Abs(filepath.Join(cfg.CacheDir, cfg.Language, bookMarkdownName))
}

func newHintError(info, path string) error {
	return fmt.Errorf("%s, pleanse check %s", info, path)
}
