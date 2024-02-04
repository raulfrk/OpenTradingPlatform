package sentiment

import (
	"fmt"
	"strings"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"
)

// NewSentimentDataTopic creates a new topic for sentiment data
func NewSentimentDataTopic(symbol string, queueID string, queueCount int) utils.Topic {
	return utils.NewDataTopic(types.SentimentAnalyzer, types.Internal, types.News, types.Sentiment, symbol, queueID, queueCount)
}

// ExtractSentimentFromLLMAnswer extracts the sentiment from the answer of the LLM
func ExtractSentimentFromLLMAnswer(answer string) (string, error) {
	lower := strings.ToLower(answer)
	lower = strings.ReplaceAll(lower, ".", "")
	lower = strings.ReplaceAll(lower, ":", "")
	lower = strings.ReplaceAll(lower, "sentiment", "")
	lower = strings.TrimSpace(lower)
	if lower == "positive" || lower == "neutral" || lower == "negative" {
		return lower, nil
	}
	return "", fmt.Errorf("no sentiment found in LLM answer")
}
