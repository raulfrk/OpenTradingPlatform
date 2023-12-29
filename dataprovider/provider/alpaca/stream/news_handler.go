package stream

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"tradingplatform/dataprovider/data"
	"tradingplatform/dataprovider/provider/alpaca"

	"tradingplatform/dataprovider/provider"
	"tradingplatform/shared/communication/producer"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"

	astream "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

// Handle a news stream request for Alpaca
func handleAlpacaNewsStreamRequest(req requests.StreamRequest) types.StreamResponse {
	account := req.GetAccount()
	// Only default account is supported for now
	// TODO: Allow account separation
	if account != requests.DefaultAccount && account != requests.AnyAccount {
		err := fmt.Errorf("account %s not supported", account)
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("handling Alpaca news stream request")
		return provider.NewStreamError(
			err,
		)
	}

	// Get the client for the account and client lock
	alpacaNewsMapLock.RLock()
	clientLock, okLock := alpacaNewsMapOfLocks[requests.DefaultAccount]
	client, ok := alpacaNewsClientMap[requests.DefaultAccount]
	alpacaNewsMapLock.RUnlock()

	if !okLock {
		logging.Log().Debug().
			RawJSON("request", req.JSON()).
			Msg("creating new Alpaca news client lock")
		clientLock = &sync.RWMutex{}
		alpacaNewsMapLock.Lock()
		alpacaNewsMapOfLocks[requests.DefaultAccount] = clientLock
		alpacaNewsMapLock.Unlock()
	}

	if !ok {
		logging.Log().Debug().
			RawJSON("request", req.JSON()).
			Msg("creating new Alpaca news client")
		client = astream.NewNewsClient(
			astream.WithCredentials(os.Getenv("ALPACA_KEY"), os.Getenv("ALPACA_SECRET")),
		)
		err := client.Connect(context.TODO())
		if err != nil {
			logging.Log().Error().
				Err(err).
				RawJSON("request", req.JSON()).
				Msg("connecting to Alpaca news stream")
			return provider.NewStreamError(err)
		}
		alpacaNewsMapLock.Lock()
		alpacaNewsClientMap[requests.DefaultAccount] = client
		alpacaNewsMapLock.Unlock()
	}

	// Handle request operation
	switch req.GetOperation() {
	case types.StreamAddOp:
		return handleAlpacaNewsStreamAddRequest(client, clientLock, req)
	case types.StreamGetOp:
		return handleAlpacaStreamGetRequest(req, types.News)
	case types.StreamRemoveOp:
		return handleAlpacaNewsStreamRemoveRequest(client, clientLock, req)
	default:
		err := fmt.Errorf("operation %s not supported", req.GetOperation())
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("handling Alpaca news stream request")
		return provider.NewStreamError(
			err,
		)
	}
}

// Handle adding a news stream for Alpaca
func handleAlpacaNewsStreamAddRequest(client *astream.NewsClient,
	clientLock *sync.RWMutex,
	req requests.StreamRequest) types.StreamResponse {

	logging.Log().Info().RawJSON("request", req.JSON()).Msg("adding news stream")
	var symbols []string
	for _, symbol := range req.GetSymbols() {
		// Remove the slash from the symbol (for Crypto)
		symbols = append(symbols, strings.Replace(symbol, "/", "", -1))
	}
	clientLock.Lock()
	err := client.SubscribeToNews(func(n astream.News) {
		for _, symbol := range req.Symbols {
			msg, err := handleOnStreamData[astream.News,
				*sharedent.News](n, types.News, types.RawText, symbol)
			if err != nil {
				js, _ := n.MarshalJSON()
				logging.Log().Error().
					RawJSON("news", js).
					Err(err).Msg("handling news")
				return
			}
			producer.GetStreamHandler(msg.Topic).Ch <- msg
		}

	}, symbols...)
	clientLock.Unlock()

	if err != nil {
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("adding news stream")
		return provider.NewStreamError(err)
	}
	data.AddDataProviderStreamForDType(req, types.RawText)

	return provider.NewStreamResponseAssetClass(
		types.Success,
		"Successfully added news stream",
		nil,
		types.News,
	)
}

// Handle removing a news stream for Alpaca
func handleAlpacaNewsStreamRemoveRequest(client *astream.NewsClient,
	clientLock *sync.RWMutex, req requests.StreamRequest) types.StreamResponse {

	logging.Log().Info().RawJSON("request", req.JSON()).Msg("removing news stream")
	var symbols []string

	for _, symbol := range req.GetSymbols() {
		// Remove the slash from the symbol (for Crypto)
		symbols = append(symbols, strings.Replace(symbol, "/", "", -1))
	}
	clientLock.Lock()
	err := client.UnsubscribeFromNews(symbols...)
	clientLock.Unlock()

	if err != nil {
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("removing news stream")
		return provider.NewStreamError(err)
	}
	for _, symbol := range symbols {
		tTopic := alpaca.NewNewsStreamTopic(types.RawText, symbol).Generate()
		producer.StopTopicHandler(tTopic)
	}
	data.RemoveDataProviderStreamForDType(req, types.RawText)

	return provider.NewStreamResponseAssetClass(
		types.Success,
		"Successfully unsubscribed from news stream",
		nil,
		types.News,
	)
}
