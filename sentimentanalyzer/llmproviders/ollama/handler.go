package ollama

import (
	"context"
	"fmt"
	"strings"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"

	"github.com/tmc/langchaingo/llms/ollama"
)

func HandleAnalysisNews(ctx context.Context, news *entities.News, req *requests.SentimentAnalysisRequest) (string, error) {
	var systemPrompt string
	var err error
	var newsText string
	systemPrompt, err = req.GetSystemPrompt()

	switch req.SentimentAnalysisProcess {
	case types.Plain:
		newsText = fmt.Sprintf("Symbol:%s\nNews:%s", news.Headline, req.Symbols[0])
	case types.Semantic:
		newsText = fmt.Sprintf("Symbols:%s\nNews:%s", news.Headline, strings.Join(news.Symbols, ","))

	default:
		return "", fmt.Errorf("sentiment analysis process %s does not have an implementation", req.SentimentAnalysisProcess)
	}
	if err != nil {
		return "", err
	}

	return HandleAnalysis(ctx, systemPrompt, newsText, req.Model)
}

func HandleAnalysis(ctx context.Context, systemPrompt string, news string, model string) (string, error) {
	llm, err := ollama.New(ollama.WithModel(model), ollama.WithSystemPrompt(systemPrompt))
	if err != nil {
		return "", fmt.Errorf("while creating ollama implementation: %v", err)
	}

	query := news
	completion, err := llm.Call(ctx, query)
	if err != nil {
		return "", fmt.Errorf("while calling ollama implementation: %v", err)

	}
	return completion, nil
}