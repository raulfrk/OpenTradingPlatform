package data

import (
	"fmt"
	"tradingplatform/dataprovider/provider/alpaca"
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

// Given a getter function that takes a symbol and a request of type T (e.g. GetCryptoBarsRequest) and
// returns a list of elements G (e.g. marketdata.Bars); a symbol; a request T, and a data type.
// Map the result of fun to a entity of type V (e.g. sharedent.Bar) (set its fignerprint).
// Return a list of messages and a data response
func handleDataFetch[T any,
	G any,
	V sharedent.Fingerprintable](
	fun func(string, T) ([]G, error),
	symbol string,
	req T,
	dtype types.DataType, assetClass types.AssetClass,
	timeFrame types.TimeFrame) (*[]*sharedent.Message, types.DataResponse) {

	// Use the getter function to get the data
	result, err := fun(symbol, req)
	if err != nil {
		logging.Log().Error().
			Err(err).
			Str("symbol", symbol).
			Str("dtype", string(dtype)).
			Str("assetClass", string(assetClass)).
			Str("timeFrame", string(timeFrame)).
			Msg("getting data from alpaca")
		return nil, types.NewDataError(err)
	}
	var responseTopic = ""
	var queueID = ""

	var messages []*sharedent.Message
	for _, element := range result {
		newEntity := alpaca.MapEntityWithReturnEntity(element, symbol)
		entity, ok := newEntity.(V)
		if !ok {
			err := fmt.Errorf("error casting to fingerprintable")
			logging.Log().Error().
				Err(err).
				Interface("entity", newEntity).
				Send()
			return nil, types.NewDataError(
				err,
			)
		}
		// Set the source of the entity
		sourceSettable, ok := newEntity.(sharedent.SourceSettable)
		if ok {
			sourceSettable.SetSource("alpaca")
		} else {
			err := fmt.Errorf("error casting to source settable")
			logging.Log().Error().
				Err(err).
				Interface("entity", newEntity).
				Send()
			return nil, types.NewDataError(
				err,
			)
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

		// Set timeframe if possible
		timeframeSettable, ok := newEntity.(sharedent.TimeframeSettable)
		if ok {
			timeframeSettable.SetTimeframe(string(timeFrame))
		}
		queueID = entity.GetFingerprint()

		payloader, ok := newEntity.(sharedent.Payloader)

		if !ok {
			err := fmt.Errorf("error casting to payloader")
			logging.Log().Error().
				Err(err).
				Interface("entity", newEntity).
				Send()
			return nil, types.NewDataError(
				err,
			)
		}

		messages = append(messages, sharedent.GenerateMessage(payloader, dtype,
			responseTopic))
	}

	if dtype == types.Bar {
		responseTopic = alpaca.NewBarDataTopic(assetClass, timeFrame, symbol,
			queueID, len(messages)).Generate()
	} else {
		responseTopic = alpaca.NewDataTopic(assetClass, dtype, symbol, queueID,
			len(messages)).Generate()
	}
	for _, message := range messages {
		message.Topic = responseTopic
	}
	return &messages, types.NewDataResponse(
		types.Success,
		"Successfully added data",
		nil,
		responseTopic,
	)
}

// Delegate a data request to the appropriate handler based on asset class
func HandleAlpacaDataRequest(dataRequest requests.DataRequest) types.DataResponse {

	switch dataRequest.GetAssetClass() {
	case types.Crypto:
		return handleAlpacaCryptoDataRequest(dataRequest)
	case types.Stock:
		return handleAlpacaStockDataRequest(dataRequest)
	case types.News:
		return handleAlpacaNewsDataRequest(dataRequest)
	default:
		return types.NewDataError(
			fmt.Errorf("invalid asset type %s", dataRequest.GetAssetClass()),
		)
	}
}
