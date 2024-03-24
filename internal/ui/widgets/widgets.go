package widgets

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"termichat/config"
	"termichat/internal/chat"
	"termichat/util"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

var chatHistory string
var chatArea tview.TextView

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

// CreateButtons returns a button box with buttons for the user to click to clear chat history
func ClearButton(chatArea *tview.TextView) *tview.Button {
	button := tview.NewButton("Clear").
		SetLabelColor(tcell.ColorWhite).
		SetSelectedFunc(func() {
			// Clear the chat history and the text in the chat area.
			chatHistory = ""
			chatArea.SetText("")
		})
	button.SetBackgroundColor(tcell.ColorBlack)
	return button
}

func CloseButton(app *tview.Application) *tview.Button {
	button := tview.NewButton("Close").
		SetLabelColor(tcell.ColorWhite).
		SetSelectedFunc(func() {
			// Close the chat terminal user interface
			app.Stop()
		})
	return button
}

func ExportButton() *tview.Button {
	button := tview.NewButton("Export").
		SetLabelColor(tcell.ColorWhite).
		SetSelectedFunc(func() {
			// Get the current chat history
			currentChatHistory := GetChatHistory()

			// Generate a filename with timestamp and a predefined topic
			topic := "MyTopic"                                // Replace with your desired topic
			timestamp := time.Now().Format("20060102-150405") // Format: YYYYMMDD-HHMMSS
			filename := fmt.Sprintf("%s_chat_history_%s_%s.txt", config.ChatHistoryPath, topic, timestamp)

			// Create the file
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			// Write the chat history to the file
			file.WriteString(currentChatHistory)
		})
	return button
}

func ConfigOptionsDropdown(app *tview.Application) *tview.DropDown {
	dropDown := tview.NewDropDown().
		SetLabel("Config Options: ").
		SetFieldBackgroundColor(tcell.ColorRed).
		SetOptions([]string{"Delete API Key", "Change Export Location", "Change GPT Model"}, nil)

	dropDown.SetSelectedFunc(func(text string, index int) {
		switch index {
		case 0: // Delete API Key
			util.DeleteAPIKey()
			app.Stop()
		case 1: // Change Export Location
			// Create an input field within a custom modal for the new path
			inputField := tview.NewInputField().
				SetLabel("New Path: ").
				SetFieldWidth(40)

			modal := tview.NewModal().
				SetText("Enter the new path for CHAT_HISTORY_PATH:").
				AddButtons([]string{"Confirm", "Cancel"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Confirm" {
						newPath := inputField.GetText()
						if util.UpdateChatHistoryPath(newPath) {
							// Successfully updated, perform any additional actions needed
						}
						// Return back to the main view

					}
					if buttonLabel == "Cancel" {
						// Return back to the main view

					}
				})

			// Construct a new Flex layout to hold the input field above the modal buttons
			form := tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(inputField, 3, 0, true).
				AddItem(modal, 0, 1, false)

			app.SetRoot(form, false)
			app.SetFocus(inputField)
		case 2: // Change GPT Model
			// Define action for Option 3
		}
	})

	return dropDown
}

// GetChatHistory returns the current chat history
func GetChatHistory() string {
	return chatHistory
}
