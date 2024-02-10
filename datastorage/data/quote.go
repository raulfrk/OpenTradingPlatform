package data

import (
	"time"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"

	"gorm.io/gorm/clause"
)

type Quote struct {
	Symbol      string
	BidExchange string
	Exchange    string
	BidPrice    float64
	BidSize     float64
	AskExchange string
	AskPrice    float64
	AskSize     float64
	Timestamp   time.Time        `gorm:"index"`
	Conditions  []QuoteCondition `gorm:"foreignKey:TradeFingerprint"`
	Tape        string
	Fingerprint string `gorm:"primaryKey"`
	Source      string
	AssetClass  string `gorm:"not null"`
}

type QuoteCondition struct {
	Condition        string `gorm:"primaryKey"`
	TradeFingerprint string `gorm:"primaryKey"`
}

func QuoteFromEntity(entity *entities.Quote) Quote {
	conditions := make([]QuoteCondition, len(entity.Conditions))
	for i, c := range entity.Conditions {
		conditions[i] = QuoteCondition{Condition: c}
	}
	return Quote{
		Symbol:      entity.Symbol,
		BidExchange: entity.BidExchange,
		Exchange:    entity.Exchange,
		BidPrice:    entity.BidPrice,
		BidSize:     entity.BidSize,
		AskExchange: entity.AskExchange,
		AskPrice:    entity.AskPrice,
		AskSize:     entity.AskSize,
		Timestamp:   time.Unix(entity.Timestamp, 0),
		Conditions:  conditions,
		Tape:        entity.Tape,
		Fingerprint: entity.Fingerprint,
		Source:      entity.Source,
		AssetClass:  entity.AssetClass,
	}
}

func QuoteToEntity(quote Quote) *entities.Quote {
	conditions := make([]string, len(quote.Conditions))
	for i, c := range quote.Conditions {
		conditions[i] = c.Condition
	}
	return &entities.Quote{
		Symbol:      quote.Symbol,
		BidExchange: quote.BidExchange,
		Exchange:    quote.Exchange,
		BidPrice:    quote.BidPrice,
		BidSize:     quote.BidSize,
		AskExchange: quote.AskExchange,
		AskPrice:    quote.AskPrice,
		AskSize:     quote.AskSize,
		Timestamp:   quote.Timestamp.Unix(),
		Conditions:  conditions,
		Tape:        quote.Tape,
		Fingerprint: quote.Fingerprint,
		Source:      quote.Source,
		AssetClass:  quote.AssetClass,
	}
}

func QuotesFromEntities(entities []*entities.Quote) []Quote {
	quote := make([]Quote, len(entities))
	for i, entity := range entities {
		quote[i] = QuoteFromEntity(entity)
	}
	return quote
}

func QuotesToEntities(quotes []Quote) []*entities.Quote {
	entities := make([]*entities.Quote, len(quotes))
	for i, quote := range quotes {
		entities[i] = QuoteToEntity(quote)
	}
	return entities
}

func InsertQuote(quote *entities.Quote) {
	dbQuote := QuoteFromEntity(quote)
	tx := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(dbQuote)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			RawJSON("quote", entities.GenerateJson(quote)).
			Msg("inserting quote")
	}
}

func GetQuoteFromRequest(symbol string, req requests.DataRequest) ([]*entities.Quote, error) {
	return GetQuote(string(req.GetSource()),
		symbol,
		string(req.AssetClass),
		req.GetStartTime(),
		req.GetEndTime()), nil
}

func GetQuote(source, symbol, assetClass string, startTime, endTime int64) []*entities.Quote {
	var quotes []Quote
	tx := DB.Preload("Conditions").Where("source = ? AND symbol = ? AND asset_class = ? AND timestamp >= ? AND timestamp <= ?",
		source,
		symbol,
		assetClass,
		time.Unix(startTime, 0),
		time.Unix(endTime, 0)).Order("Timestamp")
	quotes = PaginateRequest(tx, Quote{})
	return QuotesToEntities(quotes)
}
