package provider

import (
	"encoding/json"
	"tradingplatform/dataprovider/data"
	"tradingplatform/shared/types"
)

// New dataprovider-specific stream error
func NewStreamError(err error) types.StreamResponse {
	return NewStreamResponse(types.Failure, "", "", err)
}

// New dataprovider-specific stream response
func NewStreamResponse(status types.OpStatus, message string, topics string, err error) types.StreamResponse {
	streams := data.GetDataProviderStreams()
	streamsJSON, _ := json.Marshal(streams)
	return types.NewStreamResponse(status, message, err, string(streamsJSON), topics)
}

// New dataprovider-specific stream response for a given asset class
func NewStreamResponseAssetClass(status types.OpStatus, message string, topics string, err error,
	assetClass types.AssetClass) types.StreamResponse {

	streams := data.GetDataProviderStreamsAssetClass(assetClass)
	streamsJSON, _ := json.Marshal(streams)
	return types.NewStreamResponse(status, message, err, string(streamsJSON), topics)
}
