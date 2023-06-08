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

	CodeStartFlag4Editor = "@submit start"
	codeStartFlag        = "// @submit start\n"
	codeEndFlag          = "// @submit end\n"
)

var (
	ErrNotCached = errors.New("not cached in local yet")
)

var (
	extesionDic = map[string]string{
		"go":     ".go",
		"golang": ".go",
		"java":   ".java",
		"python": ".py",
		// TODO, support other languages
	}
)

func Read(id string) (*model.Question, error) {
	log.Trace("begin to read question in local, id:", id)
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(GetCodeFile(cfg, id))
	if err != nil {
		return nil, err
	}
	res := &model.Question{}
	err = json.Unmarshal(data, res)
	if err != nil {
		log.Trace(err)
		return nil, err
	}
	return res, nil
}

func Write(question *model.Question) error {
	log.Trace("begin to write question in local, id:", question.ID)
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	codes, err := question.ParseCodes()
	if err != nil {
		return err
	}
	id := question.ID
	err = makeDir(cfg, id)
	if err != nil {
		return err
	}
	codeFile := GetCodeFile(cfg, id)
	buf := bytes.NewBuffer(nil)
	isGo := cfg.CodeLang == config.CodeLangGo
	if isGo {
		buf.WriteString("package solution\n\n")
	}
	buf.WriteString("/*\n")
	buf.WriteString(question.MdContent)
	buf.WriteString("\n*/\n\n")
	buf.WriteString(codeStartFlag)
	for _, v := range codes {
		if v.Value == cfg.CodeLang {
			buf.WriteString(v.DefaultCode)
			buf.WriteString("\n")
			break
		}
	}
	buf.WriteString(codeEndFlag)
	err = os.WriteFile(codeFile, buf.Bytes(), 0640)
	if err != nil {
		log.Trace(err)
		return err
	}
	if !isGo {
		return nil
	}
	testPath := filepath.Join(GetDir(cfg, id), codeFileName+"_test.go")
	_ = os.Remove(testPath) // need remove the test file when update
	tests, err := gotests.GenerateTests(codeFile, nil)
	if err != nil {
		return err
	}
	if len(tests) == 0 {
		return errors.New("no tests generated")
	}
	sample := fmt.Sprintf("\n/* sample test case:\n%s\n*/\n", question.SampleTestCase)
	content := append(tests[0].Output, []byte(sample)...)
	err = os.WriteFile(testPath, content, 0640)
	log.Trace(err)
	return err
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
	return filepath.Join(cfg.Language, cfg.CodeLang, id, codeFileName+extesionDic[cfg.CodeLang])
}
