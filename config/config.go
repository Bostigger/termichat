package config

import (
	"fmt"
	"os"
)

var OpenAIKey string
var ChatHistoryPath string

func LoadConfig() error {
	OpenAIKey = os.Getenv("OPENAI_KEY")
	ChatHistoryPath = os.Getenv("CHAT_HISTORY_PATH")
	if OpenAIKey == "" {
		return fmt.Errorf("OpenAI API key not found in environment variables")
	}
	if ChatHistoryPath == "" {
		return fmt.Errorf("Chat history path not set in environment variables")
	}
	return nil
}
