package main

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// InitOpenAIClient initializes the OpenAI client with the given API key
func InitOpenAIClient(apiKey string) *openai.Client {
	return openai.NewClient(option.WithAPIKey(apiKey))
}

// GetOpenAIResponse sends messages to OpenAI and returns the assistant's reply
func GetOpenAIResponse(client *openai.Client, messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F(messages),
			Model:    openai.F(openai.ChatModelGPT3_5Turbo),
		},
	)
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
