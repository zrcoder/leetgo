package local

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cweill/gotests"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/exec"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

const (
	codeFileName = "solution"
	markdownFile = "question.md"
	metaFile     = "question.json"

	noteStartFlag = "/* @note start\n"
	noteEndFlag   = "@note end */\n"
	codeStartFlag = "// @submit start\n"
	codeEndFlag   = "// @submit end\n"
)

var (
	ErrNotCached = errors.New("not cached in local yet")
)

func Exist(id string) bool {
	cfg, err := config.Get()
	if err != nil {
		log.Debug(err)
		return false
	}
	_, err = os.Stat(GetDir(cfg, id))
	return err == nil
}

func Write(question *model.Question) error {
	log.Debug("begin to write question in local, id:", question.ID)
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	if err = makeDir(cfg, question.ID); err != nil {
		return err
	}
	if err = writeMarkdown(question, cfg); err != nil {
		return err
	}
	if err = writeCodeFile(question, cfg); err != nil {
		return err
	}

	question.MdContent = ""
	question.Stats = ""
	question.CodeDefinition = ""
	if err = writeMeta(question, cfg); err != nil {
		return err
	}

	if !config.IsGolang(cfg) {
		return nil
	}
	// write xxx_test.go and go.mod for local test
	if err = writeGoTestFile(question, cfg); err != nil {
		return err
	}
	return genGoModFile(question, cfg)
}

func GetMarkdown(id string) ([]byte, error) {
	log.Debug("begin to read question.md in local, id:", id)
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	dir := GetDir(cfg, id)
	return os.ReadFile(filepath.Join(dir, markdownFile))
}

func GetTypedCode(cfg *config.Config, id string) ([]byte, error) {
	log.Debug("begin to read typed code in local, id:", id)
	return getFromCodeFile(cfg, id, codeStartFlag, codeEndFlag)
}
func getNotes(cfg *config.Config, id string) ([]byte, error) {
	log.Debug("begin to read typed notes in local, id:", id)
	data, err := getFromCodeFile(cfg, id, noteStartFlag, noteEndFlag)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(data, []byte(codeStartFlag)) {
		return nil, nil
	}
	return data, nil
}

func getFromCodeFile(cfg *config.Config, id, start, end string) ([]byte, error) {
	path := GetCodeFile(cfg, id)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	i := bytes.Index(content, []byte(start))
	if i == -1 {
		return content, nil
	}
	content = content[i+len(start):]
	i = bytes.LastIndex(content, []byte(end))
	if i == -1 {
		return content, nil
	}
	return content[:i], nil
}

func GetQuestion(cfg *config.Config, id string) (*model.Question, error) {
	log.Debug("begin to read question.json in local, id:", id)
	path := getMetaFile(cfg, id)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := &model.Question{}
	err = json.Unmarshal(content, res)
	return res, err
}

func writeMeta(question *model.Question, cfg *config.Config) error {
	data, _ := json.MarshalIndent(question, "", "  ")
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
		line := fmt.Sprintf("package _%s\n\n", strings.ReplaceAll(question.Slug, "-", "_"))
		buf.WriteString(line)
	}
	buf.WriteString(noteStartFlag)
	buf.WriteString("\n")
	buf.WriteString(noteEndFlag)
	buf.WriteString("\n")
	codeLang := config.LeetcodeLang(cfg.CodeLang)
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

func genGoModFile(question *model.Question, cfg *config.Config) error {
	return exec.Run(GetDir(cfg, question.ID), "go", "mod", "init", cfg.Language+"-"+question.Slug)
}

func makeDir(cfg *config.Config, id string) error {
	err := os.MkdirAll(GetDir(cfg, id), 0777)
	log.Debug(err)
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
