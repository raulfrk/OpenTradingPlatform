package data

import (
	"fmt"
	"time"
	"tradingplatform/dataprovider/provider/alpaca"
	"tradingplatform/shared/communication/producer"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"

	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/types"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

// Handle a crypto data request for Alpaca
func handleAlpacaCryptoDataRequest(req requests.DataRequest) types.DataResponse {
	symbol := req.GetSymbols()[0]
	dtype := req.GetDataTypes()[0]
	var messages *[]*sharedent.Message
	var response types.DataResponse
	logging.Log().Debug().RawJSON("request", req.JSON()).Msg("handling alpaca crypto data request")
	switch dtype {
	case types.Bar:
		messages, response = handleDataFetch[marketdata.GetCryptoBarsRequest,
			marketdata.CryptoBar,
			*sharedent.Bar](marketdata.GetCryptoBars, symbol, marketdata.GetCryptoBarsRequest{
			TimeFrame: alpaca.GetAlpacaTimeFrame(req.GetTimeFrame()),
			PageLimit: 10000,
			Start:     time.Unix(req.GetStartTime(), 0),
			End:       time.Unix(req.GetEndTime(), 0),
		}, types.Bar, types.Crypto, req.GetTimeFrame())
	case types.Trades:
		messages, response = handleDataFetch[marketdata.GetCryptoTradesRequest,
			marketdata.CryptoTrade,
			*sharedent.Trade](marketdata.GetCryptoTrades, symbol, marketdata.GetCryptoTradesRequest{
			PageLimit: 10000,
			Start:     time.Unix(req.GetStartTime(), 0),
			End:       time.Unix(req.GetEndTime(), 0),
		}, types.Trades, types.Crypto, req.GetTimeFrame())
	case types.Quotes:
		messages, response = handleDataFetch[marketdata.GetCryptoQuotesRequest,
			marketdata.CryptoQuote,
			*sharedent.Quote](marketdata.GetCryptoQuotes, symbol, marketdata.GetCryptoQuotesRequest{
			PageLimit: 10000,
			Start:     time.Unix(req.GetStartTime(), 0),
			End:       time.Unix(req.GetEndTime(), 0),
		}, types.Quotes, types.Crypto, req.GetTimeFrame())
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
