package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func FormatMessage(message string) string {
	// Format message for display
	return ""
}

func DeleteAPIKey() bool {
	envFile := ".env"
	tempFile := ".env.tmp"

	file, err := os.Open(envFile)
	if err != nil {
		fmt.Println("Error opening .env file:", err)
		return false
	}
	defer file.Close()

	temp, err := os.Create(tempFile)
	if err != nil {
		fmt.Println("Error creating temporary .env file:", err)
		return false
	}
	defer temp.Close()

	scanner := bufio.NewScanner(file)
	var lineDeleted bool

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "OPENAI_KEY=") {
			lineDeleted = true
			continue // Skip the line with the API key
		}
		if _, err := temp.WriteString(line + "\n"); err != nil {
			fmt.Println("Error writing to temporary .env file:", err)
			return false
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading .env file:", err)
		return false
	}

	if !lineDeleted {
		// No need to replace the file if the line wasn't found
		os.Remove(tempFile) // Clean up temp file
		return false
	}

	// Replace the old .env file with the updated one
	if err := os.Rename(tempFile, envFile); err != nil {
		fmt.Println("Error updating .env file:", err)
		return false
	}

	return true
}

// If the CHAT_HISTORY_PATH key does not exist, it adds a new line with the specified path.
func UpdateChatHistoryPath(newPath string) bool {
	envFile := ".env"
	tempFile := ".env.tmp"

	// Open the original .env file for reading
	file, err := os.Open(envFile)
	if err != nil {
		fmt.Println("Error opening .env file:", err)
		return false
	}
	defer file.Close()

	// Create a temporary file for writing the updated content
	temp, err := os.Create(tempFile)
	if err != nil {
		fmt.Println("Error creating temporary .env file:", err)
		return false
	}
	defer temp.Close()

	scanner := bufio.NewScanner(file)
	var pathUpdated bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains CHAT_HISTORY_PATH and update it if so
		if strings.HasPrefix(line, "CHAT_HISTORY_PATH=") {
			line = fmt.Sprintf("CHAT_HISTORY_PATH=%s", newPath)
			pathUpdated = true
		}

		// Write the line (updated or original) to the temporary file
		if _, err := temp.WriteString(line + "\n"); err != nil {
			fmt.Println("Error writing to temporary .env file:", err)
			return false
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading .env file:", err)
		return false
	}

	// If CHAT_HISTORY_PATH was not found in the file, add it
	if !pathUpdated {
		if _, err := temp.WriteString(fmt.Sprintf("CHAT_HISTORY_PATH=%s\n", newPath)); err != nil {
			fmt.Println("Error adding CHAT_HISTORY_PATH to .env file:", err)
			return false
		}
	}

	// Replace the old .env file with the updated temporary file
	if err := os.Rename(tempFile, envFile); err != nil {
		fmt.Println("Error updating .env file:", err)
		return false
	}

	return true
}

func GetInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
