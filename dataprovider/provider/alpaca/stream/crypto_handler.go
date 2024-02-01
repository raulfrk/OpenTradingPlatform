package stream

import (
	"context"
	"fmt"
	"os"
	"sync"
	"tradingplatform/dataprovider/data"
	"tradingplatform/dataprovider/provider"
	"tradingplatform/dataprovider/provider/alpaca"
	"tradingplatform/shared/communication/producer"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	astream "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

// Handle a crypto stream request for Alpaca and delegate based on operation
func handleAlpacaCryptoStreamRequest(req requests.StreamRequest) types.StreamResponse {
	account := req.GetAccount()
	// Only default account is supported for now
	// TODO: Allow account separation
	if account != requests.DefaultAccount && account != requests.AnyAccount {
		err := fmt.Errorf(
			"account %s not supported, only default and any account options are supported",
			account,
		)
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("handling Alpaca crypto stream request")
		return provider.NewStreamError(
			err,
		)
	}

	// Get the client for the account and client lock
	alpacaCryptoMapLock.RLock()
	clientLock, okLock := alpacaCryptoMapOfLocks[requests.DefaultAccount]
	client, ok := alpacaCryptoClientMap[requests.DefaultAccount]
	alpacaCryptoMapLock.RUnlock()

	if !okLock {
		logging.Log().Debug().
			Str("account", string(requests.DefaultAccount)).
			RawJSON("request", req.JSON()).
			Msg("ceating new Alpaca crypto client lock")
		clientLock = &sync.RWMutex{}
		alpacaCryptoMapLock.Lock()
		alpacaCryptoMapOfLocks[requests.DefaultAccount] = clientLock
		alpacaCryptoMapLock.Unlock()
	}

	if !ok {
		logging.Log().Debug().
			Str("account", string(requests.DefaultAccount)).
			RawJSON("request", req.JSON()).
			Msg("creating new Alpaca crypto client")
		client = astream.NewCryptoClient(marketdata.US,
			astream.WithCredentials(os.Getenv("ALPACA_KEY"), os.Getenv("ALPACA_SECRET")),
		)
		err := client.Connect(context.TODO())
		if err != nil {
			logging.Log().Error().
				Err(err).
				RawJSON("request", req.JSON()).
				Msg("connecting to Alpaca using crypto client")
			return provider.NewStreamError(err)
		}
		alpacaCryptoMapLock.Lock()
		alpacaCryptoClientMap[requests.DefaultAccount] = client
		alpacaCryptoMapLock.Unlock()
	}

	// Handle request operation
	switch req.GetOperation() {
	case types.StreamAddOp:
		return handleAlpacaCryptoStreamAddRequest(client, clientLock, req)
	case types.StreamGetOp:
		return handleAlpacaStreamGetRequest(req, types.Crypto)
	case types.StreamRemoveOp:
		return handleAlpacaCryptoStreamRemoveRequest(client, clientLock, req)
	default:
		err := fmt.Errorf("operation %s not supported", req.GetOperation())
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("handling Alpaca crypto stream request")
		return provider.NewStreamError(
			err,
		)
	}
}

// Handle a crypto stream add request for Alpaca
func handleAlpacaCryptoStreamAddRequest(client *astream.CryptoClient,
	clientLock *sync.RWMutex,
	req requests.StreamRequest) types.StreamResponse {

	logging.Log().Debug().RawJSON("request", req.JSON()).Msg("adding crypto stream")
	dtypes := req.GetDataType()
	symbols := req.GetSymbol()

	dtypesHandled := []types.DataType{}

	for _, dtype := range dtypes {
		var err error
		clientLock.Lock()
		logging.Log().Debug().
			Str("dtype", string(dtype)).
			RawJSON("request", req.JSON()).
			Msg("subscribing")

		switch dtype {
		case types.Bar:
			err = client.SubscribeToBars(func(cb astream.CryptoBar) {
				msg, err := handleOnStreamData[astream.CryptoBar,
					*sharedent.Bar](cb, types.Crypto, types.Bar, cb.Symbol)
				if err != nil {
					js, _ := cb.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("cryptoBar", js).
						Msg("handling crypto bar")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ch <- msg
			}, symbols...)
		case types.Orderbook:
			err = client.SubscribeToOrderbooks(func(ob astream.CryptoOrderbook) {

				msg, err := handleOnStreamData[astream.CryptoOrderbook,
					*sharedent.Orderbook](ob, types.Crypto, types.Orderbook, ob.Symbol)
				if err != nil {
					js, _ := ob.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("cryptoOrderbook", js).
						Msg("handling crypto orderbook")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ch <- msg
			}, symbols...)
		case types.DailyBars:
			err = client.SubscribeToDailyBars(func(db astream.CryptoBar) {
				msg, err := handleOnStreamData[astream.CryptoBar,
					*sharedent.Bar](db, types.Crypto, types.DailyBars, db.Symbol)
				if err != nil {
					js, _ := db.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("cryptoBar", js).
						Msg("handling crypto daily bars")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ch <- msg
			}, symbols...)
		case types.Quotes:
			err = client.SubscribeToQuotes(func(q astream.CryptoQuote) {
				msg, err := handleOnStreamData[astream.CryptoQuote,
					*sharedent.Quote](q, types.Crypto, types.Quotes, q.Symbol)
				if err != nil {
					js, _ := q.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("cryptoQuote", js).
						Msg("handling crypto quotes")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ch <- msg
			}, symbols...)
		case types.Trades:
			err = client.SubscribeToTrades(func(t astream.CryptoTrade) {
				msg, err := handleOnStreamData[astream.CryptoTrade,
					*sharedent.Trade](t, types.Crypto, types.Trades, t.Symbol)
				if err != nil {
					js, _ := t.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("cryptoTrade", js).
						Msg("handling crypto trades")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ch <- msg
			}, symbols...)
		case types.UpdatedBars:
			err = client.SubscribeToUpdatedBars(func(ub astream.CryptoBar) {
				msg, err := handleOnStreamData[astream.CryptoBar,
					*sharedent.Bar](ub, types.Crypto, types.UpdatedBars, ub.Symbol)
				if err != nil {
					js, _ := ub.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("cryptoBar", js).
						Msg("handling crypto updated bars")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ch <- msg
			}, symbols...)
		}
		clientLock.Unlock()
		if err != nil {
			logging.Log().Error().
				Err(err).
				RawJSON("request", req.JSON()).
				Msg("subscribing to crypto stream")
			return provider.NewStreamError(err)
		}
		dtypesHandled = append(dtypesHandled, dtype)
		data.AddDataProviderStreamForDType(req, dtype)
	}

	return provider.NewStreamResponseAssetClass(
		types.Success,
		"Successfully subscribed to crypto streams",
		alpaca.GenerateJSONStreamTopicDict(types.Crypto, dtypesHandled, req.GetSymbol()),
		nil, types.Crypto)
}

func handleAlpacaCryptoStreamRemoveRequest(client *astream.CryptoClient,
	clientLock *sync.RWMutex,
	req requests.StreamRequest) types.StreamResponse {

	logging.Log().Info().RawJSON("request", req.JSON()).Msg("removing crypto stream")
	symbols := req.GetSymbol()

	dtypesHandled := []types.DataType{}

	for _, dtype := range req.GetDataType() {
		clientLock.Lock()
		var err error
		logging.Log().Debug().
			RawJSON("request", req.JSON()).
			Str("dtype", string(dtype)).
			Msg("unsubscribing")
		switch dtype {
		case types.Bar:
			err = client.UnsubscribeFromBars(symbols...)
		case types.Orderbook:
			err = client.UnsubscribeFromOrderbooks(symbols...)
		case types.DailyBars:
			err = client.UnsubscribeFromDailyBars(symbols...)
		case types.Quotes:
			err = client.UnsubscribeFromQuotes(symbols...)
		case types.Trades:
			err = client.UnsubscribeFromTrades(symbols...)
		case types.UpdatedBars:
			err = client.UnsubscribeFromUpdatedBars(symbols...)
		}
		clientLock.Unlock()
		if err != nil {
			logging.Log().Error().
				Err(err).
				RawJSON("request", req.JSON()).
				Msg("unsubscribing from crypto stream")
			return provider.NewStreamError(err)
		}
		for _, symbol := range symbols {
			tTopic := alpaca.NewCryptoStreamTopic(dtype, symbol).Generate()
			producer.StopTopicHandler(tTopic)
		}
		dtypesHandled = append(dtypesHandled, dtype)

		data.RemoveDataProviderStreamForDType(req, dtype)
	}

	return provider.NewStreamResponseAssetClass(
		types.Success,
		"Successfully unsubscribed from crypto streams",
		alpaca.GenerateJSONStreamTopicDict(types.Crypto, dtypesHandled, req.GetSymbol()),
		nil,
		types.Crypto,
	)
}
