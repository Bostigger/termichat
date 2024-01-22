package chat

import "termichat/internal/api"

func ProcessUserInput(input string) (string, error) {
	response, err := api.QueryOpenAI(input)
	if err != nil {
		return "", err
	}
	return response, nil
}
