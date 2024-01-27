package alpaca

import (
	sharedent "tradingplatform/shared/entities"
	"tradingplatform/shared/types"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	astream "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
)

// MapEntity maps a given struct instance to a shared entity and wraps it in a message
func MapEntity(i interface{}, topic string) *sharedent.Message {
	switch v := i.(type) {
	case astream.CryptoBar:
		entity, t := MapCryptoBar(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.CryptoOrderbook:
		entity, t := MapCryptoOrderbook(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.CryptoQuote:
		entity, t := MapCryptoQuote(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.CryptoTrade:
		entity, t := MapCryptoTrade(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.Bar:
		entity, t := MapStockBar(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.Trade:
		entity, t := MapStockTrade(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.Quote:
		entity, t := MapStockQuote(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.LULD:
		entity, t := MapStockLULD(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.TradingStatus:
		entity, t := MapStockTradingStatus(v)
		return sharedent.GenerateMessage(entity, t, topic)
	case astream.News:
		entity, t := MapNews(v)
		return sharedent.GenerateMessage(entity, t, topic)
	}
	return nil
}

// MapEntityWithReturnEntity maps a given struct instance to a shared entity and returns it
func MapEntityWithReturnEntity(i interface{}, symbol string) interface{} {
	switch v := i.(type) {
	case astream.CryptoBar:
		entity, _ := MapCryptoBar(v)
		return entity
	case marketdata.CryptoBar:
		entity, _ := MapMarketCryptoBar(v, symbol)
		return entity
	case astream.CryptoOrderbook:
		entity, _ := MapCryptoOrderbook(v)
		return entity
	case marketdata.CryptoTrade:
		entity, _ := MapMarketCryptoTrade(v, symbol)
		return entity
	case marketdata.CryptoQuote:
		entity, _ := MapMarketCryptoQuote(v, symbol)
		return entity
	case marketdata.Bar:
		entity, _ := MapMarketStockBar(v, symbol)
		return entity
	case marketdata.Trade:
		entity, _ := MapMarketStockTrade(v, symbol)
		return entity
	case marketdata.Quote:
		entity, _ := MapMarketStockQuote(v, symbol)
		return entity
	case marketdata.News:
		entity, _ := MapMarketNews(v)
		return entity
	case astream.CryptoQuote:
		entity, _ := MapCryptoQuote(v)
		return entity
	case astream.CryptoTrade:
		entity, _ := MapCryptoTrade(v)
		return entity
	case astream.Bar:
		entity, _ := MapStockBar(v)
		return entity
	case astream.Trade:
		entity, _ := MapStockTrade(v)
		return entity
	case astream.Quote:
		entity, _ := MapStockQuote(v)
		return entity
	case astream.LULD:
		entity, _ := MapStockLULD(v)
		return entity
	case astream.TradingStatus:
		entity, _ := MapStockTradingStatus(v)
		return entity
	case astream.News:
		entity, _ := MapNews(v)
		return entity
	}
	return nil
}

func MapCryptoBar(cb astream.CryptoBar) (*sharedent.Bar, types.DataType) {
	newBar := sharedent.Bar{
		Symbol:     cb.Symbol,
		Open:       cb.Open,
		High:       cb.High,
		Low:        cb.Low,
		Close:      cb.Close,
		Volume:     cb.Volume,
		Timestamp:  cb.Timestamp.Unix(),
		Exchange:   cb.Exchange,
		VWAP:       cb.VWAP,
		AssetClass: string(types.Crypto),
	}

	newBar.SetFingerprint()
	return &newBar, types.Bar
}

func MapMarketCryptoBar(cb marketdata.CryptoBar, symbol string) (*sharedent.Bar, types.DataType) {
	newBar := sharedent.Bar{
		Symbol:     symbol,
		Open:       cb.Open,
		High:       cb.High,
		Low:        cb.Low,
		Close:      cb.Close,
		Volume:     cb.Volume,
		Timestamp:  cb.Timestamp.Unix(),
		VWAP:       cb.VWAP,
		AssetClass: string(types.Crypto),
	}

	newBar.SetFingerprint()
	return &newBar, types.Bar
}

func MapCryptoOrderbook(ob astream.CryptoOrderbook) (*sharedent.Orderbook, types.DataType) {
	newOrderbook := sharedent.Orderbook{
		Symbol:     ob.Symbol,
		Exchange:   ob.Exchange,
		Timestamp:  ob.Timestamp.Unix(),
		AssetClass: string(types.Crypto),
	}
	for _, ask := range ob.Asks {
		newOrderbook.Asks = append(newOrderbook.Asks, &sharedent.OrderbookEntry{
			Price: ask.Price,
			Size:  ask.Size,
		})
	}

	for _, bid := range ob.Bids {
		newOrderbook.Asks = append(newOrderbook.Bids, &sharedent.OrderbookEntry{
			Price: bid.Price,
			Size:  bid.Size,
		})
	}

	newOrderbook.SetFingerprint()
	return &newOrderbook, types.Orderbook
}

func MapCryptoQuote(q astream.CryptoQuote) (*sharedent.Quote, types.DataType) {
	newQuote := sharedent.Quote{
		Symbol:     q.Symbol,
		Exchange:   q.Exchange,
		Timestamp:  q.Timestamp.Unix(),
		AskPrice:   q.AskPrice,
		AskSize:    q.AskSize,
		BidPrice:   q.BidPrice,
		BidSize:    q.BidSize,
		AssetClass: string(types.Crypto),
	}

	newQuote.SetFingerprint()
	return &newQuote, types.Quotes
}

func MapMarketCryptoQuote(q marketdata.CryptoQuote, symbol string) (*sharedent.Quote, types.DataType) {
	newQuote := sharedent.Quote{
		Symbol:     symbol,
		Timestamp:  q.Timestamp.Unix(),
		AskPrice:   q.AskPrice,
		AskSize:    q.AskSize,
		BidPrice:   q.BidPrice,
		BidSize:    q.BidSize,
		AssetClass: string(types.Crypto),
	}

	newQuote.SetFingerprint()
	return &newQuote, types.Quotes
}

func MapCryptoTrade(t astream.CryptoTrade) (*sharedent.Trade, types.DataType) {
	newTrade := sharedent.Trade{
		ID:         t.ID,
		Symbol:     t.Symbol,
		Exchange:   t.Exchange,
		Timestamp:  t.Timestamp.Unix(),
		Price:      t.Price,
		Size:       t.Size,
		TakerSide:  string(t.TakerSide),
		AssetClass: string(types.Crypto),
	}

	newTrade.SetFingerprint()
	return &newTrade, types.Trades
}

func MapMarketCryptoTrade(t marketdata.CryptoTrade, symbol string) (*sharedent.Trade, types.DataType) {
	newTrade := sharedent.Trade{
		ID:         t.ID,
		Symbol:     symbol,
		Timestamp:  t.Timestamp.Unix(),
		Price:      t.Price,
		Size:       t.Size,
		TakerSide:  string(t.TakerSide),
		AssetClass: string(types.Crypto),
	}

	newTrade.SetFingerprint()
	return &newTrade, types.Trades
}

func MapStockBar(b astream.Bar) (*sharedent.Bar, types.DataType) {
	newBar := sharedent.Bar{
		Symbol:     b.Symbol,
		Open:       b.Open,
		High:       b.High,
		Low:        b.Low,
		Close:      b.Close,
		Volume:     float64(b.Volume),
		Timestamp:  b.Timestamp.Unix(),
		VWAP:       b.VWAP,
		AssetClass: string(types.Stock),
	}

	newBar.SetFingerprint()
	return &newBar, types.Bar
}

func MapMarketStockBar(b marketdata.Bar, symbol string) (*sharedent.Bar, types.DataType) {
	newBar := sharedent.Bar{
		Symbol:     symbol,
		Open:       b.Open,
		High:       b.High,
		Low:        b.Low,
		Close:      b.Close,
		Volume:     float64(b.Volume),
		TradeCount: b.TradeCount,
		Timestamp:  b.Timestamp.Unix(),
		VWAP:       b.VWAP,
		AssetClass: string(types.Stock),
	}

	newBar.SetFingerprint()
	return &newBar, types.Bar
}

func MapStockTrade(t astream.Trade) (*sharedent.Trade, types.DataType) {
	newTrade := sharedent.Trade{
		ID:         t.ID,
		Symbol:     t.Symbol,
		Exchange:   t.Exchange,
		Price:      t.Price,
		Size:       float64(t.Size),
		Timestamp:  t.Timestamp.Unix(),
		Conditions: t.Conditions,
		Tape:       t.Tape,
		AssetClass: string(types.Stock),
	}

	newTrade.SetFingerprint()
	return &newTrade, types.Trades
}

func MapMarketStockTrade(t marketdata.Trade, symbol string) (*sharedent.Trade, types.DataType) {
	newTrade := sharedent.Trade{
		ID:         t.ID,
		Symbol:     symbol,
		Exchange:   t.Exchange,
		Update:     t.Update,
		Price:      t.Price,
		Size:       float64(t.Size),
		Timestamp:  t.Timestamp.Unix(),
		Conditions: t.Conditions,
		Tape:       t.Tape,
		AssetClass: string(types.Stock),
	}

	newTrade.SetFingerprint()
	return &newTrade, types.Trades
}

func MapStockQuote(q astream.Quote) (*sharedent.Quote, types.DataType) {
	newQuote := sharedent.Quote{
		Symbol:      q.Symbol,
		BidExchange: q.BidExchange,
		BidPrice:    q.BidPrice,
		BidSize:     float64(q.BidSize),
		AskExchange: q.AskExchange,
		AskPrice:    q.AskPrice,
		AskSize:     float64(q.AskSize),
		Timestamp:   q.Timestamp.Unix(),
		Conditions:  q.Conditions,
		Tape:        q.Tape,
		AssetClass:  string(types.Stock),
	}

	newQuote.SetFingerprint()
	return &newQuote, types.Quotes
}

func MapMarketStockQuote(q marketdata.Quote, symbol string) (*sharedent.Quote, types.DataType) {
	newQuote := sharedent.Quote{
		Symbol:      symbol,
		BidExchange: q.BidExchange,
		BidPrice:    q.BidPrice,
		BidSize:     float64(q.BidSize),
		AskExchange: q.AskExchange,
		AskPrice:    q.AskPrice,
		AskSize:     float64(q.AskSize),
		Timestamp:   q.Timestamp.Unix(),
		Conditions:  q.Conditions,
		Tape:        q.Tape,
		AssetClass:  string(types.Stock),
	}

	newQuote.SetFingerprint()
	return &newQuote, types.Quotes
}

func MapStockLULD(luld astream.LULD) (*sharedent.LULD, types.DataType) {
	newLULD := sharedent.LULD{
		Symbol:         luld.Symbol,
		LimitUpPrice:   luld.LimitUpPrice,
		LimitDownPrice: luld.LimitDownPrice,
		Indicator:      luld.Indicator,
		Timestamp:      luld.Timestamp.Unix(),
		Tape:           luld.Tape,
		AssetClass:     string(types.Stock),
	}

	newLULD.SetFingerprint()
	return &newLULD, types.LULD
}

func MapStockTradingStatus(status astream.TradingStatus) (*sharedent.TradingStatus, types.DataType) {
	newStatus := sharedent.TradingStatus{
		Symbol:     status.Symbol,
		StatusCode: status.StatusCode,
		StatusMsg:  status.StatusMsg,
		ReasonCode: status.ReasonCode,
		ReasonMsg:  status.ReasonMsg,
		Timestamp:  status.Timestamp.Unix(),
		Tape:       status.Tape,
		AssetClass: string(types.Stock),
	}

	newStatus.SetFingerprint()
	return &newStatus, types.Status
}

func MapNews(news astream.News) (*sharedent.News, types.DataType) {
	newNews := sharedent.News{
		Id:        int64(news.ID),
		Author:    news.Author,
		CreatedAt: news.CreatedAt.Unix(),
		UpdatedAt: news.UpdatedAt.Unix(),
		Headline:  news.Headline,
		Summary:   news.Summary,
		Content:   news.Content,
		URL:       news.URL,
		Symbols:   news.Symbols,
	}

	newNews.SetFingerprint()
	return &newNews, types.RawText
}

func MapMarketNews(news marketdata.News) (*sharedent.News, types.DataType) {
	newNews := sharedent.News{
		Id:        int64(news.ID),
		Author:    news.Author,
		CreatedAt: news.CreatedAt.Unix(),
		UpdatedAt: news.UpdatedAt.Unix(),
		Headline:  news.Headline,
		Summary:   news.Summary,
		Content:   news.Content,
		URL:       news.URL,
		Symbols:   news.Symbols,
	}

	newNews.SetFingerprint()
	return &newNews, types.RawText
}

func GetAlpacaTimeFrame(timeFrame types.TimeFrame) marketdata.TimeFrame {
	m := map[types.TimeFrame]marketdata.TimeFrame{
		types.OneMin:   marketdata.OneMin,
		types.OneHour:  marketdata.OneHour,
		types.OneDay:   marketdata.OneDay,
		types.OneWeek:  marketdata.OneWeek,
		types.OneMonth: marketdata.OneMonth,
	}

	return m[timeFrame]

}
