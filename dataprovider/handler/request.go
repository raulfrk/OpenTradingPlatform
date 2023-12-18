package handler

import (
	"fmt"
	"tradingplatform/dataprovider/provider"
	"tradingplatform/dataprovider/provider/alpaca/data"
	alpacaStream "tradingplatform/dataprovider/provider/alpaca/stream"
	"tradingplatform/dataprovider/requests"
	"tradingplatform/shared/types"
)

func HandleStreamRequest(req requests.StreamRequest) string {
	source := req.GetSource()

	switch source {
	case requests.Alpaca:
		return alpacaStream.HandleAlpacaStreamRequest(req).Respond()
	default:
		invalidSourceError := provider.NewStreamError(
			fmt.Errorf("invalid source %s", source),
		)
		return invalidSourceError.Respond()
	}
}

func HandleDataRequest(dataRequest requests.DataRequest, och chan types.DataResponse) {
	switch dataRequest.GetSource() {
	case requests.Alpaca:
		och <- data.HandleAlpacaDataRequest(dataRequest)
	default:
		och <- types.NewDataError(
			fmt.Errorf("invalid source %s", dataRequest.GetSource()),
		)
	}
}
