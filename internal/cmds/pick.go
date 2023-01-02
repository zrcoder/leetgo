package cmds

import (
	"errors"
	"fmt"
	"strings"

	h2md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/urfave/cli/v2"

	"github.com/zrcoder/leetgo/internal/client"
	"github.com/zrcoder/leetgo/internal/render"
	"github.com/zrcoder/leetgo/internal/trace"
)

var Pick = &cli.Command{
	Name:  "pick",
	Usage: "pick a question by id",
	Action: func(context *cli.Context) error {
		if context.Args().Len() == 0 {
			return trace.Wrap(errors.New("need question id"))
		}

		id := strings.Join(context.Args().Slice(), " ")
		question, err := client.Query(id)
		if err != nil {
			return err
		}

		content := question.TranslatedContent
		if content == "" {
			content = question.Content
		}
		content = strings.ReplaceAll(content, "<sup>", "^")
		content = strings.ReplaceAll(content, "</sup>", "")
		converter := h2md.NewConverter("", true, nil)
		content, err = converter.ConvertString(content)
		if err != nil {
			return err
		}
		md := fmt.Sprintf("## [%s. %s](%s)\n\nDifficulty: %s\n\n%s",
			question.FrontendQuestionID, question.Title, question.Referer, question.Difficulty, content) // todo content
		fmt.Println(render.MarkDown(md))

		return nil
	},
}
