package data

import (
	"time"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"

	"gorm.io/gorm/clause"
)

type Trade struct {
	ID          int64
	Symbol      string
	Exchange    string
	Price       float64
	Size        float64
	Timestamp   time.Time `gorm:"index"`
	TakerSide   string
	Conditions  []TradeCondition `gorm:"foreignKey:TradeFingerprint"`
	Tape        string
	Fingerprint string `gorm:"primaryKey"`
	Update      string
	Source      string
	AssetClass  string `gorm:"not null"`
}

// TODO: rething this
type TradeCondition struct {
	Condition        string `gorm:"primaryKey"`
	TradeFingerprint string `gorm:"primaryKey"`
}

func TradeFromEntity(entity *entities.Trade) Trade {
	conditions := make([]TradeCondition, len(entity.Conditions))
	for i, c := range entity.Conditions {
		conditions[i] = TradeCondition{Condition: c}
	}
	return Trade{
		Symbol:      entity.Symbol,
		Exchange:    entity.Exchange,
		Price:       entity.Price,
		Size:        entity.Size,
		Timestamp:   time.Unix(entity.Timestamp, 0),
		TakerSide:   entity.TakerSide,
		Conditions:  conditions,
		Tape:        entity.Tape,
		Fingerprint: entity.Fingerprint,
		Update:      entity.Update,
		Source:      entity.Source,
		AssetClass:  entity.AssetClass,
	}
}

func TradeToEntity(trade Trade) *entities.Trade {
	conditions := make([]string, len(trade.Conditions))
	for i, c := range trade.Conditions {
		conditions[i] = c.Condition
	}
	return &entities.Trade{
		Symbol:      trade.Symbol,
		Exchange:    trade.Exchange,
		Price:       trade.Price,
		Size:        trade.Size,
		Timestamp:   trade.Timestamp.Unix(),
		TakerSide:   trade.TakerSide,
		Conditions:  conditions,
		Tape:        trade.Tape,
		Fingerprint: trade.Fingerprint,
		Update:      trade.Update,
		Source:      trade.Source,
		AssetClass:  trade.AssetClass,
	}
}

func TradesFromEntities(entities []*entities.Trade) []Trade {
	trades := make([]Trade, len(entities))
	for i, entity := range entities {
		trades[i] = TradeFromEntity(entity)
	}
	return trades
}

func TradesToEntities(trades []Trade) []*entities.Trade {
	entities := make([]*entities.Trade, len(trades))
	for i, trade := range trades {
		entities[i] = TradeToEntity(trade)
	}
	return entities
}

func InsertTrade(trade *entities.Trade) {
	dbTrade := TradeFromEntity(trade)
	tx := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(dbTrade)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			RawJSON("trade", entities.GenerateJson(trade)).
			Msg("inserting trade")
	}
}

func GetTradesFromRequest(symbol string, req requests.DataRequest) ([]*entities.Trade, error) {
	return GetTrades(string(req.GetSource()),
		symbol,
		string(req.AssetClass),
		req.GetStartTime(),
		req.GetEndTime()), nil
}

func GetTrades(source, symbol, assetClass string, startTime, endTime int64) []*entities.Trade {
	var trades []Trade
	DB.Preload("Conditions").Where("source = ? AND symbol = ? AND asset_class = ? AND timestamp >= ? AND timestamp <= ?",
		source,
		symbol,
		assetClass,
		time.Unix(startTime, 0),
		time.Unix(endTime, 0)).Order("Timestamp").Find(&trades)
	return TradesToEntities(trades)
}
