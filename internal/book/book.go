package book

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	tmodel "github.com/zrcoder/tdoc/model"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
)

func GetMetaList() ([]*tmodel.DocInfo, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	ids, err := local.GetPickedQuestionIds(cfg)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("you haven't pick any question yet")
	}
	log.Debug("local ids for book:", ids)

	docs := make([]*tmodel.DocInfo, len(ids))

	for i, id := range ids {
		id := id // for doc.Getter below, captrue current id
		doc := &tmodel.DocInfo{}
		question, err := local.GetQuestion(cfg, id)
		if err != nil {
			return nil, err
		}
		doc.Title = fmt.Sprintf("%s. %s", id, question.Title)
		mdFile := local.GetMarkdownFile(cfg, id)
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
			codeData, err := local.GetTypedCode(cfg, id)
			if err != nil {
				return nil, err
			}
			noteDate, err := local.GetNotes(cfg, id)
			if err != nil {
				return nil, err
			}
			noteDate = bytes.TrimSpace(noteDate)
			mdData = bytes.TrimSpace(mdData)

			buf := bytes.NewBuffer(nil)
			buf.WriteString("\n\n## My Solution:\n\n")
			if len(noteDate) > 0 {
				buf.Write(noteDate)
				buf.WriteString("\n\n")
			}
			codeLang := config.DisplayLang(cfg.CodeLang)
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

func GetMetaListFromSolutions(solutionsResp model.SolutionListResp, question *model.Question) ([]*tmodel.DocInfo, error) {
	reqs := solutionsResp.SolutionReqs()
	docs := make([]*tmodel.DocInfo, len(reqs))
	for i, req := range reqs {
		req := req // for doc.Getter
		doc := &tmodel.DocInfo{}
		doc.Title = req.Title
		doc.Getter = func(s string) ([]byte, error) {
			rsp, err := remote.GetSolution(&req, question)
			if err != nil {
				return nil, err
			}
			content, err := rsp.RegularContent()
			if err != nil {
				return nil, err
			}
			return append([]byte(content), '\n'), nil
		}
		docs[i] = doc
	}
	return docs, nil
}
