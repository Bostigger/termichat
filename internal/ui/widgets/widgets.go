package widgets

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"regexp"
	"strings"
	"sync"
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
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetWrap(false).
		SetScrollable(false)
	return titleBox
}

// CreateSubtitle returns a  subtitle text for termichat
func CreateSubtitle() *tview.TextView {
	subtleTitleBox := tview.NewTextView().
		SetText("[brown]- Your Terminal AI Assistant - " + "[-]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	return subtleTitleBox
}

// CreateChatArea returns a configured text view for chat responses.
func CreateChatArea() *tview.TextView {
	chatArea := tview.NewTextView()
	chatArea.SetDynamicColors(true)
	chatArea.SetBorder(true)
	chatArea.SetScrollable(true)
	chatArea.SetRegions(true)
	chatArea.SetWordWrap(true)
	chatArea.SetWrap(true)
	chatArea.SetBorderColor(tcell.Color(237))
	chatArea.SetTitle("Chat Area")

	return chatArea
}

// SetupInputField configures the input field and returns it.
func SetupInputField(app *tview.Application, chatArea *tview.TextView) *tview.InputField {
	var chatHistory string
	var chatMutex sync.Mutex
	var r *glamour.TermRenderer
	var err error

	// Initialize the TermRenderer once
	r, err = glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
	)
	if err != nil {
		log.Fatalf("error creating new term renderer: %v", err) // Fatal log, as we cannot proceed without the renderer
	}

	inputField := tview.NewInputField().
		SetLabel("Ask anything: ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			chatMutex.Lock()
			userInput := inputField.GetText()
			inputField.SetText("").SetDisabled(true) // Clear the input field
			chatHistory += fmt.Sprintf("You: %s\n[green]Bot is typing...[-]\n", userInput)
			chatArea.SetText(chatHistory)
			chatMutex.Unlock()

			go func() {
				response, chatErr := chat.ProcessUserInput(userInput)
				rendered, renderErr := r.Render(response)

				// Strip ANSI sequences
				re := regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]`)
				rendered = re.ReplaceAllString(rendered, "")

				app.QueueUpdateDraw(func() {
					chatMutex.Lock()
					// Remove the "Bot is typing..." message from the chat history
					chatHistory = strings.TrimSuffix(chatHistory, "[green]Bot is typing...[-]\n")

					if chatErr != nil || renderErr != nil {
						if chatErr != nil {
							chatHistory += "Error: " + chatErr.Error() + "\n"
						}
						if renderErr != nil {
							chatHistory += "Error rendering chat response: " + renderErr.Error() + "\n"
						}
						chatArea.SetText(chatHistory)
					} else {
						chatHistory += fmt.Sprintf("[green]Bot: %s\n[-]", rendered)
						chatArea.SetText(chatHistory)
					}
					inputField.SetDisabled(false)
					app.SetFocus(inputField) // Set focus back to the input field
					chatMutex.Unlock()
				})
			}()
		}
	})

	return inputField
}
