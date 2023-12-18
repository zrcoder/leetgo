package comp

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/leetgo/internal/remote"
	"github.com/zrcoder/tdoc"
	tmodel "github.com/zrcoder/tdoc/model"
)

func NewSolutionViewer(id string) Component {
	return &solutionViewer{id: id, spinner: newSpinner("Searching")}
}

type solutionViewer struct {
	spinner *spinner.Spinner
	id      string
}

func (s *solutionViewer) Run() error {
	meta, err := queryMeta(s.id)
	if err != nil {
		return err
	}

	s.spinner.Start()
	solutionsResp, err := remote.GetSolutions(meta)
	s.spinner.Stop()
	if err != nil {
		return err
	}

	docs, err := getDocsFromSolutions(solutionsResp, meta)
	if err != nil {
		return err
	}

	s.id = meta.FrontendID // original s.id may be "today"
	err = tdoc.Run(docs, tmodel.Config{Title: fmt.Sprintf("Most Voted Solutions for %s", s.id)})
	if err != nil {
		return err
	}

	return askToCode(s.id)
}
