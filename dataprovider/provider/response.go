package provider

import (
	"encoding/json"
	"tradingplatform/dataprovider/data"
	"tradingplatform/shared/types"
)

func NewStreamError(err error) types.StreamResponse {
	return NewStreamResponse(types.Failure, "", err)
}

func NewStreamResponse(status types.OpStatus, message string, err error) types.StreamResponse {
	streams := data.GetActiveStreams()
	streamsJSON, _ := json.Marshal(streams)
	return types.NewStreamResponse(status, message, err, string(streamsJSON))
}

func NewStreamResponseAssetClass(status types.OpStatus, message string, err error, assetClass types.AssetClass) types.StreamResponse {
	streams := data.GetActiveStreamsAssetClass(assetClass)
	streamsJSON, _ := json.Marshal(streams)
	return types.NewStreamResponse(status, message, err, string(streamsJSON))
}
