package data

import (
	"fmt"
	"os"
	"time"
	"tradingplatform/dataprovider/requests"
	"tradingplatform/shared/communication/producer"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type ClientWrapper struct {
	client *marketdata.Client
}

func (c *ClientWrapper) getNewsWrapper(symbol string, req marketdata.GetNewsRequest) ([]marketdata.News, error) {
	req.Symbols = []string{symbol}
	// This is used to get around the fact that each response can have max 50 elements
	// by paginating with a page limit of 50 all the elements are returned
	// (as long as max req number is not exceeded)
	req.NoTotalLimit = true
	req.PageLimit = 50
	req.IncludeContent = true
	return c.client.GetNews(req)
}

func handleAlpacaNewsDataRequest(req requests.DataRequest) types.DataResponse {
	var response types.DataResponse
	var messages *[]*sharedent.Message
	logging.Log().Debug().RawJSON("request", req.JSON()).Msg("handling alpaca news data request")
	client := ClientWrapper{
		client: marketdata.NewClient(marketdata.ClientOpts{
			APIKey:    os.Getenv("ALPACA_KEY"),
			APISecret: os.Getenv("ALPACA_SECRET"),
		})}
	symbol := req.GetSymbols()[0]

	messages, response = handleDataFetch[marketdata.GetNewsRequest, marketdata.News, *sharedent.News](client.getNewsWrapper, symbol, marketdata.GetNewsRequest{
		PageLimit: 10000,
		Start:     time.Unix(req.GetStartTime(), 0),
		End:       time.Unix(req.GetEndTime(), 0),
	}, types.RawText, types.News, req.GetTimeFrame())
	if response.Err != "" {
		return response
	}
	if len(*messages) == 0 {
		return types.NewDataError(
			fmt.Errorf("no data found for symbol %s", symbol),
		)
	}
	responseTopic := (*messages)[0].Topic
	handler, handlerResponse := producer.GetQueueHandler(responseTopic, req.GetNoConfirm())
	if handlerResponse.Err != "" {
		return handlerResponse
	}
	handler.Ich <- messages

	return response
}
