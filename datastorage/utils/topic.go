package utils

import (
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"
)

func NewBarDataTopic(assetClass types.AssetClass,
	timeFrame types.TimeFrame,
	symbol string, queueID string, queueCount int) utils.Topic {
	return utils.NewBarDataTopic(types.DataStorage,
		types.Internal,
		assetClass,
		timeFrame,
		symbol,
		queueID,
		queueCount)
}

func NewDataTopic(assetClass types.AssetClass,
	dataType types.DataType, symbol string, queueID string, queueCount int) utils.Topic {
	return utils.NewDataTopic(types.DataStorage,
		types.Internal,
		assetClass,
		dataType,
		symbol,
		queueID,
		queueCount)
}
