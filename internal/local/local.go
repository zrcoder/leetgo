package local

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	data, _ := json.Marshal(all)
	err := os.WriteFile(filepath.Join(config.CfgDir, allFile), data, 0o640)
	if err != nil {
		log.Debug(err)
	}
	return err
}

func GetAll() (map[string]model.Meta, error) {
	path := filepath.Join(config.CfgDir, allFile)
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
	_, err := os.Stat(GetDir(id))
	return err == nil
}

func Write(question *model.Question) error {
	log.Debug("begin to write question in local", question.FrontendID, question.Title, question.TitleSlug)
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	if err = makeDir(question.FrontendID); err != nil {
		log.Debug(err)
		return err
	}
	if err = writeMarkdown(question); err != nil {
		log.Debug(err)
		return err
	}
	if err = writeCodeFile(question, cfg); err != nil {
		return err
	}

	question.MdContent = ""
	question.Stats = ""
	question.CodeDefinition = ""
	if err = writeMeta(question); err != nil {
		return err
	}

	if !config.IsGolang(cfg) {
		return nil
	}

	return judgeToGenGoModFile(cfg)
}

const TempMDFile = "tmp.md"

func WriteTempMd(id string, content []byte) error {
	log.Debug("write temp markdownfor question:", id)
	dir := GetDir(id)
	file := filepath.Join(dir, TempMDFile)
	log.Debug("file:", file)
	return os.WriteFile(file, content, 0o600)
}

func GetMarkdown(id string) ([]byte, error) {
	log.Debug("begin to read question.md in local, id:", id)
	dir := GetDir(id)
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
	path := getMetaFile(id)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := &model.Question{}
	err = json.Unmarshal(content, res)
	return res, err
}

func writeMeta(question *model.Question) error {
	data, _ := json.MarshalIndent(question, "", "  ")
	return os.WriteFile(getMetaFile(question.FrontendID), data, 0o640)
}

func writeMarkdown(question *model.Question) error {
	return os.WriteFile(GetMarkdownFile(question.FrontendID), []byte(question.MdContent), 0o640)
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
		buf.WriteString("package main\n\n")
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

func judgeToGenGoModFile(cfg *config.Config) error {
	if !config.IsGolang(cfg) {
		return nil
	}
	dir := config.CfgDir
	modFile := filepath.Join(dir, "go.mod")
	_, err := os.Stat(modFile)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return exec.Run(dir, "go", "mod", "init", "leetgo")
}

func makeDir(id string) error {
	err := os.MkdirAll(GetDir(id), 0o777)
	if err != nil {
		log.Debug(err)
	}
	return err
}

func GetDir(id string) string {
	return filepath.Join(config.CfgDir, id)
}

func GetCodeFile(cfg *config.Config, id string) string {
	return filepath.Join(config.CfgDir, id, codeFileName+config.GetCodeFileExt(cfg.CodeLang))
}

func GetMarkdownFile(id string) string {
	return filepath.Join(config.CfgDir, id, markdownFile)
}

func GetPickedQuestionIds() ([]string, error) {
	dir := config.CfgDir
	var ids []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, _ error) error {
		if path == dir {
			return nil
		}
		if d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			ids = append(ids, d.Name())
			return filepath.SkipDir
		}
		return nil
	})
	return ids, err
}

func getMetaFile(id string) string {
	return filepath.Join(config.CfgDir, id, metaFile)
}
