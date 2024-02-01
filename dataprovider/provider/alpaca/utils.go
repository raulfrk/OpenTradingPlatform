package alpaca

import (
	"encoding/json"
	"fmt"
	"tradingplatform/shared/logging"
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

func GenerateJSONStreamTopicDict(assetClass types.AssetClass, dataTypes []types.DataType, symbols []string) string {
	tmap := map[types.DataType][]string{}
	for _, dataType := range dataTypes {
		value, exists := tmap[dataType]
		if !exists {
			tmap[dataType] = []string{}
		}
		for _, symbol := range symbols {
			value = append(value, NewStreamTopic(assetClass, dataType, symbol).Generate())
		}
		tmap[dataType] = value
	}
	json, err := json.Marshal(tmap)
	dtypesStr := []string{}
	for _, dtype := range dataTypes {
		dtypesStr = append(dtypesStr, string(dtype))
	}
	if err != nil {
		logging.Log().Error().Str("assetClass", string(assetClass)).Strs("dataTypes", dtypesStr).Strs("symbols", symbols).Err(err).Msg("generating topics")
		return ""
	}
	return string(json)
}
