package subscriber

import (
	"encoding/json"
	"tradingplatform/shared/data"
	"tradingplatform/shared/types"
)

func NewStreamResponseTopic(status types.OpStatus, message string, err error) types.StreamResponse {
	topics := data.GetSubscribedTopics()
	json, _ := json.Marshal(topics)
	return types.NewStreamResponseTopic(status, message, err, string(json))
}

func NewStreamErrorResponseTopic(err error) types.StreamResponse {
	return types.NewStreamResponseTopic(types.Failure, "", err, "")
}
