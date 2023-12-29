package data

import (
	"time"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"

	"gorm.io/gorm/clause"
)

// TODO: Add more restrictions to the table
type DailyBar struct {
	Bar
}

func DailyBarFromEntity(entity *entities.Bar) DailyBar {
	return DailyBar{
		Bar: Bar{
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
		},
	}
}

func DailyBarToEntity(bar DailyBar) *entities.Bar {
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

func DailyBarsFromEntities(entities []*entities.Bar) []DailyBar {
	bars := make([]DailyBar, len(entities))
	for i, entity := range entities {
		bars[i] = DailyBarFromEntity(entity)
	}
	return bars
}

func DailyBarsToEntities(bars []DailyBar) []*entities.Bar {
	entities := make([]*entities.Bar, len(bars))
	for i, bar := range bars {
		entities[i] = BarToEntity(bar.Bar)
	}
	return entities
}

func InsertDailyBar(bar *entities.Bar) {
	dbBar := DailyBarFromEntity(bar)
	tx := DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(dbBar)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			RawJSON("daily-bar", entities.GenerateJson(bar)).
			Msg("inserting daily-bar")
	}
}

func GetDailyBarsFromRequest(symbol string, req requests.DataRequest) ([]*entities.Bar, error) {
	return GetDailyBars(string(req.GetSource()),
		symbol,
		string(req.AssetClass),
		req.GetStartTime(),
		req.GetEndTime(),
		string(req.TimeFrame)), nil
}

func GetDailyBars(source string,
	symbol string,
	assetClass string,
	startTime int64,
	endTime int64,
	timeframe string) []*entities.Bar {
	var dBars []DailyBar
	tx := DB.Where("source = ? AND symbol = ? AND timestamp >= ? AND timestamp <= ? AND asset_class = ? AND timeframe = ?",
		source,
		symbol,
		time.Unix(startTime, 0),
		time.Unix(endTime, 0),
		assetClass,
		timeframe).Order("timestamp").Find(&dBars)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Msg("getting bars")
	}

	return DailyBarsToEntities(dBars)
}
