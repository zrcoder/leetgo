package local

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/dgraph-io/badger/v3"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

var (
	ErrNotCached = errors.New("not cached in local yet")
)

func ReadList() (map[string]model.StatStatusPair, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	path, err := getListDir(cfg)
	if err != nil {
		return nil, err
	}
	db, err := getDB(cfg)
	if err != nil {
		return nil, err
	}
	defer func() { _ = db.Close() }()

	var data []byte
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(path))
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
	path, err := getListDir(cfg)
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
	key := []byte(path)
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, data)
	})
	log.Trace(err)
	return res, err
}

func Read(id string) (*model.Question, error) {
	log.Trace("begin to search question in local files, id:", id)
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	dir, err := getCacheLangDir(cfg)
	if err != nil {
		return nil, err
	}
	db, err := getDB(cfg)
	if err != nil {
		return nil, err
	}
	defer func() { _ = db.Close() }()

	var data []byte
	key := []byte(filepath.Join(dir, id))
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

	res := &model.Question{}
	err = json.Unmarshal(data, res)
	return res, err
}

func Write(id string, question *model.Question) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	dir, err := getCacheLangDir(cfg)
	if err != nil {
		return err
	}
	db, err := getDB(cfg)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	path := filepath.Join(dir, id)
	data, _ := json.Marshal(question)
	return db.Update(func(txn *badger.Txn) error {
		txn = db.NewTransaction(true)
		err = txn.Set([]byte(path), data)
		if err != nil {
			return err
		}
		return txn.Commit()
	})
}

func getCacheDir(cfg *config.Config) (string, error) {
	dir, err := filepath.Abs(cfg.CacheDir)
	log.Trace(err)
	return dir, err
}

func getListDir(cfg *config.Config) (string, error) {
	dir, err := filepath.Abs(cfg.CacheDir)
	if err != nil {
		log.Trace(err)
		return "", err
	}
	return filepath.Join(dir, cfg.Language, "list"), nil
}

func getCacheLangDir(cfg *config.Config) (string, error) {
	cacheDir, err := filepath.Abs(cfg.CacheDir)
	if err != nil {
		log.Trace(err)
		return "", err
	}
	return filepath.Join(cacheDir, cfg.Language), err
}

func getDB(cfg *config.Config) (*badger.DB, error) {
	cacheDir, err := getCacheDir(cfg)
	if err != nil {
		return nil, err
	}
	opts := badger.DefaultOptions(cacheDir)
	opts.Logger = log.Logger
	db, err := badger.Open(opts)
	log.Trace(err)
	return db, err
}
