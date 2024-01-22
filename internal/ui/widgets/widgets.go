package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	inputField := tview.NewInputField().
		SetLabel("Ask anything: ").
		SetFieldWidth(0)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			userInput := inputField.GetText()
			inputField.SetText("") // Clear the input field immediately
			chatArea.Write([]byte("You: " + userInput + "\n"))

			go func() {
				response, err := chat.ProcessUserInput(userInput)
				app.QueueUpdateDraw(func() {
					if err != nil {
						// Display the error in the chat area
						chatArea.Write([]byte("Error: " + err.Error() + "\n"))
					} else {
						// Display user input and response in the chat area
						chatArea.Write([]byte("Bot: " + response + "\n"))
					}
					//chatArea.ScrollToEnd()
					app.SetFocus(inputField) // Set focus back to the input field
				})
			}()
		}
	})

	return inputField
}
