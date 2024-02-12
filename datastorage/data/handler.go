package data

import (
	"tradingplatform/datastorage/utils"
	"tradingplatform/shared/communication/producer"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
)

// HandleDataFetch in a generic way to handle data fetches from the database
func HandleDataFetch[T any,
	V entities.FingerprintablePayloader](
	fun func(string, T) ([]V, error),
	symbol string,
	req T,
	dtype types.DataType, assetClass types.AssetClass,
	timeFrame types.TimeFrame) (*[]*entities.Message, types.DataResponse) {

	result, err := fun(symbol, req)
	if err != nil {
		logging.Log().Error().
			Err(err).
			Str("symbol", symbol).
			Str("dtype", string(dtype)).
			Str("assetClass", string(assetClass)).
			Str("timeFrame", string(timeFrame)).
			Msg("getting data from database")
		return nil, types.NewDataError(err)
	}
	var responseTopic = ""
	var queueID = ""

	var messages []*entities.Message
	queueID = producer.GenerateQueueID()
	for _, entity := range result {
		messages = append(messages,
			entities.GenerateMessage(entity,
				dtype,
				responseTopic))
	}

	if dtype == types.Bar {
		responseTopic = utils.NewBarDataTopic(assetClass,
			timeFrame, symbol, queueID, len(messages)).Generate()
	} else {
		responseTopic = utils.NewDataTopic(assetClass,
			dtype, symbol, queueID, len(messages)).Generate()
	}
	for _, message := range messages {
		message.Topic = responseTopic
	}
	return &messages, types.NewDataResponse(
		types.Success,
		"Successfully retrieved data",
		nil,
		responseTopic,
	)
}
