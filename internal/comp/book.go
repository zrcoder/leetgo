package comp

import (
	"fmt"

	"github.com/briandowns/spinner"

	"github.com/zrcoder/leetgo/internal/local"
	"github.com/zrcoder/leetgo/internal/render"
	"github.com/zrcoder/mdoc"
)

func NewBook(sortBy, port string) Component {
	return &Book{
		sortBy:  sortBy,
		port:    port,
		spinner: newSpinner("Generating"),
	}
}

type Book struct {
	sortBy string
	port   string

	spinner *spinner.Spinner
}

func (b *Book) Run() error {
	b.spinner.Start()
	docPath, err := local.Generate(b.sortBy)
	b.spinner.Stop()
	if err != nil {
		return err
	}
	fmt.Println(render.Infof("  Your book is generated in %s\n", docPath))

	cfg := mdoc.DefaultConfig
	cfg.HttpPort = b.port
	cfg.Page.HasLandingPage = false
	cfg.Page.DocsBasePath = docPath

	return mdoc.Serve(cfg)
}
