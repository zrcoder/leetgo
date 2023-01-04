package local

import (
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

var (
	ErrNotCached = errors.New("not cached in local yet")
)

const (
	listFileFullName   = "list.json"
	readmeFileFullName = "readme.md"
	codeFileName       = "solution" // for generated code file name, for example: solution.go, solution_test.go
)

/*
we will construct a project, the directory struct like:

.
├── cn
│   ├── golang
│   │   ├── 27-remove-element
│   │   │   ├── readme.md
│   │   │   ├── solution.go
│   │   │   └── solution_test.go
│   │   └── 8-string-to-integer-atoi
│   │       ├── readme.md
│   │       ├── solution.go
│   │       └── solution_test.go
│   └── list.json
└── en
    ├── golang
    │   └── 1-two-sum
    │       ├── readme.md
    │       ├── solution.go
    │       └── solution_test.go
    └── java
        ├── 1-two-sum
        │   ├── readme.md
        │   └── solution.java
        ├── 1004-max-consecutive-ones-iii
        │   ├── readme.md
        │   └── solution.java
        └── 2-add-two-numbers
            ├── readme.md
            └── solution.java
*/

func ReadList() (map[string]model.StatStatusPair, error) {
	path, err := getAllFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Dev(listFileFullName, ErrNotCached)
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

func WriteList(list *model.List) (map[string]model.StatStatusPair, error) {
	path, err := getAllFilePath()
	if err != nil {
		return nil, err
	}

	res := make(map[string]model.StatStatusPair, len(list.StatStatusPairs))
	for _, sp := range list.StatStatusPairs {
		id := sp.Stat.GetFrontendQuestionID()
		sp.Stat.CalculatedID = id
		res[id] = sp
	}

	data, _ := json.Marshal(res)
	err = os.WriteFile(path, data, 0640)
	log.Dev(err)
	return res, err
}

func Read(id string) ([]byte, string, error) {
	log.Dev("begin to search question in local files, id:", id)

	dir, err := getProjectLangCodeDir()
	if err != nil {
		return nil, "", err
	}

	matches, _ := filepath.Glob(filepath.Join(dir, id+"-*"))
	log.Dev("matches:", matches)
	if len(matches) == 0 {
		log.Dev(ErrNotCached)
		return nil, "", ErrNotCached
	}

	path := matches[0]
	data, err := os.ReadFile(filepath.Join(path, readmeFileFullName))
	log.Dev(err)
	return data, path, err
}

func Write(sp *model.StatStatusPair, question *model.Question, mdData []byte) (string, error) {
	cfg, err := config.Get()
	if err != nil {
		return "", err
	}

	dir, err := parseProjectLangCodeDir(cfg)
	if err != nil {
		return "", err
	}

	title := fmt.Sprintf("%s-%s", sp.Stat.CalculatedID, sp.Stat.QuestionTitleSlug)
	dir = filepath.Join(dir, title)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		log.Dev(err)
		return "", err
	}

	mdPath := filepath.Join(dir, readmeFileFullName)
	err = os.WriteFile(mdPath, mdData, 0640)
	if err != nil {
		log.Dev(err)
		return "", err
	}

	codes, err := question.ParseCodes()
	if err != nil {
		log.Dev(err)
		return "", err
	}

	currentCodeLang := cfg[config.CodeLangKey]
	fileExtension := config.CodeLangExtensionDic[currentCodeLang]
	codePath := filepath.Join(dir, codeFileName+fileExtension)
	for _, c := range codes {
		if c.Value == currentCodeLang {
			content := genCodeContent(currentCodeLang, c.DefaultCode)
			err = os.WriteFile(codePath, []byte(content), 0640)
			if err != nil {
				log.Dev(err)
				return "", err
			}
			if c.Value != config.CodeLangGo && c.Value != config.CodeLangGoShort {
				return dir, nil
			}

			testPath := filepath.Join(dir, codeFileName+"_test.go")
			if _, err = os.Stat(testPath); err == nil {
				err = os.Remove(testPath) // should remove the old test file when update
				if err != nil {
					log.Dev(err)
					return "", err
				}
			}
			tests, err := gotests.GenerateTests(codePath, nil)
			if err != nil {
				log.Dev(err)
				return "", err
			}
			log.Dev("generated", len(tests), "tests")
			if len(tests) == 0 {
				err = errors.New("no tests generated, may the test file has already exist")
				log.Dev(err)
				return "", err
			}
			sample := fmt.Sprintf("/* Sample test case:\n%s\n*/\n", question.SampleTestCase)
			data := append(tests[0].Output, []byte(sample)...)
			return dir, os.WriteFile(testPath, data, 0640)
		}
	}

	return dir, nil
}

func getAllFilePath() (string, error) {
	dir, err := getProjectLangDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, listFileFullName), nil
}

func getProjectLangDir() (string, error) {
	cfg, err := config.Get()
	if err != nil {
		return "", err
	}

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

func getProjectLangCodeDir() (string, error) {
	cfg, err := config.Get()
	if err != nil {
		return "", err
	}
	return parseProjectLangCodeDir(cfg)
}

func parseProjectLangCodeDir(cfg map[string]string) (string, error) {
	project, err := filepath.Abs(cfg[config.ProjectKey])
	if err != nil {
		log.Dev(err)
		return "", err
	}

	dir := filepath.Join(project, cfg[config.LangKey], cfg[config.CodeLangKey])
	err = os.MkdirAll(dir, 0777)
	log.Dev(err)
	return dir, err
}

func genCodeContent(lang, code string) string {
	const (
		beginFlag = "// [BEGIN leetcode submit region] don't modify"
		endFlag   = "// [END leetcode submit region] don't modify"
	)
	switch lang {
	case config.CodeLangGo, config.CodeLangGoShort:
		return fmt.Sprintf("package leetgo\n\n%s\n%s\n%s", beginFlag, code, endFlag)
	default:
		// TODO: adapt every language leetcode surpported
		return fmt.Sprintf("%s\n%s\n%s", beginFlag, code, endFlag)
	}
}
