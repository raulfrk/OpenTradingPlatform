package stream

import (
	"errors"
	"fmt"
	"sync"
	"tradingplatform/dataprovider/data"
	"tradingplatform/dataprovider/provider"
	"tradingplatform/dataprovider/provider/alpaca"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
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

// handleOnStreamData is a function that handles streaming data for a given entity.
// It maps the entity with the provided symbol, sets the source and exchange based on the asset class,
// generates the topic based on the asset class and data type, and returns a generated message.
// If any error occurs during the process, an error is logged and returned.
//
// Parameters:
// - entity: The entity to be handled.
// - assetClass: The asset class of the entity.
// - dtype: The data type of the entity.
// - symbol: The symbol associated with the entity.
//
// Returns:
// - *sharedent.Message: The generated message.
// - error: An error if any occurred during the process.
func handleOnStreamData[T any,
	V sharedent.Payloader](entity T,
	assetClass types.AssetClass,
	dtype types.DataType,
	symbol string) (*sharedent.Message, error) {

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
	var topic string
	switch assetClass {
	case types.Crypto:
		topic = alpaca.NewCryptoStreamTopic(dtype, symbol).Generate()
	case types.News:
		topic = alpaca.NewNewsStreamTopic(dtype, symbol).Generate()
	default:
		topic = alpaca.NewStockStreamTopic(dtype, symbol).Generate()
	}

	return sharedent.GenerateMessage(payloaderEntity, dtype, topic), nil

}

// Provide a response with the active streams
func handleAlpacaStreamGetRequest(req requests.StreamRequest, assetClass types.AssetClass) types.StreamResponse {
	streams := data.GetDataProviderStreamsAssetClass(assetClass)
	dtypes := []types.DataType{}
	symbols := []string{}
	for _, stream := range streams {
		dtypes = append(dtypes, stream.DataType)
		symbols = append(symbols, stream.Symbol)
	}
	return provider.NewStreamResponseAssetClass(
		types.Success,
		"Successfully retrieved streams",
		alpaca.GenerateJSONStreamTopicDict(types.Crypto, dtypes, symbols),
		nil, assetClass)
}

// Handle a stream request for Alpaca and delegate to the appropriate handler based on asset class
func HandleAlpacaStreamRequest(req requests.StreamRequest) types.StreamResponse {
	var response types.StreamResponse
	// Validate symbols
	for _, symbol := range req.GetSymbol() {
		if !alpaca.IsSymbolValid(symbol, req.GetAssetClass()) {
			return provider.NewStreamError(
				fmt.Errorf("symbol %s not valid for asset class %s", symbol, req.GetAssetClass()),
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
			errors.New("invalid asset class"),
		)
	}
	return response
}
