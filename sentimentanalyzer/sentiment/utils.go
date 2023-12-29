package sentiment

import (
	"fmt"
	"strings"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"
)

func NewSentimentDataTopic(symbol string, queueID string, queueCount int) utils.Topic {
	return utils.NewDataTopic(types.SentimentAnalyzer, types.Internal, types.News, types.Sentiment, symbol, queueID, queueCount)
}

func ExtractSentimentFromLLMAnswer(answer string) (string, error) {
	// pattern := `(?i)\b(positive|neutral|negative)\b`
	// regExp, err := regexp.Compile(pattern)
	// if err != nil {
	// 	return "", err
	// }

	// // Find all matches in the input string
	// matches := regExp.FindAllString(answer, -1)
	// if len(matches) == 0 {
	// 	return "", fmt.Errorf("no sentiment found in LLM answer: %s", answer)
	// }
	// return strings.ToLower(matches[0]), nil
	lower := strings.ToLower(answer)
	lower = strings.TrimSpace(lower)
	if lower == "positive" || lower == "neutral" || lower == "negative" {
		return lower, nil
	}
	return "", fmt.Errorf("no sentiment found in LLM answer: %s", answer)
}
