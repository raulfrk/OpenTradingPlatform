package alpaca

import (
	"fmt"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"
)

func NewStockStreamTopic(dataType types.DataType, symbol string) utils.Topic {
	return NewStreamTopic(types.Stock, dataType, symbol)
}

func NewCryptoStreamTopic(dataType types.DataType, symbol string) utils.Topic {
	return NewStreamTopic(types.Crypto, dataType, symbol)
}

func NewNewsStreamTopic(dataType types.DataType, symbol string) utils.Topic {
	return NewStreamTopic(types.News, dataType, symbol)
}

func NewStreamTopic(assetClass types.AssetClass, dataType types.DataType, symbol string) utils.Topic {
	return utils.NewStreamTopic(types.DataProvider, types.Alpaca, assetClass, dataType, symbol)
}

func NewDataTopic(assetClass types.AssetClass, dataType types.DataType, symbol string, queueID string, queueCount int) utils.Topic {
	return utils.NewDataTopic(types.DataProvider, types.Alpaca, assetClass, dataType, symbol, queueID, queueCount)
}

func NewBarDataTopic(assetClass types.AssetClass, timeFrame types.TimeFrame, symbol string, queueID string, queueCount int) utils.Topic {
	return utils.NewBarDataTopic(types.DataProvider, types.Alpaca, assetClass, timeFrame, symbol, queueID, queueCount)
}

func GetStreamTopicRoot(assetClass types.AssetClass) string {
	return fmt.Sprintf("%s.%s.%s.%s", types.DataProvider, types.Stream, types.Alpaca, assetClass)
}

func GetDataTopicRoot(assetClass types.AssetClass) string {
	return fmt.Sprintf("%s.%s.%s.%s", types.DataProvider, types.Data, types.Alpaca, assetClass)
}
