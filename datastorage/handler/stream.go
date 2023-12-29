package handler

import (
	"fmt"
	"tradingplatform/datastorage/subscriber"
	shsubscriber "tradingplatform/shared/communication/subscriber"
	"tradingplatform/shared/data"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

func HandleStreamRequest(req requests.StreamSubscribeRequest) types.StreamResponse {
	switch req.Operation {
	case types.StreamAddOp:
		for _, su := range req.StreamSubscribeWithAgents {
			shsubscriber.AttatchFunctionalityRoundRobin(su.Topic, subscriber.HandleStoreData, su.AgentCount)
			data.AddSubscribedTopic(su.Topic, su.AgentCount)
		}
	case types.StreamRemoveOp:
		for _, su := range req.StreamSubscribeWithAgents {
			shsubscriber.StopTopicHandler(su.Topic)
			data.RemoveSubscribedTopic(su.Topic)
		}
	default:
		return subscriber.NewStreamErrorResponseTopic(
			fmt.Errorf("operation %s not supported", req.Operation),
		)
	}

	return subscriber.NewStreamResponseTopic(
		types.Success,
		"Successfully added stream/s",
		nil,
	)
}
