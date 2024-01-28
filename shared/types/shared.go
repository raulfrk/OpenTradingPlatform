package types

type Component string
type Functionality string
type Source string
type AssetClass string
type DataType string
type OpStatus string
type SentimentAnalysisProcess string
type LLMProvider string

const (
	DataProvider      Component     = "dataprovider"
	DataStorage       Component     = "datastorage"
	SentimentAnalyzer Component     = "sentiment-analyzer"
	Command           Functionality = "command"
	Stream            Functionality = "stream"
	Data              Functionality = "data"
	Logging           Functionality = "logging"
	Alpaca            Source        = "alpaca"
	Internal          Source        = "internal"
	Crypto            AssetClass    = "crypto"
	Stock             AssetClass    = "stock"
	News              AssetClass    = "news"

	Log               DataType    = "log"
	Bar               DataType    = "bar"
	LULD              DataType    = "luld"
	Status            DataType    = "status"
	Orderbook         DataType    = "orderbook"
	DailyBars         DataType    = "daily-bars"
	Quotes            DataType    = "quotes"
	Trades            DataType    = "trades"
	UpdatedBars       DataType    = "updated-bars"
	RawText           DataType    = "raw-text"
	Sentiment         DataType    = "sentiment"
	NewsWithSentiment DataType    = "news-with-sentiment"
	Success           OpStatus    = "success"
	Failure           OpStatus    = "failure"
	Ollama            LLMProvider = "ollama"
)

const (
	Plain    SentimentAnalysisProcess = "plain"
	Semantic SentimentAnalysisProcess = "semantic"
)

func GetAssetClassMap() map[string]AssetClass {
	return map[string]AssetClass{
		"stock":  Stock,
		"crypto": Crypto,
		"news":   News,
	}
}

func GetSentimentAnalysisProcessMap() map[string]SentimentAnalysisProcess {
	return map[string]SentimentAnalysisProcess{
		"plain":    Plain,
		"semantic": Semantic,
	}
}
