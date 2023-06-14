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
	id string

	spinner *spinner.Spinner
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

	return tdoc.Run(docs, tmodel.Config{Title: fmt.Sprintf("Most Voted Solutions for %s", s.id)})
}
