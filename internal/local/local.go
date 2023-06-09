package local

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cweill/gotests"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

const (
	codeFileName = "solution"
	markdownFile = "question.md"
	metaFile     = "meta.json"

	codeStartFlag = "// @submit start\n"
	codeEndFlag   = "// @submit end\n"
)

var (
	ErrNotCached = errors.New("not cached in local yet")
)

func Exist(id string) bool {
	cfg, err := config.Get()
	if err != nil {
		log.Trace(err)
		return false
	}
	_, err = os.Stat(GetDir(cfg, id))
	return err == nil
}

func ReadMarkdown(id string) ([]byte, error) {
	log.Trace("begin to read question.md in local, id:", id)
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	dir := GetDir(cfg, id)
	return os.ReadFile(filepath.Join(dir, markdownFile))
}

func Write(question *model.Question) error {
	log.Trace("begin to write question in local, id:", question.ID)
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	if err = makeDir(cfg, question.ID); err != nil {
		return err
	}
	if err = writeMeta(question, cfg); err != nil {
		return err
	}
	if err = writeMarkdown(question, cfg); err != nil {
		return err
	}
	if err = writeCodeFile(question, cfg); err != nil {
		return err
	}
	if config.IsGolang(cfg) {
		return writeGoTestFile(question, cfg)
	}
	return nil
}

func writeMeta(question *model.Question, cfg *config.Config) error {
	data, _ := json.MarshalIndent(question.Meta, "", "  ")
	return os.WriteFile(getMetaFile(cfg, question.ID), data, 0640)
}

func writeMarkdown(question *model.Question, cfg *config.Config) error {
	return os.WriteFile(GetMarkdownFile(cfg, question.ID), []byte(question.MdContent), 0640)
}

func writeCodeFile(question *model.Question, cfg *config.Config) error {
	codes, err := question.ParseCodes()
	if err != nil {
		return err
	}
	id := question.ID
	codeFile := GetCodeFile(cfg, id)
	buf := bytes.NewBuffer(nil)
	if config.IsGolang(cfg) {
		line := fmt.Sprintf("package lc%s\n\n", id)
		buf.WriteString(line)
	}
	codeLang := cfg.CodeLang
	if codeLang == "go" {
		codeLang = "golang"
	}
	buf.WriteString(codeStartFlag)
	for _, v := range codes {
		if v.Value == codeLang {
			buf.WriteString(v.DefaultCode)
			buf.WriteString("\n")
			break
		}
	}
	buf.WriteString(codeEndFlag)
	return os.WriteFile(codeFile, buf.Bytes(), 0640)
}

func writeGoTestFile(question *model.Question, cfg *config.Config) error {
	codePath := GetCodeFile(cfg, question.ID)
	testPath := GetGoTestFile(cfg, question.ID)
	_ = os.Remove(testPath) // need remove the test file when update
	tests, err := gotests.GenerateTests(codePath, nil)
	if err != nil {
		return err
	}
	if len(tests) == 0 {
		return errors.New("no tests generated")
	}
	sample := fmt.Sprintf("\n/* sample test case:\n%s\n*/\n", question.SampleTestCase)
	content := tests[0].Output
	const todoFlag = "// TODO: Add test cases."
	content = bytes.Replace(content, []byte(todoFlag), []byte(todoFlag+sample), 1)
	return os.WriteFile(testPath, content, 0640)
}

func makeDir(cfg *config.Config, id string) error {
	err := os.MkdirAll(GetDir(cfg, id), 0777)
	log.Trace(err)
	return err
}
func GetDir(cfg *config.Config, id string) string {
	return filepath.Join(cfg.Language, cfg.CodeLang, id)
}
func GetCodeFile(cfg *config.Config, id string) string {
	return filepath.Join(cfg.Language, cfg.CodeLang, id, codeFileName+config.GetCodeFileExt(cfg.CodeLang))
}
func GetGoTestFile(cfg *config.Config, id string) string {
	return filepath.Join(cfg.Language, cfg.CodeLang, id, codeFileName+"_test"+config.GetCodeFileExt(cfg.CodeLang))
}
func GetMarkdownFile(cfg *config.Config, id string) string {
	return filepath.Join(cfg.Language, cfg.CodeLang, id, markdownFile)
}
func getMetaFile(cfg *config.Config, id string) string {
	return filepath.Join(cfg.Language, cfg.CodeLang, id, metaFile)
}
