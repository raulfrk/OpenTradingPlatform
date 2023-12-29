package types

type Component string
type Functionality string
type Source string
type AssetClass string
type DataType string
type OpStatus string

const (
	DataProvider  Component     = "dataprovider"
	DataStorage   Component     = "datastorage"
	Command       Functionality = "command"
	Stream        Functionality = "stream"
	Data          Functionality = "data"
	Logging       Functionality = "logging"
	Alpaca        Source        = "alpaca"
	Internal      Source        = "internal"
	Crypto        AssetClass    = "crypto"
	Stock         AssetClass    = "stock"
	News          AssetClass    = "news"
	NewsSentiment DataType      = "news-sentiment"
	Log           DataType      = "log"
	Bar           DataType      = "bar"
	LULD          DataType      = "luld"
	Status        DataType      = "status"
	Orderbook     DataType      = "orderbook"
	DailyBars     DataType      = "daily-bars"
	Quotes        DataType      = "quotes"
	Trades        DataType      = "trades"
	UpdatedBars   DataType      = "updated-bars"
	RawText       DataType      = "raw-text"
	Success       OpStatus      = "success"
	Failure       OpStatus      = "failure"
)

func GetAssetClassMap() map[string]AssetClass {
	return map[string]AssetClass{
		"stock":  Stock,
		"crypto": Crypto,
		"news":   News,
	}
}
