package data

import (
	"time"

	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"

	"gorm.io/gorm/clause"
)

type Bar struct {
	Symbol      string
	Exchange    string
	Source      string
	Open        float64
	High        float64
	Low         float64
	Close       float64
	Volume      float64
	VWAP        float64
	Timestamp   time.Time `gorm:"index"`
	TradeCount  uint64
	Fingerprint string `gorm:"primaryKey"`
	AssetClass  string `gorm:"not null"`
	Timeframe   string `gorm:"not null"`
}

func BarFromEntity(entity *entities.Bar) Bar {
	return Bar{
		Symbol:      entity.Symbol,
		Exchange:    entity.Exchange,
		Source:      entity.Source,
		Open:        entity.Open,
		High:        entity.High,
		Low:         entity.Low,
		Close:       entity.Close,
		Volume:      entity.Volume,
		VWAP:        entity.VWAP,
		Timestamp:   time.Unix(entity.Timestamp, 0),
		TradeCount:  entity.TradeCount,
		Fingerprint: entity.Fingerprint,
		AssetClass:  entity.AssetClass,
		Timeframe:   entity.Timeframe,
	}
}

func BarToEntity(bar Bar) *entities.Bar {
	return &entities.Bar{
		Symbol:      bar.Symbol,
		Exchange:    bar.Exchange,
		Source:      bar.Source,
		Open:        bar.Open,
		High:        bar.High,
		Low:         bar.Low,
		Close:       bar.Close,
		Volume:      bar.Volume,
		VWAP:        bar.VWAP,
		Timestamp:   bar.Timestamp.Unix(),
		TradeCount:  bar.TradeCount,
		Fingerprint: bar.Fingerprint,
		AssetClass:  bar.AssetClass,
		Timeframe:   bar.Timeframe,
	}
}

func BarsFromEntities(entities []*entities.Bar) []Bar {
	bars := make([]Bar, len(entities))
	for i, entity := range entities {
		bars[i] = BarFromEntity(entity)
	}
	return bars
}

func BarsToEntities(bars []Bar) []*entities.Bar {
	entities := make([]*entities.Bar, len(bars))
	for i, bar := range bars {
		entities[i] = BarToEntity(bar)
	}
	return entities
}

func InsertBar(bar *entities.Bar) {
	dbBar := BarFromEntity(bar)
	tx := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(dbBar)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			RawJSON("bar", entities.GenerateJson(bar)).
			Msg("inserting bar")
	}
}

func GetBarsFromRequest(symbol string, req requests.DataRequest) ([]*entities.Bar, error) {
	return GetBars(string(req.GetSource()),
		symbol,
		string(req.AssetClass),
		req.GetStartTime(),
		req.GetEndTime(),
		string(req.TimeFrame)), nil
}

func GetBars(source string,
	symbol string,
	assetClass string,
	startTime int64,
	endTime int64,
	timeframe string) []*entities.Bar {
	var bars []Bar
	tx := DB.Where("source = ? AND symbol = ? AND timestamp >= ? AND timestamp <= ? AND asset_class = ? AND timeframe = ?",
		source,
		symbol,
		time.Unix(startTime, 0),
		time.Unix(endTime, 0),
		assetClass,
		timeframe).Order("timestamp").Find(&bars)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Str("source", source).
			Str("symbol", symbol).
			Int64("startTime", startTime).
			Int64("endTime", endTime).
			Str("assetClass", assetClass).
			Str("timeframe", timeframe).
			Msg("getting bars")
	}

	return BarsToEntities(bars)
}
