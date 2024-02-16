package config

import (
	"fmt"
	"os"
)

var OpenAIKey string
var ChatHistoryPath string
var AIModelVersion string // New variable for AI model version

func LoadConfig() error {
	OpenAIKey = os.Getenv("OPENAI_KEY")
	ChatHistoryPath = os.Getenv("CHAT_HISTORY_PATH")
	AIModelVersion = os.Getenv("AI_MODEL_VERSION") // Load AI model version from environment
	if OpenAIKey == "" {
		return fmt.Errorf("OpenAI API key not found in environment variables")
	}
	if ChatHistoryPath == "" {
		return fmt.Errorf("Chat history path not set in environment variables")
	}
	if AIModelVersion == "" {
		return fmt.Errorf("AI model version not set in environment variables")
	}
	return nil
}
