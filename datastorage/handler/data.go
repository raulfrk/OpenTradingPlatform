package handler

import (
	"fmt"

	"tradingplatform/datastorage/data"

	"tradingplatform/shared/communication/producer"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

func HandleDataRequest(dataRequest requests.DataRequest, och chan types.DataResponse) {
	symbol := dataRequest.GetSymbol()
	dtype := dataRequest.GetDataType()
	var messages *[]*entities.Message
	var response types.DataResponse

	logging.Log().Debug().RawJSON("dataRequest", dataRequest.JSON()).Msg("handling data request to db")

	switch dtype {
	case types.Bar:
		messages, response = data.HandleDataFetch[requests.DataRequest,
			*entities.Bar](data.GetBarsFromRequest,
			symbol,
			dataRequest,
			types.Bar,
			dataRequest.AssetClass,
			dataRequest.GetTimeFrame())
	case types.DailyBars:
		messages, response = data.HandleDataFetch[requests.DataRequest,
			*entities.Bar](data.GetDailyBarsFromRequest,
			symbol,
			dataRequest,
			types.DailyBars,
			dataRequest.AssetClass,
			dataRequest.GetTimeFrame())
	case types.LULD:
		messages, response = data.HandleDataFetch[requests.DataRequest,
			*entities.LULD](data.GetLULDFromRequest,
			symbol,
			dataRequest,
			types.LULD,
			dataRequest.AssetClass, "")
	case types.RawText:
		messages, response = data.HandleDataFetch[requests.DataRequest,
			*entities.News](data.GetNewsFromDataRequest,
			symbol,
			dataRequest,
			types.RawText,
			dataRequest.AssetClass, "")
	case types.Orderbook:
		messages, response = data.HandleDataFetch[requests.DataRequest,
			*entities.Orderbook](data.GetOrderbookFromRequest,
			symbol,
			dataRequest,
			types.Orderbook,
			dataRequest.AssetClass,
			"")
	case types.Trades:
		messages, response = data.HandleDataFetch[requests.DataRequest,
			*entities.Trade](data.GetTradesFromRequest,
			symbol,
			dataRequest,
			types.Trades,
			dataRequest.AssetClass,
			"")
	case types.Status:
		messages, response = data.HandleDataFetch[requests.DataRequest,
			*entities.TradingStatus](data.GetTradingStatusesFromRequest,
			symbol,
			dataRequest,
			types.Status,
			dataRequest.AssetClass,
			"")
	case types.Quotes:
		messages, response = data.HandleDataFetch[requests.DataRequest,
			*entities.Quote](data.GetQuoteFromRequest,
			symbol,
			dataRequest,
			types.Quotes,
			dataRequest.AssetClass,
			"")
	default:
		och <- types.NewDataError(
			fmt.Errorf("invalid data type %s", dtype),
		)
		return
	}
	if response.Err != "" {
		och <- response
		return
	}
	if len(*messages) == 0 {
		och <- types.NewDataError(
			fmt.Errorf("no data found for %s", symbol),
		)
		return
	}
	responseTopic := (*messages)[0].Topic
	handler, handlerResponse := producer.GetQueueHandler(responseTopic, dataRequest.GetNoConfirm())
	if handlerResponse.Err != "" {
		och <- handlerResponse
		return
	}
	handler.Ch <- messages

	och <- response
}
