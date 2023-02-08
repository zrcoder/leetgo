package comp

import (
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/zrcoder/tdoc"
	"github.com/zrcoder/tdoc/docmgr"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/render"
)

func NewBook() Component {
	return &Book{
		spinner: newSpinner("Generating"),
	}
}

type Book struct {
	spinner *spinner.Spinner
}

func (b *Book) Run() error {
	b.spinner.Start()
	docPath, err := local.Generate()
	b.spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Println(render.Infof("Your book is generated in %s\n", docPath))

	mgr, err := docmgr.New(docPath)
	if err != nil {
		return err
	}
	return tdoc.Run(mgr.Docs())
}
