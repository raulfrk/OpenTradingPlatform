package gpt4all

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

func GetGPT4ALLServerURL() string {
	url := os.Getenv("GPT4ALL_SERVER_URL")
	if url == "" {
		return "http://localhost:4891"
	}
	return url
}

func HandleAnalysis(ctx context.Context, systemPrompt string, news string, model string) (string, error) {
	url := fmt.Sprintf("%s/v1", GetGPT4ALLServerURL())
	clientConfig := openai.DefaultConfig("")
	clientConfig.BaseURL = url
	client := openai.NewClientWithConfig(clientConfig)

	response, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: news,
			},
		},
	})

	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
