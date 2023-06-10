package local

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	tmodel "github.com/zrcoder/tdoc/model"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/log"
)

func GetMetaList() ([]*tmodel.DocInfo, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	ids, err := getPickedQuestionIds(cfg)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("you haven't pick any question yet")
	}
	log.Debug("local ids for book:", ids)

	docs := make([]*tmodel.DocInfo, len(ids))

	for i, id := range ids {
		doc := &tmodel.DocInfo{}
		question, err := GetQuestion(cfg, id)
		if err != nil {
			return nil, err
		}
		doc.Title = fmt.Sprintf("%s. %s", id, question.Title)
		mdFile := GetMarkdownFile(cfg, id)
		fi, err := os.Stat(mdFile)
		if err != nil {
			return nil, err
		}
		doc.ModTime = fi.ModTime()
		log.Debug(doc.Title, doc.ModTime)
		doc.Getter = func(filename string) ([]byte, error) {
			mdData, err := os.ReadFile(mdFile)
			if err != nil {
				return nil, err
			}
			codeData, err := GetTypedCode(cfg, id)
			if err != nil {
				return nil, err
			}
			mdData = bytes.TrimSpace(mdData)

			buf := bytes.NewBuffer(nil)
			buf.WriteString("\n\n## My Solution:\n\n")
			codeLang := cfg.CodeLang
			if codeLang == "golang" {
				codeLang = "go"
			}
			buf.WriteString(fmt.Sprintf("```%s\n", codeLang))
			buf.Write(codeData)
			buf.WriteString("\n```\n")

			mdData = append(mdData, buf.Bytes()...)
			return mdData, nil
		}
		docs[i] = doc
	}

	return docs, nil
}

func getPickedQuestionIds(cfg *config.Config) ([]string, error) {
	dir := filepath.Join(cfg.Language, cfg.CodeLang)
	var ids []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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
