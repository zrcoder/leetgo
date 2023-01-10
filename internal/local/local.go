package local

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cweill/gotests"
	"github.com/dgraph-io/badger/v3"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

const (
	codeFileName = "solution"

	codeStartFlag = "// [start] don't modify\n"
	codeEndFlag   = "// [end] don't modify\n"
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

func ReadList() (map[string]model.StatStatusPair, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	db, err := getDB(cfg)
	if err != nil {
		return nil, err
	}
	defer func() { _ = db.Close() }()

	key := makeListKey(cfg)
	var data []byte
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		log.Trace(err)
		if err == badger.ErrKeyNotFound {
			err = ErrNotCached
		}
		return nil, err
	}
	res := map[string]model.StatStatusPair{}
	err = json.Unmarshal(data, &res)
	log.Trace(err)
	return res, err
}

func WriteList(list *model.List) (map[string]model.StatStatusPair, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	db, err := getDB(cfg)
	if err != nil {
		return nil, err
	}
	defer func() { _ = db.Close() }()

	res := make(map[string]model.StatStatusPair, len(list.StatStatusPairs))
	for _, sp := range list.StatStatusPairs {
		id := sp.Stat.GetFrontendQuestionID()
		sp.Stat.CalculatedID = id
		res[id] = sp
	}
	data, _ := json.Marshal(res)
	key := makeListKey(cfg)
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, data)
	})
	log.Trace(err)
	return res, err
}

func Read(id string) (string, *model.Question, error) {
	log.Trace("begin to search question in local, id:", id)
	cfg, err := config.Get()
	if err != nil {
		return "", nil, err
	}
	db, err := getDB(cfg)
	if err != nil {
		return "", nil, err
	}
	defer func() { _ = db.Close() }()

	var data []byte
	key := makeQuestionKey(cfg, id)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		log.Trace(err)
		if err == badger.ErrKeyNotFound {
			err = ErrNotCached
		}
		return "", nil, err
	}

	res := &model.Question{}
	err = json.Unmarshal(data, res)
	if err != nil {
		log.Trace(err)
		return "", nil, err
	}
	dir, err := makeCodeDir(cfg, id)
	return dir, res, err
}

func Write(id string, question *model.Question) (string, error) {
	cfg, err := config.Get()
	if err != nil {
		return "", err
	}
	err = writeDB(id, question, cfg)
	if err != nil {
		log.Trace(err)
		return "", err
	}

	return writeCodeFile(id, question, cfg)
}

func writeDB(id string, question *model.Question, cfg *config.Config) error {
	db, err := getDB(cfg)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	key := makeQuestionKey(cfg, id)
	data, _ := json.Marshal(question)
	return db.Update(func(txn *badger.Txn) error {
		txn = db.NewTransaction(true)
		err = txn.Set(key, data)
		if err != nil {
			return err
		}
		return txn.Commit()
	})
}

func writeCodeFile(id string, question *model.Question, cfg *config.Config) (string, error) {
	codes, err := question.ParseCodes()
	if err != nil {
		return "", err
	}
	dir, err := makeCodeDir(cfg, id)
	if err != nil {
		return "", err
	}
	path, _ := getCodeFilePath(cfg, id)
	buf := bytes.NewBuffer(nil)
	isGo := cfg.CodeLang == config.CodeLangGoShort || cfg.CodeLang == config.CodeLangGo
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
	err = os.WriteFile(path, buf.Bytes(), 0640)
	if err != nil {
		log.Trace(err)
		return "", err
	}
	if !isGo {
		return dir, nil
	}
	testPath := filepath.Join(dir, codeFileName+"_test.go")
	_ = os.Remove(testPath) // need remove the test file when update
	tests, err := gotests.GenerateTests(path, nil)
	if err != nil {
		return "", err
	}
	if len(tests) == 0 {
		return "", errors.New("no tests generated")
	}
	sample := fmt.Sprintf("\n/* sample test case:\n%s\n*/\n", question.SampleTestCase)
	content := append(tests[0].Output, []byte(sample)...)
	err = os.WriteFile(testPath, content, 0640)
	log.Trace(err)
	return dir, err
}

func makeCodeDir(cfg *config.Config, id string) (string, error) {
	res, err := filepath.Abs(filepath.Join(cfg.CacheDir, cfg.Language, cfg.CodeLang, id))
	if err != nil {
		log.Trace(err)
		return "", err
	}
	err = os.MkdirAll(res, 0777)
	log.Trace(err)
	return res, err
}

func makeListKey(cfg *config.Config) []byte {
	return []byte(filepath.Join(cfg.Language, "list"))
}

func makeQuestionKey(cfg *config.Config, id string) []byte {
	return []byte(filepath.Join(cfg.Language, cfg.CodeLang, id))
}

func getDB(cfg *config.Config) (*badger.DB, error) {
	dir, err := filepath.Abs(filepath.Join(cfg.CacheDir, "db"))
	if err != nil {
		log.Trace(err)
		return nil, err
	}
	opts := badger.DefaultOptions(dir)
	opts.Logger = log.Logger
	db, err := badger.Open(opts)
	log.Trace(err)
	return db, err
}

func getCodeFilePath(cfg *config.Config, id string) (string, error) {
	path, err := filepath.Abs(filepath.Join(cfg.CacheDir, cfg.Language, cfg.CodeLang, id, codeFileName+extesionDic[cfg.CodeLang]))
	if err != nil {
		log.Trace(err)
		return "", err
	}
	return path, nil
}
