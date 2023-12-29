package data

import (
	"time"
	"tradingplatform/shared/entities"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
)

type LULD struct {
	Symbol         string
	LimitUpPrice   float64
	LimitDownPrice float64
	Indicator      string
	Timestamp      time.Time `gorm:"index"`
	Tape           string
	Fingerprint    string
	Source         string
	AssetClass     string `gorm:"not null"`
}

func LULDFromEntity(entity *entities.LULD) LULD {
	return LULD{
		Symbol:         entity.Symbol,
		LimitUpPrice:   entity.LimitUpPrice,
		LimitDownPrice: entity.LimitDownPrice,
		Indicator:      entity.Indicator,
		Timestamp:      time.Unix(entity.Timestamp, 0),
		Tape:           entity.Tape,
		Fingerprint:    entity.Fingerprint,
		Source:         entity.Source,
		AssetClass:     entity.AssetClass,
	}
}

func LULDToEntity(luld LULD) *entities.LULD {
	return &entities.LULD{
		Symbol:         luld.Symbol,
		LimitUpPrice:   luld.LimitUpPrice,
		LimitDownPrice: luld.LimitDownPrice,
		Indicator:      luld.Indicator,
		Timestamp:      luld.Timestamp.Unix(),
		Tape:           luld.Tape,
		Fingerprint:    luld.Fingerprint,
		Source:         luld.Source,
		AssetClass:     luld.AssetClass,
	}
}

func LULDsToEntities(lulds []LULD) []*entities.LULD {
	entities := make([]*entities.LULD, len(lulds))
	for i, luld := range lulds {
		entities[i] = LULDToEntity(luld)
	}
	return entities
}

func LULDsFromEntities(entities []*entities.LULD) []LULD {
	lulds := make([]LULD, len(entities))
	for i, entity := range entities {
		lulds[i] = LULDFromEntity(entity)
	}
	return lulds
}

func GetLULDFromRequest(symbol string, req requests.DataRequest) ([]*entities.LULD, error) {
	return GetLULD(string(req.GetSource()),
		symbol,
		string(req.AssetClass),
		req.GetStartTime(),
		req.GetEndTime()), nil
}

func InsertLULD(luld *entities.LULD) {
	dbLULD := LULDFromEntity(luld)
	tx := DB.Create(dbLULD)
	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			RawJSON("luld", entities.GenerateJson(luld)).
			Msg("inserting LULD")
	}
}

func GetLULD(source string,
	symbol string,
	assetClass string,
	startTime int64,
	endTime int64) []*entities.LULD {

	var lulds []LULD
	tx := DB.Where("source = ? AND symbol = ? AND timestamp >= ? AND timestamp <= ? AND asset_class = ?",
		source,
		symbol,
		time.Unix(startTime, 0),
		time.Unix(endTime, 0),
		assetClass,
	).Find(&lulds)

	if tx.Error != nil {
		logging.Log().Error().
			Err(tx.Error).
			Str("source", source).
			Str("symbol", symbol).
			Int64("startTime", startTime).
			Int64("endTime", endTime).
			Str("assetClass", assetClass).
			Msg("getting LULD from database")
		return nil
	}
	return LULDsToEntities(lulds)
}
