package alpaca

import (
	"os"

	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"

	astream "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

// Utility to verify whether a symbol is valid for a given asset class
func IsSymbolValid(symbol string, assetClass types.AssetClass) bool {
	client := astream.NewClient(astream.ClientOpts{
		APIKey:    os.Getenv("ALPACA_KEY"),
		APISecret: os.Getenv("ALPACA_SECRET"),
		BaseURL:   "https://paper-api.alpaca.markets",
	})

	asset, err := client.GetAsset(symbol)

	if err != nil {
		logging.Log().Error().
			Err(err).
			Msg("getting asset during symbol validation")
		return false
	}

	if asset.Class == "us_equity" && (assetClass == types.Stock || assetClass == types.News) {
		return true
	}
	if asset.Class == "crypto" && (assetClass == types.Crypto || assetClass == types.News) {
		return true
	}

	return false
}
