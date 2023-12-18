package stream

import (
	"context"
	"fmt"
	"os"
	"sync"
	"tradingplatform/shared/communication/producer"
	"tradingplatform/shared/logging"

	"tradingplatform/dataprovider/data"
	"tradingplatform/dataprovider/provider"
	"tradingplatform/dataprovider/requests"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/types"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	astream "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

func handleAlpacaStockStreamRequest(req requests.StreamRequest) types.StreamResponse {
	account := req.GetAccount()

	// Only default account is supported for now
	// TODO: Allow account separation
	if account != requests.DefaultAccount && account != requests.AnyAccount {
		err := fmt.Errorf("account %s not supported", account)
		logging.Log().Error().
			Err(err).
			Str("account", string(account)).
			RawJSON("request", req.JSON()).
			Msg("handling Alpaca stock stream request")

		return provider.NewStreamError(
			err,
		)
	}

	// Get the client for the account and client lock
	alpacaStocksMapLock.RLock()
	clientLock, okLock := alpacaStockMapOfLocks[requests.DefaultAccount]
	client, ok := alpacaStocksClientMap[requests.DefaultAccount]
	alpacaStocksMapLock.RUnlock()

	if !okLock {
		logging.Log().Debug().
			Str("account", string(requests.DefaultAccount)).
			RawJSON("request", req.JSON()).
			Msg("creating new Alpaca stocks client lock")
		clientLock = &sync.RWMutex{}
		alpacaStocksMapLock.Lock()
		alpacaStockMapOfLocks[requests.DefaultAccount] = clientLock
		alpacaStocksMapLock.Unlock()
	}

	if !ok {
		logging.Log().Debug().
			Str("account", string(requests.DefaultAccount)).
			RawJSON("request", req.JSON()).
			Msg("creating new Alpaca stocks client")
		client = astream.NewStocksClient(
			marketdata.IEX,
			astream.WithCredentials(
				os.Getenv("ALPACA_KEY"),
				os.Getenv("ALPACA_SECRET"),
			),
		)
		err := client.Connect(context.TODO())
		if err != nil {
			logging.Log().Error().
				Err(err).
				RawJSON("request", req.JSON()).
				Msg("connecting to Alpaca stocks stream")
			return provider.NewStreamError(err)
		}
		alpacaStocksMapLock.Lock()
		alpacaStocksClientMap[requests.DefaultAccount] = client
		alpacaStocksMapLock.Unlock()
	}

	// Handle the request
	switch req.GetOperation() {
	case types.StreamAddOp:
		return handleAlpacaStockStreamAddRequest(client, clientLock, req)
	case types.StreamGetOp:
		return handleAlpacaStreamGetRequest(req, types.Stock)
	case types.StreamRemoveOp:
		return handleAlpacaStockStreamRemoveRequest(client, clientLock, req)
	default:
		err := fmt.Errorf("request type %s not supported", req.GetOperation())
		logging.Log().Error().
			Err(err).
			RawJSON("request", req.JSON()).
			Msg("handling Alpaca stocks stream request")
		return provider.NewStreamError(err)
	}
}

func handleAlpacaStockStreamAddRequest(client *astream.StocksClient, clientLock *sync.RWMutex, req requests.StreamRequest) types.StreamResponse {
	dtypes := req.GetDataTypes()
	logging.Log().Debug().RawJSON("request", req.JSON()).Msg("adding stocks stream")
	for _, dtype := range dtypes {
		var err error
		logging.Log().Debug().RawJSON("request", req.JSON()).Str("dtype", string(dtype)).Msg("subscribing")
		clientLock.Lock()
		switch dtype {
		case types.Bar:
			err = client.SubscribeToBars(func(b astream.Bar) {
				msg, err := handleOnStreamData[astream.Bar, *sharedent.Bar](b, types.Stock, types.Bar, b.Symbol)
				if err != nil {
					js, _ := b.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("bar", js).
						Msg("handling bar stream")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ich <- msg
			}, req.GetSymbols()...)
		case types.DailyBars:
			err = client.SubscribeToDailyBars(func(b astream.Bar) {
				msg, err := handleOnStreamData[astream.Bar, *sharedent.Bar](b, types.Stock, types.DailyBars, b.Symbol)
				if err != nil {
					js, _ := b.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("bar", js).
						Msg("handling daily bar stream")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ich <- msg
			}, req.GetSymbols()...)
		case types.UpdatedBars:
			err = client.SubscribeToUpdatedBars(func(b astream.Bar) {
				msg, err := handleOnStreamData[astream.Bar, *sharedent.Bar](b, types.Stock, types.UpdatedBars, b.Symbol)
				if err != nil {
					js, _ := b.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("bar", js).
						Msg("handling updated bar stream")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ich <- msg
			}, req.GetSymbols()...)
		case types.Trades:
			err = client.SubscribeToTrades(func(t astream.Trade) {
				msg, err := handleOnStreamData[astream.Trade, *sharedent.Trade](t, types.Stock, types.Trades, t.Symbol)
				if err != nil {
					js, _ := t.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("trade", js).
						Msg("handling trade stream")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ich <- msg
			}, req.GetSymbols()...)

		case types.Quotes:
			err = client.SubscribeToQuotes(func(q astream.Quote) {
				msg, err := handleOnStreamData[astream.Quote, *sharedent.Quote](q, types.Stock, types.Trades, q.Symbol)
				if err != nil {
					js, _ := q.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("quote", js).
						Msg("handling quote stream")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ich <- msg
			}, req.GetSymbols()...)
		case types.LULD:
			err = client.SubscribeToLULDs(func(luld astream.LULD) {
				msg, err := handleOnStreamData[astream.LULD, *sharedent.LULD](luld, types.Stock, types.Trades, luld.Symbol)
				if err != nil {
					js, _ := luld.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("luld", js).
						Msg("handling LULD stream")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ich <- msg
			}, req.GetSymbols()...)
		case types.Status:
			err = client.SubscribeToStatuses(func(s astream.TradingStatus) {
				msg, err := handleOnStreamData[astream.TradingStatus, *sharedent.TradingStatus](s, types.Stock, types.Trades, s.Symbol)
				if err != nil {
					js, _ := s.MarshalJSON()
					logging.Log().Error().
						Err(err).
						RawJSON("status", js).
						Msg("handling status stream")
					return
				}
				producer.GetStreamHandler(msg.Topic).Ich <- msg
			}, req.GetSymbols()...)
		default:
			err = fmt.Errorf("data type %s not supported yet", dtype)
		}
		clientLock.Unlock()
		if err != nil {
			logging.Log().Error().
				Err(err).
				RawJSON("request", req.JSON()).
				Msg("subscribing to Alpaca stocks stream")
			return provider.NewStreamError(err)
		}
		data.AddActiveStreamForDType(req, dtype)
	}
	return provider.NewStreamResponseAssetClass(
		types.Success,
		"Successfully added Alpaca stocks stream",
		nil,
		types.Stock,
	)
}

func handleAlpacaStockStreamRemoveRequest(client *astream.StocksClient, clientLock *sync.RWMutex, req requests.StreamRequest) types.StreamResponse {
	logging.Log().Info().RawJSON("request", req.JSON()).Msg("removing crypto stream")
	symbols := req.GetSymbols()

	for _, dtype := range req.GetDataTypes() {
		clientLock.Lock()
		var err error
		logging.Log().Debug().RawJSON("request", req.JSON()).Str("dtype", string(dtype)).Msg("unsubscribing")
		switch dtype {
		case types.Bar:
			err = client.UnsubscribeFromBars(symbols...)
		case types.LULD:
			err = client.UnsubscribeFromLULDs(symbols...)
		case types.Status:
			err = client.UnsubscribeFromStatuses(symbols...)
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
				Msg("unsubscribing from stock stream")
			return provider.NewStreamError(err)
		}
		data.RemoveActiveStreamForDType(req, dtype)
	}

	return provider.NewStreamResponseAssetClass(
		types.Success,
		"Successfully unsubscribed from stock streams",
		nil,
		types.Stock,
	)
}
