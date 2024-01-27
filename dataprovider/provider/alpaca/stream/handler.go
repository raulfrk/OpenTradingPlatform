package stream

import (
	"errors"
	"fmt"
	"sync"
	"tradingplatform/dataprovider/provider"
	"tradingplatform/dataprovider/provider/alpaca"
	"tradingplatform/dataprovider/requests"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"

	astream "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

var alpacaCryptoClientMap map[requests.Account]*astream.CryptoClient = make(map[requests.Account]*astream.CryptoClient)
var alpacaStocksClientMap map[requests.Account]*astream.StocksClient = make(map[requests.Account]*astream.StocksClient)
var alpacaNewsClientMap map[requests.Account]*astream.NewsClient = make(map[requests.Account]*astream.NewsClient)

var alpacaCryptoMapLock sync.RWMutex
var alpacaStocksMapLock sync.RWMutex
var alpacaNewsMapLock sync.RWMutex

var alpacaCryptoMapOfLocks map[requests.Account]*sync.RWMutex = make(map[requests.Account]*sync.RWMutex)
var alpacaStockMapOfLocks map[requests.Account]*sync.RWMutex = make(map[requests.Account]*sync.RWMutex)
var alpacaNewsMapOfLocks map[requests.Account]*sync.RWMutex = make(map[requests.Account]*sync.RWMutex)

// TODO: Add description
func handleOnStreamData[T any, V sharedent.Payloader](entity T, assetClass types.AssetClass, dtype types.DataType, symbol string) (*sharedent.Message, error) {
	newEntity := alpaca.MapEntityWithReturnEntity(entity, symbol)
	if newEntity == nil {
		err := fmt.Errorf("error mapping entity %v", entity)
		logging.Log().Error().
			Err(err).
			Interface("entity", entity).
			Str("dtype", string(dtype)).
			Str("assetClass", string(assetClass)).
			Str("symbol", symbol).
			Send()
		return nil, err
	}
	payloaderEntity, ok := newEntity.(V)
	if !ok {
		err := fmt.Errorf("error casting entity %v to payloader", entity)
		logging.Log().Error().
			Err(err).
			Interface("entity", entity).
			Str("dtype", string(dtype)).
			Str("assetClass", string(assetClass)).
			Str("symbol", symbol).
			Send()
		return nil, err
	}
	sourceSettable, ok := newEntity.(sharedent.SourceSettable)

	if !ok {
		err := fmt.Errorf("error casting entity %v to source settable", entity)
		logging.Log().Error().
			Err(err).
			Interface("entity", entity).
			Str("dtype", string(dtype)).
			Str("assetClass", string(assetClass)).
			Str("symbol", symbol).
			Send()
		return nil, err
	} else {
		sourceSettable.SetSource("alpaca")
	}

	// Set exchange
	exchangeSettable, ok := newEntity.(sharedent.ExchangeSettable)
	if ok {
		switch assetClass {
		case types.Stock:
			exchangeSettable.SetExchange(alpaca.DEFAULT_EXCHANGE_STOCK)
		case types.Crypto:
			exchangeSettable.SetExchange(alpaca.DEFAULT_EXCHANGE_CRYPTO)
		}
	}

	topic := alpaca.NewStockStreamTopic(dtype, symbol).Generate()

	return alpaca.GenerateMessage(payloaderEntity, dtype, topic), nil

}

func handleAlpacaStreamGetRequest(req requests.StreamRequest, assetClass types.AssetClass) types.StreamResponse {
	return provider.NewStreamResponseAssetClass(
		types.Success,
		"Successfully retrieved streams",
		nil, assetClass)
}

// Handle a stream request for Alpaca
func HandleAlpacaStreamRequest(req requests.StreamRequest) types.StreamResponse {
	var response types.StreamResponse
	// Validate symbols

	for _, symbol := range req.GetSymbols() {
		if !alpaca.IsSymbolValid(symbol, req.GetAssetClass()) {
			return provider.NewStreamError(
				fmt.Errorf("symbol %s not valid for asset type %s", symbol, req.GetAssetClass()),
			)
		}
	}

	switch req.GetAssetClass() {
	case types.Stock:
		response = handleAlpacaStockStreamRequest(req)
	case types.Crypto:
		response = handleAlpacaCryptoStreamRequest(req)
	case types.News:
		response = handleAlpacaNewsStreamRequest(req)
	default:
		response = provider.NewStreamError(
			errors.New("invalid asset type"),
		)
	}
	return response
}
