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

	// Initialize or update .env file as needed
	initOrUpdateEnvFile()

	loadEnvVarsIntoProcess()

	// Load config
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// Create and run the application
	runApplication()
}

func initOrUpdateEnvFile() {
	var err error
	// Check if .env file exists
	if _, err = os.Stat(".env"); os.IsNotExist(err) {
		createInitialEnvFile()
	} else {
		fmt.Println(".env file already exists. Checking for missing values...")
		// Load .env file
		if err = godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		// Check and reprompt for missing environment variables
		checkAndUpdateEnvVars()
	}
}

func createInitialEnvFile() {
	fmt.Println("Welcome to Termichat setup.")
	fmt.Println("Please enter the required environment variables:")
	fmt.Println()

	envVars := map[string]string{
		"OPENAI_KEY":        promptForValue("Enter OpenAI API key"),
		"CHAT_HISTORY_PATH": promptForValue("Enter Chat History Path"),
		"AI_MODEL_VERSION":  promptForValue("Enter AI Model Version (e.g., gpt-3.5, gpt-4.0)"),
	}

	writeEnvFile(envVars)

	fmt.Println(".env file created successfully.")
	fmt.Println("You can now run Termichat.")
}

func checkAndUpdateEnvVars() {
	requiredVars := []string{"OPENAI_KEY", "CHAT_HISTORY_PATH", "AI_MODEL_VERSION"}
	envVars := make(map[string]string)

	for _, key := range requiredVars {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			envVars[key] = promptForValue(fmt.Sprintf("Enter %s", key))
		}
	}

	if len(envVars) > 0 {
		writeEnvFile(envVars)
	}
}

func promptForValue(prompt string) string {
	fmt.Print(prompt + ": ")
	return getInput()
}

func loadEnvVarsIntoProcess() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatalf("Failed to open .env file for reading: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			if err := os.Setenv(key, value); err != nil {
				log.Fatalf("Failed to set environment variable: %s", key)
			}
		}
	}
}

func writeEnvFile(envVars map[string]string) {
	// Read existing .env file into a map
	existingVars := make(map[string]string)
	file, err := os.Open(".env")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				existingVars[parts[0]] = parts[1]
			}
		}
		file.Close()
	} // If opening the file fails, we'll just create a new one

	// Update existingVars with new values
	for key, value := range envVars {
		existingVars[key] = value
	}

	// Write back to the .env file
	file, err = os.Create(".env")
	if err != nil {
		log.Fatalf("Failed to open .env file: %v", err)
	}
	defer file.Close()

	for key, value := range existingVars {
		if _, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, value)); err != nil {
			log.Fatalf("Failed to write to .env file: %v", err)
		}
	}
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func runApplication() {
	app := tview.NewApplication()
	app.EnableMouse(true)

	layout := ui.CreateMainLayout(app)

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}
