package subscriber

import (
	"tradingplatform/datastorage/data"

	"tradingplatform/shared/communication/subscriber"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"
)

func HandleStoreDataFromQueue(msg *entities.Message) {
	queue := subscriber.DrainQueue(msg.Topic)
	if len(queue) == 0 {
		return
	}
	switch msg.DataType {
	case string(types.Bar):
		utils.HandleEntityQueueWithConversion(queue,
			&entities.Bar{},
			data.InsertBatchEntity[data.Bar],
			data.BarsFromEntities)
	case string(types.DailyBars):
		utils.HandleEntityQueueWithConversion(queue,
			&entities.Bar{},
			data.InsertBatchEntity[data.DailyBar],
			data.DailyBarsFromEntities)
	case string(types.Trades):
		utils.HandleEntityQueueWithConversion(queue,
			&entities.Trade{},
			data.InsertBatchEntity[data.Trade],
			data.TradesFromEntities)
	case string(types.Quotes):
		utils.HandleEntityQueueWithConversion(queue,
			&entities.Quote{},
			data.InsertBatchEntity[data.Quote],
			data.QuotesFromEntities)
	case string(types.LULD):
		utils.HandleEntityQueueWithConversion(queue,
			&entities.LULD{},
			data.InsertBatchEntity[data.LULD],
			data.LULDsFromEntities)
	case string(types.Orderbook):
		utils.HandleEntityQueueWithConversion(queue,
			&entities.Orderbook{},
			data.InsertBatchEntity[data.Orderbook],
			data.OrderbooksFromEntities)
	case string(types.Status):
		utils.HandleEntityQueueWithConversion(queue,
			&entities.TradingStatus{},
			data.InsertBatchEntity[data.TradingStatus],
			data.TradingStatusesFromEntities)
	case string(types.RawText):
		utils.HandleEntityQueueWithConversion(queue,
			&entities.News{},
			data.InsertBatchNews,
			data.NewsFromEntities)
	}
}

func HandleStoreData(msg *entities.Message) {
	if subscriber.IsQueue(msg.Topic) {
		HandleStoreDataFromQueue(msg)
		return
	}
	switch msg.DataType {
	case string(types.Bar):
		utils.HandleEntity[*entities.Bar](msg, &entities.Bar{}, data.InsertBar)
	case string(types.DailyBars):
		utils.HandleEntity[*entities.Bar](msg, &entities.Bar{}, data.InsertDailyBar)
	case string(types.Trades):
		utils.HandleEntity[*entities.Trade](msg, &entities.Trade{}, data.InsertTrade)
	case string(types.Quotes):
		utils.HandleEntity[*entities.Quote](msg, &entities.Quote{}, data.InsertQuote)
	case string(types.LULD):
		utils.HandleEntity[*entities.LULD](msg, &entities.LULD{}, data.InsertLULD)
	case string(types.Log):
		data.InsertLog(msg.Payload)
	case string(types.Orderbook):
		utils.HandleEntity[*entities.Orderbook](msg, &entities.Orderbook{}, data.InsertOrderbook)
	case string(types.RawText):
		utils.HandleEntity[*entities.News](msg, &entities.News{}, data.InsertNews)
	case string(types.Status):
		utils.HandleEntity[*entities.TradingStatus](msg, &entities.TradingStatus{}, data.InsertTradingStatus)
	}
}
