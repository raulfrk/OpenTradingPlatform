package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tradingplatform/shared/logging"
)

func GetOllamaServerURL() string {
	url := os.Getenv("OLLAMA_SERVER_URL")
	if url == "" {
		return "http://localhost:11434"
	}
	return url
}

func HandleAnalysis(ctx context.Context, systemPrompt string, news string, model string) (string, error) {
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
