package huh

import "github.com/charmbracelet/huh"

func NewConfirm(title, description string, bind *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title(title).Description(description).Value(bind),
		),
	)
}
