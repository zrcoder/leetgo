package comp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"

	tmodel "github.com/zrcoder/tdoc/model"

	"github.com/zrcoder/leetgo/internal/config"
	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/log"
	"github.com/zrcoder/leetgo/internal/model"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/leetgo/utils/exec"
	"github.com/zrcoder/leetgo/utils/huh"
)

const (
	noQuestionFoundForErrFmt = "no question found for `%s`"
)

func queryMeta(frontendID string) (*model.Meta, error) {
	spinner := newSpinner("inquiring")
	spinner.Start()
	defer spinner.Stop()

	if frontendID == "today" {
		today, err := remote.GetToday()
		if err != nil {
			return nil, err
		}
		return today.Meta(), nil
	}

	return query(frontendID)
}

func queryMetas(key string) ([]model.Meta, error) {
	if key == "today" {
		today, err := remote.GetToday()
		if err != nil {
			return nil, err
		}
		return []model.Meta{*today.Meta()}, nil
	}

	return search(key)
}

func regualarID(id string) string {
	if id != "today" {
		return id
	}

	spinner := newSpinner("Fetching today")
	spinner.Start()
	today, err := remote.GetToday()
	spinner.Stop()
	if err != nil {
		log.Debug(err)
		return id
	}
	return today.Meta().FrontendID
}

func getDocFromLocal(id string) (*tmodel.DocInfo, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	doc := &tmodel.DocInfo{}
	question, err := local.GetQuestion(cfg, id)
	if err != nil {
		return nil, err
	}
	doc.Title = fmt.Sprintf("%s. %s", id, question.Title)
	mdFile := local.GetMarkdownFile(id)
	fi, err := os.Stat(mdFile)
	if err != nil {
		return nil, err
	}
	doc.ModTime = fi.ModTime()
	log.Debug(doc.Title, doc.ModTime)
	doc.Getter = func(_ string) ([]byte, error) {
		mdData, err := os.ReadFile(mdFile)
		if err != nil {
			return nil, err
		}
		codeData, err := local.GetTypedCode(cfg, id)
		if err != nil {
			return nil, err
		}
		noteData, err := local.GetNotes(cfg, id)
		if err != nil {
			return nil, err
		}
		noteData = bytes.TrimSpace(noteData)
		mdData = bytes.TrimSpace(mdData)

		buf := bytes.NewBuffer(nil)
		buf.WriteString("\n\n## My Solution:\n\n")
		if len(noteData) > 0 {
			buf.Write(noteData)
			buf.WriteString("\n\n")
		}
		codeLang := config.DisplayLang(cfg.CodeLang)
		fmt.Fprintf(buf, "```%s\n", codeLang)
		buf.Write(codeData)
		buf.WriteString("\n```\n")

		mdData = append(mdData, buf.Bytes()...)
		return mdData, nil
	}
	return doc, nil
}

func getDocsFromLocal() ([]*tmodel.DocInfo, error) {
	ids, err := local.GetPickedQuestionIds()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, errors.New("you haven't pick any question yet")
	}
	log.Debug("local ids for book:", ids)

	docs := make([]*tmodel.DocInfo, len(ids))

	for i, id := range ids {
		doc, err := getDocFromLocal(id)
		if err != nil {
			return nil, err
		}
		docs[i] = doc
	}

	return docs, nil
}

func getDocsFromSolutions(solutionsResp model.SolutionListResp, meta *model.Meta) ([]*tmodel.DocInfo, error) {
	reqs := solutionsResp.SolutionReqs()
	docs := make([]*tmodel.DocInfo, len(reqs))
	for i, req := range reqs {
		req := req // for doc.Getter
		doc := &tmodel.DocInfo{}
		doc.Title = req.Title
		doc.Description = req.Author + " . " + req.CreateAt.Format("2006-01-02 15:04")
		doc.Getter = func(_ string) ([]byte, error) {
			rsp, err := remote.GetSolution(&req, meta)
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

func query(frontendID string) (*model.Meta, error) {
	allMap, err := local.GetAll()
	if err == nil {
		res, ok := allMap[frontendID]
		if ok {
			return &res, nil
		}
		err = local.ErrNotCached
	}
	if err != local.ErrNotCached {
		return nil, err
	}

	all, err := remote.GetAll()
	if err != nil {
		return nil, err
	}
	allMap = make(map[string]model.Meta, len(all.StatStatusPairs))
	for _, sp := range all.StatStatusPairs {
		meta := sp.Meta()
		meta.Transform()
		allMap[meta.FrontendID] = meta
	}
	err = local.WriteAll(allMap)
	if err != nil {
		return nil, err
	}
	if res, ok := allMap[frontendID]; ok {
		return &res, nil
	}
	return nil, fmt.Errorf(noQuestionFoundForErrFmt, frontendID)
}

func search(key string) ([]model.Meta, error) {
	list, err := remote.Search(key)
	if err != nil {
		return nil, err
	}

	res := list.Data.ProblemsetQuestionList.Questions

	if len(res) == 0 {
		log.Debug("no questions found")
		return nil, fmt.Errorf(noQuestionFoundForErrFmt, key)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].FrontendID < res[j].FrontendID
	})
	return res, nil
}

func askToCode(id string) error {
	edit := true
	err := huh.NewConfirm(
		"Solve the question now?",
		"Open the local code file with your favorite editor to edit the code.",
		&edit).Run()
	if err != nil {
		return err
	}

	if edit {
		return code(id)
	}
	return nil
}

func code(id string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	cmd, ops := config.GetEditorCmdOps(cfg.Editor)
	mdFile := local.GetMarkdownFile(id)
	codeFile := local.GetCodeFile(cfg, id)
	ops = append(ops, mdFile, codeFile)
	return exec.Run(local.GetDir(id), cmd, ops...)
}
