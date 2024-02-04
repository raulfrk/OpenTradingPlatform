package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

func GetOllamaServerURL() string {
	url := os.Getenv("OLLAMA_SERVER_URL")
	if url == "" {
		return "http://localhost:11434"
	}
	return url
}

// TODO: Extract general functionality from here and leave only ollama specific things
// HandleAnalysisNews handles the semtiment analysis of news using the ollama tool
func HandleAnalysisNews(ctx context.Context, news *entities.News, req *requests.SentimentAnalysisRequest) (string, error) {
	var systemPrompt string
	var err error
	var newsText string
	systemPrompt, err = req.GetSystemPrompt()

	switch req.SentimentAnalysisProcess {
	case types.Plain:
		newsText = fmt.Sprintf("Symbol:%s News: %s", req.GetSymbol(), news.Headline)
	case types.Semantic:
		newsText = fmt.Sprintf("Symbols:%s News: %s", strings.Join(news.Symbols, ","), news.Headline)

	default:
		return "", fmt.Errorf("sentiment analysis process %s does not have an implementation", req.SentimentAnalysisProcess)
	}
	if err != nil {
		return "", err
	}

	return handleAnalysis(ctx, systemPrompt, newsText, req.Model)
}

func handleAnalysis(ctx context.Context, systemPrompt string, news string, model string) (string, error) {
	url := fmt.Sprintf("%s/api/chat", GetOllamaServerURL())

	// Define the request payload
	requestPayload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": news,
			},
		},
		"stream": false,
	}

	// Convert request payload to JSON
	requestPayloadJSON, err := json.Marshal(requestPayload)
	if err != nil {
		logging.Log().Error().Err(err).Msg("marshalling request payload to JSON")
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestPayloadJSON))
	if err != nil {
		logging.Log().Error().Err(err).Msg("creating HTTP request")
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Log().Error().Err(err).Msg("sending HTTP request")
		return "", err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		logging.Log().Error().Err(err).Msg("decoding response body")
		return "", err
	}
	if response["message"] == nil {
		return "", nil
	}
	message := response["message"].(map[string]interface{})
	assistantContent := message["content"].(string)

	return assistantContent, nil
}
