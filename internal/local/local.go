package local

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cweill/gotests"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/utils/exec"
)

const (
	allFile      = "all.json"
	codeFileName = "solution"
	markdownFile = "question.md"
	metaFile     = "question.json"
)

type CommentFlag struct {
	NoteStart, NoteEnd, CodeStart, CodeEnd string
}

var goCommentFlag = CommentFlag{
	NoteStart: "/* @note start",
	NoteEnd:   "@note end */",
	CodeStart: "// @submit start",
	CodeEnd:   "// @submit end",
}

var pythonCommentFlag = CommentFlag{
	NoteStart: `"""@note start`,
	NoteEnd:   `@note end"""`,
	CodeStart: "# @submit start",
	CodeEnd:   "# @submit end",
}

var commentFlagDict = map[string]CommentFlag{
	"go":      goCommentFlag,
	"golang":  goCommentFlag,
	"py":      pythonCommentFlag,
	"python":  pythonCommentFlag,
	"python3": pythonCommentFlag,
}

var ErrNotCached = errors.New("not cached in local yet")

func WriteAll(all map[string]model.Meta) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	languageDir := filepath.Join(config.CfgDir, cfg.Language)
	err = os.MkdirAll(languageDir, os.ModePerm)
	if err != nil {
		return err
	}
	data, _ := json.Marshal(all)
	err = os.WriteFile(filepath.Join(languageDir, allFile), data, 0o640)
	if err != nil {
		log.Debug(err)
	}
	return err
}

func GetAll() (map[string]model.Meta, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(config.CfgDir, cfg.Language, allFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = ErrNotCached
		}
		return nil, err
	}
	res := map[string]model.Meta{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Debug(err)
	}
	return res, err
}

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
	log.Debug("begin to write question in local", question.FrontendID, question.Title, question.TitleSlug)
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	if err = makeDir(cfg, question.FrontendID); err != nil {
		log.Debug(err)
		return err
	}
	if err = writeMarkdown(question, cfg); err != nil {
		log.Debug(err)
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
	commentFlag := commentFlagDict[cfg.CodeLang]
	return getFromCodeFile(cfg, id, commentFlag.CodeStart, commentFlag.CodeEnd)
}

func GetNotes(cfg *config.Config, id string) ([]byte, error) {
	log.Debug("begin to read typed notes in local, id:", id)
	commentFlag := commentFlagDict[cfg.CodeLang]

	return getFromCodeFile(cfg, id, commentFlag.NoteStart, commentFlag.NoteEnd)
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
	j := bytes.LastIndex(content, []byte(end))
	if j == -1 {
		return content, nil
	}

	return content[i+len(start) : j], nil
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
	return os.WriteFile(getMetaFile(cfg, question.FrontendID), data, 0o640)
}

func writeMarkdown(question *model.Question, cfg *config.Config) error {
	return os.WriteFile(GetMarkdownFile(cfg, question.FrontendID), []byte(question.MdContent), 0o640)
}

func writeCodeFile(question *model.Question, cfg *config.Config) error {
	codes, err := question.ParseCodes()
	if err != nil {
		return err
	}
	id := question.FrontendID
	codeFile := GetCodeFile(cfg, id)
	buf := bytes.NewBuffer(nil)
	if config.IsGolang(cfg) {
		line := fmt.Sprintf("package _%s\n\n", strings.ReplaceAll(question.TitleSlug, "-", "_"))
		buf.WriteString(line)
	}
	commentFlag := commentFlagDict[cfg.CodeLang]
	buf.WriteString(commentFlag.NoteStart)
	buf.WriteString("\n\n")
	buf.WriteString(commentFlag.NoteEnd + "\n\n")
	codeLang := config.LeetcodeLang(cfg.CodeLang)
	buf.WriteString(commentFlag.CodeStart)
	buf.WriteString("\n")
	for _, v := range codes {
		if v.Value == codeLang {
			buf.WriteString(v.DefaultCode)
			buf.WriteString("\n")
			break
		}
	}
	buf.WriteString(commentFlag.CodeEnd)
	return os.WriteFile(codeFile, buf.Bytes(), 0o640)
}

func writeGoTestFile(question *model.Question, cfg *config.Config) error {
	codePath := GetCodeFile(cfg, question.FrontendID)
	testPath := GetGoTestFile(cfg, question.FrontendID)
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
	return os.WriteFile(testPath, content, 0o640)
}

func genGoModFile(question *model.Question, cfg *config.Config) error {
	return exec.Run(GetDir(cfg, question.FrontendID), "go", "mod", "init", cfg.Language+"-"+question.TitleSlug)
}

func makeDir(cfg *config.Config, id string) error {
	err := os.MkdirAll(GetDir(cfg, id), 0o777)
	if err != nil {
		log.Debug(err)
	}
	return err
}

func GetDir(cfg *config.Config, id string) string {
	return filepath.Join(config.CfgDir, cfg.Language, cfg.CodeLang, id)
}

func GetCodeFile(cfg *config.Config, id string) string {
	return filepath.Join(config.CfgDir, cfg.Language, cfg.CodeLang, id, codeFileName+config.GetCodeFileExt(cfg.CodeLang))
}

func GetGoTestFile(cfg *config.Config, id string) string {
	return filepath.Join(config.CfgDir, cfg.Language, cfg.CodeLang, id, codeFileName+"_test"+config.GetCodeFileExt(cfg.CodeLang))
}

func GetMarkdownFile(cfg *config.Config, id string) string {
	return filepath.Join(config.CfgDir, cfg.Language, cfg.CodeLang, id, markdownFile)
}

func GetPickedQuestionIds(cfg *config.Config) ([]string, error) {
	dir := filepath.Join(config.CfgDir, cfg.Language, cfg.CodeLang)
	var ids []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, _ error) error {
		if path == dir {
			return nil
		}
		if d.IsDir() {
			ids = append(ids, d.Name())
			return nil
		}
		return nil
	})
	return ids, err
}

func getMetaFile(cfg *config.Config, id string) string {
	return filepath.Join(config.CfgDir, cfg.Language, cfg.CodeLang, id, metaFile)
}
