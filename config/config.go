package config

import (
	"fmt"
	"os"
)

var OpenAIKey string

func LoadConfig() error {
	OpenAIKey = os.Getenv("OPENAI_KEY")
	if OpenAIKey == "" {
		return fmt.Errorf("OpenAI API key not found in environment variables")
	}
	return nil
}
