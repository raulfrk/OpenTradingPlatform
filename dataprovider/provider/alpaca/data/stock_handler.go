package data

import (
	"fmt"
	"os"
	"time"
	"tradingplatform/dataprovider/provider/alpaca"
	"tradingplatform/shared/communication/producer"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

// Handle a stock data request for alpaca
func handleAlpacaStockDataRequest(req requests.DataRequest) types.DataResponse {
	var response types.DataResponse
	var messages *[]*sharedent.Message
	client := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    os.Getenv("ALPACA_KEY"),
		APISecret: os.Getenv("ALPACA_SECRET"),
	})
	symbol := req.GetSymbol()
	dtype := req.GetDataType()
	logging.Log().Debug().RawJSON("request", req.JSON()).Msg("handling alpaca stock data request")
	switch dtype {
	case types.Bar:
		messages, response = handleDataFetch[marketdata.GetBarsRequest,
			marketdata.Bar, *sharedent.Bar](client.GetBars, symbol, marketdata.GetBarsRequest{
			TimeFrame: alpaca.GetAlpacaTimeFrame(req.GetTimeFrame()),
			PageLimit: 10000,
			Start:     time.Unix(req.GetStartTime(), 0),
			End:       time.Unix(req.GetEndTime(), 0),
		}, types.Bar, types.Stock, req.GetTimeFrame())

	case types.Trades:
		messages, response = handleDataFetch[marketdata.GetTradesRequest,
			marketdata.Trade, *sharedent.Trade](client.GetTrades, symbol, marketdata.GetTradesRequest{
			Feed:      marketdata.IEX,
			PageLimit: 10000,
			Start:     time.Unix(req.GetStartTime(), 0),
			End:       time.Unix(req.GetEndTime(), 0),
		}, types.Trades, types.Stock, req.GetTimeFrame())
	case types.Quotes:
		messages, response = handleDataFetch[marketdata.GetQuotesRequest,
			marketdata.Quote, *sharedent.Quote](client.GetQuotes, symbol, marketdata.GetQuotesRequest{
			PageLimit: 10000,
			Start:     time.Unix(req.GetStartTime(), 0),
			End:       time.Unix(req.GetEndTime(), 0),
		}, types.Quotes, types.Stock, req.GetTimeFrame())

	default:
		return types.NewDataError(
			fmt.Errorf("invalid data type %s", dtype),
		)
	}
	if response.Err != "" {
		return response
	}
	if len(*messages) == 0 {
		return types.NewDataError(
			fmt.Errorf("no data found for %s", symbol),
		)
	}
	responseTopic := (*messages)[0].Topic
	handler, handlerResponse := producer.GetQueueHandler(responseTopic, req.GetNoConfirm())
	if handlerResponse.Err != "" {
		return handlerResponse
	}
	handler.Ch <- messages

	return response
}
