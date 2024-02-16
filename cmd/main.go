package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"termichat/config"
	"termichat/internal/ui"

	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

func main() {
	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// If .env file does not exist, prompt user to input environment variables
		fmt.Println("Welcome to Termichat setup.")
		fmt.Println("Please enter the required environment variables:")
		fmt.Println()

		// Prompt for OpenAI API key
		fmt.Print("Enter OpenAI API key: ")
		openAIKey := getInput()

		// Prompt for Chat History Path
		fmt.Print("Enter Chat History Path: ")
		chatHistoryPath := getInput()

		// Prompt for AI Model Version
		fmt.Print("Enter AI Model Version (e.g., gpt-3.5, gpt-4.0): ")
		aiModelVersion := getInput()

		// Create .env file and write environment variables
		err := createEnvFile(openAIKey, chatHistoryPath, aiModelVersion)
		if err != nil {
			log.Fatalf("Error creating .env file: %v", err)
		}

		fmt.Println(".env file created successfully.")
		fmt.Println("You can now run Termichat.")
	} else {
		fmt.Println(".env file already exists. Skipping setup.")
	}

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load config
	err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create application
	app := tview.NewApplication()
	app.EnableMouse(true)

	// Set up the main layout using the ui package
	layout := ui.CreateMainLayout(app)

	// Start the application with the layout
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func createEnvFile(openAIKey, chatHistoryPath, aiModelVersion string) error {
	file, err := os.Create(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "OPENAI_KEY=%s\nCHAT_HISTORY_PATH=%s\nAI_MODEL_VERSION=%s\n", openAIKey, chatHistoryPath, aiModelVersion)
	if err != nil {
		return err
	}

	return nil
}
