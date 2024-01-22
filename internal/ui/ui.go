package ui

import (
	"github.com/rivo/tview"
	"termichat/internal/ui/widgets"
)

// CreateMainLayout sets up the UI and returns the top-level layout.
func CreateMainLayout(app *tview.Application) *tview.Flex {
	titleBox := widgets.CreateTitleBox()
	chatArea := widgets.CreateChatArea()
	inputField := widgets.SetupInputField(app, chatArea)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(titleBox, 0, 5, false).
		AddItem(chatArea, 0, 10, false).
		AddItem(inputField, 3, 1, true)
	return layout

}
