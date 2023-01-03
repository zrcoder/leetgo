package local

import (
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
)

const (
	allFile = "all.json"
)

var (
	ErrNotCached = errors.New("not cached yet")
)

func ReadAll() (map[string]model.StatStatusPair, error) {
	path, err := getAllFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Dev("all.json", ErrNotCached)
			return nil, ErrNotCached
		}
		log.Dev(err)
		return nil, err
	}

	res := map[string]model.StatStatusPair{}
	err = json.Unmarshal(data, &res)
	log.Dev(err)
	return res, err
}

func WriteAll(all map[string]model.StatStatusPair) error {
	path, err := getAllFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(all)
	if err != nil {
		log.Dev(err)
		return err
	}

	err = os.WriteFile(path, data, 0640)
	log.Dev(err)
	return err
}

func Read(id string) ([]byte, string, error) {
	log.Dev("begin to search question in cache, id:", id)

	dir, err := getDir()
	if err != nil {
		return nil, "", err
	}

	log.Dev("begin to walk:", dir)
	var data []byte
	resPath := ""

	expectedPrefix := id + "-"
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err1 error) error {
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		log.Dev("current file is:", name)
		if strings.HasPrefix(name, expectedPrefix) && strings.HasSuffix(name, ".md") {
			log.Dev("found the cached markdown file for question", id)
			data, err = os.ReadFile(path)
			if err != nil {
				return err
			}
			resPath = path
			return filepath.SkipDir
		}
		return nil
	})

	if err == nil && data == nil {
		log.Dev("Not found in cache, qustion", id)
		return nil, "", ErrNotCached
	}
	log.Dev(err)
	return data, resPath, err
}

func Write(sp *model.StatStatusPair, question *model.Question, mdData []byte) (string, error) {
	cfg, err := config.Get()
	if err != nil {
		return "", err
	}

	dir, err := parseDir(cfg)
	if err != nil {
		return "", err
	}

	title := fmt.Sprintf("%s-%s", sp.Stat.CalculatedID, sp.Stat.QuestionTitleSlug)
	mdPath := filepath.Join(dir, title+".md")
	currentCodeLang := cfg[config.CodeLangKey]
	fileExtension := config.CodeLangExtensionDic[currentCodeLang]
	codeDir := filepath.Join(dir, cfg[config.CodeLangKey])
	err = os.MkdirAll(codeDir, 0777)
	if err != nil {
		log.Dev(err)
		return "", err
	}
	codePath := filepath.Join(codeDir, title+fileExtension)

	err = os.WriteFile(mdPath, mdData, 0640)
	if err != nil {
		log.Dev(err)
		return "", err
	}

	codes, err := question.ParseCodes()
	if err != nil {
		return "", err
	}

	for _, c := range codes {
		if c.Value == currentCodeLang {
			const beginFlag = "//leetcode submit region begin(Prohibit modification and deletion)"
			const endFlag = "//leetcode submit region end(Prohibit modification and deletion)\n"
			content := fmt.Sprintf("package leetgo\n\n%s\n%s\n%s", beginFlag, c.DefaultCode, endFlag)
			err = os.WriteFile(codePath, []byte(content), 0640)
			if err != nil {
				log.Dev(err)
				return "", err
			}

			if c.Value != config.CodeLangGo {
				return mdPath, nil
			}

			tess, err := gotests.GenerateTests(codePath, nil)
			if err != nil {
				log.Dev(err)
				return "", err
			}
			if len(tess) == 0 {
				return "", fmt.Errorf("failed to generate test")
			}

			sample := fmt.Sprintf("/* Sample test case:\n%s\n*/\n", question.SampleTestCase)
			data := append(tess[0].Output, []byte(sample)...)
			testPath := filepath.Join(dir, cfg[config.CodeLangKey], title+"_test.go")
			return mdPath, os.WriteFile(testPath, data, 0640)
		}
	}

	return mdPath, nil
}

func getAllFilePath() (string, error) {
	dir, err := getDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, allFile), nil
}

func getDir() (string, error) {
	cfg, err := config.Get()
	if err != nil {
		return "", err
	}
	return parseDir(cfg)
}

func parseDir(cfg map[string]string) (string, error) {
	project, err := filepath.Abs(cfg[config.ProjectKey])
	if err != nil {
		log.Dev(err)
		return "", err
	}

	dir := filepath.Join(project, cfg[config.LangKey])
	err = os.MkdirAll(dir, 0777)
	log.Dev(err)
	return dir, err
}
