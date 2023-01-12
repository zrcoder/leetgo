package local

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zrcoder/rdbook"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
)

const (
	bookMarkdownName = "_markdown"
)

func Generate(info *model.BookMeta) (string, string, error) {
	docs, err := getDocs()
	if err != nil {
		return "", "", err
	}

	markDownPath := ""
	if info.GenMarkdowns {
		markDownPath, err = writeBookMds(docs)
		if err != nil {
			return "", "", err
		}
	}

	sitePath, err := rdbook.GenerateWithDocs(docs, info.DestPath, rdbook.WithSortRule(rdbook.SortByTime))
	if err != nil {
		log.Trace(err)
		return "", "", err
	}

	return sitePath, markDownPath, nil
}

func getDocs() ([]*rdbook.Doc, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	db, err := getDB(cfg)
	if err != nil {
		return nil, err
	}
	defer func() { _ = db.Close() }()

	ids, err := getPickedQuestionIds(cfg)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("you didn't pick any question yet")
	}

	docs := make([]*rdbook.Doc, len(ids))
	for i, id := range ids {
		doc := &rdbook.Doc{}

		_, question, err := read(cfg, id, db)
		if err != nil {
			return nil, err
		}

		doc.Title = fmt.Sprintf("%s. %s", question.ID, question.Title)
		doc.MarkdownContent = []byte(strings.TrimSpace(question.MdContent))
		code, modTime, err := readAnswerCode(cfg, id)
		if err != nil {
			return nil, err
		}
		doc.Time = *modTime
		buf := bytes.NewBuffer(nil)
		buf.WriteString("\n\n## My Solution:\n\n")
		codeLang := cfg.CodeLang
		if codeLang == config.CodeLangGo {
			codeLang = config.CodeLangGoShort
		}
		buf.WriteString(fmt.Sprintf("```%s\n", codeLang))
		buf.Write(code)
		buf.WriteString("\n```\n")
		doc.MarkdownContent = append(doc.MarkdownContent, buf.Bytes()...)

		docs[i] = doc
	}

	return docs, nil
}

func getPickedQuestionIds(cfg *config.Config) ([]string, error) {
	dir, err := filepath.Abs(filepath.Join(cfg.CacheDir, cfg.Language, cfg.CodeLang))
	if err != nil {
		return nil, err
	}
	var ids []string
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if path == dir || d.Name() == bookMarkdownName {
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

func writeBookMds(docs []*rdbook.Doc) (string, error) {
	cfg, err := config.Get()
	if err != nil {
		return "", err
	}

	dir, err := getBookMarkdownDir(cfg)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(dir, 0777)
	if err != nil {
		log.Trace(err)
		return "", err
	}

	for _, doc := range docs {
		name := filepath.Join(dir, doc.Title+".md")
		err = os.WriteFile(name, doc.MarkdownContent, 0640)
		if err != nil {
			log.Trace(err)
			return "", err
		}
	}

	return dir, nil
}

func readAnswerCode(cfg *config.Config, id string) ([]byte, *time.Time, error) {
	path, err := getCodeFilePath(cfg, id)
	if err != nil {
		return nil, nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		log.Trace(err)
		return nil, nil, err
	}
	defer func() { _ = f.Close() }()

	stat, err := f.Stat()
	if err != nil {
		log.Trace(err)
		return nil, nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Trace(err)
		return nil, nil, err
	}

	index := bytes.Index(data, []byte(codeStartFlag))
	if index == -1 {
		return nil, nil, newHintError(fmt.Sprintf("start flag %s not found", codeStartFlag), path)
	}
	data = data[index+len(codeStartFlag):]
	index = bytes.Index(data, []byte(codeEndFlag))
	if index == -1 {
		return nil, nil, newHintError(fmt.Sprintf("end flag %s not found", codeEndFlag), path)
	}
	modTime := stat.ModTime()
	return bytes.TrimSpace(data[:index]), &modTime, nil
}

func getBookMarkdownDir(cfg *config.Config) (string, error) {
	return filepath.Abs(filepath.Join(cfg.CacheDir, cfg.Language, cfg.CodeLang, bookMarkdownName))
}

func newHintError(info, path string) error {
	return fmt.Errorf("%s, pleanse check %s", info, path)
}