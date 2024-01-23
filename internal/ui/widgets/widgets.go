package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
	"termichat/internal/chat"
)

// CreateTitleBox returns a configured title box
func CreateTitleBox() *tview.TextView {
	asciiArt := `
	
	████████╗███████╗██████╗░███╗░░░███╗██╗░█████╗░██╗░░██╗░█████╗░████████╗
	╚══██╔══╝██╔════╝██╔══██╗████╗░████║██║██╔══██╗██║░░██║██╔══██╗╚══██╔══╝
	░░░██║░░░█████╗░░██████╔╝██╔████╔██║██║██║░░╚═╝███████║███████║░░░██║░░░
	░░░██║░░░██╔══╝░░██╔══██╗██║╚██╔╝██║██║██║░░██╗██╔══██║██╔══██║░░░██║░░░
	░░░██║░░░███████╗██║░░██║██║░╚═╝░██║██║╚█████╔╝██║░░██║██║░░██║░░░██║░░░
	░░░╚═╝░░░╚══════╝╚═╝░░╚═╝╚═╝░░░░░╚═╝╚═╝░╚════╝░╚═╝░░╚═╝╚═╝░░╚═╝░░░╚═╝░░░
	
	`

	titleBox := tview.NewTextView().
		SetText("[green]" + asciiArt + "[-]").
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetWrap(false).
		SetScrollable(false)
	return titleBox
}

// CreateSubtitle returns a  subtitle text for termichat
func CreateSubtitle() *tview.TextView {
	subtleTitleBox := tview.NewTextView().
		SetText("[brown]- Your Terminal AI Assistance - " + "[-]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	return subtleTitleBox
}

// CreateChatArea returns a configured text view for chat responses.
func CreateChatArea() *tview.TextView {
	chatArea := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)
	return chatArea
}

// SetupInputField configures the input field and returns it.
func SetupInputField(app *tview.Application, chatArea *tview.TextView) *tview.InputField {
	var chatHistory string

	inputField := tview.NewInputField().
		SetLabel("Ask anything: ").
		SetFieldWidth(0).SetFieldBackgroundColor(tcell.ColorBlack)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			userInput := inputField.GetText()
			inputField.SetText("") // Clear the input field

			chatHistory += "You: " + userInput + "\n"
			chatHistory += "[green]Bot is typing...\n [-]"
			chatArea.SetText(chatHistory)

			go func() {
				response, err := chat.ProcessUserInput(userInput)
				app.QueueUpdateDraw(func() {

					chatHistory = strings.TrimSuffix(chatHistory, "[green]Bot is typing...\n [-]")
					if err != nil {
						chatHistory += "Error: " + err.Error() + "\n"
					} else {
						chatHistory += "[green]Bot: " + response + "\n" + "[-]"
					}
					chatArea.SetText(chatHistory)
					app.SetFocus(inputField) // Set focus back to the input field
				})
			}()
		}
	})

	return inputField
}
