package ui

import (
	"termichat/internal/ui/widgets"

	"github.com/rivo/tview"
)

// CreateMainLayout sets up the UI and returns the top-level layout.
func CreateMainLayout(app *tview.Application) *tview.Flex {
	titleBox := widgets.CreateTitleBox()
	subtitleBox := widgets.CreateSubtitle()
	chatArea := widgets.CreateChatArea()
	inputField := widgets.SetupInputField(app, chatArea)
	clearButton := widgets.ClearButton(chatArea)
	closeButton := widgets.CloseButton(app)
	exportButton := widgets.ExportButton()

     buttonLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
        AddItem(nil, 0, 1, false). // Left space for padding
        AddItem(clearButton, 10, 1, false). // Clear button with a fixed width
        AddItem(nil, 2, 1, false). // Space between buttons
        AddItem(closeButton, 10, 1, false). // Close button with a fixed width
        AddItem(nil, 2, 1, false). // Space between buttons
        AddItem(exportButton, 10, 1, false). // Export button with a fixed width
        AddItem(nil, 0, 1, false) // Right space for padding

    // Add the button layout to the main layout
    layout := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(titleBox, 0, 5, false).
        AddItem(subtitleBox, 0, 1, false).
        AddItem(buttonLayout, 0, 1, false). // Add the horizontal button layout here
        AddItem(chatArea, 0, 10, false).
        AddItem(inputField, 3, 1, true)

    return layout
}